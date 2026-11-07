package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// wide shows the sidebar beside the workspace.
func (r *Root) wide() giode.Element {
	r.navOn = false
	return giode.Row(giode.Styles{Width: -1, FlexGrow: 1, Height: -1, Wrap: giode.WrapNoWrap},
		r.sidebar,
		giode.Box(giode.Styles{Width: 1, Height: -1, Background: theme.SubtleBorder}),
		giode.Box(giode.Styles{FlexGrow: 1, Width: -1, Height: -1, Background: theme.Background}, r.main()),
	)
}

// compact is the phone and narrow-window layout: the server list
// slides in from a button.
func (r *Root) compact() giode.Element {
	bar := []giode.Element{r.navBtn, giode.Box(giode.Styles{FlexGrow: 1}, giode.Text("Servers", theme.Heading))}
	if r.bar == nil {
		// No title bar on phones: its actions move here.
		bar = append(bar, r.transBtn, r.themeBtn)
	}
	return giode.Stack(giode.Styles{Width: -1, FlexGrow: 1, Height: -1},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
			Padding: giode.SymmetricInset(theme.Space3, theme.Space2), Background: theme.Sidebar}, bar...),
		giode.Divider(giode.Styles{Color: theme.SubtleBorder}),
		giode.Box(giode.Styles{FlexGrow: 1, Width: -1, Height: -1, Background: theme.Background}, r.main()),
	)
}
