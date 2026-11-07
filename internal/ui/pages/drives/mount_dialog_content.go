package drives

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (d *MountDialog) content() giode.Element {
	rows := []giode.Element{
		widgets.FormField("Folder on the server", "Created if it doesn't exist", d.folder),
		d.overlay,
	}
	if d.err != "" {
		rows = append(rows, widgets.Banner("circle-alert", theme.Danger, d.err))
	}
	footer := []giode.Element{d.cancel, d.ok}
	if d.working {
		footer = []giode.Element{giode.Spinner(giode.Styles{Width: 20, Height: 20, Color: theme.Primary})}
	}
	body := giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1}, rows...)
	return widgets.DialogFrame("folder-input", "Mount "+d.volume.Title(), d.volume.Path, 460, body, footer...)
}
