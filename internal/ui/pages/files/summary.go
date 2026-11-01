package files

import (
	"fmt"

	"github.com/xmhrnx/bastion/internal/ui/format"
)

// itemsLabel summarizes a listing: "12 files (3.4 MiB), 4 folders".
func itemsLabel(files, dirs int, size int64) string {
	plural := func(n int, one, many string) string {
		if n == 1 {
			return one
		}
		return many
	}
	return fmt.Sprintf("%d %s (%s), %d %s", files, plural(files, "file", "files"), format.Size(size), dirs, plural(dirs, "folder", "folders"))
}
