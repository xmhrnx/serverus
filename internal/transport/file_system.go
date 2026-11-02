package transport

import (
	"io"
	"io/fs"
	"path"
	"time"
)

// FileInfo describes one directory entry.
type FileInfo struct {
	Name    string
	Size    int64
	Mode    fs.FileMode
	ModTime time.Time
	IsDir   bool
	// IsLink is set for symbolic links; IsDir then says whether the
	// link points to a directory.
	IsLink bool
	Owner  string
	Group  string
}

// FileSystem browses and transfers files. Paths use forward slashes,
// whatever the server's platform.
type FileSystem interface {
	// Home is the directory a session starts in.
	Home() (string, error)
	ReadDir(dir string) ([]FileInfo, error)
	Stat(p string) (FileInfo, error)
	Open(p string) (io.ReadCloser, error)
	// Create opens p for writing, truncating it.
	Create(p string) (io.WriteCloser, error)
	Mkdir(p string) error
	Rename(from, to string) error
	// Remove deletes a file, or a directory with everything in it.
	Remove(p string, isDir bool) error
	Chmod(p string, mode fs.FileMode) error
}

// Join joins remote path elements with forward slashes.
func Join(elem ...string) string { return path.Join(elem...) }

// Parent is the directory containing p ("/" for "/").
func Parent(p string) string { return path.Dir(p) }
