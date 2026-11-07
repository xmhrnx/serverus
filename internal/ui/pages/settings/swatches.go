package settings

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// swatch is one choice of a Swatches row.
type swatch struct {
	id, name string
	fill     giode.Color
}

// Swatches picks one color from a row of circles.
type Swatches struct {
	options  []swatch
	selected string
	buttons  []*giode.PressableElement
}

// newSwatches returns a picker over options with selected chosen;
// onPick runs with the id of a newly picked option.
func newSwatches(options []swatch, selected string, onPick func(id string)) *Swatches {
	s := &Swatches{options: options, selected: selected}
	if s.selected == "" && len(options) > 0 {
		s.selected = options[0].id
	}
	for _, o := range options {
		o := o
		s.buttons = append(s.buttons, giode.Pressable(func(st giode.PressState) giode.Element {
			ring := giode.Color("transparent")
			switch {
			case s.selected == o.id:
				ring = theme.Primary
			case st.Hovered:
				ring = theme.Border
			}
			return giode.Box(giode.Styles{Padding: giode.UniformInset(2), BorderRadius: theme.RadiusFull, BorderColor: ring, BorderWidth: 2},
				giode.Box(giode.Styles{Width: 26, Height: 26, BorderRadius: theme.RadiusFull, Background: o.fill, BorderColor: theme.Border, BorderWidth: 1}))
		}).Label(o.name).OnClick(func(giode.ClickEvent) {
			if s.selected != o.id {
				s.selected = o.id
				onPick(o.id)
			}
		}))
	}
	return s
}
