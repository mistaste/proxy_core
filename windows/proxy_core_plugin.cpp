#include "proxy_core_plugin.h"

#include <flutter/method_channel.h>
#include <flutter/plugin_registrar_windows.h>
#include <flutter/standard_method_codec.h>

#include <windows.h>
#include <shellapi.h>
#include <winsvc.h>

#include <atomic>
#include <memory>
#include <sstream>
#include <string>
#include <thread>
#include <vector>

#pragma comment(lib, "advapi32.lib")
#pragma comment(lib, "shell32.lib")

namespace proxy_core {

namespace {

// Named pipe exposed by guardexsvc.exe. Must match service_win.PipeName.
constexpr const char* kPipeName = R"(\\.\pipe\guardex-svc)";

std::atomic<uint64_t> g_next_id{1};

// JSON escape for control chars and quote/backslash. Good enough for the
// argument shapes we pass; we do not accept arbitrary strings from Dart
// that could contain raw newlines.
std::string EscapeJson(const std::string& s) {
  std::ostringstream o;
  for (char c : s) {
    switch (c) {
      case '"': o << "\\\""; break;
      case '\\': o << "\\\\"; break;
      case '\n': o << "\\n"; break;
      case '\r': o << "\\r"; break;
      case '\t': o << "\\t"; break;
      default:
        if (static_cast<unsigned char>(c) < 0x20) {
          char buf[8];
          snprintf(buf, sizeof(buf), "\\u%04x", c);
          o << buf;
        } else {
          o << c;
        }
    }
  }
  return o.str();
}

// SendRequest dials the pipe, writes one line of JSON, reads one line of
// JSON back. Returns the raw response line (including trailing newline
// stripped) on success, or an empty optional with err populated on
// failure.
struct RpcResult {
  bool ok = false;
  std::string body;   // raw JSON line from service
  std::string error;  // transport error, if any
};

RpcResult CallService(const std::string& method_json,
                      const std::string& params_json_or_empty) {
  // Each call opens its own pipe handle — no shared mutable state between
  // concurrent callers, so no mutex is needed here.
  HANDLE pipe = INVALID_HANDLE_VALUE;
  // WaitNamedPipe + CreateFile retry loop — 3 attempts, 500ms each.
  for (int i = 0; i < 3; ++i) {
    pipe = CreateFileA(kPipeName, GENERIC_READ | GENERIC_WRITE, 0, nullptr,
                       OPEN_EXISTING, 0, nullptr);
    if (pipe != INVALID_HANDLE_VALUE) break;
    if (GetLastError() != ERROR_PIPE_BUSY) break;
    WaitNamedPipeA(kPipeName, 500);
  }
  if (pipe == INVALID_HANDLE_VALUE) {
    RpcResult r;
    r.error = "cannot connect to guardex service (is GuardexVPN service running?)";
    return r;
  }

  uint64_t id = g_next_id.fetch_add(1);
  std::ostringstream req;
  req << "{\"id\":" << id << ",\"method\":\"" << method_json << "\"";
  if (!params_json_or_empty.empty()) {
    req << ",\"params\":" << params_json_or_empty;
  }
  req << "}\n";
  std::string req_s = req.str();

  DWORD written = 0;
  if (!WriteFile(pipe, req_s.data(), static_cast<DWORD>(req_s.size()),
                 &written, nullptr)) {
    CloseHandle(pipe);
    RpcResult r;
    r.error = "pipe write failed";
    return r;
  }

  // Read until newline. Responses are always one line.
  std::string resp;
  char buf[512];
  DWORD read = 0;
  while (true) {
    if (!ReadFile(pipe, buf, sizeof(buf), &read, nullptr) || read == 0) {
      break;
    }
    resp.append(buf, read);
    if (resp.find('\n') != std::string::npos) break;
  }
  CloseHandle(pipe);

  if (resp.empty()) {
    RpcResult r;
    r.error = "empty response from service";
    return r;
  }
  auto nl = resp.find('\n');
  if (nl != std::string::npos) resp.resize(nl);

  RpcResult r;
  r.ok = true;
  r.body = resp;
  return r;
}

// Minimal JSON response parser — avoids substring matches that could be
// fooled by string values containing "ok\":true" or similar literals.
// Only handles the flat {"id":N,"ok":bool,"error":"...","result":bool}
// shape returned by service_win.
struct ServiceResponse {
  bool ok = false;
  bool result_bool = false;
  std::string error;
};

static ServiceResponse ParseServiceResponse(const std::string& s) {
  ServiceResponse out;
  size_t i = 0;
  auto skipWS = [&]() {
    while (i < s.size() &&
           (s[i] == ' ' || s[i] == '\t' || s[i] == ',' || s[i] == '{' ||
            s[i] == '}'))
      ++i;
  };
  auto readString = [&]() -> std::string {
    if (i >= s.size() || s[i] != '"') return "";
    ++i;
    std::string val;
    while (i < s.size() && s[i] != '"') {
      if (s[i] == '\\') {
        ++i;
        if (i < s.size()) val += s[i];
      } else {
        val += s[i];
      }
      ++i;
    }
    if (i < s.size()) ++i;  // skip closing "
    return val;
  };

  while (i < s.size()) {
    skipWS();
    if (i >= s.size() || s[i] == '}') break;
    if (s[i] != '"') { ++i; continue; }
    std::string key = readString();
    // skip ':'
    while (i < s.size() && (s[i] == ' ' || s[i] == ':')) ++i;
    // read value
    if (i >= s.size()) break;
    if (s[i] == '"') {
      std::string val = readString();
      if (key == "error") out.error = val;
    } else if (s.size() - i >= 4 && s.substr(i, 4) == "true") {
      if (key == "ok") out.ok = true;
      if (key == "result") out.result_bool = true;
      i += 4;
    } else if (s.size() - i >= 5 && s.substr(i, 5) == "false") {
      i += 5;
    } else {
      // number or null — skip to next delimiter
      while (i < s.size() && s[i] != ',' && s[i] != '}') ++i;
    }
  }
  return out;
}

bool JsonOK(const std::string& body) {
  return ParseServiceResponse(body).ok;
}

std::string JsonError(const std::string& body) {
  return ParseServiceResponse(body).error;
}

bool JsonResultBool(const std::string& body) {
  return ParseServiceResponse(body).result_bool;
}

std::string GetStringArg(const flutter::EncodableMap* args,
                         const std::string& key,
                         const std::string& fallback = "") {
  auto it = args->find(flutter::EncodableValue(key));
  if (it == args->end()) return fallback;
  if (auto* s = std::get_if<std::string>(&it->second)) return *s;
  return fallback;
}

int GetIntArg(const flutter::EncodableMap* args,
              const std::string& key,
              int fallback = 0) {
  auto it = args->find(flutter::EncodableValue(key));
  if (it == args->end()) return fallback;
  if (auto* i = std::get_if<int32_t>(&it->second)) return *i;
  if (auto* i = std::get_if<int64_t>(&it->second))
    return static_cast<int>(*i);
  return fallback;
}

// Path to guardexsvc.exe sitting next to the running host .exe.
std::wstring GuardexSvcPath() {
  wchar_t buf[MAX_PATH];
  DWORD n = GetModuleFileNameW(nullptr, buf, MAX_PATH);
  if (n == 0 || n == MAX_PATH) return L"";
  std::wstring path(buf, n);
  auto slash = path.find_last_of(L"\\/");
  if (slash == std::wstring::npos) return L"";
  return path.substr(0, slash + 1) + L"guardexsvc.exe";
}

struct ServiceStatus {
  bool installed = false;
  bool running = false;
};

ServiceStatus QueryGuardexService() {
  ServiceStatus st;
  SC_HANDLE scm = OpenSCManagerW(nullptr, nullptr, SC_MANAGER_CONNECT);
  if (!scm) return st;
  SC_HANDLE svc = OpenServiceW(scm, L"GuardexVPN", SERVICE_QUERY_STATUS);
  if (!svc) {
    CloseServiceHandle(scm);
    return st;
  }
  st.installed = true;
  SERVICE_STATUS s{};
  if (QueryServiceStatus(svc, &s)) {
    st.running = (s.dwCurrentState == SERVICE_RUNNING ||
                  s.dwCurrentState == SERVICE_START_PENDING);
  }
  CloseServiceHandle(svc);
  CloseServiceHandle(scm);
  return st;
}

// RunElevated invokes guardexsvc.exe with the given verb/argument via
// ShellExecuteExW + "runas" so the user sees exactly one UAC prompt.
// Waits for the child to exit and returns its exit code, or -1 if the
// elevation dialog was cancelled.
int RunElevated(const std::wstring& args) {
  std::wstring exe = GuardexSvcPath();
  if (exe.empty()) return -1;

  SHELLEXECUTEINFOW sei{};
  sei.cbSize = sizeof(sei);
  sei.fMask = SEE_MASK_NOCLOSEPROCESS | SEE_MASK_NO_CONSOLE;
  sei.lpVerb = L"runas";
  sei.lpFile = exe.c_str();
  sei.lpParameters = args.c_str();
  sei.nShow = SW_HIDE;

  if (!ShellExecuteExW(&sei) || !sei.hProcess) {
    return -1;
  }
  DWORD wait = WaitForSingleObject(sei.hProcess, 30000);
  DWORD exitCode = 1;
  if (wait == WAIT_TIMEOUT) {
    TerminateProcess(sei.hProcess, 1);
    CloseHandle(sei.hProcess);
    return -2;  // caller maps this to a timeout error
  }
  GetExitCodeProcess(sei.hProcess, &exitCode);
  CloseHandle(sei.hProcess);
  return static_cast<int>(exitCode);
}

// WaitForService polls the SCM until the service is RUNNING or a timeout
// elapses. Needed because ShellExecuteEx returns as soon as the child
// starts the service, not once the service reaches RUNNING state.
bool WaitForServiceRunning(int timeout_ms) {
  const int step_ms = 200;
  for (int waited = 0; waited <= timeout_ms; waited += step_ms) {
    if (QueryGuardexService().running) return true;
    Sleep(step_ms);
  }
  return false;
}

}  // namespace

// static
void ProxyCorePlugin::RegisterWithRegistrar(
    flutter::PluginRegistrarWindows* registrar) {
  auto channel = std::make_unique<
      flutter::MethodChannel<flutter::EncodableValue>>(
      registrar->messenger(), "proxy_core/vpn",
      &flutter::StandardMethodCodec::GetInstance());

  auto plugin = std::make_unique<ProxyCorePlugin>();

  channel->SetMethodCallHandler(
      [plugin_pointer = plugin.get()](const auto& call, auto result) {
        plugin_pointer->HandleMethodCall(call, std::move(result));
      });

  registrar->AddPlugin(std::move(plugin));
}

ProxyCorePlugin::ProxyCorePlugin() {}
ProxyCorePlugin::~ProxyCorePlugin() {
  alive_->store(false);
}

void ProxyCorePlugin::HandleMethodCall(
    const flutter::MethodCall<flutter::EncodableValue>& method_call,
    std::unique_ptr<flutter::MethodResult<flutter::EncodableValue>> result) {
  const std::string& method = method_call.method_name();
  const auto* args =
      std::get_if<flutter::EncodableMap>(method_call.arguments());

  if (method == "prepare") {
    result->Success(flutter::EncodableValue(nullptr));
    return;
  }

  if (method == "serviceStatus") {
    ServiceStatus st = QueryGuardexService();
    flutter::EncodableMap out;
    out[flutter::EncodableValue("installed")] =
        flutter::EncodableValue(st.installed);
    out[flutter::EncodableValue("running")] =
        flutter::EncodableValue(st.running);
    result->Success(flutter::EncodableValue(out));
    return;
  }

  if (method == "ensureService") {
    // Run elevation + SCM polling on a worker thread so the Flutter
    // platform thread (which is also the UI thread on Windows) stays free
    // to render frames. Blocking it here — even briefly — causes visible
    // jank when the user taps Connect.
    std::shared_ptr<flutter::MethodResult<flutter::EncodableValue>>
        shared_result(std::move(result));
    std::thread([shared_result, alive = alive_]() {
      ServiceStatus st = QueryGuardexService();
      if (!st.installed) {
        int rc = RunElevated(L"install");
        if (rc == -1) {
          if (alive->load()) shared_result->Error("elevation_cancelled",
                               "user declined UAC prompt for service install");
          return;
        }
        if (rc == -2) {
          if (alive->load()) shared_result->Error("install_timeout",
                               "guardexsvc install timed out after 30s");
          return;
        }
        if (rc != 0) {
          std::ostringstream msg;
          msg << "guardexsvc install exited with " << rc;
          if (alive->load()) shared_result->Error("install_failed", msg.str());
          return;
        }
      }
      st = QueryGuardexService();
      if (!st.running) {
        int rc = RunElevated(L"start");
        if (rc == -1) {
          if (alive->load()) shared_result->Error("elevation_cancelled",
                               "user declined UAC prompt for service start");
          return;
        }
        if (rc == -2) {
          if (alive->load()) shared_result->Error("start_timeout",
                               "guardexsvc start timed out after 30s");
          return;
        }
        if (rc != 0) {
          std::ostringstream msg;
          msg << "guardexsvc start exited with " << rc;
          if (alive->load()) shared_result->Error("start_failed", msg.str());
          return;
        }
      }
      if (!WaitForServiceRunning(5000)) {
        if (alive->load()) shared_result->Error("not_running",
                             "service did not reach RUNNING within 5s");
        return;
      }
      if (alive->load()) shared_result->Success(flutter::EncodableValue(nullptr));
    }).detach();
    return;
  }

  if (method == "uninstallService") {
    std::shared_ptr<flutter::MethodResult<flutter::EncodableValue>>
        shared_result(std::move(result));
    std::thread([shared_result, alive = alive_]() {
      int rc = RunElevated(L"uninstall");
      if (rc == -1) {
        if (alive->load()) shared_result->Error("elevation_cancelled",
                             "user declined UAC prompt");
        return;
      }
      if (rc != 0) {
        std::ostringstream msg;
        msg << "guardexsvc uninstall exited with " << rc;
        if (alive->load()) shared_result->Error("uninstall_failed", msg.str());
        return;
      }
      if (alive->load()) shared_result->Success(flutter::EncodableValue(nullptr));
    }).detach();
    return;
  }

  if (method == "startVPN" || method == "startVPNWindows") {
    if (args == nullptr) {
      result->Error("bad_args", "arguments map required");
      return;
    }
    std::string adapter = GetStringArg(args, "adapterName", "Guardex");
    std::string proxy = GetStringArg(args, "proxyAddress");
    std::string server = GetStringArg(args, "serverIP");
    int mtu = GetIntArg(args, "mtu", 1500);
    if (proxy.empty()) {
      result->Error("bad_args", "proxyAddress is required");
      return;
    }

    std::ostringstream params;
    params << "{\"adapter\":\"" << EscapeJson(adapter) << "\","
           << "\"proxy\":\"" << EscapeJson(proxy) << "\","
           << "\"server\":\"" << EscapeJson(server) << "\","
           << "\"mtu\":" << mtu << "}";
    std::string params_s = params.str();

    std::shared_ptr<flutter::MethodResult<flutter::EncodableValue>>
        shared_result(std::move(result));
    std::thread([shared_result, params_s, alive = alive_]() {
      RpcResult rpc = CallService("start_vpn", params_s);
      if (!alive->load()) return;
      if (!rpc.ok) {
        shared_result->Error("rpc_failed", rpc.error);
        return;
      }
      if (!JsonOK(rpc.body)) {
        shared_result->Error("start_failed", JsonError(rpc.body));
        return;
      }
      // Preserve mixin fd return type on Android (int fd). Windows never
      // has a real fd so just return 0.
      shared_result->Success(flutter::EncodableValue(static_cast<int32_t>(0)));
    }).detach();
    return;
  }

  if (method == "stopVPN" || method == "stopVPNWindows") {
    std::shared_ptr<flutter::MethodResult<flutter::EncodableValue>>
        shared_result(std::move(result));
    std::thread([shared_result, alive = alive_]() {
      RpcResult rpc = CallService("stop_vpn", "");
      if (!alive->load()) return;
      if (!rpc.ok) {
        shared_result->Error("rpc_failed", rpc.error);
        return;
      }
      shared_result->Success(flutter::EncodableValue(nullptr));
    }).detach();
    return;
  }

  if (method == "isVPNRunning") {
    std::shared_ptr<flutter::MethodResult<flutter::EncodableValue>>
        shared_result(std::move(result));
    std::thread([shared_result, alive = alive_]() {
      RpcResult rpc = CallService("is_running", "");
      if (!alive->load()) return;
      if (!rpc.ok) {
        // If the service is not reachable, report "not running" rather
        // than surfacing a transport error — the UI polls this routinely
        // and should degrade gracefully.
        shared_result->Success(flutter::EncodableValue(false));
        return;
      }
      shared_result->Success(flutter::EncodableValue(JsonResultBool(rpc.body)));
    }).detach();
    return;
  }

  result->NotImplemented();
}

}  // namespace proxy_core
