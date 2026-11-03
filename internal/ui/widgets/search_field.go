package widgets

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/input"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// SearchField is a filter box with a magnifier icon.
type SearchField struct {
	field *input.Input
}

// NewSearchField returns an empty search field; onChange runs as the
// user types.
func NewSearchField(placeholder string, onChange func(string)) *SearchField {
	f := giode.Input.Text(placeholder, giode.Styles{
		Background: "transparent", BorderWidth: 0, BorderColor: "transparent", Padding: giode.SymmetricInset(0, 8),
		Color: theme.Text, FontSize: theme.TextBase,
	})
	if onChange != nil {
		f.OnChange(onChange)
	}
	return &SearchField{field: f}
}

// Text is what the user typed.
func (s *SearchField) Text() string { return s.field.Text() }

// Clear empties the field.
func (s *SearchField) Clear() { s.field.SetText("") }

// Focus gives the field the keyboard.
func (s *SearchField) Focus() { s.field.Focus() }

// Layout draws the field.
func (s *SearchField) Layout(gtx layout.Context) layout.Dimensions {
	border := theme.Border
	if s.field.Focused() {
		border = theme.Primary
	}
	return giode.Row(giode.Styles{
		Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Background: theme.Sunken, BorderColor: border, BorderWidth: 1, BorderRadius: theme.RadiusMd,
		Padding: giode.SymmetricInset(10, 0),
	},
		theme.Icon("search", 15, theme.Muted),
		giode.Box(giode.Styles{FlexGrow: 1}, s.field),
	).Layout(gtx)
}
