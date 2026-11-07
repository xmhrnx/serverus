package services

import (
	"fmt"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Layout draws the page.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if _, ok := p.data.Get(); !ok && !p.data.Loading() && p.data.Err() == nil {
		p.load()
	}
	all, _ := p.data.Get()
	running, failed := 0, 0
	for _, s := range all {
		if s.Running() {
			running++
		}
		if s.Failed() {
			failed++
		}
	}
	summary := []giode.Element{widgets.Pill(fmt.Sprintf("%d running", running), theme.Success)}
	if failed > 0 {
		summary = append(summary, widgets.Pill(fmt.Sprintf("%d failed", failed), theme.Danger))
	}
	var toolbar giode.Element
	pad := giode.UniformInset(theme.Space4)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space2)
		toolbar = giode.Stack(giode.Styles{Width: -1, Gap: theme.Space2},
			giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				giode.Box(giode.Styles{FlexGrow: 1, Width: -1}, p.filter), p.refresh),
			p.show,
		)
	} else {
		toolbar = giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{Width: 300}, p.filter),
			p.show,
			giode.Row(giode.Styles{FlexGrow: 1, Gap: theme.Space2, Justify: giode.JustifyEnd, Align: giode.AlignCenter}, summary...),
			p.refresh,
		)
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: pad, Gap: theme.Space3},
		toolbar,
		widgets.FlushCard(giode.Styles{Width: -1, FlexGrow: 1, Height: -1, Padding: giode.UniformInset(6)}, p.body()),
	), p.menu, p.output).Layout(gtx)
}
