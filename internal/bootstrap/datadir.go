package bootstrap

import (
	"os"
	"path/filepath"

	gioapp "gioui.org/app"
)

// dataDir is where the workspace and known hosts live: the platform's
// configuration folder (~/.config/Serverus, %AppData%\Serverus,
// ~/Library/Application Support/Serverus, the app's files on Android).
// SERVERUS_DATA_DIR overrides it, for portable installs and tests.
func dataDir() (string, error) {
	if d := os.Getenv("SERVERUS_DATA_DIR"); d != "" {
		return d, os.MkdirAll(d, 0o700)
	}
	base, err := gioapp.DataDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(base, "Serverus")
	return d, os.MkdirAll(d, 0o700)
}
