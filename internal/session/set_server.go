package session

import (
	"github.com/xmhrnx/serverus/internal/domain"
)

// SetServer updates the configuration, used on the next connection.
func (s *Session) SetServer(srv domain.Server) {
	s.mu.Lock()
	s.server = srv
	s.mu.Unlock()
}
