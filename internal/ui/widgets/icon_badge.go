package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// IconBadge is an icon on a rounded square tinted with its color.
func IconBadge(icon string, c giode.Color, size int) giode.Element {
	return giode.Box(giode.Styles{
		Width: properties.Width(size), Height: properties.Height(size),
		Background:   theme.Tint(c, 0.16),
		BorderRadius: properties.BorderRadius(size * 3 / 10),
		Align:        giode.AlignCenter,
		Justify:      giode.JustifyCenter,
	}, theme.Icon(icon, size*9/16, c))
}
