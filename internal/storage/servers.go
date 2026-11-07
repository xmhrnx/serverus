package storage

import (
	"slices"

	"github.com/xmhrnx/serverus/internal/domain"
)

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
