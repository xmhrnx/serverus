// Package overview is the server dashboard: system information, live
// CPU, memory, disk, network and temperature readings, and their recent
// history.
package overview

import (
	"errors"

	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/scroll"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Page is the dashboard of one server.
type Page struct {
	env     *env.Env
	session *session.Session
	scroll  *scroll.Scroll
}

// New returns the dashboard for s.
func New(e *env.Env, s *session.Session) *Page {
	p := &Page{env: e, session: s}
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
	return p.scroll.Layout(gtx)
}

func (p *Page) failure(err error) giode.Element {
	if errors.Is(err, monitor.ErrNotLinux) {
		return widgets.EmptyState("activity", "Monitoring isn't available here",
			"Live resource monitoring reads Linux's /proc and /sys. The terminal, files and snippets still work on this server.")
	}
	return widgets.EmptyState("circle-alert", "Couldn't read the server's resources", err.Error())
}

func (p *Page) content() giode.Element {
	snap, _ := p.session.Metrics()
	charts := p.session.Charts()

	tiles := []giode.Element{cpuTile(snap), memoryTile(snap), diskTile(snap)}
	if hot, ok := snap.HottestSensor(); ok {
		tiles = append(tiles, temperatureTile(hot))
	} else {
		total := make([]float64, len(charts.NetRx))
		for i := range total {
			total[i] = charts.NetRx[i]
			if i < len(charts.NetTx) {
				total[i] += charts.NetTx[i]
			}
		}
		tiles = append(tiles, networkTile(snap, total))
	}

	return giode.Stack(giode.Styles{Gap: theme.Space4, Padding: giode.UniformInset(theme.Space6), Width: -1},
		SystemHeader(snap, p.session.Latency()),
		widgets.Grid(230, theme.Space4, tiles...),
		widgets.Grid(420, theme.Space4, UsageChart(charts), ThroughputChart(snap, charts)),
		widgets.Grid(420, theme.Space4, CoresCard(snap), SensorsCard(snap)),
		widgets.Grid(420, theme.Space4, DisksCard(snap), NetworkCard(snap)),
	)
}
