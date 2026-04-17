import 'dart:async';
import 'dart:io';

import 'package:proxy_core/models/proxy_core_config.dart';
import 'package:proxy_core/platforms/proxy_core_interface.dart';
import 'package:proxy_core/platforms/proxy_core_base_impl.dart';
import 'package:proxy_core/platforms/proxy_core_ios_vpn_impl.dart';
import 'package:proxy_core/gen/bindings/ProxyCoreService.pbgrpc.dart';

/// ProxyCore - Multi-platform proxy and VPN service
///
/// This class provides a unified interface to the ProxyCore service across different platforms.
/// It automatically selects the appropriate implementation:
/// - **iOS**: Uses [ProxyCoreIosVpnImpl] with native method channels for VPN tunnel management
/// - **Android/Windows/Linux/macOS**: Uses [ProxyCoreBaseImpl] with gRPC and FFI bindings
///
/// Platform-specific behaviors:
/// - iOS: Manages VPN tunnel via Network Extension and receives status updates via EventChannel
/// - Other platforms: Uses persistent notifications, gRPC server, and VPN status polling
class ProxyCore {
  ProxyCore._();

  static final ProxyCore _instance = ProxyCore._();
  static ProxyCore get ins => _instance;

  final ProxyCoreInterface _proxyCoreImpl =
      Platform.isIOS ? ProxyCoreIosVpnImpl() : ProxyCoreBaseImpl();

  /// Sets the iOS tunnel information required for VPN configuration
  ///
  /// **iOS Only**: This method configures the app name and tunnel bundle identifier
  /// needed for the Network Extension VPN tunnel provider.
  ///
  /// Parameters:
  /// - [appName]: The display name of the VPN connection
  /// - [appTunnelBundle]: The bundle identifier of the PacketTunnelProvider extension
  ///
  /// This must be called before [initialize] on iOS platforms.
  void setIosTunnelInfo(String appName, String appTunnelBundle) {
    if (_proxyCoreImpl is ProxyCoreIosVpnImpl) {
      (_proxyCoreImpl).setIosTunnelInfo(appName, appTunnelBundle);
    }
  }

  /// Initializes the ProxyCore service
  ///
  /// **iOS**: Sets up method channel communication and event listeners for VPN status changes
  /// **Other platforms**: Initializes gRPC server, notification service, and FFI bindings
  ///
  /// Parameters:
  /// - [onCoreStateChanged]: Optional callback invoked when the core running state changes
  ///   - iOS: Triggered by native VPN status events via EventChannel
  ///   - Other platforms: Called after start/stop operations and during VPN status polling
  ///
  /// Must be called before using any other ProxyCore methods.
  Future<void> initialize({
    Future<void> Function(bool isRunning)? onCoreStateChanged,
  }) async {
    await _proxyCoreImpl.initialize(onCoreStateChanged: onCoreStateChanged);
  }

  /// Starts the proxy/VPN service with the given configuration
  ///
  /// **iOS**:
  /// - If tunnel is already connected, calls [simpleStart] to update configuration
  /// - Otherwise, prepares VPN manager and starts the Network Extension tunnel
  /// - Passes configuration to the PacketTunnelProvider via NETunnelProviderSession
  ///
  /// **Other platforms**:
  /// - Requests notification permissions if required
  /// - Stops any existing VPN connections
  /// - Prepares configuration for VPN mode if needed
  /// - Starts the core via gRPC
  /// - Shows persistent notification
  /// - Begins VPN status polling if in VPN mode
  ///
  /// Throws [ProxyCoreException] if the operation fails.
  Future<void> start(ProxyCoreConfig config) async {
    await _proxyCoreImpl.start(config);
  }

  /// Starts or restarts the core with new configuration without stopping VPN tunnel
  ///
  /// **iOS**:
  /// - Stops any running core via simple stop
  /// - If tunnel is not connected, prepares and starts VPN tunnel
  /// - Sends new configuration to the already-running tunnel provider
  /// - Useful for switching servers without reconnecting the entire VPN
  ///
  /// **Other platforms**:
  /// - Stops any existing VPN connections
  /// - Starts the core without notification handling
  /// - Does not require permission checks
  ///
  /// Throws [ProxyCoreException] if the operation fails.
  Future<void> simpleStart(ProxyCoreConfig config) async {
    await _proxyCoreImpl.simpleStart(config);
  }

  /// Stops the proxy/VPN service
  ///
  /// **iOS**: Sends stop command to the VPN tunnel via NETunnelProviderManager
  ///
  /// **Other platforms**:
  /// - Stops any active VPN connections
  /// - Stops the core via gRPC
  /// - Cancels persistent notification
  /// - Stops VPN status polling
  ///
  /// Throws [ProxyCoreException] if the operation fails.
  Future<void> stop() async {
    await _proxyCoreImpl.stop();
  }

  /// Checks if the core is currently running
  ///
  /// **iOS**: Queries the tunnel status via method channel
  ///
  /// **Other platforms**: Checks core status via gRPC call
  ///
  /// Returns `true` if the core is running, `false` otherwise.
  Future<bool> get isRunning => _proxyCoreImpl.isRunning;

  /// Fetches logs from the core
  ///
  /// **iOS**: Retrieves logs from the tunnel provider via method channel
  ///
  /// **Other platforms**: Fetches logs via gRPC
  ///
  /// Returns a [LogResponse] containing the log messages.
  Future<LogResponse> fetchLogs() => _proxyCoreImpl.fetchLogs();

  /// Clears all accumulated logs
  ///
  /// **iOS**: Sends clear command to tunnel provider via method channel
  ///
  /// **Other platforms**: Clears logs via gRPC
  ///
  /// Returns an empty response on success.
  Future<Empty> clearLogs() => _proxyCoreImpl.clearLogs();

  /// Measures ping latency to the specified URLs
  ///
  /// **iOS**:
  /// - Sends URLs to tunnel provider as comma-separated string
  /// - Receives comma-separated ping results
  /// - Parses results into [PingResult] objects
  ///
  /// **Other platforms**: Measures ping via gRPC and returns structured results
  ///
  /// Parameters:
  /// - [urls]: List of URLs to ping
  ///
  /// Returns a list of [PingResult] objects containing URL and delay information.
  Future<List<PingResult>> measurePing(List<String> urls) =>
      _proxyCoreImpl.measurePing(urls);

  /// Gets the version of the ProxyCore library
  ///
  /// **iOS**: Queries version via method channel
  ///
  /// **Other platforms**: Retrieves version via gRPC
  ///
  /// Returns the version string.
  Future<String> get version => _proxyCoreImpl.version;

  /// Gets the memory usage of the core
  ///
  /// **iOS**: Queries memory usage via method channel
  ///
  /// Returns the memory usage in bytes.
  Future<String> get memoryUsage => _proxyCoreImpl.getMemoryUsage();

  /// Starts the wintun-based TUN engine on Windows so default traffic
  /// is routed into the local SOCKS5 proxy that the core is serving.
  ///
  /// Requires the core to already be running in proxy mode via [start]
  /// (i.e. the xray instance is listening on [proxyPort]).
  ///
  /// - [adapterName] — wintun adapter name (e.g. "Guardex").
  /// - [proxyAddress] — "host:port" of the local SOCKS5 proxy.
  /// - [serverIP] — remote VPN server IP, excluded from the tunnel so
  ///   the encrypted link itself does not loop through the TUN.
  /// - [mtu] — link MTU; default 1500.
  ///
  /// No-op on non-Windows platforms.
  Future<void> startWindowsTun({
    required String adapterName,
    required String proxyAddress,
    required String serverIP,
    int mtu = 1500,
  }) async {
    if (!Platform.isWindows) return;
    if (_proxyCoreImpl is ProxyCoreBaseImpl) {
      await (_proxyCoreImpl as ProxyCoreBaseImpl).startVPNWindows(
        adapterName: adapterName,
        proxyAddress: proxyAddress,
        serverIP: serverIP,
        mtu: mtu,
      );
    }
  }

  /// Stops the wintun-based TUN engine on Windows and tears down the
  /// routes it installed. No-op on non-Windows platforms.
  Future<void> stopWindowsTun() async {
    if (!Platform.isWindows) return;
    if (_proxyCoreImpl is ProxyCoreBaseImpl) {
      await (_proxyCoreImpl as ProxyCoreBaseImpl).stopVPN();
    }
  }

  /// Returns whether the privileged Guardex Windows service is installed
  /// and running. Always reports both fields as `false` on non-Windows.
  Future<Map<String, bool>> windowsServiceStatus() async {
    if (!Platform.isWindows) {
      return const <String, bool>{'installed': false, 'running': false};
    }
    if (_proxyCoreImpl is ProxyCoreBaseImpl) {
      return (_proxyCoreImpl as ProxyCoreBaseImpl).windowsServiceStatus();
    }
    return const <String, bool>{'installed': false, 'running': false};
  }

  /// Installs and starts the Guardex Windows service if it is not
  /// already running. Triggers the UAC prompt on first call only;
  /// subsequent calls are fast no-ops. No-op on non-Windows platforms.
  Future<void> ensureWindowsService() async {
    if (!Platform.isWindows) return;
    if (_proxyCoreImpl is ProxyCoreBaseImpl) {
      await (_proxyCoreImpl as ProxyCoreBaseImpl).ensureWindowsService();
    }
  }

  /// Uninstalls the Guardex Windows service. Triggers a UAC prompt.
  /// No-op on non-Windows platforms.
  Future<void> uninstallWindowsService() async {
    if (!Platform.isWindows) return;
    if (_proxyCoreImpl is ProxyCoreBaseImpl) {
      await (_proxyCoreImpl as ProxyCoreBaseImpl).uninstallWindowsService();
    }
  }

  /// Gets the CPU usage of the core
  ///
  /// **iOS**: Queries CPU usage via method channel
  ///
  /// Returns the CPU usage in percentage.
  Future<String> get cpuUsage => _proxyCoreImpl.getCpuUsage();

  /// Gets cumulative traffic statistics (uplink/downlink bytes)
  ///
  /// **iOS**: Returns zeros (not supported via tunnel provider)
  ///
  /// **Other platforms**: Queries Xray stats.Manager via gRPC
  ///
  /// Requires Xray config to include `"stats":{}` and appropriate `"policy"` section.
  Future<TrafficStatsResponse> getTrafficStats() =>
      _proxyCoreImpl.getTrafficStats();
}
