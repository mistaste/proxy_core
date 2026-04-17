#include "proxy_core_plugin.h"

#include <flutter/method_channel.h>
#include <flutter/plugin_registrar_windows.h>
#include <flutter/standard_method_codec.h>

#include <windows.h>
#include <shellapi.h>
#include <winsvc.h>

#include <atomic>
#include <memory>
#include <mutex>
#include <sstream>
#include <string>
#include <vector>

#pragma comment(lib, "advapi32.lib")
#pragma comment(lib, "shell32.lib")

namespace proxy_core {

namespace {

// Named pipe exposed by guardexsvc.exe. Must match service_win.PipeName.
constexpr const char* kPipeName = R"(\\.\pipe\guardex-svc)";

std::mutex g_pipe_mutex;
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
  std::lock_guard<std::mutex> lock(g_pipe_mutex);

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

// Shallow JSON field scan — we only need a couple of primitive fields so
// we avoid pulling in a full parser.
bool JsonOK(const std::string& body) {
  return body.find("\"ok\":true") != std::string::npos;
}

std::string JsonError(const std::string& body) {
  auto k = body.find("\"error\":");
  if (k == std::string::npos) return "";
  auto q1 = body.find('"', k + 8);
  if (q1 == std::string::npos) return "";
  auto q2 = body.find('"', q1 + 1);
  if (q2 == std::string::npos) return "";
  return body.substr(q1 + 1, q2 - q1 - 1);
}

bool JsonResultBool(const std::string& body) {
  return body.find("\"result\":true") != std::string::npos;
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
  WaitForSingleObject(sei.hProcess, INFINITE);
  DWORD exitCode = 1;
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
ProxyCorePlugin::~ProxyCorePlugin() {}

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
    ServiceStatus st = QueryGuardexService();
    if (!st.installed) {
      int rc = RunElevated(L"install");
      if (rc == -1) {
        result->Error("elevation_cancelled",
                      "user declined UAC prompt for service install");
        return;
      }
      if (rc != 0) {
        std::ostringstream msg;
        msg << "guardexsvc install exited with " << rc;
        result->Error("install_failed", msg.str());
        return;
      }
    }
    st = QueryGuardexService();
    if (!st.running) {
      int rc = RunElevated(L"start");
      if (rc == -1) {
        result->Error("elevation_cancelled",
                      "user declined UAC prompt for service start");
        return;
      }
      if (rc != 0) {
        std::ostringstream msg;
        msg << "guardexsvc start exited with " << rc;
        result->Error("start_failed", msg.str());
        return;
      }
    }
    if (!WaitForServiceRunning(5000)) {
      result->Error("not_running",
                    "service did not reach RUNNING within 5s");
      return;
    }
    result->Success(flutter::EncodableValue(nullptr));
    return;
  }

  if (method == "uninstallService") {
    int rc = RunElevated(L"uninstall");
    if (rc == -1) {
      result->Error("elevation_cancelled", "user declined UAC prompt");
      return;
    }
    if (rc != 0) {
      std::ostringstream msg;
      msg << "guardexsvc uninstall exited with " << rc;
      result->Error("uninstall_failed", msg.str());
      return;
    }
    result->Success(flutter::EncodableValue(nullptr));
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

    RpcResult rpc = CallService("start_vpn", params.str());
    if (!rpc.ok) {
      result->Error("rpc_failed", rpc.error);
      return;
    }
    if (!JsonOK(rpc.body)) {
      result->Error("start_failed", JsonError(rpc.body));
      return;
    }
    // Preserve mixin fd return type on Android (int fd). Windows never
    // has a real fd so just return 0.
    result->Success(flutter::EncodableValue(static_cast<int32_t>(0)));
    return;
  }

  if (method == "stopVPN" || method == "stopVPNWindows") {
    RpcResult rpc = CallService("stop_vpn", "");
    if (!rpc.ok) {
      result->Error("rpc_failed", rpc.error);
      return;
    }
    result->Success(flutter::EncodableValue(nullptr));
    return;
  }

  if (method == "isVPNRunning") {
    RpcResult rpc = CallService("is_running", "");
    if (!rpc.ok) {
      // If the service is not reachable, report "not running" rather
      // than surfacing a transport error — the UI polls this routinely
      // and should degrade gracefully.
      result->Success(flutter::EncodableValue(false));
      return;
    }
    result->Success(flutter::EncodableValue(JsonResultBool(rpc.body)));
    return;
  }

  result->NotImplemented();
}

}  // namespace proxy_core
