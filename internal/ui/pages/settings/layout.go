package settings

import (
	"gioui.org/layout"
)

// Layout draws the page.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions { return p.scroll.Layout(gtx) }
