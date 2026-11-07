package sshx

import (
	"context"
	"time"
)

// Close ends the connection, its SFTP session and the jump host's.
func (c *Conn) Close() error {
	var err error
	c.closeOnce.Do(func() {
		close(c.done)
		c.sftpMu.Lock()
		if c.sftp != nil {
			c.sftp.Close()
		}
		c.sftpMu.Unlock()
		err = c.client.Close()
		if c.jump != nil {
			c.jump.Close()
		}
	})
	return err
}

// Ping times a keep-alive round trip.
func (c *Conn) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	errc := make(chan error, 1)
	go func() {
		_, _, err := c.client.SendRequest("keepalive@openssh.com", true, nil)
		errc <- err
	}()
	select {
	case err := <-errc:
		return time.Since(start), err
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}
