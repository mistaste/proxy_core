
package socks5

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/netip"
	"strconv"

	"github.com/xjasonlyu/tun2socks/v2/transport/internal/bufferpool"
)


type AuthMethod = uint8


const (
	MethodNoAuth   AuthMethod = 0x00
	MethodUserPass AuthMethod = 0x02
)


const Version = 0x05


type Command uint8


const (
	CmdConnect      Command = 0x01
	CmdBind         Command = 0x02
	CmdUDPAssociate Command = 0x03
)

func (c Command) String() string {
	switch c {
	case CmdConnect:
		return "CONNECT"
	case CmdBind:
		return "BIND"
	case CmdUDPAssociate:
		return "UDP ASSOCIATE"
	default:
		return "UNDEFINED"
	}
}

type Atyp = uint8


const (
	AtypIPv4       Atyp = 0x01
	AtypDomainName Atyp = 0x03
	AtypIPv6       Atyp = 0x04
)


type Reply uint8

func (r Reply) String() string {
	switch r {
	case 0x00:
		return "succeeded"
	case 0x01:
		return "general SOCKS server failure"
	case 0x02:
		return "connection not allowed by ruleset"
	case 0x03:
		return "network unreachable"
	case 0x04:
		return "host unreachable"
	case 0x05:
		return "connection refused"
	case 0x06:
		return "TTL expired"
	case 0x07:
		return "command not supported"
	case 0x08:
		return "address type not supported"
	default:
		return fmt.Sprintf("unassigned <%#02x>", uint8(r))
	}
}


const MaxAddrLen = 1 + 1 + 255 + 2


const MaxAuthLen = 255


type Addr []byte

func (a Addr) Valid() bool {
	if len(a) < 1+1+2  {
		return false
	}

	switch a[0] {
	case AtypDomainName:
		if len(a) < 1+1+int(a[1])+2 {
			return false
		}
	case AtypIPv4:
		if len(a) < 1+net.IPv4len+2 {
			return false
		}
	case AtypIPv6:
		if len(a) < 1+net.IPv6len+2 {
			return false
		}
	}
	return true
}


func (a Addr) String() string {
	if !a.Valid() {
		return ""
	}

	var host, port string
	switch a[0] {
	case AtypDomainName:
		hostLen := int(a[1])
		host = string(a[2 : 2+hostLen])
		port = strconv.Itoa(int(binary.BigEndian.Uint16(a[2+hostLen:])))
	case AtypIPv4:
		host = net.IP(a[1 : 1+net.IPv4len]).String()
		port = strconv.Itoa(int(binary.BigEndian.Uint16(a[1+net.IPv4len:])))
	case AtypIPv6:
		host = net.IP(a[1 : 1+net.IPv6len]).String()
		port = strconv.Itoa(int(binary.BigEndian.Uint16(a[1+net.IPv6len:])))
	}
	return net.JoinHostPort(host, port)
}


func (a Addr) UDPAddr() *net.UDPAddr {
	if !a.Valid() {
		return nil
	}

	var ip []byte
	var port int
	switch a[0] {
	case AtypDomainName  :
		return nil
	case AtypIPv4:
		ip = make([]byte, net.IPv4len)
		copy(ip, a[1:1+net.IPv4len])
		port = int(binary.BigEndian.Uint16(a[1+net.IPv4len:]))
	case AtypIPv6:
		ip = make([]byte, net.IPv6len)
		copy(ip, a[1:1+net.IPv6len])
		port = int(binary.BigEndian.Uint16(a[1+net.IPv6len:]))
	}
	return &net.UDPAddr{IP: ip, Port: port}
}


type User struct {
	Username string
	Password string
}


func ClientHandshake(rw io.ReadWriter, addr Addr, command Command, user *User) (Addr, error) {
	buf := make([]byte, MaxAddrLen)

	var method uint8
	if user != nil {
		method = MethodUserPass 
	} else {
		method = MethodNoAuth 
	}

	
	if _, err := rw.Write([]byte{Version, 0x01 , method}); err != nil {
		return nil, err
	}

	
	if _, err := io.ReadFull(rw, buf[:2]); err != nil {
		return nil, err
	}

	if buf[0] != Version {
		return nil, errors.New("socks version mismatched")
	}

	if buf[1] == MethodUserPass  {
		if user == nil {
			return nil, errors.New("auth required")
		}

		uLen := len(user.Username)
		pLen := len(user.Password)

		
		if uLen == 0 || pLen == 0 {
			return nil, errors.New("auth username/password empty")
		} else if uLen > MaxAuthLen || pLen > MaxAuthLen {
			return nil, errors.New("auth username/password too long")
		}

		
		authMsg := bufferpool.Get()
		defer bufferpool.Put(authMsg)
		authMsg.WriteByte(0x01 )
		authMsg.WriteByte(byte(uLen) )
		authMsg.WriteString(user.Username )
		authMsg.WriteByte(byte(pLen) )
		authMsg.WriteString(user.Password )

		if _, err := rw.Write(authMsg.Bytes()); err != nil {
			return nil, err
		}

		if _, err := io.ReadFull(rw, buf[:2]); err != nil {
			return nil, err
		}

		if buf[1] != 0x00  {
			return nil, errors.New("rejected username/password")
		}

	} else if buf[1] != MethodNoAuth  {
		return nil, errors.New("unsupported method")
	}

	
	req := bufferpool.Get()
	defer bufferpool.Put(req)
	req.Grow(3 + MaxAddrLen)
	req.WriteByte(Version)
	req.WriteByte(byte(command))
	req.WriteByte(0x00 )
	req.Write(addr)

	if _, err := rw.Write(req.Bytes()); err != nil {
		return nil, err
	}

	
	if _, err := io.ReadFull(rw, buf[:3]); err != nil {
		return nil, err
	}

	if rep := Reply(buf[1]); rep != 0x00  {
		return nil, fmt.Errorf("%s: %s", command, rep)
	}

	return ReadAddr(rw, buf)
}

func ReadAddr(r io.Reader, b []byte) (Addr, error) {
	if len(b) < MaxAddrLen {
		return nil, io.ErrShortBuffer
	}

	
	if _, err := io.ReadFull(r, b[:1]); err != nil {
		return nil, err
	}

	switch b[0]  {
	case AtypDomainName:
		
		if _, err := io.ReadFull(r, b[1:2]); err != nil {
			return nil, err
		}
		domainLength := uint16(b[1])
		_, err := io.ReadFull(r, b[2:2+domainLength+2])
		return b[:1+1+domainLength+2], err
	case AtypIPv4:
		_, err := io.ReadFull(r, b[1:1+net.IPv4len+2])
		return b[:1+net.IPv4len+2], err
	case AtypIPv6:
		_, err := io.ReadFull(r, b[1:1+net.IPv6len+2])
		return b[:1+net.IPv6len+2], err
	default:
		return nil, errors.New("invalid address type")
	}
}


func SplitAddr(b []byte) Addr {
	addrLen := 1
	if len(b) < addrLen {
		return nil
	}

	switch b[0] {
	case AtypDomainName:
		if len(b) < 2 {
			return nil
		}
		addrLen = 1 + 1 + int(b[1]) + 2
	case AtypIPv4:
		addrLen = 1 + net.IPv4len + 2
	case AtypIPv6:
		addrLen = 1 + net.IPv6len + 2
	default:
		return nil
	}

	if len(b) < addrLen {
		return nil
	}

	return b[:addrLen]
}



func SerializeAddr(domainName string, dstIP netip.Addr, dstPort uint16) Addr {
	var (
		buf  [][]byte
		port [2]byte
	)
	binary.BigEndian.PutUint16(port[:], dstPort)

	if domainName != ""  {
		length := len(domainName)
		buf = [][]byte{{AtypDomainName, uint8(length)}, []byte(domainName), port[:]}
	} else if dstIP.Is4()  {
		buf = [][]byte{{AtypIPv4}, dstIP.AsSlice(), port[:]}
	} else  {
		buf = [][]byte{{AtypIPv6}, dstIP.AsSlice(), port[:]}
	}
	return bytes.Join(buf, nil)
}



func ParseAddr(addr net.Addr) Addr {
	if v, ok := addr.(interface {
		AddrPort() netip.AddrPort
	}); ok {
		ap := v.AddrPort()
		return SerializeAddr("", ap.Addr(), ap.Port())
	}
	return ParseAddrString(addr.String())
}


func ParseAddrString(s string) Addr {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return nil
	}

	dstPort, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil
	}

	if ip, _ := netip.ParseAddr(host); ip.IsValid() {
		return SerializeAddr("", ip, uint16(dstPort))
	}
	return SerializeAddr(host, netip.Addr{}, uint16(dstPort))
}


func DecodeUDPPacket(packet []byte) (addr Addr, payload []byte, err error) {
	if len(packet) < 5 {
		err = errors.New("insufficient length of packet")
		return
	}

	
	if !bytes.Equal(packet[:2], []byte{0x00, 0x00}) {
		err = errors.New("reserved fields should be zero")
		return
	}

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	if packet[2] != 0x00  {
		err = errors.New("discarding fragmented payload")
		return
	}

	addr = SplitAddr(packet[3:])
	if addr == nil {
		err = errors.New("socks5 UDP addr is nil")
	}

	payload = packet[3+len(addr):]
	return
}

func EncodeUDPPacket(addr Addr, payload []byte) (packet []byte, err error) {
	if addr == nil {
		return nil, errors.New("address is invalid")
	}
	packet = bytes.Join([][]byte{{0x00, 0x00, 0x00}, addr, payload}, nil)
	return
}
