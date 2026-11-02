package widgets

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"github.com/mehran9675/giode/components/menu"
	"github.com/mehran9675/giode/elements"
)

// ContextMenu opens a menu at a window position. Menu.OpenAt takes
// window coordinates, but the panel it registers is then moved again by
// the containers the menu is laid out in once they place it, so a menu
// laid out anywhere but the window's origin opens that far off (see
// docs/GIODE_REQUIREMENTS.md). ContextMenu measures that shift and
// subtracts it.
type ContextMenu struct {
	*menu.Menu
	laidOutAt image.Point // where layout started, before placement
	placedAt  image.Point // where the containers finally put it
}

// NewContextMenu wraps m.
func NewContextMenu(m *menu.Menu) *ContextMenu { return &ContextMenu{Menu: m} }

// OpenAt opens the menu with its corner at the window position pos.
func (c *ContextMenu) OpenAt(pos f32.Point) {
	shift := c.placedAt.Sub(c.laidOutAt)
	c.Menu.OpenAt(pos.Sub(f32.Pt(float32(shift.X), float32(shift.Y))))
}

// Layout lays the menu out, recording how far its containers move it.
func (c *ContextMenu) Layout(gtx layout.Context) layout.Dimensions {
	c.laidOutAt = elements.AbsolutePosition()
	dims := c.Menu.Layout(gtx)
	elements.TrackPosition(&c.placedAt)
	return dims
}
