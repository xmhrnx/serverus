package monitor

func newHistory() *History {
	return &History{
		CPU: NewSeries(HistoryLength), Memory: NewSeries(HistoryLength),
		NetRx: NewSeries(HistoryLength), NetTx: NewSeries(HistoryLength),
		DiskRead: NewSeries(HistoryLength), DiskWrite: NewSeries(HistoryLength),
		Temperature: NewSeries(HistoryLength),
	}
}

func (h *History) push(s Snapshot) {
	h.CPU.Push(s.CPUPercent)
	h.Memory.Push(s.MemPercent)
	h.NetRx.Push(s.NetRxRate)
	h.NetTx.Push(s.NetTxRate)
	h.DiskRead.Push(s.DiskReadRate)
	h.DiskWrite.Push(s.DiskWriteRate)
	t := 0.0
	if hot, ok := s.HottestSensor(); ok {
		t = hot.Celsius
	}
	h.Temperature.Push(t)
}
