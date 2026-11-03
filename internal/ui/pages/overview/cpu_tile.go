package overview

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// cpuTile is the CPU headline: total use and I/O wait.
func cpuTile(s monitor.Snapshot) giode.Element {
	return kpiTile("cpu", "CPU", format.Percent(s.CPUPercent),
		fmt.Sprintf("%d cores", len(s.CorePercents)),
		fmt.Sprintf("I/O wait %s", format.Percent(s.IOWait)),
		widgets.Ring(s.CPUPercent, 64, 7, theme.Primary))
}
