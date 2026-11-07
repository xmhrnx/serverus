package tunnels

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Page lists a server's tunnels.
type Page struct {
	env     *env.Env
	session *session.Session

	add     *giode.ButtonElement
	addIcon *widgets.ToolButton
	dialog  *Dialog
	scroll  *scroll.Scroll
	toggles *giode.KeyedCache[string, *widgets.ToolButton]
	edits   *giode.KeyedCache[string, *widgets.ToolButton]
	copies  *giode.KeyedCache[string, *widgets.ToolButton]
	deletes *giode.KeyedCache[string, *widgets.ToolButton]
}

// New returns the tunnels page for s.
func New(e *env.Env, s *session.Session) *Page {
	p := &Page{env: e, session: s}
	p.dialog = NewDialog(func(t domain.Tunnel) error {
		if p.env.Tunnels.Status(t.ID).Running {
			p.env.Tunnels.Stop(t.ID)
		}
		return p.env.Repo.SaveTunnel(t)
	})
	p.add = widgets.Button(widgets.Primary, "New tunnel", "plus").OnClick(func() { p.dialog.Add(s.Server().ID) })
	p.addIcon = widgets.NewToolButton("plus", "New tunnel", func() { p.dialog.Add(s.Server().ID) }).Tint(theme.Primary)
	p.scroll = giode.Scroll(p.content, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}})
	p.toggles = giode.NewKeyedCache[string, *widgets.ToolButton]()
	p.edits = giode.NewKeyedCache[string, *widgets.ToolButton]()
	p.copies = giode.NewKeyedCache[string, *widgets.ToolButton]()
	p.deletes = giode.NewKeyedCache[string, *widgets.ToolButton]()
	return p
}

func (p *Page) find(id string) (domain.Tunnel, bool) {
	for _, t := range p.env.Repo.Tunnels(p.session.Server().ID) {
		if t.ID == id {
			return t, true
		}
	}
	return domain.Tunnel{}, false
}
