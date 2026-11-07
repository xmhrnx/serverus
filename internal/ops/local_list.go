package ops

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

// List returns processes, busiest first.
func (LocalProcesses) List(ctx context.Context) ([]Process, error) {
	ps, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	out := make([]Process, 0, len(ps))
	for _, p := range ps {
		name, err := p.NameWithContext(ctx)
		if err != nil {
			continue
		}
		proc := Process{PID: p.Pid, Name: name}
		proc.PPID, _ = p.PpidWithContext(ctx)
		proc.User, _ = p.UsernameWithContext(ctx)
		proc.CPU, _ = p.CPUPercentWithContext(ctx)
		if m, err := p.MemoryPercentWithContext(ctx); err == nil {
			proc.Memory = float64(m)
		}
		if mi, err := p.MemoryInfoWithContext(ctx); err == nil {
			proc.RSS = mi.RSS
		}
		if created, err := p.CreateTimeWithContext(ctx); err == nil {
			proc.Elapsed = now.Sub(time.UnixMilli(created)).Truncate(time.Second)
		}
		if st, err := p.StatusWithContext(ctx); err == nil && len(st) > 0 {
			proc.State = st[0]
		}
		if runtime.GOOS == "linux" {
			if b, err := os.ReadFile(fmt.Sprintf("/proc/%d/cgroup", p.Pid)); err == nil {
				proc.Unit = unitFromCgroup(string(b))
			}
		}
		proc.Command, _ = p.CmdlineWithContext(ctx)
		if proc.Command == "" {
			proc.Command = name
		}
		out = append(out, proc)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CPU > out[j].CPU })
	if len(out) > processLimit {
		out = out[:processLimit]
	}
	return out, nil
}
