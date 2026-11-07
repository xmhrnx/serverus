package monitor

// Total is the sum of every state.
func (t CPUTimes) Total() float64 {
	return t.User + t.Nice + t.System + t.Idle + t.IOWait + t.IRQ + t.SoftIRQ + t.Steal
}

// Busy is the time spent doing work.
func (t CPUTimes) Busy() float64 { return t.Total() - t.Idle - t.IOWait }

// Used is the memory taken by programs: everything but what's free or
// reclaimable.
func (m Memory) Used() uint64 {
	if m.Available > 0 && m.Available <= m.Total {
		return m.Total - m.Available
	}
	used := m.Total - m.Free - m.Buffers - m.Cached
	if used > m.Total {
		return 0
	}
	return used
}
