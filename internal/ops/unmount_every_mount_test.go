package ops

import (
	"context"
	"strings"
	"testing"

	"github.com/xmhrnx/serverus/internal/transport"
)

func TestUnmountEveryMount(t *testing.T) {
	f := &fakeMounts{mounts: []string{`/media/pi/My\040Stick`, "/mnt/data"}}
	v := Volume{Name: "sdb1", Path: "/dev/sdb1", Label: "My Stick", MountPoint: "/media/pi/My Stick"}
	if err := (Drives{Exec: f}).Unmount(context.Background(), v); err != nil {
		t.Fatal(err)
	}
	if len(f.mounts) != 0 {
		t.Fatalf("still mounted on %v", f.mounts)
	}
	var umounts []string
	for _, c := range f.ran {
		if strings.HasPrefix(c, "sudo -n") {
			umounts = append(umounts, c)
		}
	}
	if len(umounts) != 2 || !strings.Contains(umounts[0], "/mnt/data") || !strings.Contains(umounts[1], "/media/pi/My Stick") {
		t.Errorf("unmounted %v", umounts)
	}
}

func (s *stuckMounts) Run(ctx context.Context, cmd string) (transport.Output, error) {
	if strings.HasPrefix(cmd, "sudo -n sh -c 'umount ") {
		return transport.Output{}, nil // claims success, stays mounted
	}
	return s.fakeMounts.Run(ctx, cmd)
}

func TestUnmountReportsWhatStaysMounted(t *testing.T) {
	s := &stuckMounts{fakeMounts{mounts: []string{"/srv/disk"}}}
	err := (Drives{Exec: s}).Unmount(context.Background(), Volume{Name: "sdb1", Path: "/dev/sdb1"})
	if err == nil || !strings.Contains(err.Error(), "still mounted at /srv/disk") {
		t.Fatalf("got %v", err)
	}
}
