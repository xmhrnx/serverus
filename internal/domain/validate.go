package domain

import (
	"fmt"
	"net"
	"strings"
)

// Validate reports the first problem that would keep the server from
// connecting, or nil.
func (s Server) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("give the server a name")
	}
	if s.Protocol.NeedsHost() && strings.TrimSpace(s.Host) == "" {
		return fmt.Errorf("enter the host name or IP address")
	}
	if s.Port < 0 || s.Port > 65535 {
		return fmt.Errorf("the port must be between 1 and 65535")
	}
	if s.Protocol.NeedsCredentials() && strings.TrimSpace(s.Username) == "" && s.Protocol == ProtocolSSH {
		return fmt.Errorf("enter the user name to sign in as")
	}
	if s.Protocol == ProtocolSSH && s.Auth == AuthKey && strings.TrimSpace(s.KeyPath) == "" {
		return fmt.Errorf("choose the private key file")
	}
	if s.MAC != "" {
		if _, err := net.ParseMAC(s.MAC); err != nil {
			return fmt.Errorf("the MAC address should look like 00:1a:2b:3c:4d:5e")
		}
	}
	return nil
}

// CanWake reports whether the server can be woken with Wake-on-LAN.
func (s Server) CanWake() bool { return s.MAC != "" && s.Protocol.NeedsHost() }
