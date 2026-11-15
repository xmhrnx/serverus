package dashboard

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

var sizes = []struct {
	size  domain.CardSize
	label string
}{
	{domain.CardSmall, "S"}, {domain.CardMedium, "M"}, {domain.CardLarge, "L"},
}

// controls are one card's edit actions: width, move and remove.
type controls struct {
	kind                 string
	earlier, later, drop *widgets.ToolButton
	sizes                []*giode.PressableElement
}

func (b *Board) controlsFor(kind string) *controls {
	return b.controls.Get(kind, func() *controls {
		c := &controls{kind: kind}
		c.earlier = widgets.NewToolButton("chevron-left", "Move earlier", func() { b.move(kind, -1) })
		c.later = widgets.NewToolButton("chevron-right", "Move later", func() { b.move(kind, 1) })
		c.drop = widgets.NewToolButton("x", "Remove the card", func() { b.remove(kind) }).Tint(theme.Danger)
		for _, s := range sizes {
			s := s
			c.sizes = append(c.sizes, giode.Pressable(func(ps giode.PressState) giode.Element {
				bg, fg := giode.Color("transparent"), theme.Muted
				switch {
				case b.sizeOf(kind) == s.size:
					bg, fg = theme.SelectedFill, theme.Primary
				case ps.Hovered:
					bg = theme.HoverFill
				}
				return giode.Box(giode.Styles{Width: 26, Height: 24, BorderRadius: theme.RadiusSm, Background: bg, Align: giode.AlignCenter, Justify: giode.JustifyCenter},
					giode.Text(s.label, giode.Styles{Color: fg, FontSize: theme.TextXS, FontWeight: giode.FontWeightBold}))
			}).Label(map[domain.CardSize]string{domain.CardSmall: "Small", domain.CardMedium: "Half width", domain.CardLarge: "Full width"}[s.size]).
				OnClick(func(giode.ClickEvent) { b.resize(kind, s.size) }))
		}
		return c
	})
}

func (b *Board) sizeOf(kind string) domain.CardSize {
	for _, c := range b.Cards() {
		if c.Kind == kind {
			return c.Size
		}
	}
	return ""
}

// editFrame wraps a card while editing: its name and actions above it,
// inside a highlighted outline.
func (b *Board) editFrame(spec Spec, card giode.Element) giode.Element {
	c := b.controlsFor(spec.Kind)
	sizeButtons := make([]giode.Element, len(c.sizes))
	for i, s := range c.sizes {
		sizeButtons[i] = s
	}
	// The card shows its own title; the bar keeps to its icon so the
	// actions fit on a small card.
	bar := giode.Row(giode.Styles{Width: -1, Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		giode.Box(giode.Styles{FlexGrow: 1, Padding: giode.Inset{Left: 4}}, theme.Icon(spec.Icon, 14, theme.Muted)),
		giode.Row(giode.Styles{Gap: 2, Padding: giode.UniformInset(2), Background: theme.Sunken, BorderRadius: theme.RadiusMd}, sizeButtons...),
		c.earlier, c.later, c.drop,
	)
	return giode.Box(giode.Styles{
		Width: -1, FlexGrow: 1, FlexDirection: giode.FlexDirectionColumn, Gap: theme.Space2,
		Padding: giode.UniformInset(theme.Space2), BorderRadius: theme.RadiusLg + 4,
		BorderWidth: 1, BorderColor: theme.Primary, Background: theme.Tint(theme.Primary, 0.06),
	}, bar, card)
}

// unavailable stands in for a card with nothing to show here.
func unavailable(spec Spec) giode.Element {
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1, Align: giode.AlignCenter, Justify: giode.JustifyCenter},
		theme.Icon(spec.Icon, 20, theme.Faint),
		giode.Text(spec.Title, theme.Heading),
		giode.Text("Nothing to show on this server", theme.Caption),
	)
}
