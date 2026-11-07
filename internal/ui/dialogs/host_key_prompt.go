package dialogs

import (
	"sync"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"

	"github.com/xmhrnx/serverus/internal/domain"
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
