package monitor

// Sensor is a temperature in degrees Celsius. High and Critical are 0
// when the hardware doesn't report them.
type Sensor struct {
	Name           string
	Celsius        float64
	High, Critical float64
}

// Level grades the temperature: 0 normal, 1 warm, 2 hot.
func (s Sensor) Level() int {
	crit := s.Critical
	if crit <= 0 {
		crit = 95
	}
	high := s.High
	if high <= 0 || high >= crit {
		high = crit - 15
	}
	switch {
	case s.Celsius >= crit-5:
		return 2
	case s.Celsius >= high:
		return 1
	}
	return 0
}

func percent(part, whole float64) float64 {
	if whole <= 0 {
		return 0
	}
	p := part / whole * 100
	return max(0, min(100, p))
}
