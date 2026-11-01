package monitor

import (
	"sort"
	"time"
)

// Snapshot is the state of a machine derived from two readings: what
// the dashboard shows.
type Snapshot struct {
	Time   time.Time
	System SystemInfo

	CPUPercent   float64
	IOWait       float64
	CorePercents []float64

	Memory      Memory
	MemPercent  float64
	SwapPercent float64

	Load      [3]float64
	Processes int

	Disks         []DiskUsage
	DiskReadRate  float64 // bytes per second, every disk together
	DiskWriteRate float64
	Net           []NetRate
	NetRxRate     float64 // bytes per second, every interface together
	NetTxRate     float64
	Sensors       []Sensor
}

// NetRate is one interface's traffic.
type NetRate struct {
	Interface        string
	RxRate, TxRate   float64 // bytes per second
	RxTotal, TxTotal uint64  // since boot
}

// HottestSensor returns the warmest sensor, ok false when there is none.
func (s Snapshot) HottestSensor() (Sensor, bool) {
	if len(s.Sensors) == 0 {
		return Sensor{}, false
	}
	best := s.Sensors[0]
	for _, x := range s.Sensors[1:] {
		if x.Celsius > best.Celsius {
			best = x
		}
	}
	return best, true
}

// RootDisk returns the file system mounted at "/" (or the largest one
// when there's no root, as on Windows).
func (s Snapshot) RootDisk() (DiskUsage, bool) {
	var best DiskUsage
	found := false
	for _, d := range s.Disks {
		if d.Mount == "/" {
			return d, true
		}
		if !found || d.Total > best.Total {
			best, found = d, true
		}
	}
	return best, found
}

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

// cpuUsage returns the busy and I/O wait percentages between two
// readings of the same CPU.
func cpuUsage(cur, prev CPUTimes) (busy, iowait float64) {
	total := cur.Total() - prev.Total()
	if total <= 0 {
		return 0, 0
	}
	return percent(cur.Busy()-prev.Busy(), total), percent(cur.IOWait-prev.IOWait, total)
}

// rate is the per-second change of a counter; a counter that went
// backwards (a reboot, a wrapped 32-bit counter) yields 0.
func rate(cur, prev uint64, secs float64) float64 {
	if cur < prev {
		return 0
	}
	return float64(cur-prev) / secs
}
