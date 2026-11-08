package dialogs

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// sheetFooter is the sheet's buttons, pinned below the form: cancel
// and save side by side, the connection test under them.
func (d *ServerDialog) sheetFooter() giode.Element {
	half := giode.Styles{FlexGrow: 1, Width: -1}
	return giode.Stack(giode.Styles{Width: -1, Gap: theme.Space2, Padding: giode.Inset{Top: theme.Space3, Bottom: theme.Space3}},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Wrap: giode.WrapNoWrap},
			giode.Box(half, d.sheetCancel), giode.Box(half, d.sheetSave)),
		giode.Box(giode.Styles{Width: -1, Justify: giode.JustifyCenter}, d.test),
	)
}
