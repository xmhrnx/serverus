package files

import (
	"context"
	"path/filepath"

	"github.com/xmhrnx/giode"
)

// pickUploads uploads files chosen in the system picker.
func (p *Page) pickUploads() {
	fsys, err := p.fs()
	if err != nil {
		return
	}
	dir, serverID := p.cwd, p.session.Server().ID
	go func() {
		files, err := p.env.App.PickFiles(context.Background(), giode.PickOptions{Multiple: true})
		if err != nil || len(files) == 0 {
			return
		}
		for _, f := range files {
			p.env.Transfers.UploadStream(serverID, fsys, f, f.Name, f.Size, dir, p.refreshLater)
		}
		p.env.App.Post(func() { p.env.Info("Uploading to " + dir) })
	}()
}

// pickFolderUpload uploads a folder chosen in the system picker, with
// everything in it, into the current folder.
func (p *Page) pickFolderUpload() {
	fsys, err := p.fs()
	if err != nil {
		return
	}
	dir, serverID := p.cwd, p.session.Server().ID
	go func() {
		local, err := p.env.App.PickFolder(context.Background(), giode.FolderOptions{Title: "Upload a folder to " + dir})
		if err != nil {
			return
		}
		if err := p.env.Transfers.UploadPath(serverID, fsys, local, dir, p.refreshLater); err != nil {
			p.env.App.Post(func() { p.env.Fail("Upload "+filepath.Base(local), err) })
			return
		}
		p.env.App.Post(func() { p.env.Info("Uploading " + filepath.Base(local) + " to " + dir) })
	}()
}
