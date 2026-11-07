package widgets

import (
	"github.com/xmhrnx/giode"
)

// ChartSeries is one line of an AreaChart.
type ChartSeries struct {
	Values []float64
	Color  giode.Color
	// Fill shades the area under the line.
	Fill bool
}
