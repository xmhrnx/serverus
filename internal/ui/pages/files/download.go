package files

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/transport"
)

// download saves f (a folder recursively) to the download folder, asking
// for it the first time or when ask is set.
func (p *Page) download(f transport.FileInfo, ask bool) {
	fsys, err := p.fs()
	if err != nil {
		return
	}
	remote := p.pathOf(f)
	serverID := p.session.Server().ID
	dir := p.downloadDir
	go func() {
		if dir == "" || ask {
			start := dir
			if start == "" {
				if home, err := os.UserHomeDir(); err == nil {
					start = filepath.Join(home, "Downloads")
				}
			}
			chosen, err := p.env.App.PickFolder(context.Background(), giode.FolderOptions{Title: "Download " + f.Name + " to", Directory: start})
			switch {
			case errors.Is(err, giode.ErrNotSupported):
				p.saveWithDialog(fsys, f, remote, serverID)
				return
			case err != nil:
				return
			}
			dir = chosen
		}
		err := p.env.Transfers.DownloadPath(serverID, fsys, f, remote, dir)
		p.env.App.Post(func() {
			if err != nil {
				p.env.Fail("Download "+f.Name, err)
				return
			}
			p.downloadDir = dir
			p.env.Info("Downloading " + f.Name + " to " + dir)
		})
	}()
}
