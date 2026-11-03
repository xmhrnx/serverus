package settings

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/platform"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// windowSection holds the start-at-login and tray preferences, on
// desktops; nil on phones.
func (p *Page) windowSection() giode.Element {
	if !platform.Desktop {
		return nil
	}
	p.syncLaunchAtLogin()
	rows := []giode.Element{p.login}
	note := "Serverus opens normally when started at login."
	if giode.TraySupported() {
		rows = append(rows, p.minimize)
		if platform.CloseToTray {
			rows = append(rows, p.closeTray)
		}
		note = "Started at login, Serverus waits in the tray. The tray icon opens it, connects and disconnects servers, and quits."
	}
	rows = append(rows, giode.Text(note, theme.Caption.With(giode.Styles{LineHeight: 18})))
	return section("monitor", "Window", rows...)
}
