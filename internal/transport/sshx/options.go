// Package sshx connects to servers over SSH: shells, commands, SFTP and
// port forwarding, directly or through a jump host.
package sshx

import (
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/storage"
	"github.com/xmhrnx/serverus/internal/transport"
)

// Options configure a Dialer.
type Options struct {
	Credentials transport.Credentials
	HostKeys    *storage.KnownHosts
	Prompt      transport.HostKeyPrompt
	// LookupServer finds jump hosts by ID.
	LookupServer func(id string) (domain.Server, bool)
	// Timeout bounds the TCP connect and the SSH handshake.
	Timeout time.Duration
	// KeepAlive is the interval of keep-alive requests; 0 disables
	// them.
	KeepAlive time.Duration
}

// maxJumpDepth stops jump host chains that loop back on themselves.
const maxJumpDepth = 4
