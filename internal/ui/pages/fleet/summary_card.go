package fleet

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// summaryCard is the fleet in numbers.
func summaryCard(rows []row) giode.Element {
	connected, failed := 0, 0
	var cpu, mem float64
	live := liveRows(rows)
	for _, r := range rows {
		switch r.state {
		case session.Connected:
			connected++
		case session.Failed:
			failed++
		}
	}
	for _, r := range live {
		cpu += r.snap.CPUPercent
		mem += r.snap.MemPercent
	}
	avg := func(v float64) string {
		if len(live) == 0 {
			return "—"
		}
		return format.Percent(v / float64(len(live)))
	}
	alerts := len(alertsFor(rows))
	alertColor := theme.Success
	if alerts > 0 {
		alertColor = theme.Danger
	}
	stat := func(icon string, c giode.Color, label, value string) giode.Element {
		return giode.Row(giode.Styles{Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Width: -1},
			widgets.IconBadge(icon, c, 38),
			giode.Stack(giode.Styles{Gap: 2},
				giode.Text(value, theme.Metric.With(giode.Styles{FontSize: theme.TextXL})),
				giode.Text(label, theme.Caption),
			))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1},
		widgets.Grid(130, theme.Space4,
			stat("server", theme.Primary, "Servers", fmt.Sprint(len(rows))),
			stat("plug", theme.Success, "Connected", fmt.Sprint(connected)),
			stat("unplug", theme.Danger, "Connection failed", fmt.Sprint(failed)),
			stat("cpu", theme.Primary, "Average CPU", avg(cpu)),
			stat("memory-stick", theme.Accent, "Average memory", avg(mem)),
			stat("bell-ring", alertColor, "Need attention", fmt.Sprint(alerts)),
		),
	)
}
