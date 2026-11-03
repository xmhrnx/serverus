package widgets

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// cardBase is the look every card shares: a raised surface with a hairline
// border and soft corners.
var cardBase = giode.Styles{
	FlexDirection: giode.FlexDirectionColumn,
	Background:    theme.Raised,
	BorderColor:   theme.SubtleBorder,
	BorderWidth:   1,
	BorderRadius:  theme.RadiusLg,
	Padding:       giode.UniformInset(theme.Space5),
	Gap:           theme.Space4,
}

// Card is a surface grouping related content. st overrides the base
// look (padding, width, FlexGrow).
func Card(st giode.Styles, children ...giode.Element) giode.Element {
	return giode.Box(cardBase.With(st), children...)
}
