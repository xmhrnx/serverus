package dashboard

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// pickerContent lists the cards not on the board, each with Add.
func (b *Board) pickerContent() giode.Element {
	placed := map[string]bool{}
	for _, c := range b.Cards() {
		placed[c.Kind] = true
	}
	rows := []giode.Element{}
	for _, s := range b.specs {
		if placed[s.Kind] {
			continue
		}
		s := s
		add := b.picks.Get(s.Kind, func() *giode.ButtonElement {
			return widgets.Button(widgets.Secondary, "Add", "plus").OnClick(func() { b.addCard(s.Kind) })
		})
		if len(rows) > 0 {
			rows = append(rows, giode.Divider(giode.Styles{Color: theme.SubtleBorder}))
		}
		rows = append(rows, giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.SymmetricInset(0, theme.Space2)},
			widgets.IconBadge(s.Icon, theme.Primary, 34),
			giode.Stack(giode.Styles{FlexGrow: 1, Gap: 2},
				giode.Text(s.Title, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium}),
				giode.Text(s.Description, theme.Caption.With(giode.Styles{LineHeight: 17})),
			),
			add,
		))
	}
	if len(rows) == 0 {
		rows = append(rows, giode.Box(giode.Styles{Padding: giode.SymmetricInset(0, theme.Space6)},
			giode.Text("Every card is on the dashboard. Remove one to add it again elsewhere.", theme.Caption)))
	}
	return giode.Stack(giode.Styles{Width: -1, Padding: giode.Inset{Bottom: theme.Space4}}, rows...)
}
