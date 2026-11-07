package session

import (
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
)

// Capabilities reports what the server's protocol can do.
func (s *Session) Capabilities() domain.Capability {
	return s.Server().Protocol.Capabilities()
}

// Latency is the last measured round trip, 0 before the first.
func (s *Session) Latency() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.latency
}
