package bootstrap

import (
	"os"
	"path/filepath"

	gioapp "gioui.org/app"
)

// dataDir is where the workspace and known hosts live: the platform's
// configuration folder (~/.config/Bastion, %AppData%\Bastion,
// ~/Library/Application Support/Bastion, the app's files on Android).
// BASTION_DATA_DIR overrides it, for portable installs and tests.
func dataDir() (string, error) {
	if d := os.Getenv("BASTION_DATA_DIR"); d != "" {
		return d, os.MkdirAll(d, 0o700)
	}
	base, err := gioapp.DataDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(base, "Bastion")
	return d, os.MkdirAll(d, 0o700)
}
