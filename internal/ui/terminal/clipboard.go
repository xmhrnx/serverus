package terminal

// Clipboard is the system clipboard; *giode.App implements it.
type Clipboard interface {
	SetClipboard(text string)
	ReadClipboard(fn func(text string))
}
