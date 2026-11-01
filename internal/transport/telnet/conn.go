package telnet

import (
	"bytes"
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/xmhrnx/bastion/internal/domain"
	"github.com/xmhrnx/bastion/internal/transport"
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

// Close ends the connection.
func (c *Conn) Close() error {
	var err error
	c.closeOnce.Do(func() {
		err = c.nc.Close()
		close(c.done)
	})
	return err
}

// OpenShell returns the session. Telnet carries a single session, so a
// second shell can't be opened while the first is open.
func (c *Conn) OpenShell(ctx context.Context, cols, rows int) (transport.Shell, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.opened {
		return nil, errors.New("telnet carries one session: close the open terminal first")
	}
	c.opened = true
	c.neg.resize(cols, rows)
	return &shell{c: c}, nil
}

// shell reads and writes the connection, negotiating options on the way.
type shell struct {
	c   *Conn
	buf []byte
}

func (s *shell) Read(p []byte) (int, error) {
	for {
		if s.buf == nil {
			s.buf = make([]byte, 32*1024)
		}
		// Filtering only shrinks data, so reading at most len(p) bytes
		// always fits.
		n, err := s.c.nc.Read(s.buf[:min(len(s.buf), len(p))])
		if n > 0 {
			data := s.c.neg.filter(s.buf[:n])
			if len(data) > 0 {
				return copy(p, data), nil
			}
		}
		if err != nil {
			s.c.err = err
			s.c.Close()
			return 0, err
		}
	}
}

// Write sends keystrokes, doubling 0xFF bytes as Telnet requires.
func (s *shell) Write(p []byte) (int, error) {
	data := bytes.ReplaceAll(p, []byte{cmdIAC}, []byte{cmdIAC, cmdIAC})
	if _, err := s.c.nc.Write(data); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (s *shell) Resize(cols, rows int) error {
	s.c.neg.resize(cols, rows)
	return nil
}

func (s *shell) Close() error {
	s.c.mu.Lock()
	s.c.opened = false
	s.c.mu.Unlock()
	return s.c.Close()
}
