package overview

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// diskTile is the root file system headline: use and throughput.
func diskTile(s monitor.Snapshot) giode.Element {
	d, ok := s.RootDisk()
	if !ok {
		return kpiTile("hard-drive", "Disk", "—", "No file systems", "", widgets.Ring(0, 64, 7, theme.Info))
	}
	return kpiTile("hard-drive", "Disk "+d.Mount, format.Percent(d.Percent()),
		fmt.Sprintf("%s free of %s", format.Bytes(d.Free), format.Bytes(d.Total)),
		fmt.Sprintf("Read %s · Write %s", format.Rate(s.DiskReadRate), format.Rate(s.DiskWriteRate)),
		widgets.Ring(d.Percent(), 64, 7, theme.Level(d.Percent())))
}
