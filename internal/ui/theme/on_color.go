package theme

import (
	"github.com/xmhrnx/giode"
)

// onColor is white or near-black, whichever reads better on c.
func onColor(c giode.Color) giode.Color {
	v := rgb(c)
	lum := 0.2126*linear(v.R) + 0.7152*linear(v.G) + 0.0722*linear(v.B)
	if lum > 0.4 {
		return "#0b0d12"
	}
	return "#ffffff"
}
