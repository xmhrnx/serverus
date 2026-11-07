// Package wol wakes sleeping and powered-off machines with Wake-on-LAN
// magic packets.
package wol

import (
	"bytes"
	"net"
)

// Packet is the magic packet for mac: six 0xff bytes, then the address
// sixteen times.
func Packet(mac net.HardwareAddr) []byte {
	p := bytes.Repeat([]byte{0xff}, 6)
	for range 16 {
		p = append(p, mac...)
	}
	return p
}
