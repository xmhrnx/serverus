package dialogs

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (d *ServerDialog) content() giode.Element {
	title, subtitle := "Add a server", "Connect over SSH, Telnet or FTP, or manage this computer"
	if !d.isNew {
		title, subtitle = "Edit server", d.editing.Name
	}
	footer := []giode.Element{d.test, giode.Box(giode.Styles{FlexGrow: 1}), d.cancel, d.save}
	return widgets.DialogFrame("server", title, subtitle, 560,
		giode.Box(giode.Styles{Width: -1, MaxHeight: 540}, d.body), footer...)
}
