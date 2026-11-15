package fleet

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/dashboard"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// defaultCards is the dashboard until the user arranges it.
var defaultCards = []string{"summary", "servers", "alerts", "cpu-history", "top-cpu", "top-memory", "disk-space", "temperatures"}

// catalog is every card the fleet dashboard offers.
func (p *Page) catalog() []dashboard.Spec {
	return []dashboard.Spec{
		{Kind: "summary", Title: "Summary", Icon: "layout-dashboard", Size: domain.CardLarge,
			Description: "Servers, connections, average use and alerts, in numbers.",
			Render:      func() giode.Element { return summaryCard(p.rows()) }},
		{Kind: "servers", Title: "Servers", Icon: "server", Size: domain.CardLarge,
			Description: "Every server as a tile with its status, CPU, memory and disk.",
			Render:      func() giode.Element { return p.servers.Render(p.rows()) }},
		{Kind: "alerts", Title: "Needs attention", Icon: "bell-ring", Size: domain.CardMedium,
			Description: "Failed connections, busy CPUs, full memory and disks, hot sensors.",
			Render:      func() giode.Element { return alertsCard(p.rows()) }},
		{Kind: "cpu-history", Title: "CPU history", Icon: "chart-line", Size: domain.CardMedium,
			Description: "CPU use of every connected server over the last few minutes.",
			Render: func() giode.Element {
				return historyCard("chart-line", "CPU history", p.rows(), func(c monitor.Charts) []float64 { return c.CPU }, 100, cpuValue)
			}},
		{Kind: "memory-history", Title: "Memory history", Icon: "chart-line", Size: domain.CardMedium,
			Description: "Memory use of every connected server over the last few minutes.",
			Render: func() giode.Element {
				return historyCard("chart-line", "Memory history", p.rows(), func(c monitor.Charts) []float64 { return c.Memory }, 100, memValue)
			}},
		{Kind: "top-cpu", Title: "Busiest CPUs", Icon: "cpu", Size: domain.CardMedium,
			Description: "Connected servers by CPU use.",
			Render: func() giode.Element {
				return rankingCard("cpu", theme.Primary, "Busiest CPUs", "By CPU use", p.entries(func(r row) entry {
					v := r.snap.CPUPercent
					return entry{value: v, percent: v, label: format.Percent(v), detail: fmt.Sprintf("Load %.2f", r.snap.Load[0]), color: theme.Level(v)}
				}))
			}},
		{Kind: "top-memory", Title: "Most memory used", Icon: "memory-stick", Size: domain.CardMedium,
			Description: "Connected servers by memory use.",
			Render: func() giode.Element {
				return rankingCard("memory-stick", theme.Accent, "Most memory used", "By memory use", p.entries(func(r row) entry {
					v := r.snap.MemPercent
					return entry{value: v, percent: v, label: format.Percent(v), detail: format.Bytes(r.snap.Memory.Used()) + " of " + format.Bytes(r.snap.Memory.Total), color: theme.Level(v)}
				}))
			}},
		{Kind: "disk-space", Title: "Fullest disks", Icon: "hard-drive", Size: domain.CardMedium,
			Description: "Each connected server's fullest file system.",
			Render: func() giode.Element {
				return rankingCard("hard-drive", theme.Info, "Fullest disks", "Each server's fullest file system", p.entries(func(r row) entry {
					d, ok := fullestDisk(r.snap)
					if !ok {
						return entry{value: -1}
					}
					v := d.Percent()
					return entry{value: v, percent: v, label: format.Percent(v), detail: d.Mount + " · " + format.Bytes(d.Free) + " free", color: theme.Level(v)}
				}))
			}},
		{Kind: "temperatures", Title: "Hottest servers", Icon: "thermometer", Size: domain.CardMedium,
			Description: "Each connected server's warmest sensor.",
			Render: func() giode.Element {
				return rankingCard("thermometer", theme.Warning, "Hottest servers", "Each server's warmest sensor", p.entries(func(r row) entry {
					hot, ok := r.snap.HottestSensor()
					if !ok {
						return entry{value: -1}
					}
					crit := hot.Critical
					if crit <= 0 {
						crit = 100
					}
					return entry{value: hot.Celsius, percent: hot.Celsius / crit * 100, label: format.Celsius(hot.Celsius), detail: hot.Name,
						color: []giode.Color{theme.Success, theme.Warning, theme.Danger}[hot.Level()]}
				}))
			}},
		{Kind: "network", Title: "Network traffic", Icon: "arrow-up-down", Size: domain.CardMedium,
			Description: "Connected servers by network traffic right now.",
			Render: func() giode.Element {
				rows := liveRows(p.rows())
				top := 1.0
				for _, r := range rows {
					top = max(top, r.snap.NetRxRate+r.snap.NetTxRate)
				}
				return rankingCard("arrow-up-down", theme.Info, "Network traffic", "In and out, right now", p.entries(func(r row) entry {
					v := r.snap.NetRxRate + r.snap.NetTxRate
					return entry{value: v, percent: v / top * 100, label: format.Rate(v),
						detail: "↓ " + format.Rate(r.snap.NetRxRate) + "  ↑ " + format.Rate(r.snap.NetTxRate), color: theme.Info}
				}))
			}},
		{Kind: "status", Title: "Connections", Icon: "plug", Size: domain.CardSmall,
			Description: "How many servers are connected, connecting and down.",
			Render:      func() giode.Element { return statusCard(p.rows()) }},
	}
}

// entries ranks the connected servers by what of returns; an entry
// with a negative value has nothing to rank and is left out.
func (p *Page) entries(of func(row) entry) []entry {
	var out []entry
	for _, r := range liveRows(p.rows()) {
		e := of(r)
		if e.value < 0 {
			continue
		}
		e.r = r
		out = append(out, e)
	}
	return out
}

// statusCard counts servers by connection state.
func statusCard(rows []row) giode.Element {
	counts := map[session.State]int{}
	for _, r := range rows {
		counts[r.state]++
	}
	return connectionsCard(counts, len(rows))
}
