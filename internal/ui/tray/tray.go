// Package tray puts Serverus in the system tray: open the window, see
// and change which servers are connected, and quit.
package tray

import (
	"sync"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/platform"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

// Tray is the tray icon and its menu.
type Tray struct {
	env  *env.Env
	tray *giode.Tray

	mu        sync.Mutex
	signature string // what the menu shows, to rebuild it only on change
}

// Install adds the tray icon, or returns nil where there's no tray
// (phones, Wayland desktops without a status notifier host).
func Install(e *env.Env, icon, template []byte) *Tray {
	if !giode.TraySupported() {
		return nil
	}
	t := &Tray{env: e, tray: e.App.Tray()}
	if platform.Mac {
		t.tray.SetTemplateIcon(template)
	} else {
		t.tray.SetIcon(icon)
	}
	t.tray.OnClick(e.App.ShowWindow)
	e.Repo.OnChange(t.Sync)
	t.Sync()
	return t
}
