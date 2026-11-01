package widgets

import (
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/xmhrnx/bastion/internal/ui/theme"
)

// SortHeader is a sortable column title: clicking sorts by it, again
// reverses the order.
type SortHeader struct {
	label   string
	col     int
	width   int // 0 grows
	current func() (int, bool)
	press   *giode.PressableElement
}

// NewSortHeader returns the header; onSort runs with its column key.
func NewSortHeader(label string, col int, width int, current func() (int, bool), onSort func(int)) *SortHeader {
	h := &SortHeader{label: label, col: col, width: width, current: current}
	h.press = giode.Pressable(func(ps giode.PressState) giode.Element {
		by, desc := h.current()
		fg := theme.Faint
		if ps.Hovered || by == h.col {
			fg = theme.Muted
		}
		row := []giode.Element{giode.Text(h.label, theme.Overline.With(giode.Styles{Color: fg}))}
		if by == h.col {
			icon := "chevron-up"
			if desc {
				icon = "chevron-down"
			}
			row = append(row, theme.Icon(icon, 12, fg))
		}
		return giode.Row(giode.Styles{Gap: 4, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, row...)
	}).Label("Sort by " + label).OnClick(func(giode.ClickEvent) { onSort(h.col) })
	return h
}

// Cell returns the header sized for its column, for a header row.
func (h *SortHeader) Cell() giode.Element {
	st := giode.Styles{Width: properties.Width(h.width)}
	if h.width == 0 {
		st = giode.Styles{FlexGrow: 1}
	}
	return giode.Box(st, h.press)
}
