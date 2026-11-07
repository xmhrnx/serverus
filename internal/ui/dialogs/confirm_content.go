package dialogs

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (c *Confirm) content() giode.Element {
	icon, btn := "circle-alert", c.ok
	if c.req.Destructive {
		icon, btn = "triangle-alert", c.danger
	}
	body := giode.Text(c.req.Message, theme.Body.With(giode.Styles{Color: theme.Muted, LineHeight: 20}))
	return widgets.DialogFrame(icon, c.req.Title, "", 420, body, c.cancel, btn)
}

// decline runs the pending request's OnCancel.
func (c *Confirm) decline() {
	if c.pending {
		c.pending = false
		if c.req.OnCancel != nil {
			giode.AfterLayout(c.req.OnCancel)
		}
	}
}
