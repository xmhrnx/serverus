package overview

import (
	"sort"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// temperatureChartLines is how many sensors the chart draws.
const temperatureChartLines = 4

// TemperatureChart plots the warmest sensors over the last few minutes.
func TemperatureChart(s monitor.Snapshot, c monitor.Charts) giode.Element {
	sensors := append([]monitor.Sensor(nil), s.Sensors...)
	sort.Slice(sensors, func(i, j int) bool { return sensors[i].Celsius > sensors[j].Celsius })
	if len(sensors) > temperatureChartLines {
		sensors = sensors[:temperatureChartLines]
	}
	var series []widgets.ChartSeries
	var keys []giode.Element
	top := 0.0
	for i, x := range sensors {
		color := seriesColors[i%len(seriesColors)]
		values := c.Sensors[x.Name]
		for _, v := range values {
			top = max(top, v)
		}
		series = append(series, widgets.ChartSeries{Values: values, Color: color, Fill: i == 0})
		keys = append(keys, legend(color, x.Name, format.Celsius(x.Celsius)))
	}
	// A scale from 0 to a little over the warmest reading, at least
	// 60°C, so a cool machine doesn't look like it's running hot.
	ceiling := max(60, top*1.15)
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1},
		widgets.CardHeader("thermometer", theme.Warning, "Temperature history", "Warmest sensors, last few minutes"),
		widgets.AreaChart(150, monitor.HistoryLength, ceiling, series...),
		giode.Row(giode.Styles{Gap: theme.Space4, Width: -1}, keys...),
	)
}
