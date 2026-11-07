package sshx

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/transport"
)

// Dialer opens SSH connections.
type Dialer struct {
	opts Options
}

// NewDialer returns a Dialer using opts.
func NewDialer(opts Options) *Dialer {
	if opts.Timeout == 0 {
		opts.Timeout = 15 * time.Second
	}
	return &Dialer{opts: opts}
}

// Dial connects to s, through its jump host when it has one.
func (d *Dialer) Dial(ctx context.Context, s domain.Server) (transport.Conn, error) {
	return d.dial(ctx, s, 0)
}
