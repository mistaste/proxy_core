#include "proxy_core_plugin.h"

#include <flutter/method_channel.h>
#include <flutter/plugin_registrar_windows.h>
#include <flutter/standard_method_codec.h>

#include <windows.h>

#include <memory>
#include <sstream>
#include <string>

namespace proxy_core {

namespace {

using StartVPNWindowsFn = int (*)(char*, char*, char*, int);
using StopVPNWindowsFn = int (*)();
using IsVPNRunningFn = int (*)();

HMODULE g_dll = nullptr;
StartVPNWindowsFn g_start = nullptr;
StopVPNWindowsFn g_stop = nullptr;
IsVPNRunningFn g_is_running = nullptr;

bool LoadDll() {
  if (g_dll != nullptr) return true;
  g_dll = LoadLibraryA("libproxy_core.dll");
  if (g_dll == nullptr) return false;
  g_start = reinterpret_cast<StartVPNWindowsFn>(
      GetProcAddress(g_dll, "StartVPNWindows"));
  g_stop = reinterpret_cast<StopVPNWindowsFn>(
      GetProcAddress(g_dll, "StopVPNWindows"));
  g_is_running = reinterpret_cast<IsVPNRunningFn>(
      GetProcAddress(g_dll, "IsVPNRunning"));
  return g_start && g_stop && g_is_running;
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

  if (!LoadDll()) {
    result->Error("dll_load_failed",
                  "libproxy_core.dll not found or missing exports");
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

    int rc = g_start(const_cast<char*>(adapter.c_str()),
                     const_cast<char*>(proxy.c_str()),
                     const_cast<char*>(server.c_str()), mtu);
    if (rc != 0) {
      std::ostringstream msg;
      msg << "StartVPNWindows returned " << rc;
      result->Error("start_failed", msg.str());
      return;
    }
    // Return a fake fd for compatibility with mixin signature.
    result->Success(flutter::EncodableValue(static_cast<int32_t>(0)));
    return;
  }

  if (method == "stopVPN" || method == "stopVPNWindows") {
    g_stop();
    result->Success(flutter::EncodableValue(nullptr));
    return;
  }

  if (method == "isVPNRunning") {
    int running = g_is_running();
    result->Success(flutter::EncodableValue(running != 0));
    return;
  }

  result->NotImplemented();
}

}  // namespace proxy_core
