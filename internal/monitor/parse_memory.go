package monitor

import (
	"strconv"
	"strings"
)

// parseMemory reads /proc/meminfo, whose values are in kB.
func parseMemory(lines []string) Memory {
	var m Memory
	for _, line := range lines {
		key, rest, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		f := strings.Fields(rest)
		if len(f) == 0 {
			continue
		}
		kb, err := strconv.ParseUint(f[0], 10, 64)
		if err != nil {
			continue
		}
		b := kb * 1024
		switch key {
		case "MemTotal":
			m.Total = b
		case "MemAvailable":
			m.Available = b
		case "MemFree":
			m.Free = b
		case "Buffers":
			m.Buffers = b
		case "Cached":
			m.Cached = b
		case "SwapTotal":
			m.SwapTotal = b
		case "SwapFree":
			m.SwapFree = b
		}
	}
	return m
}
