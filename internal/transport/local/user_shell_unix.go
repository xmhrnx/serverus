//go:build !windows

package local

import "os"

// userShell is the user's login shell, $SHELL, or /bin/sh.
func userShell() (string, []string) {
	if sh := os.Getenv("SHELL"); sh != "" {
		return sh, []string{"-l"}
	}
	return "/bin/sh", []string{"-l"}
}
