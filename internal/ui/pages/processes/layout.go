package processes

import (
	"fmt"
	"time"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Layout draws the list, reloading it periodically. Updates pause while
// a process's menu is open, so the list doesn't move under it.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	p.seen = p.shown
	p.shown = p.visible()
	next := p.data.LoadedAt().Add(refreshEvery)
	if !p.paused && !p.menu.Opened() && !p.data.Loading() && time.Now().After(next) {
		p.load()
	}
	if p.paused {
		p.pause.SetIcon("play")
		p.pause.SetTooltip("Resume updates")
	} else {
		p.pause.SetIcon("circle-stop")
		p.pause.SetTooltip("Pause updates")
	}
	compact := giode.Compact()
	toolbar := []giode.Element{giode.Box(giode.Styles{Width: 360}, p.filter)}
	if compact {
		toolbar = []giode.Element{giode.Box(giode.Styles{FlexGrow: 1, Width: -1}, p.filter)}
	} else {
		all, _ := p.data.Get()
		toolbar = append(toolbar, giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(
			fmt.Sprintf("%d processes · refreshed every %ds · right-click a process for actions", len(all), int(refreshEvery.Seconds())),
			theme.Caption.With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}))))
	}
	toolbar = append(toolbar, p.pause)
	card := []giode.Element{}
	if !compact {
		card = append(card, p.header(), giode.Divider(giode.Styles{Color: theme.SubtleBorder}))
	}
	card = append(card, giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1, Padding: giode.UniformInset(6)}, p.body(compact)))
	pad := giode.UniformInset(theme.Space4)
	if compact {
		pad = giode.UniformInset(theme.Space2)
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: pad, Gap: theme.Space3},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, toolbar...),
		widgets.FlushCard(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, card...),
	), p.menu, widgets.RefreshAt(next)).Layout(gtx)
}
