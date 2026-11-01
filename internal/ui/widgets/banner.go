package widgets

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// Banner is an inline notice: an icon, a message and optional actions,
// tinted with c (theme.Danger for errors, theme.Warning, theme.Info).
func Banner(icon string, c giode.Color, message string, actions ...giode.Element) giode.Element {
	row := []giode.Element{
		theme.Icon(icon, 18, c),
		giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(message, theme.Body.With(giode.Styles{LineHeight: 19}))),
	}
	row = append(row, actions...)
	return giode.Row(giode.Styles{
		Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding:      giode.SymmetricInset(theme.Space4, theme.Space3),
		Background:   theme.Tint(c, 0.10),
		BorderColor:  theme.Tint(c, 0.35),
		BorderWidth:  1,
		BorderRadius: theme.RadiusMd,
	}, row...)
}
