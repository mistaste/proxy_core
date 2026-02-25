import Foundation
import NetworkExtension
import Combine
import os.log



protocol VpnServiceProtocol {
    func prepareVPN(
        appName: String, appTunnelBundle: String,
        completion: @escaping (Result<Void, Error>) -> Void)
    func startVPN(
        startMode: String,
        coreName: String?,
        config: String?,
        cacheDir: String?,
        port: Int32,
        appName: String,
        appTunnelBundle: String,
        completion: @escaping (Result<Void, Error>) -> Void
    )

    func stopVPN(completion: @escaping (Result<Void, Error>) -> Void)
    func sendTunnelMessage(_ messageDict: [String: String], completion: ((String?) -> Void)?)
    var vpnManager: NETunnelProviderManager? { get set }

}


protocol VpnStatusDelegate: AnyObject {
    func vpnStatusDidChange(_ status: NEVPNStatus)
}

enum VpnServiceError: Error, LocalizedError {
    case managerNotInitialized
    case invalidSession
    case messageEncodingFailed
    case alreadyStopped

    var errorDescription: String? {
        switch self {
        case .managerNotInitialized:
            return "VPN Manager not initialized"
        case .invalidSession:
            return "Invalid VPN connection session"
        case .messageEncodingFailed:
            return "Failed to encode tunnel message"
        case .alreadyStopped:
            return "VPN is already stopped"
        }
    }
}

class VpnService: VpnServiceProtocol {
    private var _vpnManager: NETunnelProviderManager?
    
    
    var vpnManager: NETunnelProviderManager? {
        get { return _vpnManager }
        set { _vpnManager = newValue }
    }
    
    private var cancellables = Set<AnyCancellable>()
    private var vpnStatusCancellable: AnyCancellable?
    private let logger = Logger(
        subsystem: Bundle.main.bundleIdentifier ?? "ProxyCore", category: "VpnService")

    
    private var lastLoggedStatus: NEVPNStatus?
    
    
    weak var statusDelegate: VpnStatusDelegate?

    static let shared = VpnService()

    private init() {
        
        
        
        
        
        
        
    }
    
    deinit {
        NotificationCenter.default.removeObserver(self)
    }
    
    @objc private func appWillEnterForeground() {
        
        loadManager { [weak self] success in
            guard let self = self, success, let manager = self._vpnManager else { return }
            
            
            self.observeVPNStatus(manager)
            
            
            let status = manager.connection.status
            self.updateVPNStatus(status)
            
            
            if let session = manager.connection as? NETunnelProviderSession {
                self.sendTunnelMessage(["command": "IS_CORE_RUNNING"]) { response in
                    if response?.lowercased() == "true" {
                        
                        DispatchQueue.main.async {
                            
                            if status != .connected {
                                self.logger.info("Core is running but status is not connected, forcing connected status")
                                self.statusDelegate?.vpnStatusDidChange(.connected)
                            }
                        }
                    }
                }
            }
        }
    }
    
    func loadManager(completion: @escaping (Bool) -> Void) {
        Task {
            do {
                let managers = try await NETunnelProviderManager.loadAllFromPreferences()
                if let manager = managers.first {
                    self._vpnManager = manager
                    self.logger.info("Loaded existing VPN manager")
                    completion(true)
                } else {
                    self.logger.warning("No existing VPN manager found")
                    completion(false)
                }
            } catch {
                self.logger.error("Failed to load managers: \(error.localizedDescription)")
                completion(false)
            }
        }
    }

    func prepareVPN(
        appName: String, appTunnelBundle: String,
        completion: @escaping (Result<Void, Error>) -> Void
    ) {
        Task {
            do {
                let managers = try await NETunnelProviderManager.loadAllFromPreferences()

                if managers.isEmpty {
                    self._vpnManager = NETunnelProviderManager()
                    do {
                        try configureVPNManager(
                            self._vpnManager!, appName: appName, appTunnelBundle: appTunnelBundle)
                        try await self._vpnManager?.saveToPreferences()
                    } catch {
                        completion(.failure(error))
                        return
                    }
                    os_log("❤️ manager is empty")

                } else {
                    os_log("❤️❤️ manager is not empty")
                    self._vpnManager = managers.first
                }

                guard let manager = self._vpnManager else {
                    completion(.failure(VpnServiceError.managerNotInitialized))
                    return
                }

                try await manager.loadFromPreferences()
                self.observeVPNStatus(manager)
                completion(.success(()))
            } catch {
                self.logger.error("Failed to prepare VPN: \(error.localizedDescription)")
                completion(.failure(error))
            }
        }
    }

    private func enableVPNManager(_ manager: NETunnelProviderManager, appName: String, appTunnelBundle: String) async throws {
        do {
                try await manager.loadFromPreferences()
                
                try configureVPNManager(
                    manager, appName: appName, appTunnelBundle: appTunnelBundle)
                
                try await manager.saveToPreferences()
        } catch {
            print(error.localizedDescription)
        }
    }
    
    
    func startVPN(
        startMode: String,
        coreName: String?,
        config: String?,
        cacheDir: String?,
        port: Int32,
        appName: String,
        appTunnelBundle: String,
        completion: @escaping (Result<Void, Error>) -> Void
    ) {
        Task {
            do {
                guard let manager = _vpnManager else {
                    throw VpnServiceError.managerNotInitialized
                }
                try await manager.loadFromPreferences()

                try configureVPNManager(
                    manager, appName: appName, appTunnelBundle: appTunnelBundle)
                try await manager.saveToPreferences()
                
                
                if manager.connection.status == .connected || manager.connection.status == .connecting {
                    manager.connection.stopVPNTunnel()
                    self.logger.info("Stopping existing connection before reconnecting")
                    
                }

                let options: [String: NSObject] = [
                    "startMode": startMode as NSString,
                    "coreName": (coreName ?? "") as NSString,
                    "port": NSNumber(value: port),
                    "address": "127.0.0.1" as NSString,
                    "mtu": NSNumber(value: 1500),
                    "config": (config ?? "") as NSString,
                    "cacheDir": (cacheDir ?? "") as NSString,
                ]
                try await Task.sleep(nanoseconds: 1_000_000_000) 


                try manager.connection.startVPNTunnel(options: options)
                self.logger.info("VPN tunnel initiation successful")
                completion(.success(()))
            } catch {
                self.logger.error("VPN Start Error: \(error.localizedDescription)")
                completion(.failure(error))
            }
        }
    }


    func stopVPN(completion: @escaping (Result<Void, Error>) -> Void) {
        Task {
            do {
                guard let manager = _vpnManager else {
                    throw VpnServiceError.managerNotInitialized
                }

                try await manager.loadFromPreferences()

                try await stopVpnTunnel(manager)
                completion(.success(()))
            } catch {
                self.logger.error("Failed to stop VPN: \(error.localizedDescription)")
                completion(.failure(error))
            }
        }
    }

    func sendTunnelMessage(_ messageDict: [String: String], completion: ((String?) -> Void)? = nil)
    {
        Task {
            do {
                guard let manager = _vpnManager else {
                    self.logger.error("VPN Manager is not initialized")
                    completion?(nil)
                    return
                }

                try await manager.loadFromPreferences()

                guard let session = manager.connection as? NETunnelProviderSession else {
                    self.logger.error("Invalid VPN connection")
                    completion?(nil)
                    return
                }

                guard
                    let data = try? JSONSerialization.data(withJSONObject: messageDict, options: [])
                else {
                    self.logger.error("Failed to encode message")
                    completion?(nil)
                    return
                }

                try session.sendProviderMessage(data) { responseData in
                    if let responseData = responseData,
                        let responseString = String(data: responseData, encoding: .utf8)
                    {
                        completion?(responseString)
                    } else {
                        self.logger.warning("No response or invalid response from tunnel")
                        completion?(nil)
                    }
                }
            } catch {
                self.logger.error("Error sending tunnel message: \(error.localizedDescription)")
                completion?(nil)
            }
        }
    }

    private func configureVPNManager(
        _ manager: NETunnelProviderManager, appName: String, appTunnelBundle: String
    ) throws {
        manager.localizedDescription = appName

        let protocolConfig = NETunnelProviderProtocol()
        protocolConfig.providerBundleIdentifier = appTunnelBundle
        protocolConfig.serverAddress = appName

        let configData: [String: Any] = [
            "address": "127.0.0.1",
            "port": 2080,
            "mtu": 1500,
        ]

        guard let data = try? JSONSerialization.data(withJSONObject: configData) else {
            throw VpnServiceError.messageEncodingFailed
        }

        protocolConfig.providerConfiguration = ["config": data]
        protocolConfig.excludeLocalNetworks = true

        manager.protocolConfiguration = protocolConfig
        manager.isEnabled = true
        manager.isOnDemandEnabled = false
        manager.onDemandRules = []
    }

    private func observeVPNStatus(_ manager: NETunnelProviderManager) {
        
        vpnStatusCancellable?.cancel()
        
        vpnStatusCancellable = NotificationCenter.default.publisher(for: .NEVPNStatusDidChange, object: manager.connection)
            .sink { [weak self] _ in
                guard let self = self else { return }
                let status = manager.connection.status
                self.updateVPNStatus(status)
            }
    }

    private func stopVpnTunnel(_ manager: NETunnelProviderManager) async throws {
        switch manager.connection.status {
                case .connected, .connecting, .reasserting:
                    manager.connection.stopVPNTunnel()
                    self.logger.info("Stopping existing connection before reconnecting")
                default:
                    
                    self.logger.info("VPN is already stopped. Current status: \(String(manager.connection.status.rawValue))")
                }
    }

    private func updateVPNStatus(_ status: NEVPNStatus) {
        
        if lastLoggedStatus != status {
            switch status {
            case .connected:
                logger.info("VPN Connected")
            case .connecting:
                logger.info("VPN Connecting...")
            case .disconnecting:
                logger.info("VPN Disconnecting...")
            case .disconnected:
                logger.info("VPN Disconnected")
            case .reasserting:
                logger.info("VPN Reasserting...")
            case .invalid:
                logger.error("VPN Status Invalid")
            @unknown default:
                logger.warning("VPN Unknown Status: \(String(status.rawValue))")
            }
            
            
            lastLoggedStatus = status
            
            
            statusDelegate?.vpnStatusDidChange(status)
        }
    }
}
