package overview

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// infoChip is one fact in the system header: an icon and a label.
func infoChip(icon, label string) giode.Element {
	return giode.Row(giode.Styles{
		Gap: 6, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: giode.SymmetricInset(10, 5), BorderRadius: theme.RadiusFull,
		Background: theme.Sunken, BorderColor: theme.SubtleBorder, BorderWidth: 1,
	},
		theme.Icon(icon, 13, theme.Muted),
		giode.Text(label, giode.Styles{Color: theme.Text, FontSize: theme.TextSm}),
	)
}
