package terminal

import (
	"gioui.org/layout"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"

	term "github.com/xmhrnx/serverus/internal/ui/terminal"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// keyBarKey is one key of the bar.
type keyBarKey struct {
	label string
	code  rune
	mod   uv.KeyMod
}

// keyBarKeys are the keys on-screen keyboards lack.
var keyBarKeys = []keyBarKey{
	{label: "Esc", code: uv.KeyEscape},
	{label: "Tab", code: uv.KeyTab},
	{label: "↑", code: uv.KeyUp},
	{label: "↓", code: uv.KeyDown},
	{label: "←", code: uv.KeyLeft},
	{label: "→", code: uv.KeyRight},
	{label: "Home", code: uv.KeyHome},
	{label: "End", code: uv.KeyEnd},
	{label: "PgUp", code: uv.KeyPgUp},
	{label: "PgDn", code: uv.KeyPgDown},
	{label: "^C", code: 'c', mod: uv.ModCtrl},
	{label: "^D", code: 'd', mod: uv.ModCtrl},
	{label: "^Z", code: 'z', mod: uv.ModCtrl},
}

// KeyBar is a row of the keys a touch keyboard lacks: Esc, Tab, a
// sticky Ctrl, arrows and common Ctrl combinations. It scrolls sideways
// when it doesn't fit.
type KeyBar struct {
	view   func() *term.View
	ctrl   *giode.PressableElement
	keys   []*giode.PressableElement
	scroll *scroll.Scroll
}

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

// keyCap draws one key.
func keyCap(label string, on, pressed bool) giode.Element {
	bg, fg := theme.Raised, theme.Text
	switch {
	case on:
		bg, fg = theme.Primary, theme.OnPrimary
	case pressed:
		bg = theme.PressFill
	}
	return giode.Box(giode.Styles{
		Padding: giode.SymmetricInset(14, 10), BorderRadius: theme.RadiusSm, Background: bg,
		BorderColor: theme.Border, BorderWidth: 1, MinWidth: 44, Justify: giode.JustifyCenter,
	}, giode.Text(label, giode.Styles{Color: fg, FontSize: theme.TextSm, FontFamily: theme.Mono, FontWeight: giode.FontWeightSemiBold}))
}
