package transfer

import (
	"context"
	"io"
	"sync/atomic"
)

// countingReader counts bytes and stops the copy when ctx ends. Size
// tells an SFTP file how many requests it may keep in flight.
type countingReader struct {
	r    io.Reader
	n    *atomic.Int64
	ctx  context.Context
	size int64
}

// countingWriter counts bytes and stops the copy when ctx ends.
type countingWriter struct {
	w   io.Writer
	n   *atomic.Int64
	ctx context.Context
}

func (c *countingReader) Read(p []byte) (int, error) {
	if err := c.ctx.Err(); err != nil {
		return 0, err
	}
	n, err := c.r.Read(p)
	c.n.Add(int64(n))
	return n, err
}

// Size is the bytes left to read, or -1 when unknown.
func (c *countingReader) Size() int64 {
	if c.size < 0 {
		return -1
	}
	return max(0, c.size-c.n.Load())
}

func (c *countingWriter) Write(p []byte) (int, error) {
	if err := c.ctx.Err(); err != nil {
		return 0, err
	}
	n, err := c.w.Write(p)
	c.n.Add(int64(n))
	return n, err
}
