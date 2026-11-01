package widgets

import (
	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"
)

// Grid lays children in equal columns, as many as fit at least minWidth
// dp wide, wrapping into rows: a dashboard of tiles that reflows with
// the window.
func Grid(minWidth, gap int, children ...giode.Element) giode.Element {
	return giode.Raw(func(gtx layout.Context) layout.Dimensions {
		avail := int(gtx.Metric.PxToDp(gtx.Constraints.Max.X))
		cols := max(1, min(len(children), (avail+gap)/(minWidth+gap)))
		rows := []giode.Element{}
		for i := 0; i < len(children); i += cols {
			cells := make([]giode.Element, 0, cols)
			for j := i; j < i+cols; j++ {
				var c giode.Element = giode.Box(giode.Styles{})
				if j < len(children) {
					c = children[j]
				}
				cells = append(cells, giode.Box(giode.Styles{FlexGrow: 1, Width: -1}, c))
			}
			rows = append(rows, giode.Row(giode.Styles{Width: -1, Gap: properties.Gap(gap), Wrap: giode.WrapNoWrap, Align: giode.AlignStart}, cells...))
		}
		return giode.Stack(giode.Styles{Width: -1, Gap: properties.Gap(gap)}, rows...).Layout(gtx)
	})
}
