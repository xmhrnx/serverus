package dialogs

import (
	"fmt"
	"sync"

	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/dialog"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// HostKeyPrompt asks whether to trust a server's host key the first
// time the app connects, and warns loudly when a known key changed.
type HostKeyPrompt struct {
	app  *giode.App
	turn sync.Mutex

	open        bool
	server      domain.Server
	fingerprint string
	changed     bool
	reply       chan bool

	reject *giode.ButtonElement
	trust  *giode.ButtonElement
	danger *giode.ButtonElement
	dialog *dialog.Dialog
}

// NewHostKeyPrompt returns the prompt.
func NewHostKeyPrompt(app *giode.App) *HostKeyPrompt {
	h := &HostKeyPrompt{app: app}
	h.reject = widgets.Button(widgets.Secondary, "Don't connect", "").OnClick(func() { h.answer(false) })
	h.trust = widgets.Button(widgets.Primary, "Trust and connect", "shield-check").OnClick(func() { h.answer(true) })
	h.danger = widgets.Button(widgets.Danger, "Trust the new key", "").OnClick(func() { h.answer(true) })
	h.dialog = giode.Dialog(&h.open, h.content, widgets.DialogPanel).Dismissible(false)
	return h
}

// TrustHostKey implements transport.HostKeyPrompt. It blocks until the
// user decides.
func (h *HostKeyPrompt) TrustHostKey(s domain.Server, fingerprint string, changed bool) bool {
	h.turn.Lock()
	defer h.turn.Unlock()
	reply := make(chan bool, 1)
	h.app.Post(func() {
		h.server, h.fingerprint, h.changed, h.reply = s, fingerprint, changed, reply
		h.open = true
	})
	return <-reply
}

func (h *HostKeyPrompt) answer(ok bool) {
	if h.reply == nil {
		return
	}
	h.open = false
	h.reply <- ok
	h.reply = nil
}

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
