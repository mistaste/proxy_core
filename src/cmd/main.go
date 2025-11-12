package main

import "C"
import (
	Sserver "segment/server"
)


func GRPCSERVER() bool {
	return Sserver.StartGRPCServer()
}


func ENFORCE_BINDING() {
}

func main() {
	GRPCSERVER()
	Sserver.WaitForServer()
}
