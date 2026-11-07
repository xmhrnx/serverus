package domain

import (
	"fmt"
	"net"
	"strconv"
)

// Tunnel forwards a local TCP port through an SSH server to an address
// the server can reach (ssh -L).
type Tunnel struct {
	ID         string `json:"id"`
	ServerID   string `json:"serverId"`
	Name       string `json:"name"`
	BindHost   string `json:"bindHost,omitempty"`
	LocalPort  int    `json:"localPort"`
	RemoteHost string `json:"remoteHost"`
	RemotePort int    `json:"remotePort"`
	// AutoStart opens the tunnel whenever the server connects.
	AutoStart bool `json:"autoStart,omitempty"`
}

// LocalAddress is where the tunnel listens.
func (t Tunnel) LocalAddress() string {
	host := t.BindHost
	if host == "" {
		host = "127.0.0.1"
	}
	return net.JoinHostPort(host, strconv.Itoa(t.LocalPort))
}

// RemoteAddress is where the server connects for each forwarded
// connection.
func (t Tunnel) RemoteAddress() string {
	return net.JoinHostPort(t.RemoteHost, strconv.Itoa(t.RemotePort))
}

// Describe is the "localhost:5432 → db:5432" summary line.
func (t Tunnel) Describe() string {
	return fmt.Sprintf("%s → %s", t.LocalAddress(), t.RemoteAddress())
}
