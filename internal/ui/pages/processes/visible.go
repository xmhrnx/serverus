package processes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/xmhrnx/serverus/internal/ops"
)

func (p *Page) visible() []ops.Process {
	all, _ := p.data.Get()
	q := strings.ToLower(strings.TrimSpace(p.filter.Text()))
	out := make([]ops.Process, 0, len(all))
	for _, pr := range all {
		if q != "" && !strings.Contains(strings.ToLower(pr.Name+" "+pr.User+" "+pr.Command+" "+fmt.Sprint(pr.PID)), q) {
			continue
		}
		out = append(out, pr)
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		var less bool
		switch p.sortBy {
		case byPID:
			less = a.PID < b.PID
		case byUser:
			less = a.User < b.User
		case byMemory:
			less = a.Memory < b.Memory
		case byTime:
			less = a.Elapsed < b.Elapsed
		case byName:
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		default:
			less = a.CPU < b.CPU
		}
		if p.desc {
			return !less
		}
		return less
	})
	return out
}
