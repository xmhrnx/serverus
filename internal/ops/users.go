package ops

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

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
