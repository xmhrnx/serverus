// Package fleet is the dashboard of every server at once: how many are
// up, what they're using, what needs attention, as cards the user
// arranges.
package fleet

import (
	"fmt"
	"time"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/dashboard"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Page is the fleet dashboard.
type Page struct {
	env        *env.Env
	open       func(domain.Server)
	board      *dashboard.Board
	scroll     *scroll.Scroll
	connectAll *giode.ButtonElement
	connectBtn *widgets.ToolButton
	servers    *ServersCard
	cached     []row
	rowsAt     time.Time
}

// New returns the dashboard; open shows a server's workspace.
func New(e *env.Env, open func(domain.Server)) *Page {
	p := &Page{env: e, open: open}
	p.servers = NewServersCard(e, open)
	p.board = dashboard.New(e, domain.DashboardFleet, p.catalog(), defaultCards...)
	p.connectAll = widgets.Button(widgets.Secondary, "Connect all", "plug").OnClick(p.connectEvery)
	p.connectBtn = widgets.NewToolButton("plug", "Connect all", p.connectEvery)
	p.scroll = giode.Scroll(p.content, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6, Inset: 2}})
	return p
}

// connectEvery connects the servers that aren't connected or
// connecting.
func (p *Page) connectEvery() {
	for _, srv := range p.env.Repo.Servers() {
		s := p.env.Sessions.Get(srv)
		if st, _ := s.State(); st == session.Disconnected || st == session.Failed {
			s.Connect()
		}
	}
}

func (p *Page) content() giode.Element {
	rows := p.rows()
	connected := 0
	for _, r := range rows {
		if r.state == session.Connected {
			connected++
		}
	}
	pad := giode.UniformInset(theme.Space6)
	connect := giode.Element(p.connectAll)
	if giode.Compact() {
		pad, connect = giode.UniformInset(theme.Space3), p.connectBtn
	}
	caption := fmt.Sprintf("%d servers · %d connected", len(rows), connected)
	if p.board.Editing() {
		caption = "Resize, move or remove cards, or add new ones."
	}
	tools := []giode.Element{}
	if connected < len(rows) && !p.board.Editing() {
		tools = append(tools, connect)
	}
	tools = append(tools, p.board.Controls())
	header := giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		giode.Stack(giode.Styles{FlexGrow: 1, Gap: 3},
			giode.Text("Overview", theme.Title.With(giode.Styles{FontSize: theme.Text2XL})),
			giode.Text(caption, theme.Caption),
		),
		giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, tools...),
	)
	return giode.Stack(giode.Styles{Gap: theme.Space4, Padding: pad, Width: -1}, header, p.board.Grid())
}

// Layout draws the dashboard.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if len(p.env.Repo.Servers()) == 0 {
		return widgets.EmptyState("layout-dashboard", "No servers yet", "Add a server and its readings show up here, next to every other server's.").Layout(gtx)
	}
	return widgets.WithOverlays(p.scroll, p.board.Overlays()...).Layout(gtx)
}
