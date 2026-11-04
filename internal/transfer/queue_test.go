package transfer

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xmhrnx/serverus/internal/transport/local"
)

func waitIdle(t *testing.T, q *Queue) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for q.Active() > 0 {
		if time.Now().After(deadline) {
			t.Fatal("transfers didn't finish")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestUploadAndDownloadTrees(t *testing.T) {
	src := t.TempDir()
	os.MkdirAll(filepath.Join(src, "site", "css"), 0o755)
	os.WriteFile(filepath.Join(src, "site", "index.html"), []byte("<h1>hi</h1>"), 0o644)
	os.WriteFile(filepath.Join(src, "site", "css", "app.css"), []byte("body{}"), 0o644)

	remote := t.TempDir() // the "server" is this computer's file system
	fs := local.Open()
	q := NewQueue(2, nil)
	if err := q.UploadPath("srv", fs, filepath.Join(src, "site"), filepath.ToSlash(remote), nil); err != nil {
		t.Fatal(err)
	}
	waitIdle(t, q)
	got, err := os.ReadFile(filepath.Join(remote, "site", "css", "app.css"))
	if err != nil || string(got) != "body{}" {
		t.Fatalf("uploaded file = %q, %v", got, err)
	}
	for _, j := range q.Jobs() {
		if st := j.Status(); st.State != Done || st.Progress(j.Size) != 1 {
			t.Errorf("job %s: %+v", j.Name, st)
		}
	}

	dest := t.TempDir()
	info, _ := fs.Stat(filepath.ToSlash(filepath.Join(remote, "site")))
	if err := q.DownloadPath("srv", fs, info, filepath.ToSlash(filepath.Join(remote, "site")), dest); err != nil {
		t.Fatal(err)
	}
	waitIdle(t, q)
	got, err = os.ReadFile(filepath.Join(dest, "site", "index.html"))
	if err != nil || string(got) != "<h1>hi</h1>" {
		t.Fatalf("downloaded file = %q, %v", got, err)
	}
	q.ClearFinished()
	if len(q.Jobs()) != 0 {
		t.Error("finished jobs not cleared")
	}
}

// blocked is a source that waits until release closes.
type blocked struct{ release chan struct{} }

func (b blocked) Read([]byte) (int, error) { <-b.release; return 0, io.EOF }
func (b blocked) Close() error             { return nil }

func running(q *Queue) int {
	n := 0
	for _, j := range q.Jobs() {
		if j.Status().State == Running {
			n++
		}
	}
	return n
}

func eventually(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("condition not reached")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestSetConcurrency(t *testing.T) {
	release := make(chan struct{})
	q := NewQueue(1, nil)
	for range 3 {
		q.Add(&Job{
			Size:       -1,
			OpenSource: func() (io.ReadCloser, error) { return blocked{release}, nil },
			OpenDest:   func() (io.WriteCloser, error) { return nopWriter{}, nil },
		})
	}
	eventually(t, func() bool { return running(q) == 1 })
	time.Sleep(50 * time.Millisecond)
	if n := running(q); n != 1 {
		t.Fatalf("%d running with a limit of 1", n)
	}
	q.SetConcurrency(3)
	eventually(t, func() bool { return running(q) == 3 })
	close(release)
	waitIdle(t, q)
}

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }
func (nopWriter) Close() error                { return nil }
