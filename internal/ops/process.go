package ops

import (
	"context"
	"errors"
	"strings"
	"time"
)

// ErrProcessGone is returned when a signal finds the process has
// already exited.
var ErrProcessGone = errors.New("the process has already exited")

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
	TTY     string // the terminal it runs in, "" for none
	Unit    string // the systemd unit it belongs to, "" when unknown
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

// Service is the systemd service that runs the process, "" when it
// isn't one (a login session, a user's own programs). systemd starts a
// service's process again when it's killed, so it's stopped through
// the service.
func (p Process) Service() string {
	if !strings.HasSuffix(p.Unit, ".service") || strings.HasPrefix(p.Unit, "user@") {
		return ""
	}
	return p.Unit
}

// unitFromCgroup returns the system unit in a /proc/<pid>/cgroup file:
// the first ".service" or ".scope" element of its path, as ps's unit
// column does (a user's own services sit under user@<uid>.service).
func unitFromCgroup(content string) string {
	for _, line := range strings.Split(content, "\n") {
		_, path, ok := strings.Cut(line, "::") // cgroup v2: "0::/system.slice/nginx.service"
		if !ok {
			if i := strings.Index(line, ":name=systemd:"); i >= 0 {
				path, ok = line[i+len(":name=systemd:"):], true // cgroup v1
			}
		}
		if !ok {
			continue
		}
		parts := strings.Split(strings.TrimSpace(path), "/")
		for _, part := range parts {
			if strings.HasSuffix(part, ".service") || strings.HasSuffix(part, ".scope") {
				return part
			}
		}
	}
	return ""
}
