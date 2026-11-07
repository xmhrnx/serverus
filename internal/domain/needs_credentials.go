package domain

// NeedsCredentials reports whether the protocol authenticates a user.
func (p Protocol) NeedsCredentials() bool { return p == ProtocolSSH || p == ProtocolFTP }
