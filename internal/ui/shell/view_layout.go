package shell

import (
	"slices"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Layout draws the workspace.
func (v *ServerView) Layout(gtx layout.Context) layout.Dimensions {
	specs := availablePages(v.sess.Server().Protocol)
	v.nav.SetPages(specs)
	if !slices.ContainsFunc(specs, func(s pageSpec) bool { return s.kind == v.active }) && len(specs) > 0 {
		v.active = specs[0].kind
	}
	var body giode.Element = v.panel
	if st, _ := v.sess.State(); st == session.Connected {
		body = v.page(v.active)
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1},
		v.topbar,
		v.nav,
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, body),
	), v.power.Overlays()...).Layout(gtx)
}
