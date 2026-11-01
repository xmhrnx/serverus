package ops

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xmhrnx/bastion/internal/transport"
)

// RemoteProcesses manages a Linux server's processes with ps and kill.
type RemoteProcesses struct {
	Exec transport.Executor
}

// processLimit bounds the list to the busiest processes.
const processLimit = 400

// List returns processes, busiest first.
func (p RemoteProcesses) List(ctx context.Context) ([]Process, error) {
	cmd := fmt.Sprintf("ps -eo pid=,ppid=,user:32=,pcpu=,pmem=,rss=,etimes=,stat=,comm=,args= --sort=-pcpu | head -n %d", processLimit)
	out, err := run(ctx, p.Exec, cmd)
	if err != nil {
		return nil, err
	}
	return parsePS(out), nil
}

// Signal sends sig to pid, through sudo when needed.
func (p RemoteProcesses) Signal(ctx context.Context, pid int32, sig Signal) error {
	_, err := runPrivileged(ctx, p.Exec, fmt.Sprintf("kill -%s %d", sig, pid))
	return err
}

// parsePS parses the ps output: nine fixed columns, then the command
// line with its spaces.
func parsePS(out string) []Process {
	var procs []Process
	for _, line := range strings.Split(out, "\n") {
		f := strings.Fields(line)
		if len(f) < 9 {
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
			RSS: rss * 1024, Elapsed: time.Duration(secs) * time.Second, State: f[7], Name: f[8],
		}
		if len(f) > 9 {
			p.Command = strings.Join(f[9:], " ")
		} else {
			p.Command = p.Name
		}
		procs = append(procs, p)
	}
	return procs
}
