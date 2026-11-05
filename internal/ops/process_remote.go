package ops

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xmhrnx/serverus/internal/transport"
)

// RemoteProcesses manages a Linux server's processes with ps and kill.
type RemoteProcesses struct {
	Exec transport.Executor
}

// processLimit bounds the list to the busiest processes.
const processLimit = 400

// List returns processes, busiest first, leaving out the commands this
// app runs in the background over its own connection (the listing and
// the metrics probe). They start every few seconds and, as ps averages
// CPU over a process's life, sort to the top while they live, only to
// be gone by the time someone picks them.
func (p RemoteProcesses) List(ctx context.Context) ([]Process, error) {
	cmd := "sh -c " + Quote(`echo "self $$"; ps -eo pid=,ppid=,user:32=,pcpu=,pmem=,rss=,etimes=,stat=,tty=,comm=,args= --sort=-pcpu`)
	out, err := run(ctx, p.Exec, cmd)
	if err != nil {
		return nil, err
	}
	self, rest, _ := strings.Cut(out, "\n")
	selfPID, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(self, "self")))
	procs := withoutOwnCommands(parsePS(rest), int32(selfPID))
	if len(procs) > processLimit {
		procs = procs[:processLimit]
	}
	return procs, nil
}

// withoutOwnCommands drops the processes without a terminal that run
// under the SSH connection self belongs to. Programs started in this
// app's terminal tabs have one, so they stay.
func withoutOwnCommands(procs []Process, self int32) []Process {
	byPID := make(map[int32]Process, len(procs))
	for _, p := range procs {
		byPID[p.PID] = p
	}
	conn := int32(0)
	for pid, hops := self, 0; pid > 1 && hops < 32; hops++ {
		p, ok := byPID[pid]
		if !ok {
			break
		}
		if p.Name == "sshd" || p.Name == "sshd-session" || p.Name == "dropbear" {
			conn = pid
			break
		}
		pid = p.PPID
	}
	if conn == 0 {
		return procs
	}
	under := func(p Process) bool {
		for pid, hops := p.PPID, 0; pid > 1 && hops < 64; hops++ {
			if pid == conn {
				return true
			}
			pid = byPID[pid].PPID
		}
		return false
	}
	out := procs[:0]
	for _, p := range procs {
		if p.TTY == "" && under(p) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// Signal sends sig to pid, through sudo when needed.
func (p RemoteProcesses) Signal(ctx context.Context, pid int32, sig Signal) error {
	_, err := runPrivileged(ctx, p.Exec, fmt.Sprintf("kill -%s %d", sig, pid))
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "no such process") {
		return ErrProcessGone
	}
	return err
}

// parsePS parses the ps output: ten fixed columns, then the command
// line with its spaces.
func parsePS(out string) []Process {
	var procs []Process
	for _, line := range strings.Split(out, "\n") {
		f := strings.Fields(line)
		if len(f) < 10 {
			continue
		}
		pid, err := strconv.ParseInt(f[0], 10, 32)
		if err != nil {
			continue
		}
		ppid, _ := strconv.ParseInt(f[1], 10, 32)
		cpu, _ := strconv.ParseFloat(f[3], 64)
		mem, _ := strconv.ParseFloat(f[4], 64)
		rss, _ := strconv.ParseUint(f[5], 10, 64)
		secs, _ := strconv.ParseInt(f[6], 10, 64)
		p := Process{
			PID: int32(pid), PPID: int32(ppid), User: f[2], CPU: cpu, Memory: mem,
			RSS: rss * 1024, Elapsed: time.Duration(secs) * time.Second, State: f[7], Name: f[9],
		}
		if f[8] != "?" && f[8] != "-" {
			p.TTY = f[8]
		}
		if len(f) > 10 {
			p.Command = strings.Join(f[10:], " ")
		} else {
			p.Command = p.Name
		}
		procs = append(procs, p)
	}
	return procs
}
