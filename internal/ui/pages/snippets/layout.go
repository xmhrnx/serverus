package snippets

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Layout draws the page.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	search, add := giode.Element(giode.Box(giode.Styles{Width: 320}, p.search)), giode.Element(p.add)
	pad := giode.Inset{Left: theme.Space4, Right: theme.Space4, Top: theme.Space4}
	if giode.Compact() {
		search, add = giode.Box(giode.Styles{Width: -1}, p.search), p.addIcon
		pad = giode.Inset{Left: theme.Space2, Right: theme.Space2, Top: theme.Space2}
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: pad},
			giode.Box(giode.Styles{FlexGrow: 1}, search),
			add),
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, p.scroll),
	), p.edit, p.broadcast, p.output).Layout(gtx)
}
