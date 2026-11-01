package session

import (
	"sync"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/transport"
)

// Manager owns the sessions of every server.
type Manager struct {
	dialer   transport.Dialer
	notify   func()
	interval func() time.Duration
	hooks    Hooks

	mu       sync.Mutex
	sessions map[string]*Session
}

// NewManager returns a manager dialing with dialer. notify runs after
// any session changes; interval is read when a collector starts.
func NewManager(dialer transport.Dialer, notify func(), interval func() time.Duration, hooks Hooks) *Manager {
	return &Manager{dialer: dialer, notify: notify, interval: interval, hooks: hooks, sessions: map[string]*Session{}}
}

// Get returns the server's session, creating it disconnected.
func (m *Manager) Get(srv domain.Server) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[srv.ID]; ok {
		return s
	}
	s := &Session{dialer: m.dialer, notify: m.notify, interval: m.interval, hooks: m.hooks, server: srv}
	m.sessions[srv.ID] = s
	return s
}

// Lookup returns the server's session if it has one.
func (m *Manager) Lookup(id string) (*Session, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[id]
	return s, ok
}

// Remove disconnects and forgets a server's session.
func (m *Manager) Remove(id string) {
	m.mu.Lock()
	s, ok := m.sessions[id]
	delete(m.sessions, id)
	m.mu.Unlock()
	if ok {
		s.Disconnect()
	}
}

// SetInterval changes the metrics interval of every session.
func (m *Manager) SetInterval(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range m.sessions {
		s.SetInterval(d)
	}
}

// DisconnectAll closes every connection, at shutdown.
func (m *Manager) DisconnectAll() {
	m.mu.Lock()
	all := make([]*Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		all = append(all, s)
	}
	m.mu.Unlock()
	for _, s := range all {
		s.Disconnect()
	}
}
