package overview

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// UsageChart plots CPU and memory use over the last few minutes.
func UsageChart(c monitor.Charts) giode.Element {
	return widgets.Card(giode.Styles{Width: -1},
		widgets.CardHeader("activity", theme.Primary, "CPU & memory", "Last few minutes",
			legend(theme.Primary, "CPU", format.Percent(last(c.CPU))),
			legend(theme.Accent, "RAM", format.Percent(last(c.Memory))),
		),
		widgets.AreaChart(150, monitor.HistoryLength, 100,
			widgets.ChartSeries{Values: c.Memory, Color: theme.Accent, Fill: true},
			widgets.ChartSeries{Values: c.CPU, Color: theme.Primary, Fill: true},
		),
	)
}

func last(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	return v[len(v)-1]
}
