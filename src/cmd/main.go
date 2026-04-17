package main




import "C"
import (
	Sserver "segment/server"
	"segment/libtun"
)


func GRPCSERVER() bool {
	return Sserver.StartGRPCServer()
}


func ENFORCE_BINDING() {
}





func StartVPNWindows(adapterName *C.char, proxyAddress *C.char, serverIP *C.char, mtu C.int) C.int {
	err := libtun.StartWintun(
		C.GoString(adapterName),
		C.GoString(proxyAddress),
		C.GoString(serverIP),
		int(mtu),
	)
	if err != nil {
		return 1
	}
	return 0
}




func StopVPNWindows() C.int {
	libtun.Stop()
	return 0
}



func IsVPNRunning() C.int {
	if libtun.IsStarted() {
		return 1
	}
	return 0
}

func main() {
	GRPCSERVER()
	Sserver.WaitForServer()
}
