package ops

import (
	"context"
)

// Cancel calls off a scheduled shutdown.
func (p Power) Cancel(ctx context.Context) error {
	switch p.OS {
	case "windows":
		_, err := run(ctx, p.Exec, "shutdown /a")
		return err
	case "darwin":
		_, err := runPrivileged(ctx, p.Exec, "killall shutdown")
		return err
	}
	_, err := runPrivileged(ctx, p.Exec, "shutdown -c")
	return err
}

// Pending reports a scheduled shutdown, on systemd machines, which keep
// it in /run/systemd/shutdown/scheduled. ok is false when nothing is
// scheduled or the machine can't tell.
func (p Power) Pending(ctx context.Context) (Scheduled, bool) {
	if p.OS == "windows" || p.OS == "darwin" {
		return Scheduled{}, false
	}
	out, err := run(ctx, p.Exec, "cat /run/systemd/shutdown/scheduled 2>/dev/null || true")
	if err != nil {
		return Scheduled{}, false
	}
	return parseScheduled(out)
}
