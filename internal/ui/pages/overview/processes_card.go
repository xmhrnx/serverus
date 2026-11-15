package overview

import (
	"context"
	"fmt"
	"time"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// processesShown is how many of the busiest processes the card lists.
const processesShown = 6

// ProcessesCard lists the busiest processes, refreshed every few
// seconds, with a link to the full list.
type ProcessesCard struct {
	data *live[[]ops.Process]
	open *giode.ButtonElement
}

// NewProcessesCard returns the card for s; open shows the Processes page.
func NewProcessesCard(notify func(), s *session.Session, open func()) *ProcessesCard {
	return &ProcessesCard{
		data: newLive(notify, 3*time.Second, func(ctx context.Context) ([]ops.Process, error) {
			procs, err := s.Processes()
			if err != nil {
				return nil, err
			}
			return procs.List(ctx)
		}),
		open: cardLink("All", open),
	}
}

// Render builds the card.
func (c *ProcessesCard) Render() giode.Element {
	procs, ok, err := c.data.get()
	children := []giode.Element{widgets.CardHeader("list", theme.Primary, "Top processes", "By CPU, refreshed every 3 s", c.open)}
	switch {
	case !ok && err != nil:
		children = append(children, giode.Text(err.Error(), theme.Caption))
	case !ok:
		children = append(children, giode.Text("Reading…", theme.Caption))
	}
	for i, p := range procs {
		if i == processesShown {
			break
		}
		children = append(children, giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Stack(giode.Styles{FlexGrow: 1, Gap: 2},
				giode.Text(p.Name, giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontWeight: giode.FontWeightMedium, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}),
				giode.Text(fmt.Sprintf("%d · %s", p.PID, p.User), theme.Caption.With(giode.Styles{FontSize: theme.TextXS, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
			),
			giode.Stack(giode.Styles{Width: 90, Gap: 4},
				giode.Text("CPU "+format.Percent(p.CPU), giode.Styles{Color: theme.Muted, FontSize: theme.TextXS, FontFamily: theme.Mono}),
				widgets.Meter(min(100, p.CPU), theme.Primary, 4),
			),
			giode.Stack(giode.Styles{Width: 70, Gap: 4},
				giode.Text("RAM "+format.Percent(p.Memory), giode.Styles{Color: theme.Muted, FontSize: theme.TextXS, FontFamily: theme.Mono}),
				widgets.Meter(p.Memory, theme.Accent, 4),
			),
		))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1, Gap: theme.Space3}, children...)
}
