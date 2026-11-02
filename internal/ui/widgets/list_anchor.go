package widgets

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/elements"
)

// ListAnchor finds where a VirtualList row is in the window, for opening
// a context menu at a click in it. It tracks one position, the list's,
// and works a row's out from its index, the row height and the scroll
// offset, rather than tracking every row as it's built.
type ListAnchor struct {
	list      func() *giode.VirtualListElement
	rowHeight float32 // dp
	origin    image.Point
}

// NewListAnchor returns an anchor for the list list returns, whose rows
// are rowHeight dp tall.
func NewListAnchor(list func() *giode.VirtualListElement, rowHeight float32) *ListAnchor {
	return &ListAnchor{list: list, rowHeight: rowHeight}
}

// Wrap lays the list out, tracking where it is.
func (a *ListAnchor) Wrap(list giode.Element) giode.Element {
	return giode.Raw(func(gtx layout.Context) layout.Dimensions {
		// Track after laying out: click handlers run inside Layout and
		// need the exact position from the last frame.
		dims := list.Layout(gtx)
		elements.TrackPosition(&a.origin)
		return dims
	})
}

// At is the window position of a click e in row i.
func (a *ListAnchor) At(i int, e giode.ClickEvent) f32.Point {
	top := float32(i)*a.rowHeight - a.list().OffsetDp()
	return f32.Pt(float32(a.origin.X+e.Position.X), float32(a.origin.Y+giode.DpToPx(top)+e.Position.Y))
}
