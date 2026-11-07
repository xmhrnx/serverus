package dialogs

import (
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
)

// Add opens the editor for a new server.
func (d *ServerDialog) Add() {
	id := domain.NewID()
	d.load(domain.Server{ID: id, Protocol: domain.ProtocolSSH, Auth: domain.AuthPassword, Accent: domain.AccentFor(id)}, true)
}

// Edit opens the editor for an existing server.
func (d *ServerDialog) Edit(s domain.Server) { d.load(s, false) }

// Duplicate opens the editor on a copy of s.
func (d *ServerDialog) Duplicate(s domain.Server) {
	s.ID = domain.NewID()
	s.Name += " copy"
	s.CreatedAt = time.Time{}
	d.load(s, true)
}
