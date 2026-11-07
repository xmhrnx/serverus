package drives

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// DriveCard shows one drive and its volumes, with the eject action for
// drives that can be unplugged.
type DriveCard struct {
	page      *Page
	drive     ops.Drive
	eject     *giode.ButtonElement
	ejectIcon *widgets.ToolButton
}

// NewDriveCard returns a card for the page's drives.
func NewDriveCard(p *Page) *DriveCard {
	c := &DriveCard{page: p}
	c.eject = widgets.Button(widgets.Secondary, "Eject", "power-off").OnClick(func() { p.Eject(c.drive) })
	c.ejectIcon = widgets.NewToolButton("power-off", "Eject", func() { p.Eject(c.drive) })
	return c
}

// Set shows d.
func (c *DriveCard) Set(d ops.Drive) { c.drive = d }
