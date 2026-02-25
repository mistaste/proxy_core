import 'dart:async';
import 'dart:io';

import 'package:proxy_core/models/proxy_core_config.dart';
import 'package:proxy_core/platforms/proxy_core_interface.dart';
import 'package:proxy_core/platforms/proxy_core_base_impl.dart';
import 'package:proxy_core/platforms/proxy_core_ios_vpn_impl.dart';
import 'package:proxy_core/gen/bindings/ProxyCoreService.pbgrpc.dart';











class ProxyCore {
  ProxyCore._();

  static final ProxyCore _instance = ProxyCore._();
  static ProxyCore get ins => _instance;

  final ProxyCoreInterface _proxyCoreImpl =
      Platform.isIOS ? ProxyCoreIosVpnImpl() : ProxyCoreBaseImpl();

  
  
  
  
  
  
  
  
  
  
  void setIosTunnelInfo(String appName, String appTunnelBundle) {
    if (_proxyCoreImpl is ProxyCoreIosVpnImpl) {
      (_proxyCoreImpl).setIosTunnelInfo(appName, appTunnelBundle);
    }
  }

  
  
  
  
  
  
  
  
  
  
  
  Future<void> initialize({
    Future<void> Function(bool isRunning)? onCoreStateChanged,
  }) async {
    await _proxyCoreImpl.initialize(onCoreStateChanged: onCoreStateChanged);
  }

  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  Future<void> start(ProxyCoreConfig config) async {
    await _proxyCoreImpl.start(config);
  }

  
  
  
  
  
  
  
  
  
  
  
  
  
  
  Future<void> simpleStart(ProxyCoreConfig config) async {
    await _proxyCoreImpl.simpleStart(config);
  }

  
  
  
  
  
  
  
  
  
  
  
  Future<void> stop() async {
    await _proxyCoreImpl.stop();
  }

  
  
  
  
  
  
  
  Future<bool> get isRunning => _proxyCoreImpl.isRunning;

  
  
  
  
  
  
  
  Future<LogResponse> fetchLogs() => _proxyCoreImpl.fetchLogs();

  
  
  
  
  
  
  
  Future<Empty> clearLogs() => _proxyCoreImpl.clearLogs();

  
  
  
  
  
  
  
  
  
  
  
  
  
  Future<List<PingResult>> measurePing(List<String> urls) =>
      _proxyCoreImpl.measurePing(urls);

  
  
  
  
  
  
  
  Future<String> get version => _proxyCoreImpl.version;

  
  
  
  
  
  Future<String> get memoryUsage => _proxyCoreImpl.getMemoryUsage();

  
  
  
  
  
  Future<String> get cpuUsage => _proxyCoreImpl.getCpuUsage();
}
