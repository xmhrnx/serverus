package terminal

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// keyCap draws one key.
func keyCap(label string, on, pressed bool) giode.Element {
	bg, fg := theme.Raised, theme.Text
	switch {
	case on:
		bg, fg = theme.Primary, theme.OnPrimary
	case pressed:
		bg = theme.PressFill
	}
	return giode.Box(giode.Styles{
		Padding: giode.SymmetricInset(14, 10), BorderRadius: theme.RadiusSm, Background: bg,
		BorderColor: theme.Border, BorderWidth: 1, MinWidth: 44, Justify: giode.JustifyCenter,
	}, giode.Text(label, giode.Styles{Color: fg, FontSize: theme.TextSm, FontFamily: theme.Mono, FontWeight: giode.FontWeightSemiBold}))
}
