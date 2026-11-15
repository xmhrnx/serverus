package fleet

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// alert is something on a server that needs a look.
type alert struct {
	server  domain.Server
	danger  bool
	message string
}

// Thresholds, in percent: warning, then danger.
const (
	warnAt   = 80
	dangerAt = 90
)

// alertsFor lists what needs attention: failed connections, and busy
// CPUs, full memory, full disks and hot sensors on connected servers.
func alertsFor(rows []row) []alert {
	var out []alert
	level := func(r row, what string, v float64) {
		switch {
		case v >= dangerAt:
			out = append(out, alert{r.server, true, fmt.Sprintf("%s at %s", what, format.Percent(v))})
		case v >= warnAt:
			out = append(out, alert{r.server, false, fmt.Sprintf("%s at %s", what, format.Percent(v))})
		}
	}
	for _, r := range rows {
		if r.state == session.Failed {
			msg := "Couldn't connect"
			if r.err != nil {
				msg = format.Sentence(transport.Explain(r.err).Error())
			}
			out = append(out, alert{r.server, true, msg})
			continue
		}
		if !r.live {
			continue
		}
		level(r, "CPU", r.snap.CPUPercent)
		level(r, "Memory", r.snap.MemPercent)
		for _, d := range r.snap.Disks {
			level(r, "Disk "+d.Mount, d.Percent())
		}
		if hot, ok := r.snap.HottestSensor(); ok && hot.Level() > 0 {
			out = append(out, alert{r.server, hot.Level() == 2, fmt.Sprintf("%s at %s", hot.Name, format.Celsius(hot.Celsius))})
		}
	}
	return out
}

// alertsCard lists alertsFor, the most serious first.
func alertsCard(rows []row) giode.Element {
	alerts := alertsFor(rows)
	children := []giode.Element{widgets.CardHeader("bell-ring", theme.Danger, "Needs attention", fmt.Sprintf("Above %d%% warns, above %d%% is serious", warnAt, dangerAt))}
	if len(alerts) == 0 {
		children = append(children, giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter},
			theme.Icon("circle-check", 16, theme.Success), giode.Text("All clear", theme.Body)))
	}
	shown := 0
	for _, serious := range []bool{true, false} {
		for _, a := range alerts {
			if a.danger != serious || shown == 10 {
				continue
			}
			shown++
			c := theme.Warning
			if a.danger {
				c = theme.Danger
			}
			children = append(children, giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				theme.Icon("triangle-alert", 14, c),
				giode.Text(a.server.Name, giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontWeight: giode.FontWeightMedium}),
				giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(a.message, theme.Caption.With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}))),
			))
		}
	}
	if len(alerts) > shown {
		children = append(children, giode.Text(fmt.Sprintf("and %d more", len(alerts)-shown), theme.Caption))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1, Gap: theme.Space3}, children...)
}
