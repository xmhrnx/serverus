package overview

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// LoadChart plots the 1-minute load average against the number of
// cores, the line above which work waits for a CPU.
func LoadChart(s monitor.Snapshot, c monitor.Charts) giode.Element {
	cores := float64(max(1, len(s.CorePercents)))
	top := cores
	for _, v := range c.Load {
		top = max(top, v)
	}
	limit := make([]float64, len(c.Load))
	for i := range limit {
		limit[i] = cores
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1},
		widgets.CardHeader("activity", theme.Primary, "Load average", fmt.Sprintf("1 · 5 · 15 min: %.2f · %.2f · %.2f", s.Load[0], s.Load[1], s.Load[2]),
			legend(theme.Primary, "Now", fmt.Sprintf("%.2f", s.Load[0])),
			legend(theme.Faint, "Cores", fmt.Sprintf("%.0f", cores)),
		),
		widgets.AreaChart(150, monitor.HistoryLength, top*1.2,
			widgets.ChartSeries{Values: limit, Color: theme.Faint},
			widgets.ChartSeries{Values: c.Load, Color: theme.Primary, Fill: true},
		),
	)
}
