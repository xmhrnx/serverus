package monitor

import (
	"testing"
)

func TestSeriesRing(t *testing.T) {
	s := NewSeries(3)
	for i := 1; i <= 5; i++ {
		s.Push(float64(i))
	}
	got := s.Values()
	if len(got) != 3 || got[0] != 3 || got[2] != 5 {
		t.Fatalf("got %v", got)
	}
}
