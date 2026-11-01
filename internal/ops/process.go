package ops

import (
	"context"
	"time"
)

// Process is one running program.
type Process struct {
	PID     int32
	PPID    int32
	User    string
	CPU     float64 // percent of one core
	Memory  float64 // percent of RAM
	RSS     uint64  // bytes
	Elapsed time.Duration
	State   string
	Name    string
	Command string
}

// Signal is a signal sent to a process.
type Signal string

const (
	// SignalTerm asks the process to exit.
	SignalTerm Signal = "TERM"
	// SignalKill ends it immediately.
	SignalKill Signal = "KILL"
	// SignalHup asks a daemon to reload its configuration.
	SignalHup Signal = "HUP"
)

// Processes lists and signals processes.
type Processes interface {
	List(ctx context.Context) ([]Process, error)
	Signal(ctx context.Context, pid int32, sig Signal) error
}
