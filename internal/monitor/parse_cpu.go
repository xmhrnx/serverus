package monitor

import (
	"strconv"
	"strings"
)

// parseCPU reads the cpu lines of /proc/stat: the first is every core
// together, then one per core.
func parseCPU(lines []string) (total CPUTimes, cores []CPUTimes) {
	for _, line := range lines {
		f := strings.Fields(line)
		if len(f) < 5 || !strings.HasPrefix(f[0], "cpu") {
			continue
		}
		v := make([]float64, 8)
		for i := 0; i < 8 && i+1 < len(f); i++ {
			v[i], _ = strconv.ParseFloat(f[i+1], 64)
		}
		t := CPUTimes{User: v[0], Nice: v[1], System: v[2], Idle: v[3], IOWait: v[4], IRQ: v[5], SoftIRQ: v[6], Steal: v[7]}
		if f[0] == "cpu" {
			total = t
		} else {
			cores = append(cores, t)
		}
	}
	return total, cores
}
