package drives

import (
	"fmt"
	"time"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Busy reports whether an action is running on the device at path.
func (p *Page) Busy(path string) bool { return p.busy[path] }

// Layout draws the page.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if !p.data.Loading() && (p.data.LoadedAt().IsZero() && p.data.Err() == nil || time.Since(p.data.LoadedAt()) > rescanEvery && len(p.busy) == 0) {
		p.load()
	}
	pad := giode.UniformInset(theme.Space4)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space2)
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: pad, Gap: theme.Space3},
		p.toolbar(),
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, p.scroll),
	), p.mountTo, p.inUse).Layout(gtx)
}

func (p *Page) toolbar() giode.Element {
	all, _ := p.data.Get()
	external := 0
	for _, d := range all {
		if d.External() {
			external++
		}
	}
	summary := widgets.Pill(fmt.Sprintf("%d external", external), theme.Primary)
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		p.show,
		giode.Row(giode.Styles{FlexGrow: 1, Gap: theme.Space2, Justify: giode.JustifyEnd, Align: giode.AlignCenter}, summary),
		p.refresh,
	)
}
