package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// StatusDot is a small colored dot with a soft halo, for connection and
// service states.
func StatusDot(c giode.Color, size int) giode.Element {
	return giode.Canvas(giode.Styles{Width: properties.Width(size), Height: properties.Height(size)}, func(p *giode.Painter) {
		mid := p.Width() / 2
		p.Circle(mid, mid, mid, theme.Tint(c, 0.28))
		p.Circle(mid, mid, mid*0.55, c)
	})
}
