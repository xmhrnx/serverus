package monitor

import (
	"strconv"
	"strings"
)

// parseNetDev reads /proc/net/dev, skipping the loopback interface.
func parseNetDev(lines []string) []NetCounter {
	var out []NetCounter
	for _, line := range lines {
		name, rest, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		name = strings.TrimSpace(name)
		if name == "lo" {
			continue
		}
		f := strings.Fields(rest)
		if len(f) < 9 {
			continue
		}
		rx, _ := strconv.ParseUint(f[0], 10, 64)
		tx, _ := strconv.ParseUint(f[8], 10, 64)
		out = append(out, NetCounter{Interface: name, RxBytes: rx, TxBytes: tx})
	}
	return out
}
