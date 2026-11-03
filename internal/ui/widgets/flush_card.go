package widgets

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// flushBase is a card without padding, for content that runs to its
// edges: tables, lists, terminals.
var flushBase = giode.Styles{
	FlexDirection: giode.FlexDirectionColumn,
	Background:    theme.Raised,
	BorderColor:   theme.SubtleBorder,
	BorderWidth:   1,
	BorderRadius:  theme.RadiusLg,
	Overflow:      giode.OverflowHidden,
}

// FlushCard is a Card whose content touches its edges.
func FlushCard(st giode.Styles, children ...giode.Element) giode.Element {
	return giode.Box(flushBase.With(st), children...)
}
