package monitor

import (
	"strconv"
	"strings"
)

// parseDiskStats reads /proc/diskstats, keeping whole physical disks
// (not partitions, loop devices or RAM disks). Sectors are 512 bytes.
func parseDiskStats(lines []string) []DiskCounter {
	type row struct {
		name string
		r, w uint64
	}
	var rows []row
	names := map[string]bool{}
	for _, line := range lines {
		f := strings.Fields(line)
		if len(f) < 10 {
			continue
		}
		name := f[2]
		if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") ||
			strings.HasPrefix(name, "zram") || strings.HasPrefix(name, "sr") || strings.HasPrefix(name, "fd") {
			continue
		}
		r, _ := strconv.ParseUint(f[5], 10, 64)
		w, _ := strconv.ParseUint(f[9], 10, 64)
		rows = append(rows, row{name, r * 512, w * 512})
		names[name] = true
	}
	var out []DiskCounter
	for _, r := range rows {
		if isPartition(r.name, names) {
			continue
		}
		out = append(out, DiskCounter{Device: r.name, ReadBytes: r.r, WriteBytes: r.w})
	}
	return out
}
