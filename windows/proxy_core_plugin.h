#ifndef FLUTTER_PLUGIN_PROXY_CORE_PLUGIN_H_
#define FLUTTER_PLUGIN_PROXY_CORE_PLUGIN_H_

#include <flutter/method_channel.h>
#include <flutter/plugin_registrar_windows.h>

#include <memory>

namespace proxy_core {

class ProxyCorePlugin : public flutter::Plugin {
 public:
  static void RegisterWithRegistrar(flutter::PluginRegistrarWindows* registrar);

  ProxyCorePlugin();
  virtual ~ProxyCorePlugin();

  ProxyCorePlugin(const ProxyCorePlugin&) = delete;
  ProxyCorePlugin& operator=(const ProxyCorePlugin&) = delete;

 private:
  void HandleMethodCall(
      const flutter::MethodCall<flutter::EncodableValue>& method_call,
      std::unique_ptr<flutter::MethodResult<flutter::EncodableValue>> result);
};

}  // namespace proxy_core

#endif  // FLUTTER_PLUGIN_PROXY_CORE_PLUGIN_H_
