package dialogs

import (
	"sync"

	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/checkbox"
	"github.com/mehran9675/giode/components/dialog"
	"github.com/mehran9675/giode/components/input"

	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// SecretRequest describes what the prompt asks for.
type SecretRequest struct {
	Title    string
	Message  string
	Label    string
	Hidden   bool // mask the answer (passwords)
	Remember bool // offer to save the answer in the keychain
	Retry    bool // the last answer was wrong
}

// SecretAnswer is what the user entered.
type SecretAnswer struct {
	Value    string
	Remember bool
	OK       bool
}

// SecretPrompt asks for a password, passphrase or one-time code while a
// connection waits. Ask blocks the calling (connecting) goroutine until
// the user answers; requests from several connections take turns.
type SecretPrompt struct {
	app  *giode.App
	turn sync.Mutex

	open     bool
	req      SecretRequest
	reply    chan SecretAnswer
	field    *input.Input
	remember *checkbox.Checkbox
	cancel   *giode.ButtonElement
	submit   *giode.ButtonElement
	dialog   *dialog.Dialog
}

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

func (p *SecretPrompt) answer(ok bool) {
	if !p.open || p.reply == nil {
		return
	}
	p.open = false
	p.reply <- SecretAnswer{Value: p.field.Text(), Remember: p.req.Remember && p.remember.Checked(), OK: ok}
	p.reply = nil
	p.field.SetText("")
}

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
