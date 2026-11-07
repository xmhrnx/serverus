package ops

import (
	"context"
	"fmt"
	"strings"
)

// Signal sends sig to pid, through sudo when needed.
func (p RemoteProcesses) Signal(ctx context.Context, pid int32, sig Signal) error {
	_, err := runPrivileged(ctx, p.Exec, fmt.Sprintf("kill -%s %d", sig, pid))
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "no such process") {
		return ErrProcessGone
	}
	return err
}
