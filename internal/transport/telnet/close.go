package telnet

import (
	"context"
	"errors"

	"github.com/xmhrnx/serverus/internal/transport"
)

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
