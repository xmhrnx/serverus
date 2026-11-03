// Package tray puts Serverus in the system tray: open the window, see
// and change which servers are connected, and quit.
package tray

import (
	"fmt"
	"strings"
	"sync"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/platform"
	"github.com/xmhrnx/serverus/internal/session"
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

// Sync updates the menu and tooltip when servers or their connection
// states changed. It's cheap, and safe from any goroutine: sessions call
// it as they connect and drop, which keeps the menu current while the
// window is hidden and draws no frames.
func (t *Tray) Sync() {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	servers := t.env.Repo.Servers()
	var sig strings.Builder
	connected := 0
	for _, s := range servers {
		st, _ := t.env.Sessions.Get(s).State()
		if st == session.Connected {
			connected++
		}
		fmt.Fprintf(&sig, "%s|%s|%d;", s.ID, s.Name, st)
	}
	if sig.String() == t.signature {
		return
	}
	t.signature = sig.String()

	items := make([]giode.TrayItem, 0, len(servers))
	for _, s := range servers {
		sess := t.env.Sessions.Get(s)
		st, _ := sess.State()
		items = append(items, giode.TrayCheckbox(s.Name, st == session.Connected, func() {
			if st, _ := sess.State(); st == session.Connected || st == session.Connecting {
				sess.Disconnect()
			} else {
				sess.Connect()
			}
		}))
	}
	if len(items) == 0 {
		items = append(items, giode.TrayItem{Label: "No servers yet", Disabled: true})
	}
	t.tray.SetTooltip(tooltip(connected))
	t.tray.SetMenu(
		giode.TrayAction("Open Serverus", t.env.App.ShowWindow),
		giode.TraySeparator(),
		giode.TraySubmenu("Servers", items...),
		giode.TrayAction("Disconnect all", t.env.Sessions.DisconnectAll),
		giode.TraySeparator(),
		giode.TrayAction("Quit Serverus", t.env.App.Quit),
	)
}

func tooltip(connected int) string {
	switch connected {
	case 0:
		return "Serverus"
	case 1:
		return "Serverus · 1 server connected"
	}
	return fmt.Sprintf("Serverus · %d servers connected", connected)
}
