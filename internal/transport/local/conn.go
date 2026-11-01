// Package local manages the computer the app runs on, through the same
// interfaces as a remote server: a PTY shell, commands and the file
// system.
package local

import (
	"sync"

	"github.com/xmhrnx/bastion/internal/domain"
	"github.com/xmhrnx/bastion/internal/transport"
)

// Conn is the "connection" to this computer.
type Conn struct {
	closeOnce sync.Once
	done      chan struct{}
}

var (
	_ transport.Conn       = (*Conn)(nil)
	_ transport.Terminal   = (*Conn)(nil)
	_ transport.Executor   = (*Conn)(nil)
	_ transport.FileSystem = (*Conn)(nil)
)

// Open returns a connection to this computer.
func Open() *Conn { return &Conn{done: make(chan struct{})} }

// Capabilities reports what local management supports.
func (c *Conn) Capabilities() domain.Capability { return domain.ProtocolLocal.Capabilities() }

// Wait blocks until Close.
func (c *Conn) Wait() error {
	<-c.done
	return nil
}

// Close ends the connection.
func (c *Conn) Close() error {
	c.closeOnce.Do(func() { close(c.done) })
	return nil
}
