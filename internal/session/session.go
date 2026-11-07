package session

import (
	"context"
	"sync"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/transport"
)

// Hooks are called on a session's goroutines when its connection comes
// up or goes down.
type Hooks struct {
	OnConnect    func(s *Session)
	OnDisconnect func(s *Session, err error)
	// SudoPassword asks for the password sudo wants on s's server;
	// retry is set after a wrong one. It's called off the UI goroutine
	// and may block until the user answers.
	SudoPassword func(s *Session, retry bool) (string, bool)
}

// Session is the connection to one server.
type Session struct {
	dialer   transport.Dialer
	notify   func()
	interval func() time.Duration
	hooks    Hooks

	mu          sync.RWMutex
	server      domain.Server
	state       State
	err         error
	conn        transport.Conn
	collector   *monitor.Collector
	latency     time.Duration
	connectedAt time.Time
	cancel      context.CancelFunc
	generation  int

	sudoMu       sync.Mutex
	sudoPassword string
}

// Server returns the server's current configuration.
func (s *Session) Server() domain.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.server
}
