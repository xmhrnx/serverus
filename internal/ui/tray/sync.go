package tray

import (
	"fmt"
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/session"
)

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
