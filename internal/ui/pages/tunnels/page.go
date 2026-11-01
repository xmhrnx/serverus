package tunnels

import (
	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/scroll"

	"github.com/xmhrnx/bastion/internal/domain"
	"github.com/xmhrnx/bastion/internal/session"
	"github.com/xmhrnx/bastion/internal/ui/env"
	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// Page lists a server's tunnels.
type Page struct {
	env     *env.Env
	session *session.Session

	add     *giode.ButtonElement
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

func (p *Page) toggle(t domain.Tunnel) {
	if p.env.Tunnels.Status(t.ID).Running {
		p.env.Tunnels.Stop(t.ID)
		return
	}
	fw, err := p.session.Forwarder()
	if err != nil {
		p.env.Fail("Open "+t.Name, err)
		return
	}
	if err := p.env.Tunnels.Start(t, fw); err != nil {
		p.env.Fail("Open "+t.Name, err)
		return
	}
	p.env.Success("Tunnel open on " + t.LocalAddress())
}

func (p *Page) content() giode.Element {
	tunnels := p.env.Repo.Tunnels(p.session.Server().ID)
	rows := []giode.Element{}
	for _, t := range tunnels {
		id := t.ID
		st := p.env.Tunnels.Status(id)
		toggle := p.toggles.Get(id, func() *widgets.ToolButton {
			return widgets.NewToolButton("play", "Open", func() {
				if t, ok := p.find(id); ok {
					p.toggle(t)
				}
			})
		})
		if st.Running {
			toggle.SetIcon("square")
			toggle.SetTooltip("Close")
		} else {
			toggle.SetIcon("play")
			toggle.SetTooltip("Open")
		}
		edit := p.edits.Get(id, func() *widgets.ToolButton {
			return widgets.NewToolButton("pencil", "Edit", func() {
				if t, ok := p.find(id); ok {
					p.dialog.Edit(t)
				}
			})
		})
		copyAddr := p.copies.Get(id, func() *widgets.ToolButton {
			return widgets.NewToolButton("copy", "Copy local address", func() {
				if t, ok := p.find(id); ok {
					p.env.App.SetClipboard(t.LocalAddress())
					p.env.Info("Copied " + t.LocalAddress())
				}
			})
		})
		del := p.deletes.Get(id, func() *widgets.ToolButton {
			return widgets.NewToolButton("trash-2", "Delete", func() {
				t, ok := p.find(id)
				if !ok {
					return
				}
				p.env.Confirm(env.Confirmation{
					Title: "Delete " + t.Name + "?", Message: "The tunnel " + t.Describe() + " is closed and forgotten.",
					Action: "Delete", Destructive: true,
					OnConfirm: func() {
						p.env.Tunnels.Stop(id)
						_ = p.env.Repo.DeleteTunnel(id)
					},
				})
			}).Tint(theme.Danger)
		})
		rows = append(rows, tunnelRow(t, st, toggle, copyAddr, edit, del))
	}
	return giode.Stack(giode.Styles{Gap: theme.Space3, Width: -1, Padding: giode.UniformInset(theme.Space4)}, rows...)
}

// Layout draws the page.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if len(p.env.Repo.Tunnels(p.session.Server().ID)) == 0 {
		return widgets.WithOverlays(
			widgets.EmptyState("cable", "No tunnels yet",
				"Reach a database, admin panel or any service that only listens on the server's private network, as if it ran on this computer.",
				p.add),
			p.dialog,
		).Layout(gtx)
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1},
		giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.Inset{Left: theme.Space4, Right: theme.Space4, Top: theme.Space4}},
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text("Forward local ports through this server. Tunnels close when it disconnects.", theme.Caption)),
			p.add),
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, p.scroll),
	), p.dialog).Layout(gtx)
}
