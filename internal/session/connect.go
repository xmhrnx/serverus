package session

import (
	"context"
)

// State returns the connection state and, for Failed, why.
func (s *Session) State() (State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state, s.err
}

// Connect starts connecting unless the session is already connected or
// connecting.
func (s *Session) Connect() {
	s.mu.Lock()
	if s.state == Connecting || s.state == Connected {
		s.mu.Unlock()
		return
	}
	s.generation++
	gen := s.generation
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.state, s.err = Connecting, nil
	srv := s.server
	s.mu.Unlock()
	s.notify()

	go s.connect(ctx, gen, srv)
}
