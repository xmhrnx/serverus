package settings

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/platform"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// windowSection holds the tray preferences, on desktops with a tray;
// nil elsewhere.
func (p *Page) windowSection() giode.Element {
	if !platform.Desktop || !giode.TraySupported() {
		return nil
	}
	rows := []giode.Element{p.minimize}
	if platform.CloseToTray {
		rows = append(rows, p.closeTray)
	}
	rows = append(rows, giode.Text(
		"The tray icon opens Serverus, connects and disconnects servers, and quits. These take effect the next time Serverus starts.",
		theme.Caption.With(giode.Styles{LineHeight: 18})))
	return section("monitor", "Window", rows...)
}
