package domain

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Address is host:port for dialing, with the protocol's default port
// when none is set.
func (s Server) Address() string {
	port := s.Port
	if port == 0 {
		port = s.Protocol.DefaultPort()
	}
	return net.JoinHostPort(s.Host, strconv.Itoa(port))
}

// Subtitle is the short "user@host" line shown under the name.
func (s Server) Subtitle() string {
	if !s.Protocol.NeedsHost() {
		return "Local machine"
	}
	host := s.Host
	if s.Port != 0 && s.Port != s.Protocol.DefaultPort() {
		host = fmt.Sprintf("%s:%d", s.Host, s.Port)
	}
	if s.Username != "" {
		return s.Username + "@" + host
	}
	return host
}

// DisplayGroup is the sidebar group the server is listed under.
func (s Server) DisplayGroup() string {
	if g := strings.TrimSpace(s.Group); g != "" {
		return g
	}
	return "Servers"
}
