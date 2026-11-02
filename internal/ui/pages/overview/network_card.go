package overview

import (
	"fmt"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// NetworkCard lists interfaces with their current rates and totals.
func NetworkCard(s monitor.Snapshot) giode.Element {
	rows := []giode.Element{widgets.CardHeader("network", theme.Accent, "Interfaces", fmt.Sprintf("%d active", len(s.Net)))}
	for _, n := range s.Net {
		rows = append(rows, giode.Row(giode.Styles{Gap: theme.Space3, Width: -1, Wrap: giode.WrapNoWrap, Align: giode.AlignCenter, Padding: giode.SymmetricInset(0, 4)},
			theme.Icon("wifi", 15, theme.Muted),
			giode.Stack(giode.Styles{Gap: 2, FlexGrow: 1},
				giode.Text(n.Interface, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium, FontFamily: theme.Mono}),
				giode.Text(fmt.Sprintf("Total ↓ %s  ↑ %s", format.Bytes(n.RxTotal), format.Bytes(n.TxTotal)), theme.Caption),
			),
			giode.Stack(giode.Styles{Gap: 2, Align: giode.AlignEnd},
				giode.Text("↓ "+format.Rate(n.RxRate), giode.Styles{Color: theme.Info, FontSize: theme.TextSm, FontFamily: theme.Mono}),
				giode.Text("↑ "+format.Rate(n.TxRate), giode.Styles{Color: theme.Warning, FontSize: theme.TextSm, FontFamily: theme.Mono}),
			),
		))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1}, rows...)
}
