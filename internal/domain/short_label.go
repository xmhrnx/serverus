package domain

// ShortLabel is the protocol's name where space is tight, as in the
// server editor's connection type picker.
func (p Protocol) ShortLabel() string {
	switch p {
	case ProtocolSSH:
		return "SSH"
	case ProtocolTelnet:
		return "Telnet"
	case ProtocolFTP:
		return "FTP"
	case ProtocolLocal:
		return "Local"
	}
	return string(p)
}

// DefaultPort is the well-known port for the protocol, 0 when it has
// none.
func (p Protocol) DefaultPort() int {
	switch p {
	case ProtocolSSH:
		return 22
	case ProtocolTelnet:
		return 23
	case ProtocolFTP:
		return 21
	}
	return 0
}

// NeedsHost reports whether the protocol connects to a remote address.
func (p Protocol) NeedsHost() bool { return p != ProtocolLocal }
