package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// serverItem draws one server in the sidebar: avatar, name, address,
// connection dot and, while connected, a CPU reading.
func serverItem(s domain.Server, sess *session.Session, selected bool, ps giode.PressState) giode.Element {
	bg := giode.Color("transparent")
	switch {
	case selected:
		bg = theme.SelectedFill
	case ps.Pressed:
		bg = theme.PressFill
	case ps.Hovered:
		bg = theme.HoverFill
	}
	st, _ := sess.State()
	trailing := []giode.Element{}
	if snap, ok := sess.Metrics(); ok && st == session.Connected {
		trailing = append(trailing, giode.Text(format.Percent(snap.CPUPercent), giode.Styles{
			Color: theme.Muted, FontSize: theme.TextXS, FontFamily: theme.Mono,
		}))
	}
	trailing = append(trailing, widgets.StatusDot(stateColor(st), 10))

	nameColor := theme.Text
	if !selected && st != session.Connected {
		nameColor = giode.Mix("theme.text", "theme.muted", 0.5)
	}
	return giode.Row(giode.Styles{
		Width: -1, Gap: 10, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: giode.SymmetricInset(10, 8), BorderRadius: theme.RadiusMd, Background: bg,
	},
		ServerAvatar(s, 32),
		giode.Stack(giode.Styles{Gap: 2, FlexGrow: 1},
			giode.Text(s.Name, giode.Styles{Color: nameColor, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Text(s.Subtitle(), giode.Styles{Color: theme.Faint, FontSize: theme.TextXS, TextOverflow: giode.TextOverflowEllipsis}),
		),
		giode.Row(giode.Styles{Gap: 6, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, trailing...),
	)
}
