package storage

import (
	"fmt"
	"path/filepath"
	"slices"
	"sync"

	"github.com/xmhrnx/serverus/internal/domain"
)

// Repository is the saved workspace. It's safe for concurrent use; every
// change is written to disk before the method returns.
type Repository struct {
	mu        sync.RWMutex
	path      string
	ws        Workspace
	listeners []func()
}

// OpenRepository loads the workspace from dir, creating a starter one
// when there is none yet.
func OpenRepository(dir string) (*Repository, error) {
	r := &Repository{path: filepath.Join(dir, "workspace.json")}
	ws := Workspace{}
	if err := readJSON(r.path, &ws); err != nil {
		return nil, fmt.Errorf("read workspace: %w", err)
	}
	if ws.Version == 0 {
		ws = newWorkspace()
		if err := writeJSON(r.path, ws); err != nil {
			return nil, fmt.Errorf("create workspace: %w", err)
		}
	}
	ws.Settings = ws.Settings.Normalized()
	r.ws = ws
	return r, nil
}

// OnChange registers fn to run after every saved change.
func (r *Repository) OnChange(fn func()) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listeners = append(r.listeners, fn)
}

// mutate applies fn to the workspace and saves it.
func (r *Repository) mutate(fn func(ws *Workspace)) error {
	r.mu.Lock()
	fn(&r.ws)
	err := writeJSON(r.path, r.ws)
	listeners := slices.Clone(r.listeners)
	r.mu.Unlock()
	for _, l := range listeners {
		l()
	}
	return err
}

// Servers returns a copy of every server, in the user's order.
func (r *Repository) Servers() []domain.Server {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return slices.Clone(r.ws.Servers)
}

// Server returns the server with the given ID.
func (r *Repository) Server(id string) (domain.Server, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.ws.Servers {
		if s.ID == id {
			return s, true
		}
	}
	return domain.Server{}, false
}

// SaveServer adds the server, or replaces the one with its ID.
func (r *Repository) SaveServer(s domain.Server) error {
	return r.mutate(func(ws *Workspace) {
		for i := range ws.Servers {
			if ws.Servers[i].ID == s.ID {
				ws.Servers[i] = s
				return
			}
		}
		ws.Servers = append(ws.Servers, s)
	})
}

// DeleteServer removes the server and its tunnels, and clears it as a
// jump host from the servers using it.
func (r *Repository) DeleteServer(id string) error {
	return r.mutate(func(ws *Workspace) {
		ws.Servers = slices.DeleteFunc(ws.Servers, func(s domain.Server) bool { return s.ID == id })
		ws.Tunnels = slices.DeleteFunc(ws.Tunnels, func(t domain.Tunnel) bool { return t.ServerID == id })
		for i := range ws.Servers {
			if ws.Servers[i].JumpHostID == id {
				ws.Servers[i].JumpHostID = ""
			}
		}
	})
}

// Snippets returns a copy of every snippet.
func (r *Repository) Snippets() []domain.Snippet {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return slices.Clone(r.ws.Snippets)
}

// SaveSnippet adds the snippet, or replaces the one with its ID.
func (r *Repository) SaveSnippet(s domain.Snippet) error {
	return r.mutate(func(ws *Workspace) {
		for i := range ws.Snippets {
			if ws.Snippets[i].ID == s.ID {
				ws.Snippets[i] = s
				return
			}
		}
		ws.Snippets = append(ws.Snippets, s)
	})
}

// DeleteSnippet removes the snippet.
func (r *Repository) DeleteSnippet(id string) error {
	return r.mutate(func(ws *Workspace) {
		ws.Snippets = slices.DeleteFunc(ws.Snippets, func(s domain.Snippet) bool { return s.ID == id })
	})
}

// Tunnels returns the tunnels of one server.
func (r *Repository) Tunnels(serverID string) []domain.Tunnel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Tunnel
	for _, t := range r.ws.Tunnels {
		if t.ServerID == serverID {
			out = append(out, t)
		}
	}
	return out
}

// SaveTunnel adds the tunnel, or replaces the one with its ID.
func (r *Repository) SaveTunnel(t domain.Tunnel) error {
	return r.mutate(func(ws *Workspace) {
		for i := range ws.Tunnels {
			if ws.Tunnels[i].ID == t.ID {
				ws.Tunnels[i] = t
				return
			}
		}
		ws.Tunnels = append(ws.Tunnels, t)
	})
}

// DeleteTunnel removes the tunnel.
func (r *Repository) DeleteTunnel(id string) error {
	return r.mutate(func(ws *Workspace) {
		ws.Tunnels = slices.DeleteFunc(ws.Tunnels, func(t domain.Tunnel) bool { return t.ID == id })
	})
}

// Settings returns the preferences.
func (r *Repository) Settings() domain.Settings {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ws.Settings
}

// SaveSettings replaces the preferences.
func (r *Repository) SaveSettings(s domain.Settings) error {
	return r.mutate(func(ws *Workspace) { ws.Settings = s.Normalized() })
}
