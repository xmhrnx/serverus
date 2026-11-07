package session

import (
	"time"

	"github.com/xmhrnx/serverus/internal/monitor"
)

// ConnectedAt is when the connection came up.
func (s *Session) ConnectedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connectedAt
}

// Metrics returns the latest snapshot; ok is false until one exists.
func (s *Session) Metrics() (snap monitor.Snapshot, ok bool) {
	s.mu.RLock()
	c := s.collector
	s.mu.RUnlock()
	if c == nil {
		return monitor.Snapshot{}, false
	}
	return c.Snapshot()
}

// Charts returns the metrics history.
func (s *Session) Charts() monitor.Charts {
	s.mu.RLock()
	c := s.collector
	s.mu.RUnlock()
	if c == nil {
		return monitor.Charts{}
	}
	return c.Charts()
}

// MetricsErr is the last metrics collection error.
func (s *Session) MetricsErr() error {
	s.mu.RLock()
	c := s.collector
	s.mu.RUnlock()
	if c == nil {
		return nil
	}
	return c.Err()
}
