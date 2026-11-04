package ops

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

// DriveUser is a program keeping a mounted drive busy.
type DriveUser struct {
	PID  int
	Name string
	Path string // the open file or working folder on the drive
}

// BusyError is returned by Unmount when programs still use the volume.
type BusyError struct {
	Volume  Volume
	Targets []string    // where it's still mounted
	Users   []DriveUser // the programs found; may be empty
}

func (e *BusyError) Error() string {
	if len(e.Users) == 0 {
		return e.Volume.Title() + " is in use"
	}
	names := []string{}
	for _, u := range e.Users {
		if !slices.Contains(names, u.Name) {
			names = append(names, u.Name)
		}
	}
	return e.Volume.Title() + " is in use by " + strings.Join(names, ", ")
}

func isBusy(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "busy")
}

func (d Drives) busyError(ctx context.Context, v Volume, targets []string) error {
	return &BusyError{Volume: v, Targets: targets, Users: d.users(ctx, targets)}
}

// users finds the programs with files or working folders under targets,
// from /proc. It looks as root when sudo needs no password, and at the
// user's own programs otherwise; it never asks.
func (d Drives) users(ctx context.Context, targets []string) []DriveUser {
	const list = `find /proc/[0-9]*/cwd /proc/[0-9]*/root /proc/[0-9]*/fd -maxdepth 1 -type l -printf '%p\t%l\n' 2>/dev/null; true`
	out, err := run(ctx, d.Exec, "sudo -n sh -c "+Quote(list))
	if err != nil {
		out, _ = run(ctx, d.Exec, "sh -c "+Quote(list))
	}
	seen := map[int]bool{}
	var users []DriveUser
	for _, line := range strings.Split(out, "\n") {
		link, target, ok := strings.Cut(line, "\t")
		if !ok || !underAny(target, targets) {
			continue
		}
		parts := strings.Split(link, "/") // "", "proc", pid, ...
		if len(parts) < 4 {
			continue
		}
		pid, err := strconv.Atoi(parts[2])
		if err != nil || seen[pid] {
			continue
		}
		seen[pid] = true
		users = append(users, DriveUser{PID: pid, Path: target})
	}
	if len(users) == 0 {
		return nil
	}
	pids := make([]string, len(users))
	for i, u := range users {
		pids[i] = strconv.Itoa(u.PID)
	}
	script := "for p in " + strings.Join(pids, " ") + `; do printf '%s\t%s\n' "$p" "$(cat /proc/$p/comm 2>/dev/null)"; done`
	names, _ := run(ctx, d.Exec, "sh -c "+Quote(script))
	byPID := map[int]string{}
	for _, line := range strings.Split(names, "\n") {
		if p, n, ok := strings.Cut(line, "\t"); ok {
			if pid, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
				byPID[pid] = strings.TrimSpace(n)
			}
		}
	}
	for i := range users {
		users[i].Name = byPID[users[i].PID]
		if users[i].Name == "" {
			users[i].Name = fmt.Sprintf("process %d", users[i].PID)
		}
	}
	return users
}

func underAny(p string, dirs []string) bool {
	for _, d := range dirs {
		if p == d || strings.HasPrefix(p, strings.TrimSuffix(d, "/")+"/") {
			return true
		}
	}
	return false
}

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
