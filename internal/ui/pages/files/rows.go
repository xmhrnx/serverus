package files

import (
	"github.com/xmhrnx/giode"
)

func (p *Page) row(i int) giode.Element {
	v := p.visible()
	if i >= len(v) {
		return giode.Box(giode.Styles{})
	}
	f := v[i]
	key := p.cwd + "/" + f.Name
	pr := p.rows.Get(key, func() *giode.PressableElement {
		name := f.Name
		openMenu := func(e giode.ClickEvent) {
			if cur, ok := p.lookup(name); ok {
				p.selected, p.target = name, cur
				p.menu.OpenAtClick(e)
			}
		}
		return giode.Pressable(func(ps giode.PressState) giode.Element {
			cur, ok := p.lookup(name)
			if !ok {
				return giode.Box(giode.Styles{})
			}
			if giode.Compact() {
				return compactFileRow(cur, p.selected == name, ps)
			}
			return fileRow(cur, p.selected == name, ps)
		}).OnClick(func(e giode.ClickEvent) {
			p.selected = name
			if e.Clicks == 2 {
				if cur, ok := p.lookup(name); ok {
					p.activate(cur)
				}
			}
		}).OnSecondaryClick(openMenu).OnLongPress(openMenu)
	})
	pr.Label(f.Name)
	return pr
}
