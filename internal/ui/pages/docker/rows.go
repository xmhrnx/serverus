package docker

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
	id := v[i].ID
	toggle := p.toggles.Get(id, func() *widgets.ToolButton {
		return widgets.NewToolButton("play", "Start", func() {
			if c, ok := p.lookup(id); ok {
				if c.Running() {
					p.do(c, ops.ContainerStop)
				} else {
					p.do(c, ops.ContainerStart)
				}
			}
		})
	})
	logs := p.logs.Get(id, func() *widgets.ToolButton {
		return widgets.NewToolButton("scroll-text", "Logs", func() {
			if c, ok := p.lookup(id); ok {
				p.showLogs(c)
			}
		})
	})
	pr := p.rows.Get(id, func() *giode.PressableElement {
		openMenu := func(e giode.ClickEvent) {
			if c, ok := p.lookup(id); ok {
				p.selected, p.target = id, c
				p.menu.OpenAtClick(e)
			}
		}
		return giode.Pressable(func(ps giode.PressState) giode.Element {
			c, ok := p.lookup(id)
			if !ok {
				return giode.Box(giode.Styles{})
			}
			if c.Running() {
				toggle.SetIcon("square")
				toggle.SetTooltip("Stop")
			} else {
				toggle.SetIcon("play")
				toggle.SetTooltip("Start")
			}
			return containerRow(c, p.selected == id, ps, giode.Row(giode.Styles{Gap: 2, Wrap: giode.WrapNoWrap}, toggle, logs))
		}).OnClick(func(giode.ClickEvent) { p.selected = id }).
			OnSecondaryClick(openMenu).OnLongPress(openMenu)
	})
	pr.Label(v[i].Name)
	return pr
}
