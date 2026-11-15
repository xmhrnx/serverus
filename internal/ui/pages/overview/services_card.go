package overview

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// ServicesCard counts systemd services and lists the failed ones.
type ServicesCard struct {
	data *live[[]ops.Service]
	open *giode.ButtonElement
}

// NewServicesCard returns the card for s; open shows the Services page.
func NewServicesCard(notify func(), s *session.Session, open func()) *ServicesCard {
	return &ServicesCard{
		data: newLive(notify, 15*time.Second, func(ctx context.Context) ([]ops.Service, error) {
			ex, err := s.Executor()
			if err != nil {
				return nil, err
			}
			return ops.Services{Exec: ex}.List(ctx)
		}),
		open: cardLink("All", open),
	}
}

// Available reports whether the server runs systemd.
func (c *ServicesCard) Available() bool {
	_, _, err := c.data.get()
	return !errors.Is(err, ops.ErrNoSystemd)
}

// Render builds the card.
func (c *ServicesCard) Render() giode.Element {
	all, ok, err := c.data.get()
	running, failed := 0, []ops.Service{}
	for _, s := range all {
		if s.Running() {
			running++
		}
		if s.Failed() {
			failed = append(failed, s)
		}
	}
	children := []giode.Element{widgets.CardHeader("layers", theme.Info, "Services", "systemd units", c.open)}
	switch {
	case !ok && err != nil:
		return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1}, append(children, giode.Text(err.Error(), theme.Caption))...)
	case !ok:
		return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1}, append(children, giode.Text("Reading…", theme.Caption))...)
	}
	failedColor := theme.Success
	if len(failed) > 0 {
		failedColor = theme.Danger
	}
	children = append(children, giode.Row(giode.Styles{Gap: theme.Space2, Width: -1},
		widgets.Pill(fmt.Sprintf("%d running", running), theme.Success),
		widgets.Pill(fmt.Sprintf("%d failed", len(failed)), failedColor),
		widgets.Pill(fmt.Sprintf("%d in all", len(all)), theme.Muted),
	))
	if len(failed) == 0 {
		children = append(children, giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter},
			theme.Icon("circle-check", 16, theme.Success), giode.Text("No failed services", theme.Body)))
	}
	for i, s := range failed {
		if i == 5 {
			children = append(children, giode.Text(fmt.Sprintf("and %d more", len(failed)-5), theme.Caption))
			break
		}
		children = append(children, giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			widgets.StatusDot(theme.Danger, 8),
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(s.Unit, giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontFamily: theme.Mono, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
			giode.Text(s.Sub, theme.Caption),
		))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1, Gap: theme.Space3}, children...)
}
