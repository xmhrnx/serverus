package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// refreshEvery is how often the list reloads while shown.
const refreshEvery = 5 * time.Second

// Page lists containers.
type Page struct {
	env     *env.Env
	session *session.Session
	data    *async.Resource[[]ops.Container]

	selected string
	filter   *widgets.SearchField
	list     *giode.VirtualListElement
	rows     *giode.KeyedCache[string, *giode.PressableElement]
	toggles  *giode.KeyedCache[string, *widgets.ToolButton]
	logs     *giode.KeyedCache[string, *widgets.ToolButton]
	menu     *menu.Menu
	target   ops.Container
	output   *dialogs.OutputDialog
}

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

func (p *Page) docker() (ops.Docker, bool) {
	ex, err := p.session.Executor()
	if err != nil {
		return ops.Docker{}, false
	}
	return ops.Docker{Exec: ex}, true
}

func (p *Page) load() {
	d, ok := p.docker()
	if !ok {
		return
	}
	p.data.Load(func(ctx context.Context) ([]ops.Container, error) { return d.List(ctx) })
}

func (p *Page) visible() []ops.Container {
	all, _ := p.data.Get()
	q := strings.ToLower(p.filter.Text())
	if q == "" {
		return all
	}
	var out []ops.Container
	for _, c := range all {
		if strings.Contains(strings.ToLower(c.Name+" "+c.Image+" "+c.Status), q) {
			out = append(out, c)
		}
	}
	return out
}

func (p *Page) lookup(id string) (ops.Container, bool) {
	all, _ := p.data.Get()
	for _, c := range all {
		if c.ID == id {
			return c, true
		}
	}
	return ops.Container{}, false
}

func (p *Page) row(i int) giode.Element {
	v := p.visible()
	if i >= len(v) {
		return giode.Box(giode.Styles{})
	}
	id := v[i].ID
	toggle := p.toggles.Get(id, func() *widgets.ToolButton {
		return widgets.NewToolButton("play", "Start", func() {
			if c, ok := p.lookup(id); ok {
				if c.Running() {
					p.do(c, ops.ContainerStop)
				} else {
					p.do(c, ops.ContainerStart)
				}
			}
		})
	})
	logs := p.logs.Get(id, func() *widgets.ToolButton {
		return widgets.NewToolButton("scroll-text", "Logs", func() {
			if c, ok := p.lookup(id); ok {
				p.showLogs(c)
			}
		})
	})
	pr := p.rows.Get(id, func() *giode.PressableElement {
		openMenu := func(e giode.ClickEvent) {
			if c, ok := p.lookup(id); ok {
				p.selected, p.target = id, c
				p.menu.OpenAt(widgets.ClickPoint(e))
			}
		}
		return giode.Pressable(func(ps giode.PressState) giode.Element {
			c, ok := p.lookup(id)
			if !ok {
				return giode.Box(giode.Styles{})
			}
			if c.Running() {
				toggle.SetIcon("square")
				toggle.SetTooltip("Stop")
			} else {
				toggle.SetIcon("play")
				toggle.SetTooltip("Start")
			}
			return containerRow(c, p.selected == id, ps, giode.Row(giode.Styles{Gap: 2, Wrap: giode.WrapNoWrap}, toggle, logs))
		}).OnClick(func(giode.ClickEvent) { p.selected = id }).
			OnSecondaryClick(openMenu).OnLongPress(openMenu)
	})
	pr.Label(v[i].Name)
	return pr
}

func (p *Page) do(c ops.Container, a ops.ContainerAction) {
	d, ok := p.docker()
	if !ok {
		return
	}
	run := func() {
		p.env.Async(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			return d.Do(ctx, c.ID, a)
		}, func(err error) {
			if err != nil {
				p.env.Fail(fmt.Sprintf("docker %s %s", a, c.Name), err)
			} else {
				p.env.Success(fmt.Sprintf("%s: %s done", c.Name, a))
			}
			p.data.Reload(func(ctx context.Context) ([]ops.Container, error) { return d.List(ctx) })
		})
	}
	switch a {
	case ops.ContainerStop:
		p.env.Confirm(env.Confirmation{Title: "Stop " + c.Name + "?", Message: "The container stops serving until it's started again.", Action: "Stop", Destructive: true, OnConfirm: run})
	case ops.ContainerRemove:
		p.env.Confirm(env.Confirmation{Title: "Remove " + c.Name + "?", Message: "The container and its writable layer are deleted. Volumes are kept.", Action: "Remove", Destructive: true, OnConfirm: run})
	default:
		run()
	}
}

func (p *Page) showLogs(c ops.Container) {
	d, ok := p.docker()
	if !ok {
		return
	}
	p.output.Show("scroll-text", "Container logs", c.Name, func(ctx context.Context) (string, error) { return d.Logs(ctx, c.ID, 400) })
}

// Layout draws the page, reloading periodically.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	next := p.data.LoadedAt().Add(refreshEvery)
	if !p.data.Loading() && time.Now().After(next) && p.data.Err() != ops.ErrNoDocker {
		p.load()
	}
	all, _ := p.data.Get()
	running := 0
	for _, c := range all {
		if c.Running() {
			running++
		}
	}
	filterWidth, pad := giode.Styles{Width: 300}, giode.UniformInset(theme.Space4)
	if giode.Compact() {
		filterWidth, pad = giode.Styles{FlexGrow: 1, Width: -1}, giode.UniformInset(theme.Space2)
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: pad, Gap: theme.Space3},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(filterWidth, p.filter),
			giode.Box(giode.Styles{FlexGrow: 1}),
			widgets.Pill(fmt.Sprintf("%d running", running), theme.Success),
			widgets.Pill(fmt.Sprintf("%d total", len(all)), theme.Muted),
		),
		widgets.FlushCard(giode.Styles{Width: -1, FlexGrow: 1, Height: -1, Padding: giode.UniformInset(6)}, p.body()),
	), p.menu, p.output, widgets.RefreshAt(next)).Layout(gtx)
}

func (p *Page) body() giode.Element {
	if _, ok := p.data.Get(); !ok {
		if err := p.data.Err(); err != nil {
			switch err {
			case ops.ErrNoDocker:
				return widgets.EmptyState("container", "Docker isn't installed", "Install Docker on this server to manage its containers here.")
			case ops.ErrDockerDown:
				return widgets.EmptyState("container", "Docker isn't running", "The Docker CLI is installed, but its daemon isn't answering. Start it with: sudo systemctl start docker")
			}
			return widgets.EmptyState("circle-alert", "Couldn't list containers", err.Error())
		}
		return widgets.Loading("Asking Docker…")
	}
	if len(p.visible()) == 0 {
		return widgets.EmptyState("container", "No containers", "Containers you run on this server show up here.")
	}
	return p.list
}
