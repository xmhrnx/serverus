package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
)

// tools are the icon actions beside the connect button: power, while
// connected to a server that runs commands, and edit.
func (t *TopBar) tools(st session.State) giode.Element {
	items := []giode.Element{}
	if st == session.Connected && t.sess.Capabilities().Has(domain.CapExec) {
		items = append(items, t.power.Trigger())
	}
	items = append(items, t.edit)
	return giode.Row(giode.Styles{Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, items...)
}
