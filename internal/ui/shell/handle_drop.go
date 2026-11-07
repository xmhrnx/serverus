package shell

import (
	"github.com/xmhrnx/serverus/internal/session"
)

// HandleDrop uploads dropped files when the file manager is showing.
func (v *ServerView) HandleDrop(paths []string) bool {
	if v.active != PageFiles || v.files == nil {
		return false
	}
	if st, _ := v.sess.State(); st != session.Connected {
		return false
	}
	v.files.HandleDrop(paths)
	return true
}

// Close ends the view's shells, when the server is deleted.
func (v *ServerView) Close() {
	if v.terminal != nil {
		v.terminal.Close()
	}
}
