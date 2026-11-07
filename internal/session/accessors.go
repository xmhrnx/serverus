package session

import (
	"errors"

	"github.com/xmhrnx/serverus/internal/transport"
)

func (s *Session) connection() (transport.Conn, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.conn == nil {
		return nil, errors.New("not connected")
	}
	return s.conn, nil
}

// Terminal returns the connection's shell opener.
func (s *Session) Terminal() (transport.Terminal, error) {
	c, err := s.connection()
	if err != nil {
		return nil, err
	}
	if t, ok := c.(transport.Terminal); ok {
		return t, nil
	}
	return nil, transport.ErrUnsupported
}

// Executor returns the connection's command runner.
func (s *Session) Executor() (transport.Executor, error) {
	c, err := s.connection()
	if err != nil {
		return nil, err
	}
	e, ok := c.(transport.Executor)
	if !ok {
		return nil, transport.ErrUnsupported
	}
	if in, ok := c.(transport.InputRunner); ok && s.hooks.SudoPassword != nil {
		return sudoExecutor{run: e, in: in, s: s}, nil
	}
	return e, nil
}
