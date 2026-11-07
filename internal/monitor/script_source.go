package monitor

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/xmhrnx/serverus/internal/transport"
)

// ScriptSource samples a remote Linux server by running a shell script
// over the connection.
type ScriptSource struct {
	Exec transport.Executor
}

// ErrNotLinux is returned by ScriptSource when the server isn't Linux:
// the probe reads Linux's /proc and /sys.
var ErrNotLinux = errors.New("resource monitoring needs a Linux server")

// Sample runs the probe and parses its output.
func (s ScriptSource) Sample(ctx context.Context) (Reading, error) {
	out, err := s.Exec.Run(ctx, probeScript)
	if err != nil {
		return Reading{}, err
	}
	return parseProbe(string(out.Stdout), time.Now())
}

// parseProbe turns the probe's output into a Reading.
func parseProbe(out string, now time.Time) (Reading, error) {
	sec := splitSections(out)
	if _, ok := sec["end"]; !ok || len(sec["stat"]) == 0 {
		if strings.TrimSpace(out) == "" {
			return Reading{}, errors.New("the server returned no data")
		}
		return Reading{}, ErrNotLinux
	}
	r := Reading{Time: now}
	r.System = parseSystem(sec)
	r.CPU, r.Cores = parseCPU(sec["stat"])
	r.System.Cores = len(r.Cores)
	r.Memory = parseMemory(sec["meminfo"])
	r.Load = parseLoad(sec["loadavg"])
	r.Disks = parseDF(sec["df"])
	r.DiskIO = parseDiskStats(sec["diskstats"])
	r.Net = parseNetDev(sec["netdev"])
	r.Sensors = parseSensors(sec["hwmon"], sec["thermal"])
	r.Processes = parseCount(sec["procs"])
	return r, nil
}
