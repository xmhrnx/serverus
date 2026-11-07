package ops

import (
	"strings"
)

// unitFromCgroup returns the system unit in a /proc/<pid>/cgroup file:
// the first ".service" or ".scope" element of its path, as ps's unit
// column does (a user's own services sit under user@<uid>.service).
func unitFromCgroup(content string) string {
	for _, line := range strings.Split(content, "\n") {
		_, path, ok := strings.Cut(line, "::") // cgroup v2: "0::/system.slice/nginx.service"
		if !ok {
			if i := strings.Index(line, ":name=systemd:"); i >= 0 {
				path, ok = line[i+len(":name=systemd:"):], true // cgroup v1
			}
		}
		if !ok {
			continue
		}
		parts := strings.Split(strings.TrimSpace(path), "/")
		for _, part := range parts {
			if strings.HasSuffix(part, ".service") || strings.HasSuffix(part, ".scope") {
				return part
			}
		}
	}
	return ""
}
