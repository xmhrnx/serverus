package monitor

import (
	"sort"
)

// derive computes a Snapshot from the current reading and the previous
// one (nil for the first, which yields zero rates).
func derive(cur Reading, prev *Reading) Snapshot {
	s := Snapshot{
		Time:      cur.Time,
		System:    cur.System,
		Memory:    cur.Memory,
		Load:      cur.Load,
		Processes: cur.Processes,
		Disks:     cur.Disks,
		Sensors:   cur.Sensors,
	}
	s.MemPercent = percent(float64(cur.Memory.Used()), float64(cur.Memory.Total))
	s.SwapPercent = percent(float64(cur.Memory.SwapUsed()), float64(cur.Memory.SwapTotal))
	sort.Slice(s.Disks, func(i, j int) bool { return s.Disks[i].Mount < s.Disks[j].Mount })

	if prev == nil {
		s.CorePercents = make([]float64, len(cur.Cores))
		for _, n := range cur.Net {
			s.Net = append(s.Net, NetRate{Interface: n.Interface, RxTotal: n.RxBytes, TxTotal: n.TxBytes})
		}
		return s
	}
	secs := cur.Time.Sub(prev.Time).Seconds()
	if secs <= 0 {
		secs = 1
	}

	s.CPUPercent, s.IOWait = cpuUsage(cur.CPU, prev.CPU)
	s.CorePercents = make([]float64, len(cur.Cores))
	for i, c := range cur.Cores {
		if i < len(prev.Cores) {
			s.CorePercents[i], _ = cpuUsage(c, prev.Cores[i])
		}
	}

	prevDisk := map[string]DiskCounter{}
	for _, d := range prev.DiskIO {
		prevDisk[d.Device] = d
	}
	for _, d := range cur.DiskIO {
		if p, ok := prevDisk[d.Device]; ok {
			s.DiskReadRate += rate(d.ReadBytes, p.ReadBytes, secs)
			s.DiskWriteRate += rate(d.WriteBytes, p.WriteBytes, secs)
		}
	}

	prevNet := map[string]NetCounter{}
	for _, n := range prev.Net {
		prevNet[n.Interface] = n
	}
	for _, n := range cur.Net {
		r := NetRate{Interface: n.Interface, RxTotal: n.RxBytes, TxTotal: n.TxBytes}
		if p, ok := prevNet[n.Interface]; ok {
			r.RxRate = rate(n.RxBytes, p.RxBytes, secs)
			r.TxRate = rate(n.TxBytes, p.TxBytes, secs)
		}
		s.NetRxRate += r.RxRate
		s.NetTxRate += r.TxRate
		s.Net = append(s.Net, r)
	}
	return s
}
