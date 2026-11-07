package storage

import (
	"slices"

	"github.com/xmhrnx/serverus/internal/domain"
)

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
