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

// ContainersCard lists Docker containers, running first.
type ContainersCard struct {
	data *live[[]ops.Container]
	open *giode.ButtonElement
}

// NewContainersCard returns the card for s; open shows the Docker page.
func NewContainersCard(notify func(), s *session.Session, open func()) *ContainersCard {
	return &ContainersCard{
		data: newLive(notify, 10*time.Second, func(ctx context.Context) ([]ops.Container, error) {
			ex, err := s.Executor()
			if err != nil {
				return nil, err
			}
			return ops.Docker{Exec: ex}.List(ctx)
		}),
		open: cardLink("All", open),
	}
}

// Available reports whether the server has Docker.
func (c *ContainersCard) Available() bool {
	_, _, err := c.data.get()
	return !errors.Is(err, ops.ErrNoDocker)
}

// Render builds the card.
func (c *ContainersCard) Render() giode.Element {
	all, ok, err := c.data.get()
	running := 0
	for _, x := range all {
		if x.Running() {
			running++
		}
	}
	children := []giode.Element{widgets.CardHeader("container", theme.Info, "Containers", fmt.Sprintf("%d of %d running", running, len(all)), c.open)}
	switch {
	case !ok && err != nil:
		children = append(children, giode.Text(err.Error(), theme.Caption))
	case !ok:
		children = append(children, giode.Text("Reading…", theme.Caption))
	case len(all) == 0:
		children = append(children, giode.Text("No containers", theme.Caption))
	}
	shown := 0
	for _, pass := range []bool{true, false} {
		for _, x := range all {
			if x.Running() != pass || shown == 6 {
				continue
			}
			shown++
			dot := theme.Faint
			if x.Running() {
				dot = theme.Success
			}
			detail := x.Status
			if x.CPU != "" {
				detail = "CPU " + x.CPU + " · " + x.Memory
			}
			children = append(children, giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				widgets.StatusDot(dot, 8),
				giode.Stack(giode.Styles{FlexGrow: 1, Gap: 2},
					giode.Text(x.Name, giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontWeight: giode.FontWeightMedium, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}),
					giode.Text(x.Image, theme.Caption.With(giode.Styles{FontSize: theme.TextXS, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
				),
				giode.Text(detail, theme.Caption.With(giode.Styles{FontFamily: theme.Mono, FontSize: theme.TextXS})),
			))
		}
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1, Gap: theme.Space3}, children...)
}
