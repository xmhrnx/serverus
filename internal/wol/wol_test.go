package wol

import (
	"bytes"
	"net"
	"testing"
)

func TestPacket(t *testing.T) {
	mac, _ := net.ParseMAC("00:1a:2b:3c:4d:5e")
	p := Packet(mac)
	if len(p) != 102 || !bytes.Equal(p[:6], bytes.Repeat([]byte{0xff}, 6)) || !bytes.Equal(p[96:], mac) {
		t.Fatalf("packet %x", p)
	}
}
