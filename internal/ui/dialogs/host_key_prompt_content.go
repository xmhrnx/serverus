package dialogs

import (
	"fmt"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (h *HostKeyPrompt) content() giode.Element {
	title, msg, btn := "Verify the server", fmt.Sprintf(
		"This is the first connection to %s. Check that the fingerprint below matches the server's before trusting it.", h.server.Address()), h.trust
	var banner giode.Element
	if h.changed {
		title, btn = "The host key changed", h.danger
		msg = fmt.Sprintf("%s presents a different key than last time.", h.server.Address())
		banner = widgets.Banner("triangle-alert", theme.Danger,
			"Someone could be intercepting the connection. Only continue if you know the server was reinstalled or its keys were replaced.")
	}
	body := []giode.Element{giode.Text(msg, theme.Body.With(giode.Styles{Color: theme.Muted, LineHeight: 20}))}
	if banner != nil {
		body = append(body, banner)
	}
	body = append(body, giode.Box(giode.Styles{
		Width: -1, Background: theme.Sunken, BorderRadius: theme.RadiusMd, Padding: giode.UniformInset(theme.Space3),
		BorderColor: theme.Border, BorderWidth: 1,
	}, giode.Text(h.fingerprint, theme.Code).Selectable("hostkey")))
	return widgets.DialogFrame("shield-check", title, h.server.Name, 480,
		giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1}, body...), h.reject, btn)
}

// Layout draws the prompt while it's open.
func (h *HostKeyPrompt) Layout(gtx layout.Context) layout.Dimensions { return h.dialog.Layout(gtx) }
