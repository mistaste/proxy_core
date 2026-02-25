import Foundation
import Flutter
import NetworkExtension
import Combine
import os.log


protocol ProxyCorePluginImpl {
    func handleMethodCall(_ call: FlutterMethodCall, result: @escaping FlutterResult)
}

class ProxyCorePluginIosImpl: NSObject, ProxyCorePluginImpl, FlutterStreamHandler, VpnStatusDelegate {
    private var vpnManager: NETunnelProviderManager?
    private var logTimer: Timer?
    private var logCallback: ((String) -> Void)?
    private let vpnService: VpnService
    private let logger = Logger(subsystem: Bundle.main.bundleIdentifier ?? "com.proxycore", category: "ProxyCorePlugin")
    
    
    private var eventSink: FlutterEventSink?
    
    override init() {
        self.vpnService = VpnService.shared
        super.init()
        self.vpnService.statusDelegate = self
        
        
        NotificationCenter.default.addObserver(
            self,
            selector: #selector(appDidBecomeActive),
            name: UIApplication.didBecomeActiveNotification,
            object: nil
        )
    }
    
    deinit {
        NotificationCenter.default.removeObserver(self)
    }
    
    @objc private func appDidBecomeActive() {
        
        refreshVpnStatus()
    }
    
    private func refreshVpnStatus() {
        
        vpnService.loadManager { [weak self] success in
            if success, let manager = self?.vpnService.vpnManager {
                
                let status = manager.connection.status
                self?.logger.info("App became active, current VPN status: \(String(status.rawValue))")
                
                if let session = manager.connection as? NETunnelProviderSession {
                    
                    self?.vpnService.sendTunnelMessage(["command": "IS_CORE_RUNNING"]) { response in
                        let isRunning = response?.lowercased() == "true"
                        self?.logger.info("VPN tunnel is running: \(isRunning ? "true" : "false")")
                        
                        if isRunning {
                            
                            self?.sendVpnStatus(true)
                        } else {
                            
                            self?.sendVpnStatus(status == .connected)
                        }
                    }
                } else {
                    
                    self?.sendVpnStatus(status == .connected)
                }
            } else {
                self?.sendVpnStatus(false)
            }
        }
    }

    
    
    func onListen(withArguments arguments: Any?, eventSink events: @escaping FlutterEventSink) -> FlutterError? {
        self.eventSink = events
        
        
        refreshVpnStatus()
        
        return nil
    }
    
    func onCancel(withArguments arguments: Any?) -> FlutterError? {
        self.eventSink = nil
        return nil
    }
    
    
    
    func vpnStatusDidChange(_ status: NEVPNStatus) {
        let isConnected = (status == .connected)
        sendVpnStatus(isConnected)
        
        
        if status == .disconnected {
            vpnService.sendTunnelMessage(["command": "IS_CORE_RUNNING"]) { [weak self] response in
                if response?.lowercased() == "true" {
                    
                    self?.sendVpnStatus(true)
                }
            }
        }
    }
    
    private func sendVpnStatus(_ isConnected: Bool) {
        DispatchQueue.main.async { [weak self] in
            print("VPN status swift: \(isConnected)")
            self?.eventSink?(isConnected)
        }
    }
    
    enum StartMode {
        case normal
        case simple
    }

    func handleMethodCall(_ call: FlutterMethodCall, result: @escaping FlutterResult) {
        switch call.method {
        case "prepare":
            handlePrepareVPN(call, result: result)
        case "startVPN":
            handleStartVPN(call, result: result, mode: .normal)
        case "stopVPN":
            handleStopVPN(result)
        
        
        case "simpleStartVpn":
            handleStartVPN(call, result: result, mode: .simple)
        case "simpleStart":
            handleSimpleStart(call, result: result)
        case "simpleStop":
            handleSimpleStop(result)
        

        case "isCoreRunning":
            handleIsCoreRunning(result)
        case "measurePing":
            handleMeasurePing(call, result: result)
        case "fetchLogs":
            handleFetchLogs(result)
        case "clearLogs":
            handleClearLogs(result)
        case "getTunnelStatus":
            getTunnelStatus(result)
        case "getVersion":
            handleGetVersion(result)
        case "getMemoryUsage":
            handleGetMemoryUsage(result)
        case "getCpuUsage":
            handleGetCpuUsage(result)
        default:
            result(FlutterMethodNotImplemented)
        }
    }

    private func handlePrepareVPN(_ call: FlutterMethodCall, result: @escaping FlutterResult) {
        guard let args = call.arguments as? [String: Any],
              let appName = args["appName"] as? String,
              let appTunnelBundle = args["appTunnelBundle"] as? String
        else {
            result(MethodError.invalidArguments("Missing required parameters for prepare").flutterError)
            return
        }
        
        prepareVPN(appName: appName, appTunnelBundle: appTunnelBundle, result)
    }

    private func prepareVPN(
        appName: String, appTunnelBundle: String, _ result: @escaping FlutterResult
    ) {
        vpnService.prepareVPN(appName: appName, appTunnelBundle: appTunnelBundle) { vpnResult in
            DispatchQueue.main.async {
                switch vpnResult {
                case .success:
                    result(true)
                case .failure(let error):
                    result(FlutterError(
                        code: "PREPARE_FAILED",
                        message: error.localizedDescription,
                        details: nil
                    ))
                }
            }
        }
    }

    private func handleSimpleStart(_ call: FlutterMethodCall, result: @escaping FlutterResult) {
        guard let args = call.arguments as? [String: Any],
            let coreName = args["coreName"] as? String,
            let config = args["config"] as? String,
            let cacheDir = args["cacheDir"] as? String
        else {
            result(
                MethodError.invalidArguments("Missing required parameters for startVPN")
                    .flutterError)
            return
        }

        vpnService.sendTunnelMessage([
            "command": "SIMPLE_START_CORE", 
            "coreName": coreName,
            "config": config,
            "cacheDir": cacheDir
        ]) { coreResult in
            DispatchQueue.main.async {
                result(coreResult)
            }
        }

    }
    private func handleSimpleStop(_ result: @escaping FlutterResult) {
        vpnService.sendTunnelMessage(["command": "SIMPLE_STOP_CORE"]) { response in
            DispatchQueue.main.async {
                let responseValue = response?.lowercased() == "true"
                result(responseValue)
            }
        }
        
        result(true)
    }

    private func handleStartVPN(_ call: FlutterMethodCall, result: @escaping FlutterResult, mode: StartMode) {
        guard let args = call.arguments as? [String: Any],
            let appName = args["appName"] as? String,
            let appTunnelBundle = args["appTunnelBundle"] as? String
        else {
            result(
                MethodError.invalidArguments("Missing required parameters for startVPN")
                    .flutterError)
            return
        }
        
        let coreName: String? = args["coreName"] as? String
        let config: String?   = args["config"] as? String
        let cacheDir: String? = args["cacheDir"] as? String


        let startMode = mode == .normal ? "normal" : "simple"

        logger.debug("Starting VPN with config")

        vpnService.startVPN(startMode: startMode, coreName: coreName, config: config, cacheDir: cacheDir, port: 2080, appName: appName, appTunnelBundle: appTunnelBundle) { vpnResult in
            switch vpnResult {
                case .success:
                    self.logger.info("VPN tunnel started successfully")
                    DispatchQueue.main.async {
                        result(true)
                    }
                case .failure(let error):
                    self.logger.error("Failed to start VPN: \(error.localizedDescription)")
                    DispatchQueue.main.async {
                        result(
                            FlutterError(
                                code: "START_FAILED",
                                message: error.localizedDescription,
                                details: nil))
                    }
                }
            }
    }

    private func handleStopVPN(_ result: @escaping FlutterResult) {
        vpnService.stopVPN { vpnResult in
            switch vpnResult {
            case .success:
                DispatchQueue.main.async {
                    result(true)
                }
            case .failure(let error):
                DispatchQueue.main.async {
                    result(
                        FlutterError(
                            code: "STOP_FAILED",
                            message: error.localizedDescription,
                            details: nil))
                }
            }
        }
    }

    private func handleIsCoreRunning(_ result: @escaping FlutterResult) {
        vpnService.sendTunnelMessage(["command": "IS_CORE_RUNNING"]) { response in
            DispatchQueue.main.async {
                result(response?.lowercased() == "true")
            }
        }
        
        result(false)
    }

    private func handleMeasurePing(_ call: FlutterMethodCall, result: @escaping FlutterResult) {
        guard let args = call.arguments as? [String: Any],
            let urls = args["urls"] as? String
        else {
            result(
                MethodError.invalidArguments("Missing required parameters for measurePing")
                    .flutterError)
            return
        }

        vpnService.sendTunnelMessage([
            "command": "measurePing", 
            "urls": urls
            ]) { pingResult in
            DispatchQueue.main.async {
                result(pingResult)
            }
        }
    }
    private func handleFetchLogs(_ result: @escaping FlutterResult) {
        vpnService.sendTunnelMessage(["command": "FETCH_LOGS"]) { response in
            DispatchQueue.main.async {
                if let logs = response {
                    result(logs)
                } else {
                    result(
                        FlutterError(
                            code: "FETCH_LOGS_FAILED",
                            message: "No response received from tunnel",
                            details: nil
                        )
                    )
                }
            }
        }
    }

    private func handleClearLogs(_ result: @escaping FlutterResult) {
        vpnService.sendTunnelMessage(["command": "CLEAR_LOGS"]) { response in
            DispatchQueue.main.async {
                if let logs = response {
                    result(logs)
                }else{
                    result(nil)
                }
            }
        }
    }
    private func handleGetVersion(_ result: @escaping FlutterResult) {
        vpnService.sendTunnelMessage(["command":"GET_VERSION"]){ response in
            DispatchQueue.main.async {
                if let version = response {
                    result(version)
                } else {
                    result(
                        FlutterError(
                            code: "GET_VERSION_FAILED",
                            message: "No response received from tunnel",
                            details: nil
                        )
                    )
                }
            }
        }
    }
    private func handleGetMemoryUsage(_ result: @escaping FlutterResult) {
        vpnService.sendTunnelMessage(["command":"GET_MEMORY_USAGE"]){ response in
            DispatchQueue.main.async {
                if let memoryUsage = response {
                    result(memoryUsage)
                }
                else{
                    result(0)
                }
            }
        }
    }
    private func handleGetCpuUsage(_ result: @escaping FlutterResult) {
        vpnService.sendTunnelMessage(["command":"GET_CPU_USAGE"]){ response in
            DispatchQueue.main.async {
                if let cpuUsage = response {
                    result(cpuUsage)
                }else{
                    result(0)
                }
            }
        }
    }
    
    private func getTunnelStatus(_ result: @escaping FlutterResult) {
        if VpnService.shared.vpnManager == nil {
            VpnService.shared.loadManager { success in
                if success {
                    self.returnCurrentVpnStatus(result)
                } else {
                    result("disconnected")
                }
            }
        } else {
            returnCurrentVpnStatus(result)
        }
    }

    private func returnCurrentVpnStatus(_ result: @escaping FlutterResult) {
        if let status = VpnService.shared.vpnManager?.connection.status {
            var statusString = ""
            
            switch status {
            case .connected: statusString = "connected"
            case .connecting: statusString = "connecting"
            case .disconnected: statusString = "disconnected"
            case .disconnecting: statusString = "disconnecting"
            case .invalid: statusString = "disconnected" 
            case .reasserting: statusString = "connecting" 
            @unknown default: statusString = "disconnected" 
            }
            result(statusString)
        } else {
            result("disconnected")
        }
    }
    
    
    
    enum MethodError: Error {
        case invalidArguments(String)
        
        var flutterError: FlutterError {
            switch self {
            case .invalidArguments(let message):
                return FlutterError(
                    code: "INVALID_ARGUMENTS",
                    message: message,
                    details: nil
                )
            }
        }
    }
}


enum FlutterEventStreamHandler {
    static func success(_ data: Any) -> Any {
        return data
    }
}
