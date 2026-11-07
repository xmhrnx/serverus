// Package tunnel runs local port forwards (ssh -L) through SSH
// connections.
package tunnel

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
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
