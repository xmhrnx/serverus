package domain

import (
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
	// MAC is the server's network card address, for waking it with
	// Wake-on-LAN. It's learned on the first SSH connection.
	MAC string `json:"mac,omitempty"`

	// Accent is the color of the server's avatar in the sidebar.
	Accent    string    `json:"accent,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	Favorite  bool      `json:"favorite,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
