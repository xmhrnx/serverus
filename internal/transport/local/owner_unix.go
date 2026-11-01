//go:build !windows

package local

import (
	"io/fs"
	"strconv"
	"syscall"

	"github.com/xmhrnx/serverus/internal/transport"
)

func fillOwner(fi *transport.FileInfo, info fs.FileInfo) {
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		fi.Owner = strconv.FormatUint(uint64(st.Uid), 10)
		fi.Group = strconv.FormatUint(uint64(st.Gid), 10)
	}
}
