package monitor

import (
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
