package files

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/transport"
)

func (p *Page) fs() (transport.FileSystem, error) { return p.session.FileSystem() }

func (p *Page) pathOf(f transport.FileInfo) string { return transport.Join(p.cwd, f.Name) }

// Navigate opens dir, from another page (a drive's mount point).
func (p *Page) Navigate(dir string) { p.navigate(dir) }

// navigate opens dir, remembering the current one for Back.
func (p *Page) navigate(dir string) {
	if dir == p.cwd {
		return
	}
	if p.cwd != "" {
		p.back = append(p.back, p.cwd)
	}
	p.open(dir)
}

func (p *Page) open(dir string) {
	p.cwd, p.selected = dir, ""
	p.filter.Clear()
	p.refresh()
	p.list.ScrollToIndex(0, giode.ScrollAlignStart)
}

func (p *Page) goBack() {
	if n := len(p.back); n > 0 {
		dir := p.back[n-1]
		p.back = p.back[:n-1]
		p.open(dir)
	}
}

func (p *Page) goHome() {
	fsys, err := p.fs()
	if err != nil {
		return
	}
	p.env.Async(func() error {
		home, err := fsys.Home()
		if err == nil {
			p.env.App.Post(func() { p.navigate(home) })
		}
		return err
	}, func(err error) { p.env.Fail("Open the home folder", err) })
}
