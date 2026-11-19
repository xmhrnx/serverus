package platform

import (
	"os"
	"strings"
)

// Sandboxed reports whether the app runs in the macOS App Sandbox, as
// Mac App Store builds do. Unlike the constants here it's decided at
// run time: the same binary runs inside and outside the sandbox. macOS
// doesn't always set APP_SANDBOX_CONTAINER_ID, but a sandboxed app's
// home is always its container, ~/Library/Containers/<id>/Data.
func Sandboxed() bool {
	return Mac && (os.Getenv("APP_SANDBOX_CONTAINER_ID") != "" ||
		strings.Contains(os.Getenv("HOME"), "/Library/Containers/"))
}
