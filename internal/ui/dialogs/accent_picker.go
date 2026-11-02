package dialogs

import (
	"gioui.org/layout"
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// AccentPicker chooses a server's avatar color from domain.Accents.
type AccentPicker struct {
	value  string
	swatch []*giode.PressableElement
}

// NewAccentPicker returns a picker with the first accent selected.
func NewAccentPicker() *AccentPicker {
	a := &AccentPicker{value: domain.Accents[0]}
	for _, c := range domain.Accents {
		color := c
		a.swatch = append(a.swatch, giode.Pressable(func(s giode.PressState) giode.Element {
			ring := giode.Color("transparent")
			if a.value == color {
				ring = theme.Text
			} else if s.Hovered {
				ring = theme.Border
			}
			return giode.Box(giode.Styles{Padding: giode.UniformInset(2), BorderRadius: theme.RadiusFull, BorderColor: ring, BorderWidth: 2},
				giode.Box(giode.Styles{Width: 20, Height: 20, BorderRadius: theme.RadiusFull, Background: giode.Color(color)}))
		}).Label("Accent "+color).OnClick(func(giode.ClickEvent) { a.value = color }))
	}
	return a
}

// Value is the chosen color.
func (a *AccentPicker) Value() string { return a.value }

// Set selects c.
func (a *AccentPicker) Set(c string) { a.value = c }

// Layout draws the swatches.
func (a *AccentPicker) Layout(gtx layout.Context) layout.Dimensions {
	items := make([]giode.Element, len(a.swatch))
	for i, s := range a.swatch {
		items[i] = s
	}
	return giode.Row(giode.Styles{Gap: 6, Align: giode.AlignCenter}, items...).Layout(gtx)
}
