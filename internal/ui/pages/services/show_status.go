package services

import (
	"context"

	"github.com/xmhrnx/serverus/internal/ops"
)

func (p *Page) showStatus(s ops.Service) {
	svc, ok := p.svc()
	if !ok {
		return
	}
	p.output.Show("info", "Status", s.Unit, func(ctx context.Context) (string, error) { return svc.Status(ctx, s.Unit) })
}
