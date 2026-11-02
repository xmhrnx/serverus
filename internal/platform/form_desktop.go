//go:build !android && !ios

package platform

// Desktop is true on Windows, macOS and Linux: there are windows with
// title bars, folders the app can browse, and a shell to run.
const Desktop = true
