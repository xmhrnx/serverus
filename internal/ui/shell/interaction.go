package shell

import (
	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
)

// ShowNth switches to the n-th page in the tab row, counting from 1.
func (v *ServerView) ShowNth(n int) {
	if specs := availablePages(v.sess.Server().Protocol); n >= 1 && n <= len(specs) {
		v.active = specs[n-1].kind
	}
}

// NewTerminalTab switches to the terminal and opens a tab, if the
// server has terminals and is connected.
func (v *ServerView) NewTerminalTab() {
	if !v.sess.Capabilities().Has(domain.CapTerminal) {
		return
	}
	if st, _ := v.sess.State(); st != session.Connected {
		return
	}
	v.active = PageTerminal
	v.page(PageTerminal)
	v.terminal.NewTab()
}

// openFolder shows dir in the file manager.
func (v *ServerView) openFolder(dir string) {
	v.active = PageFiles
	v.page(PageFiles)
	v.files.Navigate(dir)
}

// typeInTerminal switches to the terminal and types text into it.
func (v *ServerView) typeInTerminal(text string) {
	v.active = PageTerminal
	v.page(PageTerminal)
	v.terminal.SendText(text)
}
