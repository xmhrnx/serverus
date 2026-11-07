package ops

import (
	"context"
	"strconv"
	"strings"
)

// List returns processes, busiest first, leaving out the commands this
// app runs in the background over its own connection (the listing and
// the metrics probe). They start every few seconds and, as ps averages
// CPU over a process's life, sort to the top while they live, only to
// be gone by the time someone picks them.
func (p RemoteProcesses) List(ctx context.Context) ([]Process, error) {
	const columns = "pid=,ppid=,user:32=,pcpu=,pmem=,rss=,etimes=,stat=,tty=,unit=,comm=,args="
	list := func(cols string) (string, error) {
		return run(ctx, p.Exec, "sh -c "+Quote(`echo "self $$"; ps -eo `+cols+` --sort=-pcpu`))
	}
	out, err := list(columns)
	if err != nil {
		// ps built without systemd has no unit column: say "-" instead.
		out, err = list(strings.Replace(columns, "unit=", "unit:1=", 1))
		if err != nil {
			out, err = list(strings.Replace(columns, "unit=,", "tty=,", 1))
		}
	}
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
