package dialogs

import (
	"github.com/xmhrnx/serverus/internal/domain"
)

// protocolHint completes the line under the connection type picker.
func protocolHint(p domain.Protocol) string {
	switch p {
	case domain.ProtocolSSH:
		return ": terminal, monitoring, files and tunnels"
	case domain.ProtocolLocal:
		return ": manage the computer Serverus runs on"
	case domain.ProtocolTelnet:
		return ": a terminal for network devices and old systems"
	case domain.ProtocolFTP:
		return ": file management, optionally over TLS"
	}
	return ""
}
