package services

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// New returns the services page for s.
func New(e *env.Env, s *session.Session) *Page {
	p := &Page{env: e, session: s}
	p.data = async.NewResource[[]ops.Service](e.Invalidate)
	p.filter = widgets.NewSearchField("Filter services", nil)
	p.show = giode.Input.Segmented("All", "Running", "Failed", "Stopped").Styles(giode.Styles{FontSize: theme.TextSm})
	p.rows = giode.NewKeyedCache[string, *giode.PressableElement]()
	p.restarts = giode.NewKeyedCache[string, *widgets.ToolButton]()
	p.logsBtns = giode.NewKeyedCache[string, *widgets.ToolButton]()
	p.list = giode.VirtualList(func() int { return len(p.visible()) }, 52, p.row,
		giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}}).ClickFocusesRow(false)
	p.refresh = widgets.NewToolButton("refresh-cw", "Refresh", p.load)
	p.output = dialogs.NewOutputDialog(e)
	act := func(a ops.ServiceAction) func() { return func() { p.do(p.target, a) } }
	p.menu = giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd}).
		IconItem(theme.IconSource("play"), "Start", act(ops.ServiceStart)).
		IconItem(theme.IconSource("square"), "Stop", act(ops.ServiceStop)).
		IconItem(theme.IconSource("rotate-cw"), "Restart", act(ops.ServiceRestart)).
		IconItem(theme.IconSource("refresh-cw"), "Reload configuration", act(ops.ServiceReload)).
		Separator().
		IconItem(theme.IconSource("power"), "Start at boot", act(ops.ServiceEnable)).
		IconItem(theme.IconSource("circle-stop"), "Don't start at boot", act(ops.ServiceDisable)).
		Separator().
		IconItem(theme.IconSource("scroll-text"), "Logs", func() { p.showLogs(p.target) }).
		IconItem(theme.IconSource("info"), "Status", func() { p.showStatus(p.target) })
	return p
}
