package fleet

import (
	"sort"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// entry is one server's value on a ranking.
type entry struct {
	r       row
	value   float64 // what's sorted on
	percent float64 // the bar's fill, 0-100
	label   string  // the value as shown
	detail  string  // optional, under the name
	color   giode.Color
}

// rankingCard lists servers by a value, highest first, as bars.
func rankingCard(icon string, accent giode.Color, title, caption string, entries []entry) giode.Element {
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].value > entries[j].value })
	children := []giode.Element{widgets.CardHeader(icon, accent, title, caption)}
	if len(entries) == 0 {
		children = append(children, giode.Text("Connect servers to compare them here.", theme.Caption))
	}
	for i, e := range entries {
		if i == 8 {
			break
		}
		name := []giode.Element{giode.Text(e.r.server.Name, giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontWeight: giode.FontWeightMedium, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})}
		if e.detail != "" {
			name = append(name, giode.Text(e.detail, theme.Caption.With(giode.Styles{FontSize: theme.TextXS, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})))
		}
		children = append(children, giode.Stack(giode.Styles{Width: -1, Gap: 5},
			giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				widgets.ServerAvatar(e.r.server, 20),
				giode.Stack(giode.Styles{FlexGrow: 1, Gap: 1}, name...),
				giode.Text(e.label, giode.Styles{Color: e.color, FontSize: theme.TextSm, FontWeight: giode.FontWeightSemiBold, FontFamily: theme.Mono}),
			),
			widgets.Meter(e.percent, e.color, 5),
		))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1, Gap: theme.Space3}, children...)
}
