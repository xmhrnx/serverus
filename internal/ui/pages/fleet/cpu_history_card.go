package fleet

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// historyLines is how many servers a history chart draws.
const historyLines = 8

// historyCard draws one metric of every connected server over the last
// few minutes, each in the server's color.
func historyCard(icon, title string, rows []row, pick func(monitor.Charts) []float64, ceiling float64, value func(row) string) giode.Element {
	var series []widgets.ChartSeries
	var keys []giode.Element
	for i, r := range liveRows(rows) {
		if i == historyLines {
			break
		}
		c := giode.Color(r.server.Accent)
		if c == "" {
			c = giode.Color(domain.AccentFor(r.server.ID))
		}
		series = append(series, widgets.ChartSeries{Values: pick(r.charts), Color: c})
		keys = append(keys, giode.Row(giode.Styles{Gap: 6, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{Width: 10, Height: 3, BorderRadius: 2, Background: c}),
			giode.Text(r.server.Name, theme.Caption),
			giode.Text(value(r), giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontWeight: giode.FontWeightSemiBold, FontFamily: theme.Mono}),
		))
	}
	children := []giode.Element{widgets.CardHeader(icon, theme.Primary, title, "Every connected server, last few minutes")}
	if len(series) == 0 {
		children = append(children, giode.Text("Connect servers to see them here.", theme.Caption))
	} else {
		children = append(children, widgets.AreaChart(150, monitor.HistoryLength, ceiling, series...),
			giode.Row(giode.Styles{Gap: theme.Space4, Width: -1}, keys...))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1}, children...)
}

func cpuValue(r row) string { return format.Percent(r.snap.CPUPercent) }
func memValue(r row) string { return format.Percent(r.snap.MemPercent) }
