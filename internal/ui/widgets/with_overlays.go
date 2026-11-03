package widgets

import "github.com/xmhrnx/giode"

// WithOverlays places dialogs, menus and other overlays next to content.
// They take no space there and paint above the window when open, so
// content keeps its layout.
func WithOverlays(content giode.Element, overlays ...giode.Element) giode.Element {
	return giode.Box(giode.Styles{}, append([]giode.Element{content}, overlays...)...)
}
