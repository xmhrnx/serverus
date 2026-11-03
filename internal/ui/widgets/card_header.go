package widgets

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// CardHeader is the title row of a card: an icon badge, a title with an
// optional caption, and trailing elements (a value, actions) on the
// right.
func CardHeader(icon string, accent giode.Color, title, caption string, trailing ...giode.Element) giode.Element {
	text := []giode.Element{giode.Text(title, theme.Heading)}
	if caption != "" {
		text = append(text, giode.Text(caption, theme.Caption.With(giode.Styles{TextOverflow: giode.TextOverflowEllipsis})))
	}
	row := []giode.Element{}
	if icon != "" {
		row = append(row, IconBadge(icon, accent, 32))
	}
	row = append(row, giode.Stack(giode.Styles{Gap: 2, FlexGrow: 1}, text...))
	row = append(row, trailing...)
	return giode.Row(giode.Styles{Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Width: -1}, row...)
}
