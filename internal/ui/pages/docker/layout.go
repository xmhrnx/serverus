package docker

import (
	"fmt"
	"time"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Layout draws the page, reloading periodically.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	next := p.data.LoadedAt().Add(refreshEvery)
	if !p.data.Loading() && time.Now().After(next) && p.data.Err() != ops.ErrNoDocker {
		p.load()
	}
	all, _ := p.data.Get()
	running := 0
	for _, c := range all {
		if c.Running() {
			running++
		}
	}
	filterWidth, pad := giode.Styles{Width: 300}, giode.UniformInset(theme.Space4)
	if giode.Compact() {
		filterWidth, pad = giode.Styles{FlexGrow: 1, Width: -1}, giode.UniformInset(theme.Space2)
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: pad, Gap: theme.Space3},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(filterWidth, p.filter),
			giode.Box(giode.Styles{FlexGrow: 1}),
			widgets.Pill(fmt.Sprintf("%d running", running), theme.Success),
			widgets.Pill(fmt.Sprintf("%d total", len(all)), theme.Muted),
		),
		widgets.FlushCard(giode.Styles{Width: -1, FlexGrow: 1, Height: -1, Padding: giode.UniformInset(6)}, p.body()),
	), p.menu, p.output, widgets.RefreshAt(next)).Layout(gtx)
}
