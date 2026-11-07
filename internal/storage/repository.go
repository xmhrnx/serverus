package storage

import (
	"fmt"
	"path/filepath"
	"sync"
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
