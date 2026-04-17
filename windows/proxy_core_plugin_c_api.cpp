#include "include/proxy_core/proxy_core_plugin_c_api.h"

#include <flutter/plugin_registrar_windows.h>

#include "proxy_core_plugin.h"

void ProxyCorePluginCApiRegisterWithRegistrar(
    FlutterDesktopPluginRegistrarRef registrar) {
  proxy_core::ProxyCorePlugin::RegisterWithRegistrar(
      flutter::PluginRegistrarManager::GetInstance()
          ->GetRegistrar<flutter::PluginRegistrarWindows>(registrar));
}
