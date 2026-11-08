package services

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (p *Page) row(i int) giode.Element {
	v := p.visible()
	if i >= len(v) {
		return giode.Box(giode.Styles{})
	}
	unit := v[i].Unit
	restart := p.restarts.Get(unit, func() *widgets.ToolButton {
		return widgets.NewToolButton("rotate-cw", "Restart", func() {
			if s, ok := p.lookup(unit); ok {
				p.do(s, ops.ServiceRestart)
			}
		})
	})
	logs := p.logsBtns.Get(unit, func() *widgets.ToolButton {
		return widgets.NewToolButton("scroll-text", "Logs", func() {
			if s, ok := p.lookup(unit); ok {
				p.showLogs(s)
			}
		})
	})
	pr := p.rows.Get(unit, func() *giode.PressableElement {
		openMenu := func(e giode.ClickEvent) {
			if s, ok := p.lookup(unit); ok {
				p.selected, p.target = unit, s
				p.menu.OpenAtClick(e)
			}
		}
		return giode.Pressable(func(ps giode.PressState) giode.Element {
			s, ok := p.lookup(unit)
			if !ok {
				return giode.Box(giode.Styles{})
			}
			return serviceRow(s, p.selected == unit, ps, giode.Row(giode.Styles{Gap: 2, Wrap: giode.WrapNoWrap}, restart, logs))
		}).OnClick(func(e giode.ClickEvent) {
			p.selected = unit
			if e.Clicks == 2 {
				if s, ok := p.lookup(unit); ok {
					p.showStatus(s)
				}
			}
		}).OnSecondaryClick(openMenu).OnLongPress(openMenu)
	})
	pr.Label(unit)
	return pr
}
