package local

import (
	"io/fs"

	"github.com/xmhrnx/serverus/internal/transport"
)

func toFileInfo(info fs.FileInfo) transport.FileInfo {
	return transport.FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
		IsLink:  info.Mode()&fs.ModeSymlink != 0,
	}
}
