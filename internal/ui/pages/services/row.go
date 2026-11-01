// Package services manages systemd services: state, start, stop,
// restart, boot enablement and logs.
package services

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// stateTone picks the color of a service state.
func stateTone(s ops.Service) giode.Color {
	switch {
	case s.Failed():
		return theme.Danger
	case s.Running() && s.Sub == "running":
		return theme.Success
	case s.Running():
		return theme.Info
	case s.Active == "activating" || s.Active == "deactivating" || s.Active == "reloading":
		return theme.Warning
	}
	return theme.Faint
}

// serviceRow draws one unit.
func serviceRow(s ops.Service, selected bool, ps giode.PressState, actions giode.Element) giode.Element {
	bg := giode.Color("transparent")
	switch {
	case selected:
		bg = theme.SelectedFill
	case ps.Hovered:
		bg = theme.HoverFill
	}
	tone := stateTone(s)
	enabled := giode.Element(giode.Box(giode.Styles{}))
	if s.Enabled != "" {
		c := theme.Muted
		if s.Enabled == "enabled" {
			c = theme.Primary
		}
		enabled = widgets.Pill(s.Enabled, c)
	}
	return giode.Row(giode.Styles{
		Width: -1, Height: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: giode.SymmetricInset(theme.Space4, 0), Background: bg, BorderRadius: theme.RadiusSm,
	},
		widgets.StatusDot(tone, 10),
		giode.Stack(giode.Styles{FlexGrow: 1, Gap: 2},
			giode.Text(s.Unit, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium, FontFamily: theme.Mono, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Text(s.Description, theme.Caption.With(giode.Styles{TextOverflow: giode.TextOverflowEllipsis})),
		),
		giode.Box(giode.Styles{Width: 140}, widgets.Pill(s.Active+" · "+s.Sub, tone)),
		giode.Box(giode.Styles{Width: 90}, enabled),
		actions,
	)
}
