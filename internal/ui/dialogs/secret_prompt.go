package dialogs

import (
	"sync"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/checkbox"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"
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
