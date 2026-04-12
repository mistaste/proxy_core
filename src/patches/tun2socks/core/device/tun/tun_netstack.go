

package tun

import (
	"fmt"

	"golang.org/x/sys/unix"
	"gvisor.dev/gvisor/pkg/rawfile"
	"gvisor.dev/gvisor/pkg/tcpip/link/fdbased"
	"gvisor.dev/gvisor/pkg/tcpip/link/tun"
	"gvisor.dev/gvisor/pkg/tcpip/stack"

	"github.com/xjasonlyu/tun2socks/v2/core/device"
)

type TUN struct {
	stack.LinkEndpoint

	fd   int
	mtu  uint32
	name string
}

func Open(name string, mtu uint32) (device.Device, error) {
	t := &TUN{name: name, mtu: mtu}

	if len(t.name) >= unix.IFNAMSIZ {
		return nil, fmt.Errorf("interface name too long: %s", t.name)
	}

	fd, err := tun.Open(t.name)
	if err != nil {
		return nil, fmt.Errorf("create tun: %w", err)
	}
	t.fd = fd

	if t.mtu > 0 {
		if err := setMTU(t.name, t.mtu); err != nil {
			return nil, fmt.Errorf("set mtu: %w", err)
		}
	}

	_mtu, err := rawfile.GetMTU(t.name)
	if err != nil {
		return nil, fmt.Errorf("get mtu: %w", err)
	}
	t.mtu = _mtu

	ep, err := fdbased.New(&fdbased.Options{
		FDs: []int{fd},
		MTU: t.mtu,
		
		EthernetHeader: false,
		
		PacketDispatchMode: fdbased.Readv,
		
		
		
		
		
		
		
		
		MaxSyscallHeaderBytes: 0x00,
	})
	if err != nil {
		return nil, fmt.Errorf("create endpoint: %w", err)
	}
	t.LinkEndpoint = ep

	return t, nil
}

func (t *TUN) Name() string {
	return t.name
}

func (t *TUN) Close() {
	defer t.LinkEndpoint.Close()
	_ = unix.Close(t.fd)
}

func setMTU(name string, n uint32) error {
	
	fd, err := unix.Socket(
		unix.AF_INET,
		unix.SOCK_DGRAM,
		0,
	)
	if err != nil {
		return err
	}

	defer unix.Close(fd)

	ifr, err := unix.NewIfreq(name)
	if err != nil {
		return err
	}
	ifr.SetUint32(n)
	return unix.IoctlIfreq(fd, unix.SIOCSIFMTU, ifr)
}
