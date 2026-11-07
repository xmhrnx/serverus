package terminal

// ended reports whether the tab's shell is over.
func (t *tab) ended() bool { return t.shell != nil && t.shell.Ended() }

// close ends the shell.
func (t *tab) close() {
	if t.shell != nil {
		t.shell.Close()
	}
}
