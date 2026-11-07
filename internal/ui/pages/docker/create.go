package docker

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

// New returns the containers page for s.
func New(e *env.Env, s *session.Session) *Page {
	p := &Page{env: e, session: s}
	p.data = async.NewResource[[]ops.Container](e.Invalidate)
	p.filter = widgets.NewSearchField("Filter containers", nil)
	p.rows = giode.NewKeyedCache[string, *giode.PressableElement]()
	p.toggles = giode.NewKeyedCache[string, *widgets.ToolButton]()
	p.logs = giode.NewKeyedCache[string, *widgets.ToolButton]()
	p.list = giode.VirtualList(func() int { return len(p.visible()) }, 58, p.row,
		giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}}).ClickFocusesRow(false)
	p.output = dialogs.NewOutputDialog(e)
	act := func(a ops.ContainerAction) func() { return func() { p.do(p.target, a) } }
	p.menu = giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd}).
		IconItem(theme.IconSource("play"), "Start", act(ops.ContainerStart)).
		IconItem(theme.IconSource("square"), "Stop", act(ops.ContainerStop)).
		IconItem(theme.IconSource("rotate-cw"), "Restart", act(ops.ContainerRestart)).
		IconItem(theme.IconSource("circle-stop"), "Pause", act(ops.ContainerPause)).
		IconItem(theme.IconSource("play"), "Resume", act(ops.ContainerUnpause)).
		Separator().
		IconItem(theme.IconSource("scroll-text"), "Logs", func() { p.showLogs(p.target) }).
		IconItem(theme.IconSource("copy"), "Copy ID", func() { p.env.App.SetClipboard(p.target.ID) }).
		Separator().
		IconItem(theme.IconSource("trash-2"), "Remove", act(ops.ContainerRemove))
	return p
}
