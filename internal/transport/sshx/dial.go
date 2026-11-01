package sshx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

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

func (d *Dialer) dial(ctx context.Context, s domain.Server, depth int) (*Conn, error) {
	if depth > maxJumpDepth {
		return nil, errors.New("too many jump hosts: check that they don't loop")
	}
	cfg, err := d.clientConfig(s)
	if err != nil {
		return nil, err
	}

	var jump *Conn
	var raw net.Conn
	if s.JumpHostID != "" {
		js, ok := d.opts.LookupServer(s.JumpHostID)
		if !ok {
			return nil, errors.New("the jump host no longer exists")
		}
		jump, err = d.dial(ctx, js, depth+1)
		if err != nil {
			return nil, fmt.Errorf("jump host %s: %w", js.Name, err)
		}
		raw, err = jump.client.DialContext(ctx, "tcp", s.Address())
		if err != nil {
			jump.Close()
			return nil, fmt.Errorf("reach %s from the jump host: %w", s.Address(), err)
		}
	} else {
		nd := net.Dialer{Timeout: d.opts.Timeout, KeepAlive: 30 * time.Second}
		raw, err = nd.DialContext(ctx, "tcp", s.Address())
		if err != nil {
			return nil, err
		}
	}

	// The handshake has no context: bound it with a deadline instead.
	if dl, ok := ctx.Deadline(); ok {
		_ = raw.SetDeadline(dl)
	} else {
		_ = raw.SetDeadline(time.Now().Add(d.opts.Timeout + time.Minute)) // prompts may wait on the user
	}
	cc, chans, reqs, err := ssh.NewClientConn(raw, s.Address(), cfg)
	if err != nil {
		raw.Close()
		if jump != nil {
			jump.Close()
		}
		return nil, cleanAuthError(err)
	}
	_ = raw.SetDeadline(time.Time{})
	c := newConn(ssh.NewClient(cc, chans, reqs), jump)
	if d.opts.KeepAlive > 0 {
		go c.keepAlive(d.opts.KeepAlive)
	}
	return c, nil
}

func (d *Dialer) clientConfig(s domain.Server) (*ssh.ClientConfig, error) {
	auth, err := d.authMethods(s)
	if err != nil {
		return nil, err
	}
	return &ssh.ClientConfig{
		User:            s.Username,
		Auth:            auth,
		HostKeyCallback: d.hostKeyCallback(s),
		Timeout:         d.opts.Timeout,
		ClientVersion:   "SSH-2.0-Serverus",
	}, nil
}

// cleanAuthError turns x/crypto's terse handshake errors into ones a
// user can act on.
func cleanAuthError(err error) error {
	var hk *hostKeyRejected
	if errors.As(err, &hk) {
		return hk
	}
	msg := err.Error()
	if strings.Contains(msg, "unable to authenticate") || strings.Contains(msg, "no supported methods remain") {
		return errors.New("the server rejected the credentials")
	}
	return err
}
