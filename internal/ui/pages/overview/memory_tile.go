package overview

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// memoryTile is the memory headline: use and swap.
func memoryTile(s monitor.Snapshot) giode.Element {
	m := s.Memory
	swap := "No swap"
	if m.SwapTotal > 0 {
		swap = fmt.Sprintf("Swap %s of %s", format.Bytes(m.SwapUsed()), format.Bytes(m.SwapTotal))
	}
	return kpiTile("memory-stick", "Memory", format.Percent(s.MemPercent),
		fmt.Sprintf("%s of %s", format.Bytes(m.Used()), format.Bytes(m.Total)),
		swap, widgets.Ring(s.MemPercent, 64, 7, theme.Accent))
}
