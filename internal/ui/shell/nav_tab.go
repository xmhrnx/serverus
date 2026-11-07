package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

func navTab(s pageSpec, active, hovered bool) giode.Element {
	fg, underline := theme.Muted, giode.Color("transparent")
	if hovered {
		fg = theme.Text
	}
	if active {
		fg, underline = theme.Text, theme.Primary
	}
	return giode.Box(giode.Styles{},
		giode.Row(giode.Styles{Gap: 7, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.SymmetricInset(12, 12)},
			theme.Icon(s.icon, 15, fg),
			giode.Text(s.label, giode.Styles{Color: fg, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium}),
		),
		giode.Box(giode.Styles{
			Position: giode.PositionAbsolute, Left: 8, Right: 8, Bottom: giode.Flush,
			Height: 2, BorderRadius: 2, Background: underline,
		}),
	)
}
