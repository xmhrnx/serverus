package tunnels

import (
	"fmt"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/bastion/internal/domain"
	"github.com/xmhrnx/bastion/internal/tunnel"
	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// tunnelRow draws one tunnel with its live state and actions.
func tunnelRow(t domain.Tunnel, st tunnel.Status, actions ...giode.Element) giode.Element {
	tone, label := theme.Faint, "Stopped"
	switch {
	case st.Running:
		tone, label = theme.Success, fmt.Sprintf("Open · %d active · %d total", st.Connections, st.Total)
	case st.Err != nil:
		tone, label = theme.Danger, st.Err.Error()
	}
	auto := giode.Element(giode.Box(giode.Styles{}))
	if t.AutoStart {
		auto = widgets.Pill("auto", theme.Primary)
	}
	row := []giode.Element{
		widgets.IconBadge("cable", tone, 36),
		giode.Stack(giode.Styles{FlexGrow: 1, Gap: 3},
			giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter},
				giode.Text(t.Name, theme.Heading), auto),
			giode.Text(t.Describe(), theme.Code.With(giode.Styles{Color: theme.Muted})),
			giode.Text(label, giode.Styles{Color: tone, FontSize: theme.TextXS, TextOverflow: giode.TextOverflowEllipsis}),
		),
	}
	row = append(row, actions...)
	return widgets.Card(giode.Styles{FlexDirection: giode.FlexDirectionRow, Align: giode.AlignCenter, Gap: theme.Space4, Padding: giode.SymmetricInset(theme.Space4, theme.Space3), Width: -1}, row...)
}
