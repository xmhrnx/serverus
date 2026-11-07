package monitor

import (
	"strconv"
	"strings"
)

// ignoredFS are file system types that aren't storage.
var ignoredFS = map[string]bool{
	"tmpfs": true, "devtmpfs": true, "squashfs": true, "overlay": true, "proc": true,
	"sysfs": true, "cgroup": true, "cgroup2": true, "devpts": true, "efivarfs": true,
	"ramfs": true, "autofs": true, "fuse.lxcfs": true, "nsfs": true, "tracefs": true,
}

// parseDF reads "df -PTk": device, type, 1K-blocks, used, available,
// capacity, mount point (which may contain spaces).
func parseDF(lines []string) []DiskUsage {
	var out []DiskUsage
	seen := map[string]bool{}
	for i, line := range lines {
		f := strings.Fields(line)
		if i == 0 && len(f) > 0 && f[0] == "Filesystem" {
			continue
		}
		if len(f) < 7 {
			continue
		}
		if ignoredFS[f[1]] || strings.HasPrefix(f[0], "/dev/loop") {
			continue
		}
		total, err1 := strconv.ParseUint(f[2], 10, 64)
		used, err2 := strconv.ParseUint(f[3], 10, 64)
		free, err3 := strconv.ParseUint(f[4], 10, 64)
		if err1 != nil || err2 != nil || err3 != nil || total == 0 {
			continue
		}
		// One entry per device: bind mounts repeat it.
		if seen[f[0]] {
			continue
		}
		seen[f[0]] = true
		out = append(out, DiskUsage{
			Device: f[0], FSType: f[1], Mount: strings.Join(f[6:], " "),
			Total: total * 1024, Used: used * 1024, Free: free * 1024,
		})
	}
	return out
}
