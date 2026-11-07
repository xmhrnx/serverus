package env

import (
	"errors"
	"time"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
)

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
