package storage

import (
	"slices"

	"github.com/xmhrnx/serverus/internal/domain"
)

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
