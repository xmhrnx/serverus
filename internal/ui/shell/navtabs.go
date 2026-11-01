package shell

import (
	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/scroll"

	"github.com/xmhrnx/bastion/internal/ui/theme"
)

// NavTabs is the row of page tabs under the server's header.
type NavTabs struct {
	scroll *scroll.Scroll
	items  map[PageKind]*giode.PressableElement
	active func() PageKind
	choose func(PageKind)
	specs  []pageSpec
}

// NewNavTabs returns the tab row; choose runs when a tab is clicked.
func NewNavTabs(active func() PageKind, choose func(PageKind)) *NavTabs {
	n := &NavTabs{items: map[PageKind]*giode.PressableElement{}, active: active, choose: choose}
	n.scroll = giode.Scroll(n.row, giode.Styles{ScrollBar: giode.ScrollBar{Hidden: true}}).Axis(giode.Horizontal)
	for _, s := range pageSpecs {
		spec := s
		n.items[spec.kind] = giode.Pressable(func(ps giode.PressState) giode.Element {
			return navTab(spec, n.active() == spec.kind, ps.Hovered)
		}).Label(spec.label).OnClick(func(giode.ClickEvent) { n.choose(spec.kind) })
	}
	return n
}

// SetPages sets which tabs show.
func (n *NavTabs) SetPages(specs []pageSpec) { n.specs = specs }

func (n *NavTabs) row() giode.Element {
	tabs := make([]giode.Element, 0, len(n.specs))
	for _, s := range n.specs {
		tabs = append(tabs, n.items[s.kind])
	}
	return giode.Row(giode.Styles{Gap: 2, Padding: giode.Inset{Left: theme.Space4, Right: theme.Space4}, Wrap: giode.WrapNoWrap}, tabs...)
}

// Layout draws the tabs; they scroll sideways when the window is too
// narrow for all of them.
func (n *NavTabs) Layout(gtx layout.Context) layout.Dimensions {
	return giode.Stack(giode.Styles{Width: -1},
		giode.Box(giode.Styles{Width: -1}, n.scroll),
		giode.Divider(giode.Styles{Color: theme.SubtleBorder}),
	).Layout(gtx)
}

func navTab(s pageSpec, active, hovered bool) giode.Element {
	fg, underline := theme.Muted, giode.Color("transparent")
	if hovered {
		fg = theme.Text
	}
	if active {
		fg, underline = theme.Text, theme.Primary
	}
	return giode.Box(giode.Styles{},
		giode.Row(giode.Styles{Gap: 7, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.SymmetricInset(12, 12)},
			theme.Icon(s.icon, 15, fg),
			giode.Text(s.label, giode.Styles{Color: fg, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium}),
		),
		giode.Box(giode.Styles{
			Position: giode.PositionAbsolute, Left: 8, Right: 8, Bottom: giode.Flush,
			Height: 2, BorderRadius: 2, Background: underline,
		}),
	)
}
