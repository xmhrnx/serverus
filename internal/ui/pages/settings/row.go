package settings

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// row is one preference: its title and hint on the left, the control
// on the right.
func row(title, hint string, control giode.Element) giode.Element {
	text := []giode.Element{giode.Text(title, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium})}
	if hint != "" {
		text = append(text, giode.Text(hint, theme.Caption))
	}
	if giode.Compact() {
		// On a phone the control goes under its title, full width.
		return giode.Stack(giode.Styles{Width: -1, Gap: theme.Space2, Padding: giode.SymmetricInset(0, 6)},
			giode.Stack(giode.Styles{Gap: 3}, text...),
			giode.Box(giode.Styles{Width: -1}, control))
	}
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space4, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.SymmetricInset(0, 6)},
		giode.Stack(giode.Styles{FlexGrow: 1, Gap: 3}, text...), control)
}
