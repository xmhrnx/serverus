package tunnel

import (
	"errors"
	"net"
)

// Status returns a tunnel's state.
func (m *Manager) Status(id string) Status {
	m.mu.Lock()
	defer m.mu.Unlock()
	if a, ok := m.running[id]; ok {
		return Status{Running: true, Connections: a.open.Load(), Total: a.total.Load()}
	}
	return Status{Err: m.errs[id]}
}

func (m *Manager) setErr(id string, err error) {
	if errors.Is(err, net.ErrClosed) {
		return
	}
	m.mu.Lock()
	m.errs[id] = err
	m.mu.Unlock()
	m.notify()
}
