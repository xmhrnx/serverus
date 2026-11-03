package overview

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// legend is a chart key: a colored bar, a label and the current value.
func legend(c giode.Color, label, value string) giode.Element {
	return giode.Row(giode.Styles{Gap: 6, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		giode.Box(giode.Styles{Width: 10, Height: 3, BorderRadius: 2, Background: c}),
		giode.Text(label, theme.Caption),
		giode.Text(value, giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontWeight: giode.FontWeightSemiBold, FontFamily: theme.Mono}),
	)
}
