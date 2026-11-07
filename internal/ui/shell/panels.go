package shell

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (s *Sidebar) content() giode.Element {
	groups, byGroup := s.filtered()
	rows := []giode.Element{}
	for _, g := range groups {
		rows = append(rows, groupHeader(g, len(byGroup[g])))
		for _, srv := range byGroup[g] {
			rows = append(rows, s.item(srv))
		}
	}
	if len(rows) == 0 {
		rows = append(rows, giode.Box(giode.Styles{Padding: giode.UniformInset(theme.Space4)},
			giode.Text("No servers match the search.", theme.Caption)))
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 2, Padding: giode.Inset{Left: theme.Space3, Right: theme.Space3, Bottom: theme.Space3}}, rows...)
}

// Layout draws the sidebar.
func (s *Sidebar) Layout(gtx layout.Context) layout.Dimensions {
	return widgets.WithOverlays(giode.Stack(giode.Styles{
		Width: theme.SidebarWidth, Height: -1, Background: theme.Sidebar,
	},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.Inset{Left: theme.Space4, Right: theme.Space3, Top: theme.Space4, Bottom: theme.Space2}},
			giode.Box(giode.Styles{FlexGrow: 1}, s.search),
			s.add,
		),
		giode.Box(giode.Styles{FlexGrow: 1, Width: -1, Height: -1}, s.list),
		giode.Divider(giode.Styles{Color: theme.SubtleBorder}),
		giode.Box(giode.Styles{Width: -1, Padding: giode.UniformInset(theme.Space3)}, s.settings),
	), s.menu).Layout(gtx)
}
