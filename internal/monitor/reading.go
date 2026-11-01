// Package monitor collects a server's resource usage: CPU, memory,
// disks, network, temperatures and system information. A Source takes
// raw samples; a Collector turns consecutive samples into rates and
// keeps a short history for charts.
package monitor

import (
	"context"
	"time"
)

// Source takes one Reading of a machine.
type Source interface {
	Sample(ctx context.Context) (Reading, error)
}

// Reading is one raw sample. Counters (CPU times, bytes) are cumulative
// since boot; the Collector turns two readings into rates.
type Reading struct {
	Time      time.Time
	System    SystemInfo
	CPU       CPUTimes   // all cores together
	Cores     []CPUTimes // one per core
	Memory    Memory
	Load      [3]float64
	Disks     []DiskUsage
	DiskIO    []DiskCounter
	Net       []NetCounter
	Sensors   []Sensor
	Processes int
}

// SystemInfo describes the machine.
type SystemInfo struct {
	Hostname string
	OS       string // "Ubuntu 24.04.1 LTS"
	Kernel   string // "Linux 6.8.0-45-generic x86_64"
	CPUModel string
	Cores    int
	Uptime   time.Duration
}

// CPUTimes are cumulative CPU times in any consistent unit.
type CPUTimes struct {
	User, Nice, System, Idle, IOWait, IRQ, SoftIRQ, Steal float64
}

// Total is the sum of every state.
func (t CPUTimes) Total() float64 {
	return t.User + t.Nice + t.System + t.Idle + t.IOWait + t.IRQ + t.SoftIRQ + t.Steal
}

// Busy is the time spent doing work.
func (t CPUTimes) Busy() float64 { return t.Total() - t.Idle - t.IOWait }

// Memory is in bytes.
type Memory struct {
	Total, Available, Free, Buffers, Cached uint64
	SwapTotal, SwapFree                     uint64
}

// Used is the memory taken by programs: everything but what's free or
// reclaimable.
func (m Memory) Used() uint64 {
	if m.Available > 0 && m.Available <= m.Total {
		return m.Total - m.Available
	}
	used := m.Total - m.Free - m.Buffers - m.Cached
	if used > m.Total {
		return 0
	}
	return used
}

// SwapUsed is the swap in use.
func (m Memory) SwapUsed() uint64 {
	if m.SwapFree > m.SwapTotal {
		return 0
	}
	return m.SwapTotal - m.SwapFree
}

// DiskUsage is one mounted file system, in bytes.
type DiskUsage struct {
	Device, FSType, Mount string
	Total, Used, Free     uint64
}

// Percent is the share of the file system in use, 0-100.
func (d DiskUsage) Percent() float64 { return percent(float64(d.Used), float64(d.Used+d.Free)) }

// DiskCounter is the bytes a block device has read and written.
type DiskCounter struct {
	Device                string
	ReadBytes, WriteBytes uint64
}

// NetCounter is the bytes an interface has received and sent.
type NetCounter struct {
	Interface        string
	RxBytes, TxBytes uint64
}

// Sensor is a temperature in degrees Celsius. High and Critical are 0
// when the hardware doesn't report them.
type Sensor struct {
	Name           string
	Celsius        float64
	High, Critical float64
}

// Level grades the temperature: 0 normal, 1 warm, 2 hot.
func (s Sensor) Level() int {
	crit := s.Critical
	if crit <= 0 {
		crit = 95
	}
	high := s.High
	if high <= 0 || high >= crit {
		high = crit - 15
	}
	switch {
	case s.Celsius >= crit-5:
		return 2
	case s.Celsius >= high:
		return 1
	}
	return 0
}

func percent(part, whole float64) float64 {
	if whole <= 0 {
		return 0
	}
	p := part / whole * 100
	return max(0, min(100, p))
}
