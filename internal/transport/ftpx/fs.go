package ftpx

import (
	"io"
	"io/fs"
	"path"
	"sync"

	"github.com/jlaffaye/ftp"

	"github.com/xmhrnx/bastion/internal/transport"
)

// Home is the directory the server signs the user into.
func (c *Conn) Home() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sc.CurrentDir()
}

// ReadDir lists dir.
func (c *Conn) ReadDir(dir string) ([]transport.FileInfo, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entries, err := c.sc.List(dir)
	if err != nil {
		return nil, err
	}
	out := make([]transport.FileInfo, 0, len(entries))
	for _, e := range entries {
		if e.Name == "." || e.Name == ".." {
			continue
		}
		out = append(out, toFileInfo(e))
	}
	return out, nil
}

// Stat describes p.
func (c *Conn) Stat(p string) (transport.FileInfo, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, err := c.sc.GetEntry(p)
	if err != nil {
		return transport.FileInfo{}, err
	}
	fi := toFileInfo(e)
	fi.Name = path.Base(p)
	return fi, nil
}

// Open streams p. The connection is busy until the reader is closed.
func (c *Conn) Open(p string) (io.ReadCloser, error) {
	c.mu.Lock()
	r, err := c.sc.Retr(p)
	if err != nil {
		c.mu.Unlock()
		return nil, err
	}
	return &unlockingReader{r: r, unlock: c.mu.Unlock}, nil
}

// Create uploads what's written to p; Close waits for the upload to
// finish and reports its error.
func (c *Conn) Create(p string) (io.WriteCloser, error) {
	pr, pw := io.Pipe()
	w := &uploadWriter{pw: pw, done: make(chan error, 1)}
	c.mu.Lock()
	go func() {
		err := c.sc.Stor(p, pr)
		c.mu.Unlock()
		pr.CloseWithError(err)
		w.done <- err
	}()
	return w, nil
}

// Mkdir creates the directory p.
func (c *Conn) Mkdir(p string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sc.MakeDir(p)
}

// Rename moves from to to.
func (c *Conn) Rename(from, to string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sc.Rename(from, to)
}

// Remove deletes p, recursively for a directory.
func (c *Conn) Remove(p string, isDir bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if isDir {
		return c.sc.RemoveDirRecur(p)
	}
	return c.sc.Delete(p)
}

// Chmod is not part of FTP.
func (c *Conn) Chmod(string, fs.FileMode) error { return transport.ErrUnsupported }

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

// unlockingReader frees the connection when the download is closed.
type unlockingReader struct {
	r      *ftp.Response
	once   sync.Once
	unlock func()
}

func (u *unlockingReader) Read(p []byte) (int, error) { return u.r.Read(p) }

func (u *unlockingReader) Close() error {
	err := u.r.Close()
	u.once.Do(u.unlock)
	return err
}

// uploadWriter feeds Stor through a pipe.
type uploadWriter struct {
	pw   *io.PipeWriter
	done chan error
}

func (w *uploadWriter) Write(p []byte) (int, error) { return w.pw.Write(p) }

func (w *uploadWriter) Close() error {
	w.pw.Close()
	return <-w.done
}
