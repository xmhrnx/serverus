package files

import (
	"github.com/xmhrnx/serverus/internal/transport"
)

func (p *Page) lookup(name string) (transport.FileInfo, bool) {
	l, _ := p.data.Get()
	for _, e := range l.entries {
		if e.Name == name {
			return e, true
		}
	}
	return transport.FileInfo{}, false
}

// activate opens a folder, or edits or downloads a file.
func (p *Page) activate(f transport.FileInfo) {
	if f.IsDir {
		p.navigate(p.pathOf(f))
		return
	}
	if looksTextual(f.Name, f.Size) {
		p.edit(f)
		return
	}
	p.download(f, false)
}
