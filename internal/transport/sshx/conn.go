package sshx

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/transport"
)

// Conn is an open SSH connection. It implements transport.Terminal,
// Executor, FileSystem (over SFTP), Forwarder and Pinger.
type Conn struct {
	client *ssh.Client
	jump   *Conn

	sftpMu sync.Mutex
	sftp   *sftp.Client

	closeOnce sync.Once
	done      chan struct{}
}

var (
	_ transport.Conn       = (*Conn)(nil)
	_ transport.Terminal   = (*Conn)(nil)
	_ transport.Executor   = (*Conn)(nil)
	_ transport.FileSystem = (*Conn)(nil)
	_ transport.Forwarder  = (*Conn)(nil)
	_ transport.Pinger     = (*Conn)(nil)
)

func newConn(client *ssh.Client, jump *Conn) *Conn {
	return &Conn{client: client, jump: jump, done: make(chan struct{})}
}

// Capabilities reports everything SSH can do.
func (c *Conn) Capabilities() domain.Capability { return domain.ProtocolSSH.Capabilities() }

// Wait blocks until the connection drops.
func (c *Conn) Wait() error { return c.client.Wait() }

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

// keepAlive pings the server so NATs and firewalls keep the connection
// open, and closes it when the server stops answering.
func (c *Conn) keepAlive(every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	missed := 0
	for {
		select {
		case <-c.done:
			return
		case <-t.C:
			ctx, cancel := context.WithTimeout(context.Background(), every)
			_, err := c.Ping(ctx)
			cancel()
			if err != nil {
				missed++
				if missed >= 3 {
					c.Close()
					return
				}
				continue
			}
			missed = 0
		}
	}
}

// sftpClient opens the SFTP subsystem on first use.
func (c *Conn) sftpClient() (*sftp.Client, error) {
	c.sftpMu.Lock()
	defer c.sftpMu.Unlock()
	if c.sftp != nil {
		return c.sftp, nil
	}
	sc, err := sftp.NewClient(c.client, sftp.UseConcurrentReads(true), sftp.UseConcurrentWrites(true))
	if err != nil {
		return nil, err
	}
	c.sftp = sc
	go func() {
		_ = sc.Wait()
		c.sftpMu.Lock()
		if c.sftp == sc {
			c.sftp = nil
		}
		c.sftpMu.Unlock()
	}()
	return sc, nil
}
