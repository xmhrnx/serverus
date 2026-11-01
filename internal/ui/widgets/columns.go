package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// Columns lays children side by side, sharing the width by weights (nil
// for equal columns).
func Columns(weights []int, children ...giode.Element) giode.Element {
	cells := make([]giode.Element, len(children))
	for i, c := range children {
		w := 1
		if i < len(weights) {
			w = weights[i]
		}
		cells[i] = giode.Box(giode.Styles{FlexGrow: properties.FlexGrow(w), Width: -1}, c)
	}
	return giode.Row(giode.Styles{Gap: theme.Space3, Width: -1, Wrap: giode.WrapNoWrap}, cells...)
}
