package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/input"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// CodeField returns a monospaced text field, for commands and paths.
func CodeField(placeholder string) *input.Input {
	return giode.Input.Text(placeholder, fieldStyle.With(giode.Styles{FontFamily: theme.Mono, FontSize: theme.TextSm}))
}
