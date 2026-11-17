package widgets

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// CardHeader is the title row of a card: an icon badge, a title with an
// optional caption, and trailing elements (a value, legends, actions) on
// the right. On a narrow card the trailing elements move to a line of
// their own under the title, so neither squeezes the other.
func CardHeader(icon string, accent giode.Color, title, caption string, trailing ...giode.Element) giode.Element {
	oneLine := giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}
	text := []giode.Element{giode.Text(title, theme.Heading.With(oneLine))}
	if caption != "" {
		text = append(text, giode.Text(caption, theme.Caption.With(oneLine)))
	}
	head := []giode.Element{}
	if icon != "" {
		head = append(head, IconBadge(icon, accent, 32))
	}
	head = append(head, giode.Stack(giode.Styles{Gap: 2, FlexGrow: 1}, text...))
	row := giode.Styles{Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Width: -1}
	if len(trailing) == 0 {
		return giode.Row(row, head...)
	}
	return giode.Raw(func(gtx layout.Context) layout.Dimensions {
		width := int(gtx.Metric.PxToDp(gtx.Constraints.Max.X))
		if width >= stackBelow(len(trailing)) {
			return giode.Row(row, append(head, trailing...)...).Layout(gtx)
		}
		return giode.Stack(giode.Styles{Width: -1, Gap: theme.Space2},
			giode.Row(row, head...),
			giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter}, trailing...),
		).Layout(gtx)
	})
}

// stackBelow is the width, in dp, under which n trailing elements go on
// their own line: room for the title and about 120dp for each.
func stackBelow(n int) int { return 200 + 120*n }
