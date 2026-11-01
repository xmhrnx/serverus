package overview

import (
	"fmt"
	"sort"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/bastion/internal/monitor"
	"github.com/xmhrnx/bastion/internal/ui/format"
	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// maxSensorsShown keeps the card a reasonable height on servers with
// dozens of probes.
const maxSensorsShown = 10

// SensorsCard lists temperatures, warmest first.
func SensorsCard(s monitor.Snapshot) giode.Element {
	sensors := append([]monitor.Sensor(nil), s.Sensors...)
	sort.Slice(sensors, func(i, j int) bool { return sensors[i].Celsius > sensors[j].Celsius })
	rows := []giode.Element{}
	for i, x := range sensors {
		if i == maxSensorsShown {
			rows = append(rows, giode.Text(fmt.Sprintf("+ %d more", len(sensors)-maxSensorsShown), theme.Caption))
			break
		}
		c := []giode.Color{theme.Success, theme.Warning, theme.Danger}[x.Level()]
		crit := x.Critical
		if crit <= 0 {
			crit = 100
		}
		rows = append(rows, giode.Stack(giode.Styles{Gap: 5, Width: -1},
			giode.Row(giode.Styles{Justify: giode.JustifySpaceBetween, Width: -1, Wrap: giode.WrapNoWrap, Gap: theme.Space2},
				giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(x.Name, theme.Caption.With(giode.Styles{TextOverflow: giode.TextOverflowEllipsis}))),
				giode.Text(format.Celsius(x.Celsius), giode.Styles{Color: c, FontSize: theme.TextSm, FontWeight: giode.FontWeightSemiBold, FontFamily: theme.Mono}),
			),
			widgets.Meter(x.Celsius/crit*100, c, 5),
		))
	}
	if len(rows) == 0 {
		rows = append(rows, giode.Text("This machine doesn't expose temperature sensors (common on virtual machines).", theme.Caption.With(giode.Styles{LineHeight: 18})))
	}
	return widgets.Card(giode.Styles{Width: -1},
		append([]giode.Element{widgets.CardHeader("thermometer", theme.Warning, "Temperatures", fmt.Sprintf("%d sensors", len(sensors)))}, rows...)...,
	)
}
