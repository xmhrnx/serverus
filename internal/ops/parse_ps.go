package ops

import (
	"strconv"
	"strings"
	"time"
)

// parsePS parses the ps output: eleven fixed columns, then the command
// line with its spaces.
func parsePS(out string) []Process {
	var procs []Process
	for _, line := range strings.Split(out, "\n") {
		f := strings.Fields(line)
		if len(f) < 11 {
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
			RSS: rss * 1024, Elapsed: time.Duration(secs) * time.Second, State: f[7], Name: f[10],
		}
		if strings.Contains(f[9], ".") {
			p.Unit = f[9]
		}
		if f[8] != "?" && f[8] != "-" {
			p.TTY = f[8]
		}
		if len(f) > 11 {
			p.Command = strings.Join(f[11:], " ")
		} else {
			p.Command = p.Name
		}
		procs = append(procs, p)
	}
	return procs
}
