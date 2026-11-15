// Package overview is the server dashboard: system information, live
// CPU, memory, disk, network and temperature readings, their recent
// history, and summaries of processes, services and containers, as
// cards the user arranges.
package overview

import (
	"errors"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/dashboard"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Links open other parts of the app from cards.
type Links struct {
	// Open shows another page of the server: "processes", "services"
	// or "docker".
	Open func(page string)
	// Edit opens the server editor.
	Edit func()
}

// Page is the dashboard of one server.
type Page struct {
	env     *env.Env
	session *session.Session
	board   *dashboard.Board
	scroll  *scroll.Scroll
}

// New returns the dashboard for s. Its cards are arranged once for
// every server.
func New(e *env.Env, s *session.Session, links Links) *Page {
	p := &Page{env: e, session: s}
	p.board = dashboard.New(e, domain.DashboardServer, p.catalog(links), defaultCards...)
	p.scroll = giode.Scroll(p.content, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6, Inset: 2}})
	return p
}

// Layout draws the dashboard.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if _, ok := p.session.Metrics(); !ok {
		if err := p.session.MetricsErr(); err != nil {
			return p.failure(err).Layout(gtx)
		}
		return widgets.Loading("Collecting the first readings…").Layout(gtx)
	}
	return widgets.WithOverlays(p.scroll, p.board.Overlays()...).Layout(gtx)
}

func (p *Page) failure(err error) giode.Element {
	if errors.Is(err, monitor.ErrNotLinux) {
		return widgets.EmptyState("activity", "Monitoring isn't available here",
			"Live resource monitoring reads Linux's /proc and /sys. The terminal, files and snippets still work on this server.")
	}
	return widgets.EmptyState("circle-alert", "Couldn't read the server's resources", err.Error())
}
