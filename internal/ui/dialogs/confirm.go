// Package dialogs holds the app-wide modal dialogs: confirmations,
// credential and host key prompts, and the server editor.
package dialogs

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"

	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Confirm asks the user to confirm an action.
type Confirm struct {
	open   bool
	// pending is set while a shown request has neither run nor been
	// declined; closing the dialog any other way counts as declining.
	pending bool
	req    env.Confirmation
	dialog *dialog.Dialog
	cancel *giode.ButtonElement
	ok     *giode.ButtonElement
	danger *giode.ButtonElement
}

// NewConfirm returns the confirmation dialog.
func NewConfirm() *Confirm {
	c := &Confirm{}
	c.cancel = widgets.Button(widgets.Secondary, "Cancel", "").OnClick(func() { c.open = false })
	run := func() {
		c.open, c.pending = false, false
		if c.req.OnConfirm != nil {
			giode.AfterLayout(c.req.OnConfirm)
		}
	}
	c.ok = widgets.Button(widgets.Primary, "Confirm", "").OnClick(run)
	c.danger = widgets.Button(widgets.Danger, "Delete", "").OnClick(run)
	c.dialog = giode.Dialog(&c.open, c.content, widgets.DialogPanel)
	return c
}

// Show opens the dialog for req.
func (c *Confirm) Show(req env.Confirmation) {
	if req.Action == "" {
		req.Action = "Confirm"
	}
	c.decline()
	c.req, c.pending = req, true
	c.ok.Label(req.Action)
	c.danger.Label(req.Action)
	c.open = true
}

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

// Layout draws the dialog while it's open.
func (c *Confirm) Layout(gtx layout.Context) layout.Dimensions {
	dims := c.dialog.Layout(gtx)
	if !c.open {
		c.decline() // canceled, or dismissed with Escape or a click outside
	}
	return dims
}
