package monitor

// Series is a fixed-size ring of the most recent values of one metric.
type Series struct {
	values []float64
	next   int
	full   bool
}

// NewSeries returns an empty series holding up to n values.
func NewSeries(n int) *Series { return &Series{values: make([]float64, n)} }

// Push appends v, dropping the oldest value when full.
func (s *Series) Push(v float64) {
	s.values[s.next] = v
	s.next = (s.next + 1) % len(s.values)
	if s.next == 0 {
		s.full = true
	}
}

// Values returns the values oldest first.
func (s *Series) Values() []float64 {
	if !s.full {
		return append([]float64(nil), s.values[:s.next]...)
	}
	out := make([]float64, 0, len(s.values))
	out = append(out, s.values[s.next:]...)
	return append(out, s.values[:s.next]...)
}

// Capacity is how many values the series keeps.
func (s *Series) Capacity() int { return len(s.values) }

// History keeps the recent values the dashboard charts.
type History struct {
	CPU, Memory, NetRx, NetTx, DiskRead, DiskWrite, Temperature *Series
}

// HistoryLength is how many samples each series keeps.
const HistoryLength = 90

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

// Charts is a copy of the history's values, safe to read anywhere.
type Charts struct {
	CPU, Memory, NetRx, NetTx, DiskRead, DiskWrite, Temperature []float64
}

func (h *History) charts() Charts {
	return Charts{
		CPU: h.CPU.Values(), Memory: h.Memory.Values(),
		NetRx: h.NetRx.Values(), NetTx: h.NetTx.Values(),
		DiskRead: h.DiskRead.Values(), DiskWrite: h.DiskWrite.Values(),
		Temperature: h.Temperature.Values(),
	}
}
