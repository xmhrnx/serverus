package ops

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/process"

	"github.com/xmhrnx/serverus/internal/transport"
)

// LocalProcesses manages this computer's processes, on any platform.
// Exec, when set, sends signals through sudo to processes of other
// users on Linux and macOS.
type LocalProcesses struct {
	Exec transport.Executor
}

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

// Signal sends sig to pid. Windows has no signals: every signal ends
// the process there.
func (l LocalProcesses) Signal(ctx context.Context, pid int32, sig Signal) error {
	err := l.signal(pid, sig)
	switch {
	case errors.Is(err, os.ErrProcessDone) || errors.Is(err, syscall.ESRCH):
		return ErrProcessGone
	case errors.Is(err, syscall.EPERM) && l.Exec != nil && runtime.GOOS != "windows":
		_, err = runAsRoot(ctx, l.Exec, fmt.Sprintf("kill -%s %d", sig, pid))
		if err != nil && strings.Contains(strings.ToLower(err.Error()), "no such process") {
			return ErrProcessGone
		}
	}
	return err
}

func (LocalProcesses) signal(pid int32, sig Signal) error {
	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}
	switch sig {
	case SignalKill:
		return proc.Kill()
	case SignalHup:
		return signalOrKill(proc, syscall.Signal(1))
	case SignalTerm:
		return signalOrKill(proc, syscall.Signal(15))
	}
	return fmt.Errorf("unknown signal %s", sig)
}

func signalOrKill(proc *os.Process, sig os.Signal) error {
	if err := proc.Signal(sig); err != nil {
		if err.Error() == "not supported by windows" {
			return proc.Kill()
		}
		return err
	}
	return nil
}
