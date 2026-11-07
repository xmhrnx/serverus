package files

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

// Offer asks whether to do what as root when err is a permission error,
// and runs root then, followed by done with its result. It reports
// whether it took over; when it didn't, the caller shows err itself.
func (s *Sudo) Offer(what string, err error, root func(context.Context, ops.RootFiles) error, done func(error)) bool {
	if root == nil || !ops.IsPermissionError(err) {
		return false
	}
	rf, ok := s.files()
	if !ok {
		return false
	}
	s.env.Confirm(env.Confirmation{
		Title:    "Permission denied",
		Message:  "You aren't allowed to " + what + " as the user you signed in with. Do it as root with sudo? You may be asked for the sudo password.",
		Action:   "Use sudo",
		OnCancel: func() { done(err) },
		OnConfirm: func() {
			s.env.Async(func() error {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				defer cancel()
				return root(ctx, rf)
			}, done)
		},
	})
	return true
}
