package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// maxKeySize bounds a private key copy; real keys are a few KB.
const maxKeySize = 1 << 20

// SaveKey keeps a private copy of an SSH private key in the data folder's
// keys/ subfolder, readable only by the user, and returns its path. It's
// for where the app can't reopen the user's own file later: Android's
// picker gives no path, and the Mac App Store's sandbox forgets a picked
// file when the app quits. The copy's name carries a hash of its content,
// so keys with the same file name from different folders don't collide.
func (r *Repository) SaveKey(name string, src io.Reader) (string, error) {
	data, err := io.ReadAll(io.LimitReader(src, maxKeySize+1))
	if err != nil {
		return "", fmt.Errorf("read the private key: %w", err)
	}
	if len(data) == 0 || len(data) > maxKeySize {
		return "", fmt.Errorf("%s doesn't look like a private key", name)
	}
	sum := sha256.Sum256(data)
	base := strings.TrimSpace(filepath.Base(name))
	if base == "" || base == "." || base == string(filepath.Separator) {
		base = "key"
	}
	dir := filepath.Join(filepath.Dir(r.path), "keys")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	p := filepath.Join(dir, base+"-"+hex.EncodeToString(sum[:4]))
	if old, err := os.ReadFile(p); err == nil && bytes.Equal(old, data) {
		return p, nil
	}
	if err := os.WriteFile(p, data, 0o600); err != nil {
		return "", fmt.Errorf("save the private key: %w", err)
	}
	return p, nil
}
