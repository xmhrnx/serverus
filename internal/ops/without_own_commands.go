package ops

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
