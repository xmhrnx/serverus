package monitor

func (h *History) charts() Charts {
	c := Charts{
		CPU: h.CPU.Values(), Memory: h.Memory.Values(),
		NetRx: h.NetRx.Values(), NetTx: h.NetTx.Values(),
		DiskRead: h.DiskRead.Values(), DiskWrite: h.DiskWrite.Values(),
		Temperature: h.Temperature.Values(),
		Load:        h.Load.Values(), Swap: h.Swap.Values(), IOWait: h.IOWait.Values(),
		Sensors: make(map[string][]float64, len(h.Sensors)),
	}
	for name, s := range h.Sensors {
		c.Sensors[name] = s.Values()
	}
	return c
}
