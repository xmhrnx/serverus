package domain

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Server is one machine the user manages.
type Server struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Group    string   `json:"group,omitempty"`
	Protocol Protocol `json:"protocol"`
	Host     string   `json:"host,omitempty"`
	Port     int      `json:"port,omitempty"`
	Username string   `json:"username,omitempty"`

	Auth AuthMethod `json:"auth,omitempty"`
	// KeyPath is the private key file for AuthKey.
	KeyPath string `json:"keyPath,omitempty"`
	// JumpHostID is the ID of another SSH server to tunnel through
	// (a serverus host), empty to connect directly.
	JumpHostID string `json:"jumpHostId,omitempty"`
	// FTPSecure upgrades FTP connections with explicit TLS (FTPS).
	FTPSecure bool `json:"ftpSecure,omitempty"`

	// Accent is the color of the server's avatar in the sidebar.
	Accent    string    `json:"accent,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	Favorite  bool      `json:"favorite,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

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
	return nil
}
