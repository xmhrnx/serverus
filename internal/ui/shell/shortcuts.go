package shell

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/shortcuts"
)

// registerShortcuts registers the app-wide keyboard shortcuts for this
// frame. They get only keys nothing else took: a focused terminal keeps
// Ctrl combinations for the shell.
func (r *Root) registerShortcuts() {
	list := []giode.KeyShortcut{
		giode.Shortcut(shortcuts.Search, func() {
			r.navOn = giode.Compact()
			r.sidebar.FocusSearch()
		}),
		giode.Shortcut(shortcuts.AddServer, r.editor.Add),
		giode.Shortcut(shortcuts.Settings, func() { r.showSettings = true }),
	}
	if v, ok := r.views[r.selected]; ok && !r.showSettings {
		list = append(list, giode.Shortcut(shortcuts.NewTerminal, v.NewTerminalTab))
		for n := 1; n <= 9; n++ {
			n := n
			list = append(list, giode.Shortcut(fmt.Sprintf("%s%d", shortcuts.Page, n), func() { v.ShowNth(n) }))
		}
	}
	giode.Shortcuts(list...)
}
