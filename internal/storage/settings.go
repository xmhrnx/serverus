package storage

import (
	"github.com/xmhrnx/serverus/internal/domain"
)

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
