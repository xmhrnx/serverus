package domain

// Capability is something a connection can do. Pages are shown only for
// the capabilities a server's protocol provides.
type Capability uint

const (
	// CapTerminal opens interactive shells.
	CapTerminal Capability = 1 << iota
	// CapExec runs commands and reads their output: metrics,
	// processes, services, containers and snippets need it.
	CapExec
	// CapFiles browses and transfers files.
	CapFiles
	// CapForward forwards TCP ports.
	CapForward
)

// Has reports whether every capability in c is in s.
func (s Capability) Has(c Capability) bool { return s&c == c }

// Capabilities is what connections over p can do.
func (p Protocol) Capabilities() Capability {
	switch p {
	case ProtocolSSH:
		return CapTerminal | CapExec | CapFiles | CapForward
	case ProtocolLocal:
		return CapTerminal | CapExec | CapFiles
	case ProtocolTelnet:
		return CapTerminal
	case ProtocolFTP:
		return CapFiles
	}
	return 0
}
