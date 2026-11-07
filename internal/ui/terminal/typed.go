package terminal

// typed returns to the live screen and clears the selection, as the
// user typing means they're done looking back.
func (v *View) typed() {
	v.offset = 0
	v.sel = selection{}
}
