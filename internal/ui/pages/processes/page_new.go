package processes

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// New returns the process list for s.
func New(e *env.Env, s *session.Session) *Page {
	p := &Page{env: e, session: s, sortBy: byCPU, desc: true}
	p.data = async.NewResource[[]ops.Process](e.Invalidate)
	p.filter = widgets.NewSearchField("Filter by name, user, command or PID", nil)
	cur := func() (int, bool) { return p.sortBy, p.desc }
	p.headers = []*widgets.SortHeader{
		widgets.NewSortHeader("PID", byPID, pidWidth, cur, p.sort),
		widgets.NewSortHeader("User", byUser, userWidth, cur, p.sort),
		widgets.NewSortHeader("CPU", byCPU, cpuWidth, cur, p.sort),
		widgets.NewSortHeader("Memory", byMemory, memWidth, cur, p.sort),
		widgets.NewSortHeader("Running", byTime, timeWidth, cur, p.sort),
		widgets.NewSortHeader("Process", byName, 0, cur, p.sort),
	}
	p.rows = giode.NewKeyedCache[int, *giode.PressableElement]()
	count := func() int { return len(p.shown) }
	listStyle := giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}}
	p.list = giode.VirtualList(count, 46, p.row, listStyle).ClickFocusesRow(false)
	p.compactList = giode.VirtualList(count, 58, p.row, listStyle).ClickFocusesRow(false)
	p.pause = widgets.NewToolButton("circle-stop", "Pause updates", func() { p.paused = !p.paused })
	p.menu = p.menuFor(ops.Process{})
	return p
}
