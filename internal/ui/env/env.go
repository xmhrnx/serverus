// Package env holds the services every view needs, created once at
// startup and passed down: the app window, storage, sessions, transfers,
// tunnels, notifications and confirmations.
package env

import (
	"errors"
	"time"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
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

// Async runs work off the UI goroutine and then done on it with work's
// error. done may be nil.
func (e *Env) Async(work func() error, done func(err error)) {
	go func() {
		err := work()
		if done != nil {
			e.App.Post(func() { done(err) })
		}
	}()
}

// Invalidate requests a redraw; safe from any goroutine.
func (e *Env) Invalidate() { e.App.Invalidate() }

// Success shows a confirmation toast.
func (e *Env) Success(msg string) { e.Toasts.Show(giode.ToastSuccess, msg) }

// Info shows an informational toast.
func (e *Env) Info(msg string) { e.Toasts.Show(giode.ToastInfo, msg) }

// Fail shows an error toast prefixed with what failed.
func (e *Env) Fail(what string, err error) {
	if err == nil || errors.Is(err, ops.ErrCanceled) {
		return
	}
	e.Toasts.Show(giode.ToastError, what+": "+err.Error())
}

// MetricsInterval is the configured sampling interval.
func (e *Env) MetricsInterval() time.Duration {
	return time.Duration(e.Repo.Settings().RefreshSeconds) * time.Second
}
