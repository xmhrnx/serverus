package theme

import "github.com/mehran9675/giode"

// Text styles for the recurring roles. Merge overrides with With.
var (
	// Title is a page or dialog heading.
	Title = giode.Styles{Color: Text, FontSize: TextXL, FontWeight: giode.FontWeightSemiBold}
	// Heading is a card or section heading.
	Heading = giode.Styles{Color: Text, FontSize: TextMd, FontWeight: giode.FontWeightSemiBold}
	// Body is regular text.
	Body = giode.Styles{Color: Text, FontSize: TextBase}
	// Caption is secondary text under or beside a label.
	Caption = giode.Styles{Color: Muted, FontSize: TextSm}
	// Overline labels a group: small, uppercase, spaced.
	Overline = giode.Styles{Color: Faint, FontSize: TextXS, FontWeight: giode.FontWeightSemiBold, TextTransform: giode.TextTransformUppercase}
	// Code is monospaced text.
	Code = giode.Styles{Color: Text, FontSize: TextSm, FontFamily: Mono}
	// Metric is a large number on a dashboard tile.
	Metric = giode.Styles{Color: Text, FontSize: 24, FontWeight: giode.FontWeightSemiBold, MaxLines: 1}
)
