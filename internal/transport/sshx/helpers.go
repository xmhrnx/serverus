package sshx

import (
	"io/fs"
	"strconv"

	"github.com/pkg/sftp"

	"github.com/xmhrnx/serverus/internal/transport"
)

func toFileInfo(fi fs.FileInfo) transport.FileInfo {
	out := transport.FileInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Mode:    fi.Mode(),
		ModTime: fi.ModTime(),
		IsDir:   fi.IsDir(),
		IsLink:  fi.Mode()&fs.ModeSymlink != 0,
	}
	if st, ok := fi.Sys().(*sftp.FileStat); ok {
		out.Owner = strconv.FormatUint(uint64(st.UID), 10)
		out.Group = strconv.FormatUint(uint64(st.GID), 10)
	}
	return out
}
