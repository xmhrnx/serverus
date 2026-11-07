package files

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// Layout draws the segments.
func (b *Breadcrumb) Layout(gtx layout.Context) layout.Dimensions {
	labels, paths := crumbs(b.path())
	items := []giode.Element{}
	for i := range labels {
		label, target := labels[i], paths[i]
		last := i == len(labels)-1
		seg := b.segments.Get(target, func() *giode.PressableElement {
			return giode.Pressable(func(ps giode.PressState) giode.Element {
				fg, bg := theme.Muted, giode.Color("transparent")
				if ps.Hovered {
					fg, bg = theme.Text, theme.HoverFill
				}
				return giode.Box(giode.Styles{Padding: giode.SymmetricInset(6, 3), BorderRadius: theme.RadiusSm, Background: bg},
					giode.Text(label, giode.Styles{Color: fg, FontSize: theme.TextSm, FontFamily: theme.Mono}))
			}).Label(label).OnClick(func(giode.ClickEvent) { b.navigate(target) })
		})
		if last {
			items = append(items, giode.Box(giode.Styles{Padding: giode.SymmetricInset(6, 3)},
				giode.Text(label, giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontFamily: theme.Mono, FontWeight: giode.FontWeightSemiBold})))
		} else {
			items = append(items, seg)
		}
		if !last && label != "/" {
			items = append(items, theme.Icon("chevron-right", 13, theme.Faint))
		}
	}
	return giode.Row(giode.Styles{Gap: 0, Align: giode.AlignCenter}, items...).Layout(gtx)
}
