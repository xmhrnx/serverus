package files

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (p *Page) body() giode.Element {
	if err := p.data.Err(); err != nil {
		if _, ok := p.data.Get(); !ok || !p.data.Loading() {
			return widgets.EmptyState("circle-alert", "Couldn't open the folder", err.Error())
		}
	}
	l, ok := p.data.Get()
	if !ok || (p.data.Loading() && l.dir != p.cwd && p.cwd != "") {
		return widgets.Loading("Reading the folder…")
	}
	if len(p.visible()) == 0 {
		if p.filter.Text() != "" {
			return widgets.EmptyState("search", "No matches", "Nothing in this folder matches the filter.")
		}
		return widgets.EmptyState("folder-open", "This folder is empty", "Drop files on the window or use Upload to add some.", p.upload)
	}
	if giode.Compact() {
		return p.phoneList
	}
	return p.list
}
