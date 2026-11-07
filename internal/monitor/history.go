package monitor

// Series is a fixed-size ring of the most recent values of one metric.
type Series struct {
	values []float64
	next   int
	full   bool
}

// History keeps the recent values the dashboard charts.
type History struct {
	CPU, Memory, NetRx, NetTx, DiskRead, DiskWrite, Temperature *Series
}

// Charts is a copy of the history's values, safe to read anywhere.
type Charts struct {
	CPU, Memory, NetRx, NetTx, DiskRead, DiskWrite, Temperature []float64
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

// HistoryLength is how many samples each series keeps.
const HistoryLength = 90
