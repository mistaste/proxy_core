package adapter



type TransportHandler interface {
	HandleTCP(TCPConn)
	HandleUDP(UDPConn)
}
