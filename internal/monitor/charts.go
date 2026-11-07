package monitor

func (h *History) charts() Charts {
	return Charts{
		CPU: h.CPU.Values(), Memory: h.Memory.Values(),
		NetRx: h.NetRx.Values(), NetTx: h.NetTx.Values(),
		DiskRead: h.DiskRead.Values(), DiskWrite: h.DiskWrite.Values(),
		Temperature: h.Temperature.Values(),
	}
}
