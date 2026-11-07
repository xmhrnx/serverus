package overview

import (
	"slices"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (p *Page) content() giode.Element {
	snap, _ := p.session.Metrics()
	charts := p.session.Charts()

	tiles := []giode.Element{cpuTile(snap), memoryTile(snap), diskTile(snap)}
	if hot, ok := snap.HottestSensor(); ok {
		tiles = append(tiles, temperatureTile(hot))
	} else {
		total := make([]float64, len(charts.NetRx))
		for i := range total {
			total[i] = charts.NetRx[i]
			if i < len(charts.NetTx) {
				total[i] += charts.NetTx[i]
			}
		}
		tiles = append(tiles, networkTile(snap, total))
	}

	pad := giode.UniformInset(theme.Space6)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space3)
	}
	sections := []giode.Element{
		SystemHeader(snap, p.session.Latency()),
		widgets.Grid(230, theme.Space4, tiles...),
		widgets.Grid(420, theme.Space4, UsageChart(charts), ThroughputChart(snap, charts)),
		CoresCard(snap),
		widgets.Grid(420, theme.Space4, DisksCard(snap), NetworkCard(snap)),
	}
	if len(snap.Sensors) > 0 {
		// Next to the cores; machines without sensors (most virtual
		// ones) don't get an empty card.
		sections = slices.Insert(sections, 4, SensorsCard(snap))
	}
	return giode.Stack(giode.Styles{Gap: theme.Space4, Padding: pad, Width: -1}, sections...)
}
