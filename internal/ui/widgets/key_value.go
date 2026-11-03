package widgets

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// KeyValue is a label on the left and its value on the right, for
// property lists.
func KeyValue(key, value string) giode.Element {
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Wrap: giode.WrapNoWrap, Align: giode.AlignCenter, Padding: giode.SymmetricInset(0, 6)},
		giode.Box(giode.Styles{Width: 110}, giode.Text(key, theme.Caption)),
		giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(value, theme.Body.With(giode.Styles{TextOverflow: giode.TextOverflowEllipsis}))),
	)
}
