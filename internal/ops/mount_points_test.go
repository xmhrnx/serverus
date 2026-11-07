package ops

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/xmhrnx/serverus/internal/transport/local"
)

// The mount table's last line is another device: the script must still
// succeed. This broke unmounting on machines whose newest mount was a
// different disk.
func TestMountPointsWhenLastLineIsAnotherDevice(t *testing.T) {
	table := filepath.Join(t.TempDir(), "mounts")
	os.WriteFile(table, []byte("proc /proc proc rw 0 0\n/dev/null /srv/my\\040stick ext4 rw 0 0\n/dev/zero /mnt/other ext4 rw 0 0\n"), 0o644)
	d := Drives{Exec: local.Open()}
	got, err := d.mountPointsIn(context.Background(), "/dev/null", table)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "/srv/my stick" {
		t.Fatalf("got %q", got)
	}
	none, err := d.mountPointsIn(context.Background(), "/dev/full", table)
	if err != nil || len(none) != 0 {
		t.Fatalf("got %q, %v", none, err)
	}
}
