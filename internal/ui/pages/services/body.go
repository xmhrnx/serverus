package services

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (p *Page) body() giode.Element {
	if _, ok := p.data.Get(); !ok {
		if err := p.data.Err(); err != nil {
			if err == ops.ErrNoSystemd {
				return widgets.EmptyState("layers", "No systemd here", "This server manages services some other way (OpenRC, runit, SysV init). Use the terminal to manage them.")
			}
			return widgets.EmptyState("circle-alert", "Couldn't list services", err.Error())
		}
		return widgets.Loading("Asking systemd…")
	}
	if len(p.visible()) == 0 {
		return widgets.EmptyState("layers", "No services match", "Try another filter.")
	}
	return p.list
}
