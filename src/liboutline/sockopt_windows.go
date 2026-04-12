//go:build windows

package liboutline

import (
	"syscall"
)

func reuseAddrControl(_, _ string, r syscall.RawConn) error {
	var serr error
	r.Control(func(fd uintptr) {
		serr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	})
	return serr
}
