package overview

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

func (p *Page) content() giode.Element {
	pad := giode.UniformInset(theme.Space6)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space3)
	}
	hint := "Live readings and their last few minutes"
	if p.board.Editing() {
		hint = "Resize, move or remove cards, or add new ones. The layout is shared by every server."
	}
	header := giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(hint, theme.Caption.With(giode.Styles{MaxLines: 2}))),
		p.board.Controls(),
	)
	return giode.Stack(giode.Styles{Gap: theme.Space4, Padding: pad, Width: -1}, header, p.board.Grid())
}
