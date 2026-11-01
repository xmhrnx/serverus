//go:build windows

package local

import (
	"io/fs"

	"github.com/xmhrnx/serverus/internal/transport"
)

// fillOwner does nothing on Windows, whose files have ACLs rather than
// an owner and group the file browser could show.
func fillOwner(*transport.FileInfo, fs.FileInfo) {}
