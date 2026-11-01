package widgets

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// EmptyState fills a page that has nothing to show yet: an icon, a
// headline, a sentence of explanation and the actions that get the user
// going.
func EmptyState(icon, title, message string, actions ...giode.Element) giode.Element {
	children := []giode.Element{
		IconBadge(icon, theme.Primary, 56),
		giode.Text(title, theme.Title.With(giode.Styles{FontSize: theme.TextLg, TextAlign: giode.TextAlignCenter})),
	}
	if message != "" {
		children = append(children, giode.Box(giode.Styles{MaxWidth: 420},
			giode.Text(message, theme.Caption.With(giode.Styles{FontSize: theme.TextBase, TextAlign: giode.TextAlignCenter, LineHeight: 20}))))
	}
	if len(actions) > 0 {
		children = append(children, giode.Row(giode.Styles{Gap: theme.Space2, Margin: giode.Inset{Top: theme.Space2}, Justify: giode.JustifyCenter}, actions...))
	}
	return giode.Stack(giode.Styles{
		Width: -1, Height: -1, Align: giode.AlignCenter, Justify: giode.JustifyCenter,
		Gap: theme.Space3, Padding: giode.UniformInset(theme.Space8),
	}, children...)
}
