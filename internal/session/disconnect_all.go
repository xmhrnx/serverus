package session

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
