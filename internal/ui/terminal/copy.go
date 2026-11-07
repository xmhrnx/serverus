package terminal

// Copy puts the selection on the clipboard; it reports whether there
// was one.
func (v *View) Copy() bool {
	if v.sel.empty() {
		return false
	}
	a, b := v.sel.ordered()
	v.clipboard.SetClipboard(v.screen.text(a, b))
	return true
}

// Paste types the clipboard's text into the terminal.
func (v *View) Paste() {
	v.clipboard.ReadClipboard(func(t string) {
		if t != "" {
			v.offset = 0
			v.screen.Paste(t)
		}
	})
}
