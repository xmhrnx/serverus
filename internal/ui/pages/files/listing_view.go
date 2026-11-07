package files

import (
	"context"
)

func (p *Page) refresh() {
	fsys, err := p.fs()
	if err != nil {
		return
	}
	dir := p.cwd
	p.data.Reload(func(ctx context.Context) (listing, error) {
		if dir == "" {
			home, err := fsys.Home()
			if err != nil {
				home = "/"
			}
			dir = home
			p.env.App.Post(func() {
				if p.cwd == "" {
					p.cwd = home
				}
			})
		}
		entries, err := fsys.ReadDir(dir)
		return listing{dir: dir, entries: entries}, err
	})
}

func (p *Page) sort(c column) {
	if p.sortBy == c {
		p.desc = !p.desc
	} else {
		p.sortBy, p.desc = c, c != byName
	}
}

func (p *Page) toggleHidden() {
	s := p.env.Repo.Settings()
	s.ShowHiddenFiles = !s.ShowHiddenFiles
	_ = p.env.Repo.SaveSettings(s)
}
