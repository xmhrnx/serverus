package overview

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// DiskIOChart plots disk reads and writes over the last few minutes.
func DiskIOChart(s monitor.Snapshot, c monitor.Charts) giode.Element {
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1},
		widgets.CardHeader("hard-drive", theme.Info, "Disk activity", "Every disk, last few minutes",
			legend(theme.Accent, "Read", format.Rate(s.DiskReadRate)),
			legend(theme.Danger, "Write", format.Rate(s.DiskWriteRate)),
		),
		widgets.AreaChart(150, monitor.HistoryLength, 0,
			widgets.ChartSeries{Values: c.DiskRead, Color: theme.Accent, Fill: true},
			widgets.ChartSeries{Values: c.DiskWrite, Color: theme.Danger},
		),
	)
}
