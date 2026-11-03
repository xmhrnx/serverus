package widgets

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// DialogPanel is the style of every dialog's panel.
var DialogPanel = giode.Styles{
	Background:   theme.Surface,
	BorderColor:  theme.Border,
	BorderWidth:  1,
	BorderRadius: theme.RadiusLg,
	Padding:      giode.UniformInset(theme.Space6),
	BoxShadow:    giode.BoxShadow{Y: 18, Blur: 48, Color: "#00000070"},
}
