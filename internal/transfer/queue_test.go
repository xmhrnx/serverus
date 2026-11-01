package transfer

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xmhrnx/bastion/internal/transport/local"
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
