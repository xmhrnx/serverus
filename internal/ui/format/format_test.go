package format

import (
	"io/fs"
	"testing"
	"time"
)

func TestBytes(t *testing.T) {
	cases := map[uint64]string{0: "0 B", 1023: "1023 B", 1536: "1.5 KiB", 5 << 30: "5.0 GiB", 200 << 20: "200 MiB"}
	for in, want := range cases {
		if got := Bytes(in); got != want {
			t.Errorf("Bytes(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestDuration(t *testing.T) {
	if got := Duration(50*time.Hour + 3*time.Minute); got != "2d 2h" {
		t.Errorf("got %q", got)
	}
	if got := Duration(90 * time.Second); got != "1m 30s" {
		t.Errorf("got %q", got)
	}
}

func TestMode(t *testing.T) {
	if got := Mode(fs.ModeDir | 0o755); got != "drwxr-xr-x" {
		t.Errorf("got %q", got)
	}
	if got := Mode(0o640); got != "-rw-r-----" {
		t.Errorf("got %q", got)
	}
}
