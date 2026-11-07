package ops

import (
	"encoding/json"
	"strings"
)

func parseContainers(list, stats string) []Container {
	usage := map[string]statsLine{}
	for _, line := range strings.Split(stats, "\n") {
		var s statsLine
		if json.Unmarshal([]byte(line), &s) == nil && s.Name != "" {
			usage[s.Name] = s
		}
	}
	var out []Container
	for _, line := range strings.Split(list, "\n") {
		var p psLine
		if json.Unmarshal([]byte(line), &p) != nil || p.ID == "" {
			continue
		}
		c := Container{
			ID: p.ID, Name: strings.TrimPrefix(p.Names, "/"), Image: p.Image,
			State: p.State, Status: p.Status, Ports: p.Ports, Created: p.CreatedAt,
		}
		if c.State == "" {
			c.State = stateFromStatus(p.Status)
		}
		if u, ok := usage[c.Name]; ok {
			c.CPU, c.Memory = u.CPUPerc, u.MemUsage
		}
		out = append(out, c)
	}
	return out
}
