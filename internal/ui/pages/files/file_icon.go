// Package files is the file manager: browse the server's file system,
// upload and download with progress, edit text files in place, rename,
// delete and change permissions.
package files

import (
	"path"
	"strings"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// fileIcon picks an icon and color for an entry.
func fileIcon(f transport.FileInfo) (string, giode.Color) {
	if f.IsDir {
		return "folder", theme.Primary
	}
	switch strings.ToLower(path.Ext(f.Name)) {
	case ".sh", ".bash", ".zsh", ".py", ".go", ".js", ".ts", ".rb", ".php", ".rs", ".c", ".h", ".java":
		return "code", theme.Accent
	case ".conf", ".cfg", ".ini", ".toml", ".yaml", ".yml", ".json", ".env", ".service":
		return "settings", theme.Warning
	case ".log", ".txt", ".md":
		return "file-text", theme.Info
	case ".pem", ".key", ".crt", ".pub":
		return "key-round", theme.Danger
	}
	return "file", theme.Muted
}
