package overview

import (
	"fmt"
	"time"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// SystemHeader shows who the machine is: host name, OS and kernel, and
// chips with uptime, load, processes and the CPU.
func SystemHeader(s monitor.Snapshot, latency time.Duration) giode.Element {
	chips := []giode.Element{
		chip("clock", "Up "+format.Duration(s.System.Uptime)),
		chip("activity", fmt.Sprintf("Load %.2f · %.2f · %.2f", s.Load[0], s.Load[1], s.Load[2])),
		chip("list", fmt.Sprintf("%d processes", s.Processes)),
	}
	if s.System.CPUModel != "" {
		chips = append(chips, chip("cpu", fmt.Sprintf("%d × %s", s.System.Cores, s.System.CPUModel)))
	}
	if latency > 0 {
		chips = append(chips, chip("zap", format.Latency(latency)+" round trip"))
	}
	sub := s.System.OS
	if s.System.Kernel != "" {
		sub += "  ·  " + s.System.Kernel
	}
	return widgets.Card(giode.Styles{
		BackgroundGradient: giode.Linear(120, theme.Tint(theme.Primary, 0.14), theme.Raised, theme.Raised),
	},
		giode.Row(giode.Styles{Gap: theme.Space4, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Width: -1},
			widgets.IconBadge("monitor", theme.Primary, 48),
			giode.Stack(giode.Styles{Gap: 4, FlexGrow: 1},
				giode.Text(s.System.Hostname, theme.Title),
				giode.Text(sub, theme.Caption.With(giode.Styles{TextOverflow: giode.TextOverflowEllipsis})),
			),
		),
		giode.Row(giode.Styles{Gap: theme.Space2, Width: -1}, chips...),
	)
}

func chip(icon, label string) giode.Element {
	return giode.Row(giode.Styles{
		Gap: 6, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: giode.SymmetricInset(10, 5), BorderRadius: theme.RadiusFull,
		Background: theme.Sunken, BorderColor: theme.SubtleBorder, BorderWidth: 1,
	},
		theme.Icon(icon, 13, theme.Muted),
		giode.Text(label, giode.Styles{Color: theme.Text, FontSize: theme.TextSm}),
	)
}
