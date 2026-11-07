package terminal

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"
)

func (b *KeyBar) row() giode.Element {
	keys := []giode.Element{b.ctrl}
	for _, k := range b.keys {
		keys = append(keys, k)
	}
	return giode.Row(giode.Styles{Gap: 6, Wrap: giode.WrapNoWrap}, keys...)
}

// Layout draws the bar.
func (b *KeyBar) Layout(gtx layout.Context) layout.Dimensions {
	return giode.Box(giode.Styles{Width: -1}, b.scroll).Layout(gtx)
}
