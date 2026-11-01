package files

import (
	"sort"
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
)

// column is what the listing is sorted by.
type column int

const (
	byName column = iota
	bySize
	byModified
)

// sortEntries orders directories first, then by the column.
func sortEntries(entries []transport.FileInfo, by column, desc bool) {
	sort.SliceStable(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		if a.IsDir != b.IsDir {
			return a.IsDir
		}
		var less bool
		switch by {
		case bySize:
			less = a.Size < b.Size
		case byModified:
			less = a.ModTime.Before(b.ModTime)
		default:
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		}
		if desc {
			return !less
		}
		return less
	})
}
