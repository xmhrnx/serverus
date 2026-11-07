package terminal

import (
	"slices"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (p *Page) tabRow() giode.Element {
	chips := make([]giode.Element, 0, len(p.tabs)+1)
	for _, t := range p.tabs {
		t := t
		closer := p.closers.Get(t.id, func() *widgets.ToolButton {
			return widgets.NewToolButton("x", "Close tab", func() { giode.AfterLayout(func() { p.closeTab(t.id) }) })
		})
		chip := p.chips.Get(t.id, func() *giode.PressableElement {
			return giode.Pressable(func(s giode.PressState) giode.Element {
				return tabChip(t.title(), p.current() == t, s.Hovered, t.ended(), closer)
			}).Label(t.title()).OnClick(func(giode.ClickEvent) {
				if j := slices.Index(p.tabs, t); j >= 0 {
					p.active = j
					t.view.Focus()
				}
			})
		})
		chips = append(chips, chip)
	}
	chips = append(chips, p.newTab)
	return giode.Row(giode.Styles{Gap: 6, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, chips...)
}
