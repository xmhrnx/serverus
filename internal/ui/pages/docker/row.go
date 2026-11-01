// Package docker manages a server's Docker containers.
package docker

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/bastion/internal/ops"
	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// stateTone is the color of a container state.
func stateTone(c ops.Container) giode.Color {
	switch c.State {
	case "running":
		return theme.Success
	case "paused":
		return theme.Warning
	case "restarting":
		return theme.Info
	case "dead":
		return theme.Danger
	}
	return theme.Faint
}

// containerRow draws one container.
func containerRow(c ops.Container, selected bool, ps giode.PressState, actions giode.Element) giode.Element {
	bg := giode.Color("transparent")
	switch {
	case selected:
		bg = theme.SelectedFill
	case ps.Hovered:
		bg = theme.HoverFill
	}
	usage := ""
	if c.CPU != "" {
		usage = "CPU " + c.CPU + "  ·  " + c.Memory
	}
	ports := c.Ports
	if ports == "" {
		ports = "no published ports"
	}
	return giode.Row(giode.Styles{
		Width: -1, Height: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: giode.SymmetricInset(theme.Space4, 0), Background: bg, BorderRadius: theme.RadiusSm,
	},
		widgets.IconBadge("container", stateTone(c), 34),
		giode.Stack(giode.Styles{FlexGrow: 2, Gap: 2, Width: -1},
			giode.Text(c.Name, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Text(c.Image, theme.Caption.With(giode.Styles{FontFamily: theme.Mono, TextOverflow: giode.TextOverflowEllipsis})),
		),
		giode.Stack(giode.Styles{FlexGrow: 2, Gap: 2, Width: -1},
			giode.Text(c.Status, giode.Styles{Color: stateTone(c), FontSize: theme.TextSm, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Text(ports, theme.Caption.With(giode.Styles{FontSize: theme.TextXS, FontFamily: theme.Mono, TextOverflow: giode.TextOverflowEllipsis})),
		),
		giode.Box(giode.Styles{Width: 200}, giode.Text(usage, giode.Styles{Color: theme.Muted, FontSize: theme.TextXS, FontFamily: theme.Mono, TextOverflow: giode.TextOverflowEllipsis})),
		actions,
	)
}
