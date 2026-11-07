package session

import (
	"github.com/xmhrnx/serverus/internal/transport"
)

// FileSystem returns the connection's file system.
func (s *Session) FileSystem() (transport.FileSystem, error) {
	c, err := s.connection()
	if err != nil {
		return nil, err
	}
	if f, ok := c.(transport.FileSystem); ok {
		return f, nil
	}
	return nil, transport.ErrUnsupported
}

// Forwarder returns the connection's port forwarder.
func (s *Session) Forwarder() (transport.Forwarder, error) {
	c, err := s.connection()
	if err != nil {
		return nil, err
	}
	if f, ok := c.(transport.Forwarder); ok {
		return f, nil
	}
	return nil, transport.ErrUnsupported
}
