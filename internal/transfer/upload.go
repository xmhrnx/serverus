package transfer

import (
	"io"
	"os"
	"path"

	"github.com/xmhrnx/serverus/internal/transport"
)

func uploadJob(serverID string, fsys transport.FileSystem, local, remote string, size int64, onDone func()) *Job {
	return &Job{
		ServerID: serverID, Direction: Upload, Name: path.Base(remote),
		From: local, To: remote, Size: size,
		OpenSource: func() (io.ReadCloser, error) { return os.Open(local) },
		OpenDest:   func() (io.WriteCloser, error) { return fsys.Create(remote) },
		OnDone:     onDone,
	}
}

// UploadStream queues the upload of an already open stream (a file
// picked on a phone, which has no path).
func (q *Queue) UploadStream(serverID string, fsys transport.FileSystem, r io.ReadCloser, name string, size int64, dir string, onDone func()) {
	remote := path.Join(dir, name)
	q.Add(&Job{
		ServerID: serverID, Direction: Upload, Name: name, From: name, To: remote, Size: size,
		OpenSource: func() (io.ReadCloser, error) { return r, nil },
		OpenDest:   func() (io.WriteCloser, error) { return fsys.Create(remote) },
		OnDone:     onDone,
	})
}
