package files

import (
	"path"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

func (p *Page) footer() giode.Element {
	v := p.visible()
	dirs := 0
	var size int64
	for _, f := range v {
		if f.IsDir {
			dirs++
		} else {
			size += f.Size
		}
	}
	note := "Drop files here to upload"
	if _, dragging := p.env.App.FileDrag(); dragging {
		note = "Release to upload into " + path.Base(p.cwd)
	}
	return giode.Row(giode.Styles{Width: -1, Justify: giode.JustifySpaceBetween, Padding: giode.SymmetricInset(theme.Space4, 10), Background: theme.Sunken},
		giode.Text(itemsLabel(len(v)-dirs, dirs, size), theme.Caption),
		giode.Text(note, theme.Caption.With(giode.Styles{Color: theme.Faint})),
	)
}
