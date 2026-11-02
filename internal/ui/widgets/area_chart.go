package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// ChartSeries is one line of an AreaChart.
type ChartSeries struct {
	Values []float64
	Color  giode.Color
	// Fill shades the area under the line.
	Fill bool
}

// AreaChart draws series over capacity samples, newest on the right,
// scaled to ceiling (or to the largest value when ceiling is 0), with
// faint horizontal grid lines.
func AreaChart(height, capacity int, ceiling float64, series ...ChartSeries) giode.Element {
	top := ceiling
	if top <= 0 {
		for _, s := range series {
			for _, v := range s.Values {
				top = max(top, v)
			}
		}
		top *= 1.15
		if top <= 0 {
			top = 1
		}
	}
	return giode.Canvas(giode.Styles{Width: -1, Height: properties.Height(height)}, func(p *giode.Painter) {
		w, h := p.Width(), p.Height()
		for i := 1; i < 4; i++ {
			y := h * float32(i) / 4
			p.Line(0, y, w, y, 1, theme.SubtleBorder)
		}
		if capacity < 2 {
			return
		}
		step := w / float32(capacity-1)
		for _, s := range series {
			n := len(s.Values)
			if n < 2 {
				continue
			}
			x0 := w - float32(n-1)*step
			yOf := func(v float64) float32 {
				f := float32(v / top)
				return h - max(0, min(1, f))*(h-2) - 1
			}
			if s.Fill {
				area := p.Path()
				area.MoveTo(x0, h)
				for i, v := range s.Values {
					area.LineTo(x0+float32(i)*step, yOf(v))
				}
				area.LineTo(w, h)
				area.Close()
				area.Fill(theme.Tint(s.Color, 0.18))
			}
			line := p.Path()
			for i, v := range s.Values {
				x, y := x0+float32(i)*step, yOf(v)
				if i == 0 {
					line.MoveTo(x, y)
				} else {
					line.LineTo(x, y)
				}
			}
			line.Stroke(2, s.Color)
			lx, ly := w, yOf(s.Values[n-1])
			p.Circle(lx-3, ly, 3.5, s.Color)
		}
	})
}
