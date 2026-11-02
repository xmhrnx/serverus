package shell

import (
	"fmt"

	"gioui.org/layout"
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// TopBar is the selected server's header: who it is, how the connection
// is doing, and the connect and edit actions.
type TopBar struct {
	sess       *session.Session
	connect    *giode.ButtonElement
	disconnect *giode.ButtonElement
	edit       *widgets.ToolButton
	// Icon-only versions of the connect actions, for phones.
	connectIcon    *widgets.ToolButton
	disconnectIcon *widgets.ToolButton
}

// NewTopBar returns the header for sess; onEdit opens the editor.
func NewTopBar(sess *session.Session, onEdit func()) *TopBar {
	return &TopBar{
		sess:       sess,
		connect:    widgets.Button(widgets.Primary, "Connect", "plug").OnClick(sess.Connect),
		disconnect: widgets.Button(widgets.Secondary, "Disconnect", "unplug").OnClick(sess.Disconnect),
		edit:       widgets.NewToolButton("pencil", "Edit server", onEdit),

		connectIcon:    widgets.NewToolButton("plug", "Connect", sess.Connect).Tint(theme.Primary),
		disconnectIcon: widgets.NewToolButton("unplug", "Disconnect", sess.Disconnect),
	}
}

// Layout draws the header.
func (t *TopBar) Layout(gtx layout.Context) layout.Dimensions {
	s := t.sess.Server()
	st, _ := t.sess.State()
	status := widgets.Pill(st.Label(), stateColor(st))
	meta := []giode.Element{
		giode.Text(s.Subtitle(), theme.Caption),
	}
	if s.Protocol != domain.ProtocolLocal {
		meta = append(meta, widgets.Pill(s.Protocol.Label(), theme.Muted))
	}
	if st == session.Connected {
		if lat := t.sess.Latency(); lat > 0 {
			meta = append(meta, giode.Text(format.Latency(lat), theme.Caption.With(giode.Styles{FontFamily: theme.Mono})))
		}
		meta = append(meta, giode.Text(fmt.Sprintf("connected %s", format.Ago(t.sess.ConnectedAt())), theme.Caption))
	}
	for _, tag := range s.Tags {
		meta = append(meta, widgets.Pill(tag, theme.Primary))
	}

	if giode.Compact() {
		return t.compact(s, st, status).Layout(gtx)
	}
	var action giode.Element = t.connect
	switch st {
	case session.Connected:
		action = t.disconnect
	case session.Connecting:
		action = giode.Spinner(giode.Styles{Width: 20, Height: 20, Color: theme.Primary})
	}
	return giode.Row(giode.Styles{
		Width: -1, Gap: theme.Space4, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: giode.Inset{Left: theme.Space6, Right: theme.Space5, Top: theme.Space5, Bottom: theme.Space3},
	},
		ServerAvatar(s, 44),
		giode.Stack(giode.Styles{Gap: 5, FlexGrow: 1},
			giode.Row(giode.Styles{Gap: theme.Space3, Align: giode.AlignCenter},
				giode.Text(s.Name, theme.Title.With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
				status,
			),
			giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter}, meta...),
		),
		t.edit,
		action,
	).Layout(gtx)
}

// compact is the header on phones: one line of name and status, the
// address under it, and icon-only actions.
func (t *TopBar) compact(s domain.Server, st session.State, status giode.Element) giode.Element {
	var action giode.Element = t.connectIcon
	switch st {
	case session.Connected:
		action = t.disconnectIcon
	case session.Connecting:
		action = giode.Spinner(giode.Styles{Width: 18, Height: 18, Color: theme.Primary})
	}
	return giode.Row(giode.Styles{
		Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: giode.Inset{Left: theme.Space4, Right: theme.Space2, Top: theme.Space3, Bottom: theme.Space2},
	},
		ServerAvatar(s, 36),
		giode.Stack(giode.Styles{Gap: 3, FlexGrow: 1},
			giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(s.Name, theme.Heading.With(giode.Styles{FontSize: theme.TextLg, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}))),
			),
			giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				status,
				giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(s.Subtitle(), theme.Caption.With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}))),
			),
		),
		t.edit,
		action,
	)
}
