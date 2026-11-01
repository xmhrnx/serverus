package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// Sparkline is a small filled line of recent values, scaled to their
// maximum, for tiles.
func Sparkline(values []float64, width, height int, c giode.Color) giode.Element {
	top := 0.0
	for _, v := range values {
		top = max(top, v)
	}
	return giode.Canvas(giode.Styles{Width: properties.Width(width), Height: properties.Height(height)}, func(p *giode.Painter) {
		w, h := p.Width(), p.Height()
		p.Line(0, h-0.5, w, h-0.5, 1, theme.Tint(c, 0.25))
		if len(values) < 2 || top <= 0 {
			return
		}
		step := w / float32(len(values)-1)
		y := func(v float64) float32 { return h - float32(v/top)*(h-3) - 1 }
		area := p.Path()
		area.MoveTo(0, h)
		for i, v := range values {
			area.LineTo(float32(i)*step, y(v))
		}
		area.LineTo(w, h)
		area.Close()
		area.Fill(theme.Tint(c, 0.2))
		line := p.Path()
		for i, v := range values {
			if i == 0 {
				line.MoveTo(0, y(v))
			} else {
				line.LineTo(float32(i)*step, y(v))
			}
		}
		line.Stroke(1.6, c)
	})
}
