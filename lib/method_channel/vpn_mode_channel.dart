import 'package:flutter/services.dart';
import 'package:proxy_core/constants/vpn_mode_methods.dart';
import 'package:proxy_core/models/proxy_core_config.dart';

mixin class VpnModeChannelMixin {
  static const MethodChannel _methCh = MethodChannel('proxy_core/vpn');

  bool _isVpnRunning = false;

  /// To check if the VPN is currently running after the last toggling action
  bool get isVPNRunning => _isVpnRunning;

  /// Prepare the config for VPN mode if needed
  Future<ProxyCoreConfig> prepareConfigForVpnIfNeeded(
      ProxyCoreConfig config) async {
    if (config.vpnMode) {
      await _prepareVpnProfile();
      final fd = await _startVPN(config.blockedApps);
      config.parcelFileId = fd;
    }
    return config;
  }

  /// Start the VPN and return the file descriptor
  Future<int> _startVPN(List<String>? blockedApps) async {
    final int fd = await _methCh.invokeMethod(
      VpnModeMethods.startVPN.name,
      <String, dynamic>{
        if (blockedApps != null && blockedApps.isNotEmpty)
          'blockedApps': blockedApps,
      },
    );
    await setVpnStatus();
    return fd;
  }

  /// Stop the VPN
  Future stopVPN() async {
    await _methCh.invokeMethod(VpnModeMethods.stopVPN.name);
    await setVpnStatus();
  }

  /// Check if VPN is currently running
  Future<bool> setVpnStatus() async {
    final result = await _methCh.invokeMethod(VpnModeMethods.isVPNRunning.name);
    // Accept both int and bool from native side
    _isVpnRunning = result == true || result == 1;
    return _isVpnRunning;
  }

  //. Call this before trying to turn on vpn mode
  Future<void> _prepareVpnProfile() async =>
      await _methCh.invokeMethod(VpnModeMethods.prepare.name);

  /// To have a clean state before _simpleStart (test purposes)
  Future<void> stopAnyRunningVPN() async {
    await _prepareVpnProfile();
    await stopVPN();
  }

  /// Queries whether the privileged Guardex Windows service is
  /// installed and currently running. Returns a map with keys
  /// `installed` and `running`, both `false` on non-Windows platforms.
  Future<Map<String, bool>> windowsServiceStatus() async {
    final dynamic raw = await _methCh.invokeMethod('serviceStatus');
    if (raw is Map) {
      return <String, bool>{
        'installed': raw['installed'] == true,
        'running': raw['running'] == true,
      };
    }
    return <String, bool>{'installed': false, 'running': false};
  }

  /// Ensures the Guardex Windows service is installed and running.
  /// Triggers one or two UAC prompts (install, then start) the first
  /// time only; subsequent calls are fast no-ops.
  ///
  /// Throws a [PlatformException] with code `elevation_cancelled` if
  /// the user declines the UAC dialog.
  Future<void> ensureWindowsService() async {
    await _methCh.invokeMethod('ensureService');
  }

  /// Uninstalls the Guardex Windows service. Prompts for UAC.
  Future<void> uninstallWindowsService() async {
    await _methCh.invokeMethod('uninstallService');
  }

  /// Start the VPN on Windows using the wintun-based TUN engine.
  ///
  /// [adapterName] — name for the wintun adapter (e.g. "Guardex").
  /// [proxyAddress] — "host:port" of local SOCKS5 proxy.
  /// [serverIP] — public IP of the remote VPN server; added as /32
  /// exception through the original gateway so the encrypted tunnel
  /// to the server itself does not loop through the TUN.
  /// [mtu] — link MTU; pass 0 for default 1500.
  Future<void> startVPNWindows({
    required String adapterName,
    required String proxyAddress,
    required String serverIP,
    int mtu = 1500,
  }) async {
    await _methCh.invokeMethod(
      'startVPNWindows',
      <String, dynamic>{
        'adapterName': adapterName,
        'proxyAddress': proxyAddress,
        'serverIP': serverIP,
        'mtu': mtu,
      },
    );
    await setVpnStatus();
  }
}
