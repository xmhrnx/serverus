package sshx

import (
	"sync"

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
