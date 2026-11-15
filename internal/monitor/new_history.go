package monitor

func newHistory() *History {
	return &History{
		CPU: NewSeries(HistoryLength), Memory: NewSeries(HistoryLength),
		NetRx: NewSeries(HistoryLength), NetTx: NewSeries(HistoryLength),
		DiskRead: NewSeries(HistoryLength), DiskWrite: NewSeries(HistoryLength),
		Temperature: NewSeries(HistoryLength),
		Load:        NewSeries(HistoryLength), Swap: NewSeries(HistoryLength), IOWait: NewSeries(HistoryLength),
		Sensors: map[string]*Series{},
	}
}

func (h *History) push(s Snapshot) {
	h.CPU.Push(s.CPUPercent)
	h.Memory.Push(s.MemPercent)
	h.NetRx.Push(s.NetRxRate)
	h.NetTx.Push(s.NetTxRate)
	h.DiskRead.Push(s.DiskReadRate)
	h.DiskWrite.Push(s.DiskWriteRate)
	h.Load.Push(s.Load[0])
	h.Swap.Push(s.SwapPercent)
	h.IOWait.Push(s.IOWait)
	t := 0.0
	if hot, ok := s.HottestSensor(); ok {
		t = hot.Celsius
	}
	h.Temperature.Push(t)
	for _, x := range s.Sensors {
		series, ok := h.Sensors[x.Name]
		if !ok {
			if len(h.Sensors) >= maxSensorSeries {
				continue
			}
			series = NewSeries(HistoryLength)
			h.Sensors[x.Name] = series
		}
		series.Push(x.Celsius)
	}
}
