package services

import (
	"github.com/xmhrnx/serverus/internal/ops"
)

func (p *Page) lookup(unit string) (ops.Service, bool) {
	all, _ := p.data.Get()
	for _, s := range all {
		if s.Unit == unit {
			return s, true
		}
	}
	return ops.Service{}, false
}
