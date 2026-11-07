package terminal

import (
	"github.com/xmhrnx/giode"
	term "github.com/xmhrnx/serverus/internal/ui/terminal"
)

// NewKeyBar returns a bar typing into the view view returns.
func NewKeyBar(view func() *term.View) *KeyBar {
	b := &KeyBar{view: view}
	b.ctrl = giode.Pressable(func(s giode.PressState) giode.Element {
		v := b.view()
		return keyCap("Ctrl", v != nil && v.CtrlLatched(), s.Pressed)
	}).Label("Ctrl").Focusable(false).OnClick(func(giode.ClickEvent) {
		if v := b.view(); v != nil {
			v.SetCtrlLatch(!v.CtrlLatched())
		}
	})
	for _, k := range keyBarKeys {
		k := k
		b.keys = append(b.keys, giode.Pressable(func(s giode.PressState) giode.Element {
			return keyCap(k.label, false, s.Pressed)
		}).Label(k.label).Focusable(false).OnClick(func(giode.ClickEvent) {
			if v := b.view(); v != nil {
				v.SendKey(k.code, k.mod)
			}
		}))
	}
	b.scroll = giode.Scroll(b.row, giode.Styles{ScrollBar: giode.ScrollBar{Hidden: true}}).Axis(giode.Horizontal)
	return b
}
