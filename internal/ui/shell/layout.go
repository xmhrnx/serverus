package shell

import (
	"fmt"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

func (t *TransfersPanel) list() giode.Element {
	jobs := t.env.Transfers.Jobs()
	if len(jobs) == 0 {
		return giode.Box(giode.Styles{Padding: giode.SymmetricInset(0, theme.Space6)},
			giode.Text("Uploads and downloads appear here.", theme.Caption))
	}
	rows := make([]giode.Element, 0, len(jobs)*2)
	for i, j := range jobs {
		if i > 0 {
			rows = append(rows, giode.Divider(giode.Styles{Color: theme.SubtleBorder}))
		}
		rows = append(rows, t.job(j))
	}
	return giode.Stack(giode.Styles{Width: -1, Padding: giode.Inset{Right: theme.Space2}}, rows...)
}

// Layout draws the panel.
func (t *TransfersPanel) Layout(gtx layout.Context) layout.Dimensions {
	return giode.Stack(giode.Styles{Width: -1, Height: -1, Gap: theme.Space4, Padding: giode.UniformInset(theme.Space5)},
		giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Stack(giode.Styles{FlexGrow: 1, Gap: 2},
				giode.Text("Transfers", theme.Title.With(giode.Styles{FontSize: theme.TextLg})),
				giode.Text(fmt.Sprintf("%d active", t.env.Transfers.Active()), theme.Caption)),
			t.clear),
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, t.scroll),
	).Layout(gtx)
}
