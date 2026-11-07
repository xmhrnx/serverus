package services

import (
	"context"
	"strings"

	"github.com/xmhrnx/serverus/internal/ops"
)

func (p *Page) svc() (ops.Services, bool) {
	ex, err := p.session.Executor()
	if err != nil {
		return ops.Services{}, false
	}
	return ops.Services{Exec: ex}, true
}

func (p *Page) load() {
	s, ok := p.svc()
	if !ok {
		return
	}
	p.data.Reload(func(ctx context.Context) ([]ops.Service, error) { return s.List(ctx) })
}

func (p *Page) visible() []ops.Service {
	all, _ := p.data.Get()
	q := strings.ToLower(p.filter.Text())
	out := make([]ops.Service, 0, len(all))
	for _, s := range all {
		switch p.show.Selected() {
		case showRunning:
			if !s.Running() {
				continue
			}
		case showFailed:
			if !s.Failed() {
				continue
			}
		case showStopped:
			if s.Running() || s.Failed() {
				continue
			}
		}
		if q != "" && !strings.Contains(strings.ToLower(s.Unit+" "+s.Description), q) {
			continue
		}
		out = append(out, s)
	}
	return out
}
