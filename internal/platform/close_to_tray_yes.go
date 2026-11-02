//go:build windows || (darwin && !ios)

package platform

// CloseToTray is true where a closed window can live on in the tray.
const CloseToTray = true
