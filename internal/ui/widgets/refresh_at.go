package widgets

import (
	"time"

	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode"
)

// RefreshAt schedules a redraw at t, for views that change with time
// (polling lists, "5 min ago" labels). It takes no space.
func RefreshAt(t time.Time) giode.Element {
	return giode.Raw(func(gtx layout.Context) layout.Dimensions {
		gtx.Execute(op.InvalidateCmd{At: t})
		return layout.Dimensions{}
	})
}
