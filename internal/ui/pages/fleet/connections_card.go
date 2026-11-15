package fleet

import (
	"fmt"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// connectionsCard is the servers by connection state, as a stacked bar
// with a key.
func connectionsCard(counts map[session.State]int, total int) giode.Element {
	states := []session.State{session.Connected, session.Connecting, session.Failed, session.Disconnected}
	bar := []giode.Element{}
	keys := []giode.Element{}
	for _, st := range states {
		n := counts[st]
		c := widgets.StateColor(st)
		if n > 0 {
			bar = append(bar, giode.Box(giode.Styles{FlexGrow: properties.FlexGrow(n), Height: 10, Background: c}))
		}
		keys = append(keys, giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			widgets.StatusDot(c, 8),
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(st.Label(), theme.Caption)),
			giode.Text(fmt.Sprint(n), giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontFamily: theme.Mono}),
		))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1, Gap: theme.Space3},
		widgets.CardHeader("plug", theme.Success, "Connections", fmt.Sprintf("%d servers", total)),
		giode.Row(giode.Styles{Width: -1, Wrap: giode.WrapNoWrap, BorderRadius: theme.RadiusSm, Overflow: giode.OverflowHidden, Gap: 2}, bar...),
		giode.Stack(giode.Styles{Width: -1, Gap: 6}, keys...),
	)
}
