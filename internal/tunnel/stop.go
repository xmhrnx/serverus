package tunnel

// Stop closes the tunnel.
func (m *Manager) Stop(id string) { m.stop(id, nil) }

func (m *Manager) stop(id string, err error) {
	m.mu.Lock()
	a, ok := m.running[id]
	delete(m.running, id)
	if err != nil {
		m.errs[id] = err
	}
	m.mu.Unlock()
	if ok {
		a.cancel()
		a.ln.Close()
	}
	m.notify()
}

// StopServer closes every tunnel of a server, when it disconnects.
func (m *Manager) StopServer(serverID string) {
	m.mu.Lock()
	var ids []string
	for id, a := range m.running {
		if a.serverID == serverID {
			ids = append(ids, id)
		}
	}
	m.mu.Unlock()
	for _, id := range ids {
		m.stop(id, nil)
	}
}
