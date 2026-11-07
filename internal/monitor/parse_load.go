package monitor

import (
	"strconv"
	"strings"
)

// parseLoad reads /proc/loadavg.
func parseLoad(lines []string) (load [3]float64) {
	if len(lines) == 0 {
		return load
	}
	f := strings.Fields(lines[0])
	for i := 0; i < 3 && i < len(f); i++ {
		load[i], _ = strconv.ParseFloat(f[i], 64)
	}
	return load
}

// parseCount reads a single integer line.
func parseCount(lines []string) int {
	if len(lines) == 0 {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
	return n
}
