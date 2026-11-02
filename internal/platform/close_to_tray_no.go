//go:build !windows && (!darwin || ios)

package platform

// CloseToTray is false on Linux, where Gio can't keep a closed window
// (closing quits), and on phones.
const CloseToTray = false
