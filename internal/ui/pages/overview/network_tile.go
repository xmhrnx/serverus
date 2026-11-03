package overview

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// networkTile is the network headline, for machines without sensors.
func networkTile(s monitor.Snapshot, history []float64) giode.Element {
	return kpiTile("network", "Network", format.Rate(s.NetRxRate+s.NetTxRate),
		fmt.Sprintf("↓ %s  ↑ %s", format.Rate(s.NetRxRate), format.Rate(s.NetTxRate)),
		fmt.Sprintf("%d interfaces", len(s.Net)), widgets.Sparkline(history, 88, 44, theme.Info))
}
