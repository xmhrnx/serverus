package terminal

// applyChrome updates p from the theme when it changed since version;
// it returns the version applied.
func applyChrome(p *Palette, version int) (int, bool) {
	chrome.Lock()
	defer chrome.Unlock()
	if chrome.version == version {
		return version, false
	}
	p.Background, p.Cursor, p.Selection = chrome.background, chrome.cursor, chrome.selection
	return chrome.version, true
}
