package ops

import (
	"strings"
)

// skippedDisk reports devices that aren't drives a person manages:
// loop devices, RAM disks, zram swap and eMMC boot areas.
func skippedDisk(b blockDevice) bool {
	if b.Type != "disk" {
		return true
	}
	for _, p := range []string{"zram", "ram", "mtdblock"} {
		if strings.HasPrefix(b.Name, p) {
			return true
		}
	}
	return strings.HasPrefix(b.Name, "mmcblk") && (strings.Contains(b.Name, "boot") || strings.Contains(b.Name, "rpmb"))
}
