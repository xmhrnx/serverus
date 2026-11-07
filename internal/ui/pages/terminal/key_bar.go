package terminal

import (
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"
	term "github.com/xmhrnx/serverus/internal/ui/terminal"
)

// keyBarKey is one key of the bar.
type keyBarKey struct {
	label string
	code  rune
	mod   uv.KeyMod
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
