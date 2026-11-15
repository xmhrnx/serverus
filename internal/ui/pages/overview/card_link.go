package overview

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// cardLink is the "open the full page" action in a card's header.
func cardLink(label string, open func()) *giode.ButtonElement {
	return widgets.Button(widgets.Ghost, label, "chevron-right").OnClick(open)
}
