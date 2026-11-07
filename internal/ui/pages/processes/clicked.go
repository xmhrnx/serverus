package processes

import (
	"github.com/xmhrnx/serverus/internal/ops"
)

// clicked returns the process that was in row i when the user clicked
// it: last frame's row, before this frame's re-sort.
func (p *Page) clicked(i int) (ops.Process, bool) {
	if i < len(p.seen) {
		return p.seen[i], true
	}
	if i < len(p.shown) {
		return p.shown[i], true
	}
	return ops.Process{}, false
}

func (p *Page) lookup(pid int32) (ops.Process, bool) {
	all, _ := p.data.Get()
	for _, pr := range all {
		if pr.PID == pid {
			return pr, true
		}
	}
	return ops.Process{}, false
}
