package monitor

import (
	"strconv"
	"strings"
	"time"
)

// parseSystem reads the host, os, uptime and cpuinfo sections.
func parseSystem(sec map[string][]string) SystemInfo {
	var si SystemInfo
	if l := sec["host"]; len(l) > 0 {
		si.Hostname = strings.TrimSpace(l[0])
	}
	if l := sec["os"]; len(l) > 0 {
		si.OS = strings.TrimSpace(l[0])
		if len(l) > 1 {
			si.Kernel = strings.TrimSpace(l[1])
		} else if si.OS != "" && strings.HasPrefix(si.OS, "Linux") {
			si.Kernel, si.OS = si.OS, ""
		}
		if si.OS == "" {
			si.OS = "Linux"
		}
	}
	if l := sec["uptime"]; len(l) > 0 {
		if f := strings.Fields(l[0]); len(f) > 0 {
			if secs, err := strconv.ParseFloat(f[0], 64); err == nil {
				si.Uptime = time.Duration(secs * float64(time.Second))
			}
		}
	}
	if l := sec["cpuinfo"]; len(l) > 0 {
		if _, v, ok := strings.Cut(l[0], ":"); ok {
			si.CPUModel = strings.Join(strings.Fields(v), " ")
		}
	}
	return si
}
