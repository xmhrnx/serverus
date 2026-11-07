// Package storage persists the workspace (servers, snippets, tunnels,
// settings), known host keys and secrets.
package storage

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

// readJSON decodes the file at path into v. A missing file leaves v
// untouched and is not an error.
func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
