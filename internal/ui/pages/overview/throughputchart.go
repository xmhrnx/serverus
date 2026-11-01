package overview

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// ThroughputChart plots network traffic in and out.
func ThroughputChart(s monitor.Snapshot, c monitor.Charts) giode.Element {
	return widgets.Card(giode.Styles{Width: -1},
		widgets.CardHeader("arrow-up-down", theme.Info, "Network traffic", "All interfaces",
			legend(theme.Info, "In", format.Rate(s.NetRxRate)),
			legend(theme.Warning, "Out", format.Rate(s.NetTxRate)),
		),
		widgets.AreaChart(150, monitor.HistoryLength, 0,
			widgets.ChartSeries{Values: c.NetRx, Color: theme.Info, Fill: true},
			widgets.ChartSeries{Values: c.NetTx, Color: theme.Warning},
		),
	)
}
