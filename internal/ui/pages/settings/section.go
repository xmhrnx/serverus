package settings

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// section groups preferences in a card under a heading.
func section(icon, title string, rows ...giode.Element) giode.Element {
	children := []giode.Element{widgets.CardHeader(icon, theme.Primary, title, "")}
	for i, r := range rows {
		if i > 0 {
			children = append(children, giode.Divider(giode.Styles{Color: theme.SubtleBorder}))
		}
		children = append(children, r)
	}
	return widgets.Card(giode.Styles{Width: -1, Gap: theme.Space3}, children...)
}
