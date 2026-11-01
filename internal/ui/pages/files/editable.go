package files

import (
	"path"
	"strings"
	"unicode/utf8"
)

// maxEditSize is the largest file the editor opens.
const maxEditSize = 2 << 20

// looksTextual guesses from the name whether a file is worth opening in
// the editor; the content is checked once loaded.
func looksTextual(name string, size int64) bool {
	if size > maxEditSize {
		return false
	}
	switch strings.ToLower(path.Ext(name)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".pdf", ".zip", ".gz", ".tgz", ".xz", ".bz2",
		".7z", ".tar", ".deb", ".rpm", ".iso", ".img", ".bin", ".so", ".o", ".a", ".exe", ".dll", ".mp3",
		".mp4", ".mkv", ".mov", ".woff", ".woff2", ".ttf", ".otf", ".db", ".sqlite":
		return false
	}
	return true
}

// isText reports whether data is UTF-8 text without NUL bytes.
func isText(data []byte) bool {
	return utf8.Valid(data) && !strings.ContainsRune(string(data), 0)
}
