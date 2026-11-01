package processes

import (
	"context"
	"fmt"
	"image"
	"sort"
	"strings"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/menu"
	"github.com/mehran9675/giode/elements"

	"github.com/xmhrnx/bastion/internal/ops"
	"github.com/xmhrnx/bastion/internal/session"
	"github.com/xmhrnx/bastion/internal/ui/async"
	"github.com/xmhrnx/bastion/internal/ui/env"
	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// refreshEvery is how often the list reloads while shown.
const refreshEvery = 3 * time.Second

// Sort columns.
const (
	byPID = iota
	byUser
	byCPU
	byMemory
	byTime
	byName
)

// Page lists processes.
type Page struct {
	env     *env.Env
	session *session.Session
	data    *async.Resource[[]ops.Process]

	sortBy   int
	desc     bool
	selected int32
	paused   bool

	filter  *widgets.SearchField
	headers []*widgets.SortHeader
	list    *giode.VirtualListElement
	rows    *giode.KeyedCache[int32, *giode.PressableElement]
	anchors map[int32]*image.Point
	menu    *menu.Menu
	target  ops.Process
	pause   *widgets.ToolButton
}

// New returns the process list for s.
func New(e *env.Env, s *session.Session) *Page {
	p := &Page{env: e, session: s, sortBy: byCPU, desc: true, anchors: map[int32]*image.Point{}}
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
	p.rows = giode.NewKeyedCache[int32, *giode.PressableElement]()
	p.list = giode.VirtualList(func() int { return len(p.visible()) }, 46, p.row,
		giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}}).ClickFocusesRow(false)
	p.pause = widgets.NewToolButton("circle-stop", "Pause updates", func() { p.paused = !p.paused })
	p.menu = giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd}).
		IconItem(theme.IconSource("x"), "Terminate (SIGTERM)", func() { p.signal(p.target, ops.SignalTerm) }).
		IconItem(theme.IconSource("skull"), "Kill (SIGKILL)", func() { p.signal(p.target, ops.SignalKill) }).
		IconItem(theme.IconSource("rotate-cw"), "Reload (SIGHUP)", func() { p.signal(p.target, ops.SignalHup) }).
		Separator().
		IconItem(theme.IconSource("copy"), "Copy command line", func() { p.env.App.SetClipboard(p.target.Command) }).
		IconItem(theme.IconSource("copy"), "Copy PID", func() { p.env.App.SetClipboard(fmt.Sprint(p.target.PID)) })
	return p
}

func (p *Page) sort(c int) {
	if p.sortBy == c {
		p.desc = !p.desc
	} else {
		p.sortBy, p.desc = c, c == byCPU || c == byMemory || c == byTime
	}
}

func (p *Page) load() {
	procs, err := p.session.Processes()
	if err != nil {
		return
	}
	p.data.Load(func(ctx context.Context) ([]ops.Process, error) { return procs.List(ctx) })
}

func (p *Page) visible() []ops.Process {
	all, _ := p.data.Get()
	q := strings.ToLower(strings.TrimSpace(p.filter.Text()))
	out := make([]ops.Process, 0, len(all))
	for _, pr := range all {
		if q != "" && !strings.Contains(strings.ToLower(pr.Name+" "+pr.User+" "+pr.Command+" "+fmt.Sprint(pr.PID)), q) {
			continue
		}
		out = append(out, pr)
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		var less bool
		switch p.sortBy {
		case byPID:
			less = a.PID < b.PID
		case byUser:
			less = a.User < b.User
		case byMemory:
			less = a.Memory < b.Memory
		case byTime:
			less = a.Elapsed < b.Elapsed
		case byName:
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		default:
			less = a.CPU < b.CPU
		}
		if p.desc {
			return !less
		}
		return less
	})
	return out
}

func (p *Page) row(i int) giode.Element {
	v := p.visible()
	if i >= len(v) {
		return giode.Box(giode.Styles{})
	}
	pid := v[i].PID
	anchor, ok := p.anchors[pid]
	if !ok {
		anchor = new(image.Point)
		p.anchors[pid] = anchor
	}
	pr := p.rows.Get(pid, func() *giode.PressableElement {
		return giode.Pressable(func(ps giode.PressState) giode.Element {
			cur, ok := p.lookup(pid)
			if !ok {
				return giode.Box(giode.Styles{})
			}
			return processRow(cur, p.selected == pid, ps)
		}).OnClick(func(giode.ClickEvent) { p.selected = pid }).
			OnSecondaryClick(func(e giode.ClickEvent) {
				if cur, ok := p.lookup(pid); ok {
					p.selected, p.target = pid, cur
					p.menu.OpenAt(f32.Pt(float32(anchor.X+e.Position.X), float32(anchor.Y+e.Position.Y)))
				}
			})
	})
	pr.Label(v[i].Name)
	return giode.Raw(func(gtx layout.Context) layout.Dimensions {
		// Track after laying out: click handlers run inside Layout and
		// need the exact position from the last frame.
		dims := pr.Layout(gtx)
		elements.TrackPosition(anchor)
		return dims
	})
}

func (p *Page) lookup(pid int32) (ops.Process, bool) {
	all, _ := p.data.Get()
	for _, pr := range all {
		if pr.PID == pid {
			return pr, true
		}
	}
	return ops.Process{}, false
}

func (p *Page) signal(pr ops.Process, sig ops.Signal) {
	run := func() {
		procs, err := p.session.Processes()
		if err != nil {
			return
		}
		p.env.Async(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			return procs.Signal(ctx, pr.PID, sig)
		}, func(err error) {
			if err != nil {
				p.env.Fail(fmt.Sprintf("Signal %s (%d)", pr.Name, pr.PID), err)
				return
			}
			p.env.Success(fmt.Sprintf("Sent SIG%s to %s (%d)", sig, pr.Name, pr.PID))
			p.load()
		})
	}
	if sig == ops.SignalHup {
		run()
		return
	}
	p.env.Confirm(env.Confirmation{
		Title:       fmt.Sprintf("Stop %s?", pr.Name),
		Message:     fmt.Sprintf("Send SIG%s to process %d (%s), run by %s.", sig, pr.PID, pr.Command, pr.User),
		Action:      map[ops.Signal]string{ops.SignalTerm: "Terminate", ops.SignalKill: "Kill"}[sig],
		Destructive: true,
		OnConfirm:   run,
	})
}

// Layout draws the list, reloading it periodically.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	next := p.data.LoadedAt().Add(refreshEvery)
	if !p.paused && !p.data.Loading() && time.Now().After(next) {
		p.load()
	}
	if p.paused {
		p.pause.SetIcon("play")
		p.pause.SetTooltip("Resume updates")
	} else {
		p.pause.SetIcon("circle-stop")
		p.pause.SetTooltip("Pause updates")
	}
	all, _ := p.data.Get()
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: giode.UniformInset(theme.Space4), Gap: theme.Space3},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{Width: 360}, p.filter),
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(fmt.Sprintf("%d processes · refreshed every %ds · right-click a process for actions", len(all), int(refreshEvery.Seconds())), theme.Caption)),
			p.pause,
		),
		widgets.FlushCard(giode.Styles{Width: -1, FlexGrow: 1, Height: -1},
			p.header(),
			giode.Divider(giode.Styles{Color: theme.SubtleBorder}),
			giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1, Padding: giode.UniformInset(6)}, p.body()),
		),
	), p.menu, widgets.RefreshAt(next)).Layout(gtx)
}

func (p *Page) header() giode.Element {
	cells := make([]giode.Element, len(p.headers))
	for i, h := range p.headers {
		cells[i] = h.Cell()
	}
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Wrap: giode.WrapNoWrap, Align: giode.AlignCenter, Padding: giode.Inset{Left: 22, Right: 22, Top: 10, Bottom: 10}}, cells...)
}

func (p *Page) body() giode.Element {
	if _, ok := p.data.Get(); !ok {
		if err := p.data.Err(); err != nil {
			return widgets.EmptyState("circle-alert", "Couldn't list processes", err.Error())
		}
		return widgets.Loading("Listing processes…")
	}
	return p.list
}
