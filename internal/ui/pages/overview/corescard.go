package overview

import (
	"fmt"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/bastion/internal/monitor"
	"github.com/xmhrnx/bastion/internal/ui/format"
	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// CoresCard shows each core's load as a bar.
func CoresCard(s monitor.Snapshot) giode.Element {
	cells := make([]giode.Element, len(s.CorePercents))
	for i, p := range s.CorePercents {
		cells[i] = giode.Stack(giode.Styles{Gap: 6, Width: 132},
			giode.Row(giode.Styles{Justify: giode.JustifySpaceBetween, Width: -1, Wrap: giode.WrapNoWrap},
				giode.Text(fmt.Sprintf("Core %d", i), theme.Caption),
				giode.Text(format.Percent(p), giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontFamily: theme.Mono}),
			),
			widgets.Meter(p, coreColor(p), 6),
		)
	}
	return widgets.Card(giode.Styles{Width: -1},
		widgets.CardHeader("cpu", theme.Primary, "Cores", fmt.Sprintf("%d logical CPUs", len(s.CorePercents))),
		giode.Row(giode.Styles{Gap: theme.Space4, Width: -1}, cells...),
	)
}

// coreColor keeps cores in the primary color until they run hot.
func coreColor(p float64) giode.Color {
	if p >= 75 {
		return theme.Level(p)
	}
	return theme.Primary
}
