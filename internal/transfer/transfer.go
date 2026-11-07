package transfer

import (
	"context"
	"io"
	"os"
)

// run copies source to destination, counting bytes as they go.
func (j *Job) run(ctx context.Context) error {
	src, err := j.OpenSource()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := j.OpenDest()
	if err != nil {
		return err
	}
	err = j.copy(ctx, dst, src)
	if cerr := dst.Close(); err == nil {
		err = cerr
	}
	return err
}

// copy moves src to dst through whichever side can pipeline: an SFTP
// file sends many requests at once from ReadFrom and WriteTo, where a
// plain 32 KB copy waits a round trip for every chunk.
func (j *Job) copy(ctx context.Context, dst io.Writer, src io.Reader) error {
	_, srcIsLocal := src.(*os.File)
	_, dstIsLocal := dst.(*os.File)
	if wt, ok := src.(io.WriterTo); ok && !srcIsLocal {
		_, err := wt.WriteTo(&countingWriter{w: dst, n: &j.done, ctx: ctx})
		return err
	}
	if rf, ok := dst.(io.ReaderFrom); ok && !dstIsLocal {
		_, err := rf.ReadFrom(&countingReader{r: src, n: &j.done, ctx: ctx, size: j.Size})
		return err
	}
	_, err := io.Copy(&countingWriter{w: dst, n: &j.done, ctx: ctx}, src)
	return err
}
