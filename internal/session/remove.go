package session

import (
	"time"
)

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
