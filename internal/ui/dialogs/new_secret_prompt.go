package dialogs

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// NewSecretPrompt returns the prompt.
func NewSecretPrompt(app *giode.App) *SecretPrompt {
	p := &SecretPrompt{app: app}
	p.field = widgets.TextField("").Submit(func(string) { p.answer(true) })
	p.remember = giode.Input.Checkbox("Remember in the system keychain", giode.Styles{Color: theme.Text, Background: theme.Primary, BorderColor: theme.Border, FontSize: theme.TextSm})
	p.cancel = widgets.Button(widgets.Secondary, "Cancel", "").OnClick(func() { p.answer(false) })
	p.submit = widgets.Button(widgets.Primary, "Continue", "").OnClick(func() { p.answer(true) })
	p.dialog = giode.Dialog(&p.open, p.content, widgets.DialogPanel).Dismissible(false)
	return p
}

// Ask shows req and waits for the answer. Call it off the UI goroutine.
func (p *SecretPrompt) Ask(req SecretRequest) SecretAnswer {
	p.turn.Lock()
	defer p.turn.Unlock()
	reply := make(chan SecretAnswer, 1)
	p.app.Post(func() {
		p.req, p.reply = req, reply
		p.field.SetText("")
		p.field.Password(req.Hidden)
		p.field.Placeholder(req.Label)
		p.remember.SetChecked(true)
		p.open = true
		p.field.Focus()
	})
	return <-reply
}
