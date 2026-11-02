package overview

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// sensorReading is one sensor: its name and temperature over a bar
// filled toward its critical point.
func sensorReading(x monitor.Sensor) giode.Element {
	c := []giode.Color{theme.Success, theme.Warning, theme.Danger}[x.Level()]
	crit := x.Critical
	if crit <= 0 {
		crit = 100
	}
	return giode.Stack(giode.Styles{Gap: 5, Width: -1},
		giode.Row(giode.Styles{Justify: giode.JustifySpaceBetween, Width: -1, Wrap: giode.WrapNoWrap, Gap: theme.Space2},
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(x.Name, theme.Caption.With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}))),
			giode.Text(format.Celsius(x.Celsius), giode.Styles{Color: c, FontSize: theme.TextSm, FontWeight: giode.FontWeightSemiBold, FontFamily: theme.Mono}),
		),
		widgets.Meter(x.Celsius/crit*100, c, 5),
	)
}
