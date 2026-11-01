package widgets

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/bastion/internal/ui/theme"
)

// Pill is a small rounded label tinted with c: states, tags, counts.
func Pill(label string, c giode.Color) giode.Element {
	return giode.Box(giode.Styles{
		Background:   theme.Tint(c, 0.16),
		BorderRadius: theme.RadiusFull,
		Padding:      giode.SymmetricInset(8, 2),
	}, giode.Text(label, giode.Styles{Color: c, FontSize: theme.TextXS, FontWeight: giode.FontWeightSemiBold}))
}
