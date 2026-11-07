package shell

import (
	"fmt"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

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
	if pill := t.power.Pill(); pill != nil && st == session.Connected {
		meta = append(meta, pill)
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
		t.tools(st),
		action,
	).Layout(gtx)
}
