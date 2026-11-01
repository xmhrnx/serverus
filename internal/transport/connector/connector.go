// Package connector opens a connection to any server, picking the
// transport by its protocol.
package connector

import (
	"context"
	"fmt"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/transport/ftpx"
	"github.com/xmhrnx/serverus/internal/transport/local"
	"github.com/xmhrnx/serverus/internal/transport/sshx"
	"github.com/xmhrnx/serverus/internal/transport/telnet"
)

// Connector implements transport.Dialer for every protocol.
type Connector struct {
	ssh   *sshx.Dialer
	creds transport.Credentials
}

var _ transport.Dialer = (*Connector)(nil)

// New returns a Connector. opts configure SSH; their Credentials are
// used for FTP too.
func New(opts sshx.Options) *Connector {
	return &Connector{ssh: sshx.NewDialer(opts), creds: opts.Credentials}
}

// Dial connects to s.
func (c *Connector) Dial(ctx context.Context, s domain.Server) (transport.Conn, error) {
	switch s.Protocol {
	case domain.ProtocolSSH:
		return c.ssh.Dial(ctx, s)
	case domain.ProtocolLocal:
		return local.Open(), nil
	case domain.ProtocolTelnet:
		return telnet.Dial(ctx, s)
	case domain.ProtocolFTP:
		return ftpx.Dial(ctx, s, c.creds)
	}
	return nil, fmt.Errorf("unknown protocol %q", s.Protocol)
}
