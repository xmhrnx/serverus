package widgets

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/bastion/internal/ui/theme"
)

// FormField labels a control, with an optional hint underneath.
func FormField(label, hint string, control giode.Element) giode.Element {
	children := []giode.Element{
		giode.Text(label, giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontWeight: giode.FontWeightMedium}),
		control,
	}
	if hint != "" {
		children = append(children, giode.Text(hint, theme.Caption.With(giode.Styles{FontSize: theme.TextXS})))
	}
	return giode.Stack(giode.Styles{Gap: 6, Width: -1}, children...)
}
