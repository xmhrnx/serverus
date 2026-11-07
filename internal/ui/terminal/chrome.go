package terminal

import (
	"image/color"
	"sync"
)

// chrome is the background, cursor and selection the app's theme asks
// for; screens pick up a change the next time they're drawn.
var chrome struct {
	sync.Mutex
	version                       int
	background, cursor, selection color.NRGBA
}

// SetChrome recolors every terminal's background, cursor and selection
// to match the app's tint and accent.
func SetChrome(background, cursor, selection color.NRGBA) {
	chrome.Lock()
	defer chrome.Unlock()
	chrome.version++
	chrome.background, chrome.cursor, chrome.selection = background, cursor, selection
}
