package local

import (
	"os"
	"path/filepath"

	"github.com/xmhrnx/serverus/internal/transport"
)

// Home is the user's home directory, with forward slashes.
func (c *Conn) Home() (string, error) {
	h, err := os.UserHomeDir()
	return filepath.ToSlash(h), err
}

// ReadDir lists dir.
func (c *Conn) ReadDir(dir string) ([]transport.FileInfo, error) {
	native := filepath.FromSlash(dir)
	entries, err := os.ReadDir(native)
	if err != nil {
		return nil, err
	}
	out := make([]transport.FileInfo, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		fi := toFileInfo(info)
		if fi.IsLink {
			if target, err := os.Stat(filepath.Join(native, e.Name())); err == nil {
				fi.IsDir = target.IsDir()
			}
		}
		fillOwner(&fi, info)
		out = append(out, fi)
	}
	return out, nil
}

// Stat describes p.
func (c *Conn) Stat(p string) (transport.FileInfo, error) {
	info, err := os.Stat(filepath.FromSlash(p))
	if err != nil {
		return transport.FileInfo{}, err
	}
	fi := toFileInfo(info)
	fillOwner(&fi, info)
	return fi, nil
}
