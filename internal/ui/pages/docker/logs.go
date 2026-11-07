package docker

import (
	"context"

	"github.com/xmhrnx/serverus/internal/ops"
)

func (p *Page) showLogs(c ops.Container) {
	d, ok := p.docker()
	if !ok {
		return
	}
	p.output.Show("scroll-text", "Container logs", c.Name, func(ctx context.Context) (string, error) { return d.Logs(ctx, c.ID, 400) })
}
