package ops

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
)

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
