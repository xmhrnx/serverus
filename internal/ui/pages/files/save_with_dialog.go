package files

import (
	"context"

	"github.com/xmhrnx/serverus/internal/transport"
)

// saveWithDialog downloads through the system's save dialog, on
// platforms without folder access (phones).
func (p *Page) saveWithDialog(fsys transport.FileSystem, f transport.FileInfo, remote, serverID string) {
	if f.IsDir {
		p.env.App.Post(func() { p.env.Info("Download folders one file at a time on this device") })
		return
	}
	w, err := p.env.App.SaveFile(context.Background(), f.Name)
	if err != nil {
		return
	}
	p.env.Transfers.DownloadStream(serverID, fsys, f, remote, w)
}
