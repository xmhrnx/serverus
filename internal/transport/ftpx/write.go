package ftpx

import (
	"io"
)

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
