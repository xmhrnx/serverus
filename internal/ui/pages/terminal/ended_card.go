package terminal

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// endedCard offers a new shell once one exits; phones stack the message
// above the buttons.
func (p *Page) endedCard(msg string) giode.Element {
	style := giode.Styles{
		FlexDirection: giode.FlexDirectionRow, Align: giode.AlignCenter, Gap: theme.Space3,
		Padding: giode.SymmetricInset(theme.Space4, theme.Space3), BoxShadow: giode.BoxShadow{Y: 8, Blur: 24, Color: "#00000080"},
	}
	if giode.Compact() {
		style.FlexDirection = giode.FlexDirectionColumn
		return widgets.Card(style, giode.Text(msg, theme.Body),
			giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, p.reopen, p.closeEnded))
	}
	return widgets.Card(style, giode.Text(msg, theme.Body), p.reopen, p.closeEnded)
}
