// Package wol wakes sleeping and powered-off machines with Wake-on-LAN
// magic packets.
package wol

import (
	"bytes"
	"errors"
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

// Wake broadcasts the magic packet for mac on the local network, on the
// two ports network cards listen on. host, the machine's address, adds
// a broadcast to its /24 for networks that drop the global one.
func Wake(mac, host string) error {
	addr, err := net.ParseMAC(mac)
	if err != nil {
		return err
	}
	if len(addr) != 6 {
		return errors.New("Wake-on-LAN needs a 6-byte MAC address")
	}
	targets := []string{"255.255.255.255"}
	if ips, err := net.LookupIP(host); err == nil {
		for _, ip := range ips {
			if v4 := ip.To4(); v4 != nil {
				targets = append(targets, net.IPv4(v4[0], v4[1], v4[2], 255).String())
				break
			}
		}
	}
	packet := Packet(addr)
	var sent bool
	var lastErr error
	for _, t := range targets {
		for _, port := range []string{"9", "7"} {
			c, err := net.Dial("udp4", net.JoinHostPort(t, port))
			if err != nil {
				lastErr = err
				continue
			}
			if _, err := c.Write(packet); err != nil {
				lastErr = err
			} else {
				sent = true
			}
			c.Close()
		}
	}
	if !sent {
		return lastErr
	}
	return nil
}
