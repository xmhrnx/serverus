package processes

import (
	"context"

	"github.com/xmhrnx/serverus/internal/ops"
)

func (p *Page) sort(c int) {
	if p.sortBy == c {
		p.desc = !p.desc
	} else {
		p.sortBy, p.desc = c, c == byCPU || c == byMemory || c == byTime
	}
}

func (p *Page) load() {
	procs, err := p.session.Processes()
	if err != nil {
		return
	}
	p.data.Load(func(ctx context.Context) ([]ops.Process, error) { return procs.List(ctx) })
}
