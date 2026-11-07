package ops

import (
	"strings"
)

func parseServices(units, files string) []Service {
	enabled := map[string]string{}
	for _, line := range strings.Split(files, "\n") {
		f := strings.Fields(line)
		if len(f) >= 2 {
			enabled[f[0]] = f[1]
		}
	}
	var out []Service
	for _, line := range strings.Split(units, "\n") {
		f := strings.Fields(line)
		if len(f) < 4 || !strings.HasSuffix(f[0], ".service") {
			continue
		}
		svc := Service{Unit: f[0], Load: f[1], Active: f[2], Sub: f[3], Enabled: enabled[f[0]]}
		if len(f) > 4 {
			svc.Description = strings.Join(f[4:], " ")
		}
		out = append(out, svc)
	}
	return out
}
