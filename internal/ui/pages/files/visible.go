package files

import (
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
)

// visible returns the entries shown: filtered, without dotfiles unless
// asked, sorted.
func (p *Page) visible() []transport.FileInfo {
	l, ok := p.data.Get()
	if !ok {
		return nil
	}
	hidden := p.env.Repo.Settings().ShowHiddenFiles
	q := strings.ToLower(p.filter.Text())
	out := make([]transport.FileInfo, 0, len(l.entries))
	for _, e := range l.entries {
		if !hidden && strings.HasPrefix(e.Name, ".") {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(e.Name), q) {
			continue
		}
		out = append(out, e)
	}
	sortEntries(out, p.sortBy, p.desc)
	return out
}
