// Package ftpx manages files over FTP, or FTPS with explicit TLS.
package ftpx

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/transport"
)

// Conn is an FTP control connection. FTP runs one command at a time, so
// every operation holds mu, transfers until their stream is closed.
type Conn struct {
	mu sync.Mutex
	sc *ftp.ServerConn

	closeOnce sync.Once
	done      chan struct{}
}

var (
	_ transport.Conn       = (*Conn)(nil)
	_ transport.FileSystem = (*Conn)(nil)
	_ transport.Pinger     = (*Conn)(nil)
)

// Dial connects and signs in, anonymously when the server has no user
// name.
func Dial(ctx context.Context, s domain.Server, creds transport.Credentials) (*Conn, error) {
	opts := []ftp.DialOption{ftp.DialWithContext(ctx), ftp.DialWithTimeout(15 * time.Second)}
	if s.FTPSecure {
		opts = append(opts, ftp.DialWithExplicitTLS(&tls.Config{ServerName: s.Host, MinVersion: tls.VersionTLS12}))
	}
	sc, err := ftp.Dial(s.Address(), opts...)
	if err != nil {
		return nil, err
	}
	user, pass := "anonymous", "anonymous@"
	if s.Username != "" {
		user = s.Username
		p, ok := creds.Password(s, false)
		if !ok {
			sc.Quit()
			return nil, context.Canceled
		}
		pass = p
	}
	if err := sc.Login(user, pass); err != nil {
		sc.Quit()
		return nil, err
	}
	return &Conn{sc: sc, done: make(chan struct{})}, nil
}

// Capabilities reports that FTP only manages files.
func (c *Conn) Capabilities() domain.Capability { return domain.ProtocolFTP.Capabilities() }

// Wait blocks until Close.
func (c *Conn) Wait() error {
	<-c.done
	return nil
}

// Close signs out.
func (c *Conn) Close() error {
	c.closeOnce.Do(func() {
		close(c.done)
		c.mu.Lock()
		_ = c.sc.Quit()
		c.mu.Unlock()
	})
	return nil
}

// Ping times a NOOP.
func (c *Conn) Ping(ctx context.Context) (time.Duration, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	start := time.Now()
	err := c.sc.NoOp()
	return time.Since(start), err
}
