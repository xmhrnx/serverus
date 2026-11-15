// Package shortcuts lists the app's keyboard shortcuts, so the shell
// that registers them and the settings page that explains them agree.
package shortcuts

// Entry is a shortcut and what it does, for the settings page.
type Entry struct {
	Spec string
	Does string
}

// Key specs, in giode.Shortcut's format. Mod is Command on macOS and Ctrl
// elsewhere.
const (
	Search      = "Mod+K"
	AddServer   = "Mod+N"
	Settings    = "Mod+,"
	Overview    = "Mod+0"
	NewTerminal = "Mod+Shift+T"
	// Page is Mod plus the page's position in the tab row, 1 to 9.
	Page = "Mod+"
)

// Help lists every shortcut.
var Help = []Entry{
	{Overview, "Open the overview of every server"},
	{Search, "Search the servers"},
	{AddServer, "Add a server"},
	{Settings, "Open the settings"},
	{NewTerminal, "Open a terminal tab on the server"},
	{Page + "1", "Go to the server's first page; 2 to 9 for the others"},
}
