package shell

import (
	"fmt"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/bastion/internal/ui/theme"
)

// groupHeader labels a group of servers in the sidebar.
func groupHeader(name string, count int) giode.Element {
	return giode.Row(giode.Styles{Width: -1, Justify: giode.JustifySpaceBetween, Padding: giode.Inset{Left: 10, Right: 10, Top: theme.Space4, Bottom: 6}},
		giode.Text(name, theme.Overline),
		giode.Text(fmt.Sprint(count), theme.Overline),
	)
}
