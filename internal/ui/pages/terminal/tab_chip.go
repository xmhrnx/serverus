package terminal

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// tabChip draws one tab of the strip. The page owns the pressables and
// passes their state in.
func tabChip(title string, active, hovered, ended bool, closeBtn giode.Element) giode.Element {
	bg, fg, border := giode.Color("transparent"), theme.Muted, giode.Color("transparent")
	switch {
	case active:
		bg, fg, border = theme.Surface, theme.Text, theme.Border
	case hovered:
		bg = theme.HoverFill
	}
	dot := theme.Success
	if ended {
		dot = theme.Faint
	}
	// Phones get a slimmer chip, so the strip keeps room for a second tab.
	pad, label := giode.Inset{Left: 12, Right: 6, Top: 6, Bottom: 6}, giode.Styles{MaxWidth: 170}
	if giode.Compact() {
		pad, label = giode.Inset{Left: 10, Right: 2, Top: 1, Bottom: 1}, giode.Styles{MaxWidth: 110}
	}
	return giode.Row(giode.Styles{
		Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: pad, BorderRadius: theme.RadiusMd,
		Background: bg, BorderColor: border, BorderWidth: 1,
	},
		giode.Box(giode.Styles{Width: 7, Height: 7, BorderRadius: theme.RadiusFull, Background: dot}),
		giode.Box(label, giode.Text(title, giode.Styles{Color: fg, FontSize: theme.TextSm, FontWeight: giode.FontWeightMedium, TextOverflow: giode.TextOverflowEllipsis})),
		closeBtn,
	)
}
