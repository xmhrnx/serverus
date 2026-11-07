package theme

import (
	"image/color"
	"math"
	"strconv"

	"github.com/xmhrnx/giode"
)

func linear(c uint8) float64 {
	f := float64(c) / 255
	if f <= 0.04045 {
		return f / 12.92
	}
	return math.Pow((f+0.055)/1.055, 2.4)
}

// rgb parses a "#rrggbb" color.
func rgb(c giode.Color) color.NRGBA {
	v, _ := strconv.ParseUint(string(c)[1:], 16, 32)
	return color.NRGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xff}
}
