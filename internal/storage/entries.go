package storage

import (
	"slices"

	"github.com/xmhrnx/serverus/internal/domain"
)

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
