// Package domain holds the app's core types: servers, snippets, tunnels
// and settings. It has no dependencies on the UI or on any transport.
package domain

// Protocol is the way the app reaches a server.
type Protocol string

const (
	// ProtocolSSH is a full SSH connection: terminal, metrics, SFTP
	// file management and port forwarding.
	ProtocolSSH Protocol = "ssh"
	// ProtocolTelnet opens a plain-text terminal only.
	ProtocolTelnet Protocol = "telnet"
	// ProtocolFTP manages files over FTP or explicit FTPS.
	ProtocolFTP Protocol = "ftp"
	// ProtocolLocal manages the computer the app runs on.
	ProtocolLocal Protocol = "local"
)

// Protocols lists every protocol in the order forms offer them.
var Protocols = []Protocol{ProtocolSSH, ProtocolLocal, ProtocolTelnet, ProtocolFTP}

// Label is the protocol's display name.
func (p Protocol) Label() string {
	switch p {
	case ProtocolSSH:
		return "SSH"
	case ProtocolTelnet:
		return "Telnet"
	case ProtocolFTP:
		return "FTP / FTPS"
	case ProtocolLocal:
		return "This computer"
	}
	return string(p)
}
