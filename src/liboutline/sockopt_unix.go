//go:build !windows

package liboutline

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func reuseAddrControl(_, _ string, r syscall.RawConn) error {
	var serr error
	r.Control(func(fd uintptr) {
		serr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if serr != nil {
			return
		}
		_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1)
	})
	return serr
}
