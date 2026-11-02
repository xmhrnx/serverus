package overview

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// kpiTile is a headline metric: label, big value, detail lines and a
// ring gauge.
func kpiTile(icon, label, value, detail, sub string, visual giode.Element) giode.Element {
	lines := []giode.Element{
		giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			theme.Icon(icon, 14, theme.Muted),
			giode.Text(label, theme.Caption.With(giode.Styles{FontWeight: giode.FontWeightMedium})),
		),
		giode.Text(value, theme.Metric),
		giode.Text(detail, theme.Caption.With(giode.Styles{TextOverflow: giode.TextOverflowEllipsis})),
	}
	if sub != "" {
		lines = append(lines, giode.Text(sub, theme.Caption.With(giode.Styles{Color: theme.Faint, FontSize: theme.TextXS, TextOverflow: giode.TextOverflowEllipsis})))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1, Padding: giode.UniformInset(theme.Space5)},
		giode.Row(giode.Styles{Gap: theme.Space4, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Width: -1},
			giode.Stack(giode.Styles{Gap: 4, FlexGrow: 1}, lines...),
			visual,
		),
	)
}
