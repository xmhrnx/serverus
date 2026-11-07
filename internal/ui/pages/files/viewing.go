package files

import (
	"github.com/xmhrnx/serverus/internal/transport"
)

func (p *Page) edit(f transport.FileInfo) {
	if f.IsDir {
		p.env.Info("Folders can't be edited")
		return
	}
	fsys, err := p.fs()
	if err != nil {
		return
	}
	p.editor.Edit(fsys, p.pathOf(f), p.refresh)
}
