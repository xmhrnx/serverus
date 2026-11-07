package processes

import (
	"fmt"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// menuFor builds the actions for pr. A process run by a systemd service
// gets the service's stop and restart first: killing it alone just has
// systemd start it again.
func (p *Page) menuFor(pr ops.Process) *menu.Menu {
	m := giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd})
	if svc := pr.Service(); svc != "" {
		m.IconItem(theme.IconSource("square"), "Stop service "+svc, func() { p.service(p.target, ops.ServiceStop) }).
			IconItem(theme.IconSource("rotate-cw"), "Restart service "+svc, func() { p.service(p.target, ops.ServiceRestart) }).
			Separator()
	}
	return m.
		IconItem(theme.IconSource("x"), "Terminate (SIGTERM)", func() { p.signal(p.target, ops.SignalTerm) }).
		IconItem(theme.IconSource("skull"), "Kill (SIGKILL)", func() { p.signal(p.target, ops.SignalKill) }).
		IconItem(theme.IconSource("refresh-cw"), "Reload (SIGHUP)", func() { p.signal(p.target, ops.SignalHup) }).
		Separator().
		IconItem(theme.IconSource("copy"), "Copy command line", func() { p.env.App.SetClipboard(p.target.Command) }).
		IconItem(theme.IconSource("copy"), "Copy PID", func() { p.env.App.SetClipboard(fmt.Sprint(p.target.PID)) })
}
