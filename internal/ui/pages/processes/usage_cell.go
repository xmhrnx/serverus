package processes

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// usageCell is a percentage with a small bar under it.
func usageCell(v float64, width int, c giode.Color) giode.Element {
	return giode.Box(giode.Styles{Width: properties.Width(width), Padding: giode.Inset{Right: theme.Space4}},
		giode.Stack(giode.Styles{Gap: 4, Width: -1},
			giode.Text(format.Percent(v), giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontFamily: theme.Mono}),
			widgets.Meter(v, c, 3),
		))
}
