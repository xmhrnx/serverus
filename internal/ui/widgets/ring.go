package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// Ring is a circular gauge filled to percent (0-100), with an optional
// label centered inside.
func Ring(percent float64, size, thickness int, c giode.Color, center ...giode.Element) giode.Element {
	p := float32(max(0, min(100, percent)))
	gauge := giode.Canvas(giode.Styles{Width: properties.Width(size), Height: properties.Height(size)}, func(pt *giode.Painter) {
		mid := pt.Width() / 2
		r := mid - float32(thickness)/2
		pt.Arc(mid, mid, r, float32(thickness), 0, 360, theme.Tint(c, 0.14))
		if p > 0 {
			pt.Arc(mid, mid, r, float32(thickness), 0, 360*p/100, c)
		}
	})
	if len(center) == 0 {
		return gauge
	}
	return giode.Box(giode.Styles{Width: properties.Width(size), Height: properties.Height(size)},
		gauge,
		giode.Box(giode.Styles{
			Position: giode.PositionAbsolute, Top: giode.Flush, Left: giode.Flush, Right: giode.Flush, Bottom: giode.Flush,
			Align: giode.AlignCenter, Justify: giode.JustifyCenter, FlexDirection: giode.FlexDirectionColumn,
		}, center...),
	)
}
