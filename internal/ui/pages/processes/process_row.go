// Package processes lists a server's processes, busiest first, and
// signals them.
package processes

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// Column widths in dp.
const (
	pidWidth  = 76
	userWidth = 110
	cpuWidth  = 120
	memWidth  = 120
	timeWidth = 90
)

// processRow draws one process.
func processRow(p ops.Process, selected bool, ps giode.PressState) giode.Element {
	bg := giode.Color("transparent")
	switch {
	case selected:
		bg = theme.SelectedFill
	case ps.Hovered:
		bg = theme.HoverFill
	}
	mono := giode.Styles{Color: theme.Muted, FontSize: theme.TextSm, FontFamily: theme.Mono}
	return giode.Row(giode.Styles{
		Width: -1, Height: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: giode.SymmetricInset(theme.Space4, 0), Background: bg, BorderRadius: theme.RadiusSm,
	},
		giode.Box(giode.Styles{Width: pidWidth}, giode.Text(fmt.Sprint(p.PID), mono)),
		giode.Box(giode.Styles{Width: userWidth}, giode.Text(p.User, giode.Styles{Color: theme.Text, FontSize: theme.TextSm, TextOverflow: giode.TextOverflowEllipsis})),
		usageCell(p.CPU, cpuWidth, theme.Primary),
		usageCell(p.Memory, memWidth, theme.Accent),
		giode.Box(giode.Styles{Width: timeWidth}, giode.Text(format.Duration(p.Elapsed), mono)),
		giode.Stack(giode.Styles{FlexGrow: 1, Gap: 1},
			giode.Text(p.Name, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Text(p.Command, giode.Styles{Color: theme.Faint, FontSize: theme.TextXS, FontFamily: theme.Mono, TextOverflow: giode.TextOverflowEllipsis}),
		),
	)
}
