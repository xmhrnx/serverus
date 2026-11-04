package settings

import (
	"gioui.org/layout"
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

// Name is the selected option's name.
func (s *Swatches) Name() string {
	for _, o := range s.options {
		if o.id == s.selected {
			return o.name
		}
	}
	return ""
}

// Layout draws the circles, wrapping on narrow screens.
func (s *Swatches) Layout(gtx layout.Context) layout.Dimensions {
	items := make([]giode.Element, len(s.buttons))
	for i, b := range s.buttons {
		items[i] = b
	}
	return giode.Row(giode.Styles{Gap: 4, Align: giode.AlignCenter}, items...).Layout(gtx)
}

// toneSwatches lists the background tones.
func toneSwatches() []swatch {
	out := make([]swatch, len(theme.Tones))
	for i, t := range theme.Tones {
		out[i] = swatch{id: t.ID, name: t.Name, fill: t.Swatch}
	}
	return out
}

// highlightSwatches lists the accent colors.
func highlightSwatches() []swatch {
	out := make([]swatch, len(theme.Highlights))
	for i, h := range theme.Highlights {
		out[i] = swatch{id: h.ID, name: h.Name, fill: h.Dark}
	}
	return out
}
