package ops

import (
	"context"
	"strconv"
	"strings"
	"time"
)

// StopUsersAndUnmount asks the programs using v to quit (SIGTERM), then
// unmounts it.
func (d Drives) StopUsersAndUnmount(ctx context.Context, v Volume, users []DriveUser) error {
	if len(users) > 0 {
		pids := make([]string, len(users))
		for i, u := range users {
			pids[i] = strconv.Itoa(u.PID)
		}
		if _, err := runPrivileged(ctx, d.Exec, "kill -TERM "+strings.Join(pids, " ")); err != nil && !strings.Contains(err.Error(), "No such process") {
			return err
		}
		select {
		case <-time.After(1500 * time.Millisecond): // let them close their files
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return d.Unmount(ctx, v)
}

// UnmountLazily detaches v from every folder right away; the kernel
// finishes unmounting when the programs using it let go. Their open
// files keep working, but nothing new can be opened on the drive.
func (d Drives) UnmountLazily(ctx context.Context, v Volume) error {
	targets, err := d.mountPoints(ctx, v.Path)
	if err != nil {
		return err
	}
	for i := len(targets) - 1; i >= 0; i-- {
		cmd := "umount -l " + Quote(targets[i])
		if targets[i] == autoMountDir(v) {
			cmd += " && { rmdir " + Quote(targets[i]) + " 2>/dev/null; true; }"
		}
		if _, err := runPrivileged(ctx, d.Exec, cmd); err != nil {
			return err
		}
	}
	return nil
}
