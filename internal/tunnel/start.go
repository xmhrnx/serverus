package tunnel

import (
	"context"
	"net"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/transport"
)

// Start listens on the tunnel's local address and forwards every
// connection through fw.
func (m *Manager) Start(t domain.Tunnel, fw transport.Forwarder) error {
	m.mu.Lock()
	if _, ok := m.running[t.ID]; ok {
		m.mu.Unlock()
		return nil
	}
	m.mu.Unlock()

	ln, err := net.Listen("tcp", t.LocalAddress())
	if err != nil {
		m.setErr(t.ID, err)
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	a := &active{serverID: t.ServerID, ln: ln, cancel: cancel}
	m.mu.Lock()
	m.running[t.ID] = a
	delete(m.errs, t.ID)
	m.mu.Unlock()
	m.notify()

	go m.accept(ctx, t, a, fw)
	return nil
}
