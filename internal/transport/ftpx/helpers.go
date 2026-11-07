package ftpx

import (
	"io"
	"io/fs"
	"sync"

	"github.com/jlaffaye/ftp"

	"github.com/xmhrnx/serverus/internal/transport"
)

// unlockingReader frees the connection when the download is closed.
type unlockingReader struct {
	r      *ftp.Response
	once   sync.Once
	unlock func()
}

// uploadWriter feeds Stor through a pipe.
type uploadWriter struct {
	pw   *io.PipeWriter
	done chan error
}

func toFileInfo(e *ftp.Entry) transport.FileInfo {
	fi := transport.FileInfo{
		Name:    e.Name,
		Size:    int64(e.Size),
		ModTime: e.Time,
		IsDir:   e.Type == ftp.EntryTypeFolder,
		IsLink:  e.Type == ftp.EntryTypeLink,
		Mode:    0o644,
	}
	if fi.IsDir {
		fi.Mode = fs.ModeDir | 0o755
	}
	return fi
}

func (u *unlockingReader) Read(p []byte) (int, error) { return u.r.Read(p) }

func (u *unlockingReader) Close() error {
	err := u.r.Close()
	u.once.Do(u.unlock)
	return err
}

func (w *uploadWriter) Write(p []byte) (int, error) { return w.pw.Write(p) }

func (w *uploadWriter) Close() error {
	w.pw.Close()
	return <-w.done
}
