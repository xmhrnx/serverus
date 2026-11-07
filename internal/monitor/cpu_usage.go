package monitor

// cpuUsage returns the busy and I/O wait percentages between two
// readings of the same CPU.
func cpuUsage(cur, prev CPUTimes) (busy, iowait float64) {
	total := cur.Total() - prev.Total()
	if total <= 0 {
		return 0, 0
	}
	return percent(cur.Busy()-prev.Busy(), total), percent(cur.IOWait-prev.IOWait, total)
}

// rate is the per-second change of a counter; a counter that went
// backwards (a reboot, a wrapped 32-bit counter) yields 0.
func rate(cur, prev uint64, secs float64) float64 {
	if cur < prev {
		return 0
	}
	return float64(cur-prev) / secs
}
