package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/elements"
)

// WithOverlays places dialogs, menus and other overlays over content
// without letting them take part in its layout. Dialogs, drawers and
// sheets already take no space; anything else (a context menu, which
// fills its parent while open to catch clicks outside it) is laid out
// as an absolutely positioned layer covering content.
func WithOverlays(content giode.Element, overlays ...giode.Element) giode.Element {
	children := []giode.Element{content}
	for _, o := range overlays {
		if _, ok := o.(elements.OutOfFlow); ok {
			children = append(children, o)
			continue
		}
		children = append(children, giode.Box(giode.Styles{
			Position: giode.PositionAbsolute,
			Top:      giode.Flush, Left: giode.Flush, Right: giode.Flush, Bottom: giode.Flush,
		}, o))
	}
	return giode.Box(giode.Styles{}, children...)
}
