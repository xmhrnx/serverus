// Package env holds the services every view needs, created once at
// startup and passed down: the app window, storage, sessions, transfers,
// tunnels, notifications and confirmations.
package env

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/storage"
	"github.com/xmhrnx/serverus/internal/transfer"
	"github.com/xmhrnx/serverus/internal/tunnel"
)

// Env is the dependency container of the UI.
type Env struct {
	App       *giode.App
	Repo      *storage.Repository
	Secrets   storage.Secrets
	Sessions  *session.Manager
	Tunnels   *tunnel.Manager
	Transfers *transfer.Queue
	Toasts    *giode.ToasterElement

	// confirm opens the shared confirmation dialog; set by the shell.
	confirm func(Confirmation)
}

// Confirmation describes a question to the user before an action.
type Confirmation struct {
	Title   string
	Message string
	Action  string // the confirming button's label
	// Destructive styles the button as dangerous; it's skipped when
	// the user turned confirmations off.
	Destructive bool
	OnConfirm   func()
	// OnCancel, when set, runs if the user declines or dismisses it.
	OnCancel func()
}

// SetConfirmer installs the function showing confirmations.
func (e *Env) SetConfirmer(fn func(Confirmation)) { e.confirm = fn }

// Confirm asks before running c.OnConfirm.
func (e *Env) Confirm(c Confirmation) {
	if c.Destructive && !e.Repo.Settings().ConfirmDestructive {
		c.OnConfirm()
		return
	}
	if e.confirm == nil {
		c.OnConfirm()
		return
	}
	e.confirm(c)
}
