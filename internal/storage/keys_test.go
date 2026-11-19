package storage

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSaveKeyKeepsAPrivateCopy(t *testing.T) {
	dir := t.TempDir()
	r, err := OpenRepository(dir)
	if err != nil {
		t.Fatal(err)
	}
	p, err := r.SaveKey("id_ed25519", strings.NewReader("first key"))
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(p) != filepath.Join(dir, "keys") || !strings.HasPrefix(filepath.Base(p), "id_ed25519-") {
		t.Fatalf("saved at %s", p)
	}
	if b, _ := os.ReadFile(p); string(b) != "first key" {
		t.Fatalf("content %q", b)
	}
	if st, _ := os.Stat(p); runtime.GOOS != "windows" && st.Mode().Perm() != 0o600 {
		t.Fatalf("mode %v", st.Mode().Perm())
	}
	again, err := r.SaveKey("id_ed25519", strings.NewReader("first key"))
	if err != nil || again != p {
		t.Fatalf("same key saved again at %s (%v), want %s", again, err, p)
	}
	other, err := r.SaveKey("id_ed25519", strings.NewReader("another key"))
	if err != nil || other == p {
		t.Fatalf("a different key with the same name overwrote the first: %s (%v)", other, err)
	}
	if _, err := r.SaveKey("empty", strings.NewReader("")); err == nil {
		t.Fatal("an empty file was saved as a key")
	}
}
