package dialogs

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (p *SecretPrompt) content() giode.Element {
	body := []giode.Element{}
	if p.req.Retry {
		body = append(body, widgets.Banner("circle-alert", theme.Danger, "That didn't work. Try again."))
	}
	body = append(body, widgets.FormField(p.req.Label, "", p.field))
	if p.req.Remember {
		body = append(body, p.remember)
	}
	return widgets.DialogFrame("key-round", p.req.Title, p.req.Message, 420,
		giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1}, body...),
		p.cancel, p.submit)
}

// Layout draws the prompt while it's open.
func (p *SecretPrompt) Layout(gtx layout.Context) layout.Dimensions { return p.dialog.Layout(gtx) }
