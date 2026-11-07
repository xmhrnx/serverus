package drives

import (
	"context"

	"github.com/xmhrnx/serverus/internal/ops"
)

func (p *Page) load() {
	d, ok := p.drives()
	if !ok {
		return
	}
	p.data.Reload(func(ctx context.Context) ([]ops.Drive, error) { return d.List(ctx) })
}
