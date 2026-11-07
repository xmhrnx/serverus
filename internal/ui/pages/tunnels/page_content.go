package tunnels

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

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
	pad := giode.UniformInset(theme.Space4)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space2)
	}
	return giode.Stack(giode.Styles{Gap: theme.Space3, Width: -1, Padding: pad}, rows...)
}
