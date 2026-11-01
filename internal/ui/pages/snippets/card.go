package snippets

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/bastion/internal/domain"
	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// snippetCard draws a snippet with its actions.
func snippetCard(s domain.Snippet, actions ...giode.Element) giode.Element {
	text := []giode.Element{giode.Text(s.Name, theme.Heading)}
	if s.Description != "" {
		text = append(text, giode.Text(s.Description, theme.Caption))
	}
	return widgets.Card(giode.Styles{Width: -1, Padding: giode.UniformInset(theme.Space4), Gap: theme.Space3},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			widgets.IconBadge("command", theme.Primary, 34),
			giode.Stack(giode.Styles{FlexGrow: 1, Gap: 3}, text...),
			giode.Row(giode.Styles{Gap: 2, Wrap: giode.WrapNoWrap}, actions...),
		),
		giode.Box(giode.Styles{Width: -1, Background: theme.Sunken, BorderRadius: theme.RadiusSm, Padding: giode.SymmetricInset(theme.Space3, 10)},
			giode.Text(s.Command, theme.Code.With(giode.Styles{Color: theme.Accent, MaxLines: 4, TextOverflow: giode.TextOverflowEllipsis}))),
	)
}
