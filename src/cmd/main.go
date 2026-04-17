package main

/*
#include <stdlib.h>
*/
import "C"
import (
	Sserver "segment/server"
	"segment/libtun"
)

//export GRPCSERVER
func GRPCSERVER() bool {
	return Sserver.StartGRPCServer()
}

//export ENFORCE_BINDING
func ENFORCE_BINDING() {
}

// StartVPNWindows launches tun2socks bound to a wintun adapter and installs
// Windows routing so traffic flows through the local SOCKS5 proxy.
// Returns 0 on success, non-zero error code otherwise.
//export StartVPNWindows
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

// StopVPNWindows stops the tun2socks engine and tears down routes added
// by StartVPNWindows.
//export StopVPNWindows
func StopVPNWindows() C.int {
	libtun.Stop()
	return 0
}

// IsVPNRunning reports whether tun2socks is currently active.
//export IsVPNRunning
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
