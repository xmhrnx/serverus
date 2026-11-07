package files

import (
	"path/filepath"
)

// HandleDrop uploads files and folders dropped on the window into the
// current folder.
func (p *Page) HandleDrop(paths []string) {
	fsys, err := p.fs()
	if err != nil || p.cwd == "" {
		return
	}
	dir, serverID := p.cwd, p.session.Server().ID
	go func() {
		for _, lp := range paths {
			if err := p.env.Transfers.UploadPath(serverID, fsys, lp, dir, p.refreshLater); err != nil {
				p.env.App.Post(func() { p.env.Fail("Upload "+filepath.Base(lp), err) })
			}
		}
	}()
	p.env.Info("Uploading to " + dir)
}

// refreshLater reloads the listing from a transfer's goroutine.
func (p *Page) refreshLater() { p.env.App.Post(p.refresh) }
