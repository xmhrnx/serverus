package docker

import (
	"context"
	"strings"

	"github.com/xmhrnx/serverus/internal/ops"
)

func (p *Page) docker() (ops.Docker, bool) {
	ex, err := p.session.Executor()
	if err != nil {
		return ops.Docker{}, false
	}
	return ops.Docker{Exec: ex}, true
}

func (p *Page) load() {
	d, ok := p.docker()
	if !ok {
		return
	}
	p.data.Load(func(ctx context.Context) ([]ops.Container, error) { return d.List(ctx) })
}

func (p *Page) visible() []ops.Container {
	all, _ := p.data.Get()
	q := strings.ToLower(p.filter.Text())
	if q == "" {
		return all
	}
	var out []ops.Container
	for _, c := range all {
		if strings.Contains(strings.ToLower(c.Name+" "+c.Image+" "+c.Status), q) {
			out = append(out, c)
		}
	}
	return out
}

func (p *Page) lookup(id string) (ops.Container, bool) {
	all, _ := p.data.Get()
	for _, c := range all {
		if c.ID == id {
			return c, true
		}
	}
	return ops.Container{}, false
}
