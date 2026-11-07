package env

import (
	"github.com/xmhrnx/giode"
)

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
