package overview

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/dashboard"
)

// defaultCards is the dashboard until the user arranges it.
var defaultCards = []string{
	"system", "cpu", "memory", "disk", "temperature",
	"usage-chart", "network-chart", "temperature-chart", "processes",
	"cores", "sensors", "disks", "interfaces",
}

// catalog is every card the server dashboard offers.
func (p *Page) catalog(links Links) []dashboard.Spec {
	snap := func() monitor.Snapshot { s, _ := p.session.Metrics(); return s }
	charts := p.session.Charts
	hasSensors := func() bool { return len(snap().Sensors) > 0 }
	open := func(page string) func() { return func() { links.Open(page) } }
	procs := NewProcessesCard(p.env.Invalidate, p.session, open("processes"))
	svcs := NewServicesCard(p.env.Invalidate, p.session, open("services"))
	boxes := NewContainersCard(p.env.Invalidate, p.session, open("docker"))
	notes := NewNotesCard(p.session, links.Edit)
	return []dashboard.Spec{
		{Kind: "system", Title: "System", Icon: "monitor", Size: domain.CardLarge,
			Description: "Host name, OS, kernel, uptime, load and CPU model.",
			Render:      func() giode.Element { return SystemHeader(snap(), p.session.Latency()) }},
		{Kind: "cpu", Title: "CPU", Icon: "cpu", Size: domain.CardSmall,
			Description: "Total CPU use and I/O wait, as a gauge.",
			Render:      func() giode.Element { return cpuTile(snap()) }},
		{Kind: "memory", Title: "Memory", Icon: "memory-stick", Size: domain.CardSmall,
			Description: "Memory and swap in use, as a gauge.",
			Render:      func() giode.Element { return memoryTile(snap()) }},
		{Kind: "disk", Title: "Root disk", Icon: "hard-drive", Size: domain.CardSmall,
			Description: "Space used on the system disk, as a gauge.",
			Render:      func() giode.Element { return diskTile(snap()) }},
		{Kind: "temperature", Title: "Hottest sensor", Icon: "thermometer", Size: domain.CardSmall,
			Description: "The warmest temperature reading right now.",
			Available:   hasSensors,
			Render: func() giode.Element {
				hot, _ := snap().HottestSensor()
				return temperatureTile(hot)
			}},
		{Kind: "network", Title: "Network", Icon: "network", Size: domain.CardSmall,
			Description: "Traffic in and out right now, with a sparkline.",
			Render: func() giode.Element {
				c := charts()
				total := make([]float64, len(c.NetRx))
				for i := range total {
					total[i] = c.NetRx[i]
					if i < len(c.NetTx) {
						total[i] += c.NetTx[i]
					}
				}
				return networkTile(snap(), total)
			}},
		{Kind: "usage-chart", Title: "CPU & memory history", Icon: "chart-line", Size: domain.CardMedium,
			Description: "CPU and memory use over the last few minutes.",
			Render:      func() giode.Element { return UsageChart(charts()) }},
		{Kind: "network-chart", Title: "Network history", Icon: "arrow-up-down", Size: domain.CardMedium,
			Description: "Traffic in and out over the last few minutes.",
			Render:      func() giode.Element { return ThroughputChart(snap(), charts()) }},
		{Kind: "temperature-chart", Title: "Temperature history", Icon: "thermometer", Size: domain.CardMedium,
			Description: "The warmest sensors over the last few minutes.",
			Available:   hasSensors,
			Render:      func() giode.Element { return TemperatureChart(snap(), charts()) }},
		{Kind: "load-chart", Title: "Load average", Icon: "activity", Size: domain.CardMedium,
			Description: "The load average over time, against the number of cores.",
			Render:      func() giode.Element { return LoadChart(snap(), charts()) }},
		{Kind: "disk-io-chart", Title: "Disk activity", Icon: "hard-drive", Size: domain.CardMedium,
			Description: "Disk reads and writes over the last few minutes.",
			Render:      func() giode.Element { return DiskIOChart(snap(), charts()) }},
		{Kind: "memory-breakdown", Title: "Memory breakdown", Icon: "memory-stick", Size: domain.CardMedium,
			Description: "Used, cached and free memory, and swap with its history.",
			Render:      func() giode.Element { return MemoryCard(snap(), charts()) }},
		{Kind: "processes", Title: "Top processes", Icon: "list", Size: domain.CardMedium,
			Description: "The busiest processes, refreshed every few seconds.",
			Render:      procs.Render},
		{Kind: "services", Title: "Services", Icon: "layers", Size: domain.CardMedium,
			Description: "Running and failed systemd services.",
			Available:   svcs.Available, Render: svcs.Render},
		{Kind: "containers", Title: "Containers", Icon: "container", Size: domain.CardMedium,
			Description: "Docker containers and what they use.",
			Available:   boxes.Available, Render: boxes.Render},
		{Kind: "cores", Title: "Cores", Icon: "cpu", Size: domain.CardLarge,
			Description: "Use of each CPU core, as bars.",
			Render:      func() giode.Element { return CoresCard(snap()) }},
		{Kind: "sensors", Title: "Temperatures", Icon: "thermometer", Size: domain.CardLarge,
			Description: "Every temperature sensor, as bars toward its critical point.",
			Available:   hasSensors,
			Render:      func() giode.Element { return SensorsCard(snap()) }},
		{Kind: "disks", Title: "Storage", Icon: "hard-drive", Size: domain.CardMedium,
			Description: "Every mounted file system and how full it is.",
			Render:      func() giode.Element { return DisksCard(snap()) }},
		{Kind: "interfaces", Title: "Network interfaces", Icon: "network", Size: domain.CardMedium,
			Description: "Traffic on each network interface.",
			Render:      func() giode.Element { return NetworkCard(snap()) }},
		{Kind: "connection", Title: "Connection", Icon: "plug", Size: domain.CardSmall,
			Description: "Protocol, address, round trip and how long it's been connected.",
			Render:      func() giode.Element { return ConnectionCard(p.session, snap()) }},
		{Kind: "notes", Title: "Notes", Icon: "notebook-pen", Size: domain.CardSmall,
			Description: "The notes kept with the server.",
			Render:      notes.Render},
	}
}
