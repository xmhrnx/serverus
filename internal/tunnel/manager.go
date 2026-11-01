// Package tunnel runs local port forwards (ssh -L) through SSH
// connections.
package tunnel

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/xmhrnx/bastion/internal/domain"
	"github.com/xmhrnx/bastion/internal/transport"
)

// Status is a tunnel's live state.
type Status struct {
	Running     bool
	Err         error
	Connections int64 // open forwarded connections
	Total       int64 // connections forwarded since it started
}

// Manager runs tunnels. It's safe for concurrent use.
type Manager struct {
	notify func()

	mu      sync.Mutex
	running map[string]*active
	errs    map[string]error
}

type active struct {
	serverID string
	ln       net.Listener
	cancel   context.CancelFunc
	open     atomic.Int64
	total    atomic.Int64
}

// NewManager returns a manager; notify runs when a tunnel changes.
func NewManager(notify func()) *Manager {
	return &Manager{notify: notify, running: map[string]*active{}, errs: map[string]error{}}
}

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

func (m *Manager) accept(ctx context.Context, t domain.Tunnel, a *active, fw transport.Forwarder) {
	for {
		c, err := a.ln.Accept()
		if err != nil {
			if ctx.Err() == nil {
				m.stop(t.ID, err)
			}
			return
		}
		go m.forward(ctx, t, a, fw, c)
	}
}

func (m *Manager) forward(ctx context.Context, t domain.Tunnel, a *active, fw transport.Forwarder, local net.Conn) {
	defer local.Close()
	remote, err := fw.DialRemote(ctx, t.RemoteAddress())
	if err != nil {
		m.setErr(t.ID, err)
		return
	}
	defer remote.Close()
	a.open.Add(1)
	a.total.Add(1)
	m.notify()
	defer func() {
		a.open.Add(-1)
		m.notify()
	}()

	done := make(chan struct{}, 2)
	go func() { _, _ = io.Copy(remote, local); done <- struct{}{} }()
	go func() { _, _ = io.Copy(local, remote); done <- struct{}{} }()
	select {
	case <-done:
	case <-ctx.Done():
	}
}

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
