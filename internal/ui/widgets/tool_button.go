package widgets

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/iconbutton"
	"github.com/xmhrnx/giode/components/tooltip"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// ToolButton is a compact icon button with a tooltip, for toolbars.
type ToolButton struct {
	btn  *iconbutton.Button
	tip  *tooltip.Tooltip
	icon string
	tint giode.Color
}

// NewToolButton returns a tool button running onClick.
func NewToolButton(icon, tooltipText string, onClick func()) *ToolButton {
	t := &ToolButton{icon: icon, tint: theme.Muted}
	t.btn = giode.IconButton(theme.Icon(icon, 16, theme.Muted), giode.Styles{
		Background: theme.Text, Padding: giode.UniformInset(7), BorderRadius: theme.RadiusSm,
	}).Variant(giode.IconButtonGhost).OnClick(onClick).Label(tooltipText)
	t.tip = giode.Tooltip(t.btn, tooltipText, giode.Styles{FontSize: theme.TextSm})
	return t
}

// Tint changes the icon color (theme.Danger for a destructive tool).
func (t *ToolButton) Tint(c giode.Color) *ToolButton {
	t.tint = c
	t.btn.Icon(theme.Icon(t.icon, 16, c))
	return t
}

// SetIcon swaps the icon (a play button turning into a stop button).
func (t *ToolButton) SetIcon(icon string) {
	if icon != t.icon {
		t.icon = icon
		t.btn.Icon(theme.Icon(icon, 16, t.tint))
	}
}

// SetTooltip replaces the tooltip text.
func (t *ToolButton) SetTooltip(s string) { t.tip.SetContent(s) }

// Click runs the button's action.
func (t *ToolButton) Click() { t.btn.Click() }

// Layout draws the button.
func (t *ToolButton) Layout(gtx layout.Context) layout.Dimensions { return t.tip.Layout(gtx) }
