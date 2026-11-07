package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

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
		t.tools(st),
		action,
	)
}
