import 'dart:async';

import 'package:proxy_core/models/proxy_core_config.dart';
import 'package:proxy_core/gen/bindings/ProxyCoreService.pbgrpc.dart';





abstract interface class ProxyCoreInterface {
  
  
  
  Future<void> initialize(
      {Future<void> Function(bool isRunning)? onCoreStateChanged});

  
  Future<bool> get isRunning;

  
  Future<String> get version;

  
  
  
  Future<void> start(ProxyCoreConfig config);

  
  
  
  
  
  Future<void> simpleStart(ProxyCoreConfig config);

  
  
  
  Future<void> stop();

  
  
  
  
  Future<List<PingResult>> measurePing(List<String> urls);

  
  
  
  Future<LogResponse> fetchLogs();

  
  
  
  Future<Empty> clearLogs();

  
  
  
  Future<String> getMemoryUsage();

  
  
  
  Future<String> getCpuUsage();

  
  void ensureInitialized();
}
