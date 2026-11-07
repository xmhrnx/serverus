package files

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// card is the listing's content: column headers (not on phones), the
// rows and the summary.
func (p *Page) card() []giode.Element {
	body := giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1, Padding: giode.UniformInset(6)}, p.body())
	if giode.Compact() {
		return []giode.Element{body, p.footer()}
	}
	return []giode.Element{p.header(), giode.Divider(giode.Styles{Color: theme.SubtleBorder}), body, p.footer()}
}

func (p *Page) header() giode.Element {
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Wrap: giode.WrapNoWrap, Align: giode.AlignCenter, Padding: giode.Inset{Left: 22 + 26, Right: 22, Top: 10, Bottom: 10}},
		p.headers[0].Cell(), p.headers[1].Cell(), p.headers[2].Cell(),
		giode.Box(giode.Styles{Width: modeWidth}, giode.Text("Permissions", theme.Overline)),
	)
}
