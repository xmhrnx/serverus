package settings

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"
)

// Name is the selected option's name.
func (s *Swatches) Name() string {
	for _, o := range s.options {
		if o.id == s.selected {
			return o.name
		}
	}
	return ""
}

// Layout draws the circles, wrapping on narrow screens.
func (s *Swatches) Layout(gtx layout.Context) layout.Dimensions {
	items := make([]giode.Element, len(s.buttons))
	for i, b := range s.buttons {
		items[i] = b
	}
	return giode.Row(giode.Styles{Gap: 4, Align: giode.AlignCenter}, items...).Layout(gtx)
}
