package fleet

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// ServersCard shows every server as a tile: its status and, while
// connected, its CPU, memory and disk with a CPU sparkline. A tile
// opens the server; one that's down offers to connect.
type ServersCard struct {
	env      *env.Env
	open     func(domain.Server)
	tiles    *giode.KeyedCache[string, *giode.PressableElement]
	connects *giode.KeyedCache[string, *giode.ButtonElement]
	current  map[string]row
}

// NewServersCard returns the card; open shows a server's workspace.
func NewServersCard(e *env.Env, open func(domain.Server)) *ServersCard {
	return &ServersCard{
		env: e, open: open,
		tiles:    giode.NewKeyedCache[string, *giode.PressableElement](),
		connects: giode.NewKeyedCache[string, *giode.ButtonElement](),
		current:  map[string]row{},
	}
}

// Render builds the card for rows.
func (c *ServersCard) Render(rows []row) giode.Element {
	tiles := make([]giode.Element, 0, len(rows))
	for _, r := range rows {
		id := r.server.ID
		c.current[id] = r
		tiles = append(tiles, c.tiles.Get(id, func() *giode.PressableElement {
			return giode.Pressable(func(ps giode.PressState) giode.Element { return c.tile(c.current[id], ps) }).
				Label(r.server.Name).
				OnClick(func(giode.ClickEvent) { c.open(c.current[id].server) })
		}))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1},
		widgets.CardHeader("server", theme.Primary, "Servers", "Click one to open it"),
		widgets.Grid(250, theme.Space3, tiles...),
	)
}

func (c *ServersCard) tile(r row, ps giode.PressState) giode.Element {
	bg := theme.Sunken
	if ps.Hovered {
		bg = theme.HoverFill
	}
	head := giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		widgets.ServerAvatar(r.server, 34),
		giode.Stack(giode.Styles{FlexGrow: 1, Gap: 2},
			giode.Text(r.server.Name, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightSemiBold, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Text(r.server.Subtitle(), theme.Caption.With(giode.Styles{FontSize: theme.TextXS, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
		),
		widgets.StatusDot(widgets.StateColor(r.state), 9),
	)
	body := []giode.Element{head}
	if r.live {
		disk := 0.0
		if d, ok := fullestDisk(r.snap); ok {
			disk = d.Percent()
		}
		meter := func(label string, v float64, col giode.Color) giode.Element {
			return giode.Stack(giode.Styles{FlexGrow: 1, Width: -1, Gap: 4},
				giode.Row(giode.Styles{Width: -1, Justify: giode.JustifySpaceBetween, Wrap: giode.WrapNoWrap},
					giode.Text(label, theme.Caption.With(giode.Styles{FontSize: theme.TextXS})),
					giode.Text(format.Percent(v), giode.Styles{Color: theme.Text, FontSize: theme.TextXS, FontFamily: theme.Mono})),
				widgets.Meter(v, col, 4))
		}
		body = append(body,
			giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Wrap: giode.WrapNoWrap},
				meter("CPU", r.snap.CPUPercent, theme.Primary),
				meter("RAM", r.snap.MemPercent, theme.Accent),
				meter("Disk", disk, theme.Level(disk)),
			),
			widgets.AreaChart(36, 90, 100, widgets.ChartSeries{Values: r.charts.CPU, Color: theme.Primary, Fill: true}),
		)
	} else {
		msg := r.state.Label()
		if r.state == session.Failed && r.err != nil {
			msg = format.Sentence(transport.Explain(r.err).Error())
		}
		status := giode.Text(msg, theme.Caption.With(giode.Styles{MaxLines: 2, TextOverflow: giode.TextOverflowEllipsis}))
		var action giode.Element
		switch r.state {
		case session.Disconnected, session.Failed:
			s := r.session
			action = c.connects.Get(r.server.ID, func() *giode.ButtonElement {
				return widgets.Button(widgets.Secondary, "Connect", "plug").OnClick(func() { s.Connect() })
			})
		case session.Connecting:
			action = giode.Spinner(giode.Styles{Width: 16, Height: 16, Color: theme.Primary})
		default:
			action = giode.Box(giode.Styles{})
		}
		body = append(body, giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{FlexGrow: 1}, status), action))
	}
	return giode.Box(giode.Styles{
		Width: -1, FlexGrow: 1, FlexDirection: giode.FlexDirectionColumn, Gap: theme.Space3,
		Padding: giode.UniformInset(theme.Space4), Background: bg, BorderRadius: theme.RadiusLg,
		BorderWidth: 1, BorderColor: theme.SubtleBorder,
	}, body...)
}
