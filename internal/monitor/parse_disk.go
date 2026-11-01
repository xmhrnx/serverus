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
