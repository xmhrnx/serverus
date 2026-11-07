package monitor

import (
	"strconv"
	"strings"
)

// milliCelsius parses a sysfs temperature in thousandths of a degree,
// rejecting the bogus values some drivers report.
func milliCelsius(s string) (float64, bool) {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, false
	}
	c := v / 1000
	if c <= -40 || c >= 200 {
		return 0, false
	}
	return c, true
}
