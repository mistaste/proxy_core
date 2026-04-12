package shadowstream

import (
	"crypto/rand"
	"errors"
	"io"
	"net"

	"github.com/xjasonlyu/tun2socks/v2/buffer"
)


var ErrShortPacket = errors.New("short packet")




func Pack(dst, plaintext []byte, s Cipher) ([]byte, error) {
	if len(dst) < s.IVSize()+len(plaintext) {
		return nil, io.ErrShortBuffer
	}
	iv := dst[:s.IVSize()]
	_, err := rand.Read(iv)
	if err != nil {
		return nil, err
	}
	s.Encrypter(iv).XORKeyStream(dst[len(iv):], plaintext)
	return dst[:len(iv)+len(plaintext)], nil
}



func Unpack(dst, pkt []byte, s Cipher) ([]byte, error) {
	if len(pkt) < s.IVSize() {
		return nil, ErrShortPacket
	}
	if len(dst) < len(pkt)-s.IVSize() {
		return nil, io.ErrShortBuffer
	}
	iv := pkt[:s.IVSize()]
	s.Decrypter(iv).XORKeyStream(dst, pkt[len(iv):])
	return dst[:len(pkt)-len(iv)], nil
}

type PacketConn struct {
	net.PacketConn
	Cipher
}


func NewPacketConn(c net.PacketConn, ciph Cipher) *PacketConn {
	return &PacketConn{PacketConn: c, Cipher: ciph}
}

const maxPacketSize = 64 * 1024

func (c *PacketConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	buf := buffer.Get(maxPacketSize)
	defer buffer.Put(buf)
	buf, err := Pack(buf, b, c.Cipher)
	if err != nil {
		return 0, err
	}
	_, err = c.PacketConn.WriteTo(buf, addr)
	return len(b), err
}

func (c *PacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	n, addr, err := c.PacketConn.ReadFrom(b)
	if err != nil {
		return n, addr, err
	}
	bb, err := Unpack(b[c.IVSize():], b[:n], c.Cipher)
	if err != nil {
		return n, addr, err
	}
	copy(b, bb)
	return len(bb), addr, err
}
