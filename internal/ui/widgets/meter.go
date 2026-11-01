package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/xmhrnx/bastion/internal/ui/theme"
)

// Meter is a horizontal bar filled to percent (0-100) in color c.
func Meter(percent float64, c giode.Color, height int) giode.Element {
	p := float32(max(0, min(100, percent)))
	return giode.Canvas(giode.Styles{Width: -1, Height: properties.Height(height)}, func(pt *giode.Painter) {
		w, h := pt.Width(), pt.Height()
		pt.Rect(0, 0, w, h, h/2, theme.Tint(c, 0.14))
		if fill := w * p / 100; fill > 0 {
			pt.Rect(0, 0, max(fill, h), h, h/2, c)
		}
	})
}
