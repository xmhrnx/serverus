package drives

import (
	"github.com/xmhrnx/serverus/internal/ops"
)

// visible is the drives the filter shows.
func (p *Page) visible() []ops.Drive {
	all, _ := p.data.Get()
	if p.show.Selected() == 1 {
		return all
	}
	var out []ops.Drive
	for _, d := range all {
		if d.External() {
			out = append(out, d)
		}
	}
	return out
}
