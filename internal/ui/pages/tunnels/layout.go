package tunnels

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

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
		p.header(),
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, p.scroll),
	), p.dialog).Layout(gtx)
}

// header explains tunnels next to the add button; on a phone the button
// shrinks to an icon so the sentence keeps its width.
func (p *Page) header() giode.Element {
	add, pad := giode.Element(p.add), giode.Inset{Left: theme.Space4, Right: theme.Space4, Top: theme.Space4}
	if giode.Compact() {
		add, pad = p.addIcon, giode.Inset{Left: theme.Space3, Right: theme.Space2, Top: theme.Space2}
	}
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: pad},
		giode.Box(giode.Styles{FlexGrow: 1}, giode.Text("Forward local ports through this server. Tunnels close when it disconnects.", theme.Caption)),
		add)
}
