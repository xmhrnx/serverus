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

// Memory is in bytes.
type Memory struct {
	Total, Available, Free, Buffers, Cached uint64
	SwapTotal, SwapFree                     uint64
}
