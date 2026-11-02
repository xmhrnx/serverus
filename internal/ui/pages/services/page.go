package services

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/segmented"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Filters of the segmented control.
const (
	showAll = iota
	showRunning
	showFailed
	showStopped
)

// Page lists systemd services.
type Page struct {
	env     *env.Env
	session *session.Session
	data    *async.Resource[[]ops.Service]

	selected string
	filter   *widgets.SearchField
	show     *segmented.Control
	list     *giode.VirtualListElement
	rows     *giode.KeyedCache[string, *giode.PressableElement]
	restarts *giode.KeyedCache[string, *widgets.ToolButton]
	logsBtns *giode.KeyedCache[string, *widgets.ToolButton]
	anchor   *widgets.ListAnchor
	menu     *widgets.ContextMenu
	target   ops.Service
	refresh  *widgets.ToolButton
	output   *dialogs.OutputDialog
}

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
	p.anchor = widgets.NewListAnchor(func() *giode.VirtualListElement { return p.list }, 52)
	p.refresh = widgets.NewToolButton("refresh-cw", "Refresh", p.load)
	p.output = dialogs.NewOutputDialog(e)
	act := func(a ops.ServiceAction) func() { return func() { p.do(p.target, a) } }
	p.menu = widgets.NewContextMenu(giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd}).
		IconItem(theme.IconSource("play"), "Start", act(ops.ServiceStart)).
		IconItem(theme.IconSource("square"), "Stop", act(ops.ServiceStop)).
		IconItem(theme.IconSource("rotate-cw"), "Restart", act(ops.ServiceRestart)).
		IconItem(theme.IconSource("refresh-cw"), "Reload configuration", act(ops.ServiceReload)).
		Separator().
		IconItem(theme.IconSource("power"), "Start at boot", act(ops.ServiceEnable)).
		IconItem(theme.IconSource("circle-stop"), "Don't start at boot", act(ops.ServiceDisable)).
		Separator().
		IconItem(theme.IconSource("scroll-text"), "Logs", func() { p.showLogs(p.target) }).
		IconItem(theme.IconSource("info"), "Status", func() { p.showStatus(p.target) }))
	return p
}

func (p *Page) svc() (ops.Services, bool) {
	ex, err := p.session.Executor()
	if err != nil {
		return ops.Services{}, false
	}
	return ops.Services{Exec: ex}, true
}

func (p *Page) load() {
	s, ok := p.svc()
	if !ok {
		return
	}
	p.data.Reload(func(ctx context.Context) ([]ops.Service, error) { return s.List(ctx) })
}

func (p *Page) visible() []ops.Service {
	all, _ := p.data.Get()
	q := strings.ToLower(p.filter.Text())
	out := make([]ops.Service, 0, len(all))
	for _, s := range all {
		switch p.show.Selected() {
		case showRunning:
			if !s.Running() {
				continue
			}
		case showFailed:
			if !s.Failed() {
				continue
			}
		case showStopped:
			if s.Running() || s.Failed() {
				continue
			}
		}
		if q != "" && !strings.Contains(strings.ToLower(s.Unit+" "+s.Description), q) {
			continue
		}
		out = append(out, s)
	}
	return out
}

func (p *Page) lookup(unit string) (ops.Service, bool) {
	all, _ := p.data.Get()
	for _, s := range all {
		if s.Unit == unit {
			return s, true
		}
	}
	return ops.Service{}, false
}

func (p *Page) row(i int) giode.Element {
	v := p.visible()
	if i >= len(v) {
		return giode.Box(giode.Styles{})
	}
	unit := v[i].Unit
	restart := p.restarts.Get(unit, func() *widgets.ToolButton {
		return widgets.NewToolButton("rotate-cw", "Restart", func() {
			if s, ok := p.lookup(unit); ok {
				p.do(s, ops.ServiceRestart)
			}
		})
	})
	logs := p.logsBtns.Get(unit, func() *widgets.ToolButton {
		return widgets.NewToolButton("scroll-text", "Logs", func() {
			if s, ok := p.lookup(unit); ok {
				p.showLogs(s)
			}
		})
	})
	pr := p.rows.Get(unit, func() *giode.PressableElement {
		openMenu := func(e giode.ClickEvent) {
			if s, ok := p.lookup(unit); ok {
				p.selected, p.target = unit, s
				p.menu.OpenAt(p.anchor.At(slices.IndexFunc(p.visible(), func(x ops.Service) bool { return x.Unit == unit }), e))
			}
		}
		return giode.Pressable(func(ps giode.PressState) giode.Element {
			s, ok := p.lookup(unit)
			if !ok {
				return giode.Box(giode.Styles{})
			}
			return serviceRow(s, p.selected == unit, ps, giode.Row(giode.Styles{Gap: 2, Wrap: giode.WrapNoWrap}, restart, logs))
		}).OnClick(func(e giode.ClickEvent) {
			p.selected = unit
			if e.Clicks == 2 {
				if s, ok := p.lookup(unit); ok {
					p.showStatus(s)
				}
			}
		}).OnSecondaryClick(openMenu).OnLongPress(openMenu)
	})
	pr.Label(unit)
	return pr
}

func (p *Page) do(s ops.Service, a ops.ServiceAction) {
	svc, ok := p.svc()
	if !ok {
		return
	}
	run := func() {
		p.env.Async(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			return svc.Do(ctx, s.Unit, a)
		}, func(err error) {
			if err != nil {
				p.env.Fail(fmt.Sprintf("%s %s", strings.ToUpper(string(a[:1]))+string(a[1:]), s.Unit), err)
			} else {
				p.env.Success(fmt.Sprintf("%s: %s done", s.Unit, a))
			}
			p.load()
		})
	}
	if a == ops.ServiceStop || a == ops.ServiceDisable {
		p.env.Confirm(env.Confirmation{
			Title: fmt.Sprintf("%s %s?", map[ops.ServiceAction]string{ops.ServiceStop: "Stop", ops.ServiceDisable: "Disable"}[a], s.Unit),
			Message: map[ops.ServiceAction]string{
				ops.ServiceStop:    "Whatever depends on " + s.Unit + " will stop working until it's started again.",
				ops.ServiceDisable: s.Unit + " won't start the next time the server boots.",
			}[a],
			Action: map[ops.ServiceAction]string{ops.ServiceStop: "Stop", ops.ServiceDisable: "Disable"}[a], Destructive: true, OnConfirm: run,
		})
		return
	}
	run()
}

func (p *Page) showLogs(s ops.Service) {
	svc, ok := p.svc()
	if !ok {
		return
	}
	p.output.Show("scroll-text", "Logs", s.Unit, func(ctx context.Context) (string, error) { return svc.Logs(ctx, s.Unit, 400) })
}

func (p *Page) showStatus(s ops.Service) {
	svc, ok := p.svc()
	if !ok {
		return
	}
	p.output.Show("info", "Status", s.Unit, func(ctx context.Context) (string, error) { return svc.Status(ctx, s.Unit) })
}

// Layout draws the page.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if _, ok := p.data.Get(); !ok && !p.data.Loading() && p.data.Err() == nil {
		p.load()
	}
	all, _ := p.data.Get()
	running, failed := 0, 0
	for _, s := range all {
		if s.Running() {
			running++
		}
		if s.Failed() {
			failed++
		}
	}
	summary := []giode.Element{widgets.Pill(fmt.Sprintf("%d running", running), theme.Success)}
	if failed > 0 {
		summary = append(summary, widgets.Pill(fmt.Sprintf("%d failed", failed), theme.Danger))
	}
	var toolbar giode.Element
	pad := giode.UniformInset(theme.Space4)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space2)
		toolbar = giode.Stack(giode.Styles{Width: -1, Gap: theme.Space2},
			giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				giode.Box(giode.Styles{FlexGrow: 1, Width: -1}, p.filter), p.refresh),
			p.show,
		)
	} else {
		toolbar = giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{Width: 300}, p.filter),
			p.show,
			giode.Row(giode.Styles{FlexGrow: 1, Gap: theme.Space2, Justify: giode.JustifyEnd, Align: giode.AlignCenter}, summary...),
			p.refresh,
		)
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: pad, Gap: theme.Space3},
		toolbar,
		widgets.FlushCard(giode.Styles{Width: -1, FlexGrow: 1, Height: -1, Padding: giode.UniformInset(6)}, p.body()),
	), p.menu, p.output).Layout(gtx)
}

func (p *Page) body() giode.Element {
	if _, ok := p.data.Get(); !ok {
		if err := p.data.Err(); err != nil {
			if err == ops.ErrNoSystemd {
				return widgets.EmptyState("layers", "No systemd here", "This server manages services some other way (OpenRC, runit, SysV init). Use the terminal to manage them.")
			}
			return widgets.EmptyState("circle-alert", "Couldn't list services", err.Error())
		}
		return widgets.Loading("Asking systemd…")
	}
	if len(p.visible()) == 0 {
		return widgets.EmptyState("layers", "No services match", "Try another filter.")
	}
	return p.anchor.Wrap(p.list)
}
