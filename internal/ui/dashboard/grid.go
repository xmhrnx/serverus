package dashboard

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

const columns = 12

// span is how many twelfths a card of size takes in a row width dp
// wide: four small cards across a wide window, two on a medium one, and
// one per row on a phone.
func span(size domain.CardSize, width int) int {
	s := size.Span()
	switch {
	case width >= 920:
		return s
	case width >= 520:
		if size == domain.CardSmall {
			return 6
		}
		return columns
	}
	return columns
}

type cell struct {
	span int
	el   giode.Element
}

// grid lays cells in rows of twelve columns, a cell starting a new row
// when it doesn't fit the current one. Cells in a row are as tall as
// the tallest.
func grid(gap int, cells func(width int) []cell) giode.Element {
	return giode.Raw(func(gtx layout.Context) layout.Dimensions {
		width := int(gtx.Metric.PxToDp(gtx.Constraints.Max.X))
		all := cells(width)
		unit := float32(width-(columns-1)*gap) / columns
		var rows []giode.Element
		var row []giode.Element
		used := 0
		flush := func() {
			if len(row) > 0 {
				rows = append(rows, giode.Row(giode.Styles{Width: -1, Gap: properties.Gap(gap), Wrap: giode.WrapNoWrap, Align: giode.AlignStretch}, row...))
			}
			row, used = nil, 0
		}
		for _, c := range all {
			if used+c.span > columns {
				flush()
			}
			w := unit*float32(c.span) + float32((c.span-1)*gap)
			row = append(row, giode.Box(giode.Styles{Width: properties.Width(w), FlexDirection: giode.FlexDirectionColumn}, c.el))
			used += c.span
		}
		flush()
		return giode.Stack(giode.Styles{Width: -1, Gap: properties.Gap(gap)}, rows...).Layout(gtx)
	})
}

// Grid is the board's cards, laid out to the width it's given.
func (b *Board) Grid() giode.Element {
	return grid(theme.Space4, func(width int) []cell {
		var out []cell
		for _, c := range b.Cards() {
			spec := b.byKind[c.Kind]
			var el giode.Element
			switch {
			case spec.available():
				el = spec.Render()
			case b.editing:
				el = unavailable(spec)
			default:
				continue
			}
			if b.editing {
				el = b.editFrame(spec, el)
			}
			out = append(out, cell{span: span(c.Size, width), el: el})
		}
		if len(out) == 0 {
			out = append(out, cell{span: columns, el: b.empty()})
		}
		return out
	})
}

// empty fills a board with no cards.
func (b *Board) empty() giode.Element {
	return giode.Box(giode.Styles{Width: -1, Padding: giode.SymmetricInset(0, theme.Space8), Align: giode.AlignCenter, FlexDirection: giode.FlexDirectionColumn, Gap: theme.Space3},
		theme.Icon("layout-dashboard", 28, theme.Faint),
		giode.Text("No cards on this dashboard", theme.Heading),
		b.add,
	)
}

// Controls are the board's actions for the page header: Customize, or
// while editing, Add, Reset and Done.
func (b *Board) Controls() giode.Element {
	compact := giode.Compact()
	switch {
	case !b.editing && compact:
		return b.editTool
	case !b.editing:
		return b.edit
	case compact:
		return giode.Row(giode.Styles{Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, b.addTool, b.resetTool, b.doneTool)
	}
	return giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, b.reset, b.add, b.done)
}
