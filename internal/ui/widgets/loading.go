package widgets

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// Loading is a centered spinner with a message.
func Loading(message string) giode.Element {
	return giode.Stack(giode.Styles{Width: -1, Height: -1, Align: giode.AlignCenter, Justify: giode.JustifyCenter, Gap: theme.Space3},
		giode.Spinner(giode.Styles{Width: 28, Height: 28, Color: theme.Primary}),
		giode.Text(message, theme.Caption),
	)
}
