package telnet

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/transport"
)

// Conn is a Telnet connection. It has one shell: the connection itself.
type Conn struct {
	nc  net.Conn
	neg *negotiator

	mu     sync.Mutex
	opened bool

	closeOnce sync.Once
	done      chan struct{}
	err       error
}

var (
	_ transport.Conn     = (*Conn)(nil)
	_ transport.Terminal = (*Conn)(nil)
)

// Dial connects to s.
func Dial(ctx context.Context, s domain.Server) (*Conn, error) {
	d := net.Dialer{Timeout: 15 * time.Second}
	nc, err := d.DialContext(ctx, "tcp", s.Address())
	if err != nil {
		return nil, err
	}
	c := &Conn{nc: nc, done: make(chan struct{})}
	c.neg = &negotiator{w: nc}
	return c, nil
}

// Capabilities reports that Telnet only provides a terminal.
func (c *Conn) Capabilities() domain.Capability { return domain.ProtocolTelnet.Capabilities() }

// Wait blocks until the connection closes.
func (c *Conn) Wait() error {
	<-c.done
	return c.err
}
