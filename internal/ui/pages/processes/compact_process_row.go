package processes

import (
	"fmt"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// compactProcessRow draws a process on a phone: name and user over the
// command line on the left, CPU and memory on the right.
func compactProcessRow(p ops.Process, selected bool, ps giode.PressState) giode.Element {
	bg := giode.Color("transparent")
	switch {
	case selected:
		bg = theme.SelectedFill
	case ps.Pressed || ps.Hovered:
		bg = theme.HoverFill
	}
	mono := giode.Styles{Color: theme.Muted, FontSize: theme.TextXS, FontFamily: theme.Mono}
	return giode.Row(giode.Styles{
		Width: -1, Height: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: giode.SymmetricInset(theme.Space3, 0), Background: bg, BorderRadius: theme.RadiusSm,
	},
		giode.Stack(giode.Styles{FlexGrow: 1, Gap: 2},
			giode.Text(p.Name, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Text(fmt.Sprintf("%d · %s · %s", p.PID, p.User, p.Command), mono.With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
		),
		giode.Stack(giode.Styles{Gap: 2, Align: giode.AlignEnd},
			giode.Text("CPU "+format.Percent(p.CPU), giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontFamily: theme.Mono}),
			giode.Text("MEM "+format.Percent(p.Memory), mono),
		),
	)
}
