package monitor

import (
	"strings"
)

// isPartition reports whether name is a partition of another listed
// device: sda1 of sda, nvme0n1p2 of nvme0n1, mmcblk0p1 of mmcblk0.
func isPartition(name string, names map[string]bool) bool {
	base := strings.TrimRight(name, "0123456789")
	if base == name {
		return false
	}
	if names[base] {
		return true
	}
	if b, ok := strings.CutSuffix(base, "p"); ok && names[b] {
		return true
	}
	return false
}
