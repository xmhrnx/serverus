package storage

import (
	"slices"
)

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
