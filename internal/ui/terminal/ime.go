package terminal

import (
	"strings"
	"unicode/utf8"

	"gioui.org/io/key"
	"gioui.org/layout"
	uv "github.com/charmbracelet/ultraviolet"
)

// imeSentinel is the text the terminal tells input methods it holds,
// with the caret at its end. A terminal has no editable text of its own,
// but on-screen keyboards edit the text before the caret: a backspace
// deletes a character of it, and a composing keyboard replaces the word
// it's building as the user types. Reading those edits against the
// sentinel turns them back into keystrokes: each deleted character is a
// Backspace, each inserted one is typed. The sentinel is put back after
// every edit.
const imeSentinel = "                "

var imeLen = utf8.RuneCountInString(imeSentinel)

// resetIME tells the input method the terminal holds the sentinel, with
// the caret at its end.
func (v *View) resetIME(gtx layout.Context) {
	gtx.Execute(key.SnippetCmd{Tag: v, Snippet: key.Snippet{Range: key.Range{Start: 0, End: imeLen}, Text: imeSentinel}})
	gtx.Execute(key.SelectionCmd{Tag: v, Range: key.Range{Start: imeLen, End: imeLen}})
}

// handleEdit turns an input method's edit of the sentinel into
// keystrokes.
func (v *View) handleEdit(gtx layout.Context, e key.EditEvent) {
	v.typed()
	for i := 0; i < e.Range.End-e.Range.Start; i++ {
		v.screen.SendKey(uv.KeyPressEvent{Code: uv.KeyBackspace})
	}
	text := strings.ReplaceAll(e.Text, "\n", "\r")
	if v.ctrlLatch && utf8.RuneCountInString(text) == 1 {
		// The key bar's Ctrl applies to the next character typed.
		r, _ := utf8.DecodeRuneInString(strings.ToLower(text))
		v.ctrlLatch = false
		v.screen.SendKey(uv.KeyPressEvent{Code: r, Mod: uv.ModCtrl})
	} else if text != "" {
		v.screen.SendText(text)
	}
	v.resetIME(gtx)
}
