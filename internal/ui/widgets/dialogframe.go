package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// DialogFrame lays out a dialog's content: an icon and title with a
// subtitle, the body, and the buttons right-aligned in the footer.
func DialogFrame(icon, title, subtitle string, width int, body giode.Element, footer ...giode.Element) giode.Element {
	head := []giode.Element{}
	if icon != "" {
		head = append(head, IconBadge(icon, theme.Primary, 40))
	}
	titles := []giode.Element{giode.Text(title, theme.Title.With(giode.Styles{FontSize: theme.TextLg}))}
	if subtitle != "" {
		titles = append(titles, giode.Text(subtitle, theme.Caption))
	}
	head = append(head, giode.Stack(giode.Styles{Gap: 3, FlexGrow: 1}, titles...))

	children := []giode.Element{
		giode.Row(giode.Styles{Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Width: -1}, head...),
	}
	if body != nil {
		children = append(children, body)
	}
	if len(footer) > 0 {
		children = append(children, giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Justify: giode.JustifyEnd, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Margin: giode.Inset{Top: theme.Space2}}, footer...))
	}
	return giode.Stack(giode.Styles{Gap: theme.Space5, Width: properties.Width(width)}, children...)
}
