package processes

import (
	"github.com/xmhrnx/giode"
)

func (p *Page) row(i int) giode.Element {
	if i >= len(p.shown) {
		return giode.Box(giode.Styles{})
	}
	openMenu := func(e giode.ClickEvent) {
		pr, ok := p.clicked(i)
		if !ok {
			return
		}
		p.selected, p.target = pr.PID, pr
		p.menu = p.menuFor(pr)
		p.menu.OpenAtClick(e)
	}
	pr := p.rows.Get(i, func() *giode.PressableElement {
		return giode.Pressable(func(ps giode.PressState) giode.Element {
			if i >= len(p.shown) {
				return giode.Box(giode.Styles{})
			}
			cur := p.shown[i]
			if giode.Compact() {
				return compactProcessRow(cur, p.selected == cur.PID, ps)
			}
			return processRow(cur, p.selected == cur.PID, ps)
		}).OnClick(func(giode.ClickEvent) {
			if pr, ok := p.clicked(i); ok {
				p.selected = pr.PID
			}
		}).OnSecondaryClick(openMenu).OnLongPress(openMenu)
	})
	pr.Label(p.shown[i].Name)
	return pr
}
