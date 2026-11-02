package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/input"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// fieldStyle is the look of every text field.
var fieldStyle = giode.Styles{
	Background:   theme.Sunken,
	Color:        theme.Text,
	BorderColor:  theme.Border,
	BorderWidth:  1,
	BorderRadius: theme.RadiusMd,
	Padding:      giode.SymmetricInset(12, 9),
	FontSize:     theme.TextBase,
}

// TextField returns a themed single-line text field.
func TextField(placeholder string) *input.Input {
	return giode.Input.Text(placeholder, fieldStyle)
}
