package overview

import (
	"fmt"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// MemoryCard breaks memory down into what programs use, what the
// kernel caches (and gives back when asked) and what's free, with swap
// and its history.
func MemoryCard(s monitor.Snapshot, c monitor.Charts) giode.Element {
	m := s.Memory
	total := float64(max(1, m.Total))
	used := float64(m.Used())
	cache := float64(m.Buffers + m.Cached)
	free := max(0, total-used-cache)
	parts := []struct {
		label string
		value float64
		color giode.Color
	}{
		{"Used", used, theme.Accent},
		{"Cache & buffers", cache, theme.Info},
		{"Free", free, theme.SubtleBorder},
	}
	bar := []giode.Element{}
	keys := []giode.Element{}
	for _, p := range parts {
		if p.value > 0 {
			bar = append(bar, giode.Box(giode.Styles{FlexGrow: properties.FlexGrow(p.value / total * 1000), Height: 12, Background: p.color}))
		}
		keys = append(keys, legend(p.color, p.label, format.Bytes(uint64(p.value))))
	}
	children := []giode.Element{
		widgets.CardHeader("memory-stick", theme.Accent, "Memory", fmt.Sprintf("%s total", format.Bytes(m.Total))),
		giode.Row(giode.Styles{Width: -1, Wrap: giode.WrapNoWrap, BorderRadius: theme.RadiusSm, Overflow: giode.OverflowHidden, Gap: 2}, bar...),
		giode.Row(giode.Styles{Gap: theme.Space4, Width: -1}, keys...),
	}
	if m.SwapTotal > 0 {
		children = append(children,
			giode.Divider(giode.Styles{Color: theme.SubtleBorder}),
			giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				giode.Stack(giode.Styles{Gap: 6, FlexGrow: 1},
					giode.Text(fmt.Sprintf("Swap  %s of %s", format.Bytes(m.SwapUsed()), format.Bytes(m.SwapTotal)), theme.Caption),
					widgets.Meter(s.SwapPercent, theme.Level(s.SwapPercent), 6),
				),
				widgets.Sparkline(c.Swap, 96, 32, theme.Warning),
			))
	} else {
		children = append(children, giode.Text("No swap", theme.Caption))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1}, children...)
}
