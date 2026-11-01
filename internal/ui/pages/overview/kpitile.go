package overview

import (
	"fmt"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/bastion/internal/monitor"
	"github.com/xmhrnx/bastion/internal/ui/format"
	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// kpiTile is a headline metric: label, big value, detail lines and a
// ring gauge.
func kpiTile(icon, label, value, detail, sub string, visual giode.Element) giode.Element {
	lines := []giode.Element{
		giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			theme.Icon(icon, 14, theme.Muted),
			giode.Text(label, theme.Caption.With(giode.Styles{FontWeight: giode.FontWeightMedium})),
		),
		giode.Text(value, theme.Metric),
		giode.Text(detail, theme.Caption.With(giode.Styles{TextOverflow: giode.TextOverflowEllipsis})),
	}
	if sub != "" {
		lines = append(lines, giode.Text(sub, theme.Caption.With(giode.Styles{Color: theme.Faint, FontSize: theme.TextXS, TextOverflow: giode.TextOverflowEllipsis})))
	}
	return widgets.Card(giode.Styles{Width: -1, Padding: giode.UniformInset(theme.Space5)},
		giode.Row(giode.Styles{Gap: theme.Space4, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Width: -1},
			giode.Stack(giode.Styles{Gap: 4, FlexGrow: 1}, lines...),
			visual,
		),
	)
}

func cpuTile(s monitor.Snapshot) giode.Element {
	return kpiTile("cpu", "CPU", format.Percent(s.CPUPercent),
		fmt.Sprintf("%d cores", len(s.CorePercents)),
		fmt.Sprintf("I/O wait %s", format.Percent(s.IOWait)),
		widgets.Ring(s.CPUPercent, 64, 7, theme.Primary))
}

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

func temperatureTile(hot monitor.Sensor) giode.Element {
	crit := hot.Critical
	if crit <= 0 {
		crit = 100
	}
	c := []giode.Color{theme.Success, theme.Warning, theme.Danger}[hot.Level()]
	return kpiTile("thermometer", "Temperature", format.Celsius(hot.Celsius), hot.Name,
		fmt.Sprintf("Critical at %s", format.Celsius(crit)), widgets.Ring(hot.Celsius/crit*100, 64, 7, c))
}

func networkTile(s monitor.Snapshot, history []float64) giode.Element {
	return kpiTile("network", "Network", format.Rate(s.NetRxRate+s.NetTxRate),
		fmt.Sprintf("↓ %s  ↑ %s", format.Rate(s.NetRxRate), format.Rate(s.NetTxRate)),
		fmt.Sprintf("%d interfaces", len(s.Net)), widgets.Sparkline(history, 88, 44, theme.Info))
}
