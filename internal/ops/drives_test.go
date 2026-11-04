package ops

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/transport/local"
)

func TestParseLsblk(t *testing.T) {
	out := `{"blockdevices": [
	  {"name":"loop0","path":"/dev/loop0","size":4096,"type":"loop","fstype":"squashfs","mountpoint":"/snap/core/1","rm":false,"hotplug":false},
	  {"name":"zram0","path":"/dev/zram0","size":0,"type":"disk","rm":false,"hotplug":false},
	  {"name":"nvme0n1","path":"/dev/nvme0n1","size":512110190592,"type":"disk","rm":false,"hotplug":false,"tran":"nvme","model":"Samsung SSD 980 ",
	   "children":[
	     {"name":"nvme0n1p1","path":"/dev/nvme0n1p1","size":536870912,"type":"part","fstype":"vfat","mountpoint":"/boot/efi","fsused":6291456,"fsavail":530579456},
	     {"name":"nvme0n1p2","path":"/dev/nvme0n1p2","size":511571394560,"type":"part","fstype":"ext4","mountpoint":"/","fsused":100,"fsavail":200}
	   ]},
	  {"name":"sda","path":"/dev/sda","size":64021856256,"type":"disk","rm":true,"hotplug":true,"tran":"usb","model":"SanDisk Ultra",
	   "children":[{"name":"sda1","path":"/dev/sda1","size":64020807680,"type":"part","fstype":"exfat","label":"STICK","uuid":"1234-ABCD","mountpoint":null}]}
	]}`
	drives, err := parseLsblk(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(drives) != 2 {
		t.Fatalf("got %d drives: %+v", len(drives), drives)
	}
	usb := drives[0]
	if usb.Name != "sda" || !usb.External() || usb.Title() != "SanDisk Ultra" || len(usb.Volumes) != 1 {
		t.Fatalf("external drive first, got %+v", usb)
	}
	if v := usb.Volumes[0]; v.Mounted() || !v.Mountable() || v.Title() != "STICK" || v.Path != "/dev/sda1" {
		t.Errorf("stick volume %+v", v)
	}
	sys := drives[1]
	if sys.Title() != "Samsung SSD 980" || !sys.System() || sys.Volumes[1].Used != 100 {
		t.Errorf("system drive %+v", sys)
	}
}

func TestParseLsblkOldVersion(t *testing.T) {
	// util-linux 2.31: strings for numbers and booleans, no path.
	out := `{"blockdevices": [{"name":"sdb","size":"8004304896","type":"disk","fstype":"vfat","label":"CAM","mountpoint":"/media/pi/CAM","rm":"1","hotplug":"1","tran":"usb","model":"Card Reader"}]}`
	drives, err := parseLsblk(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(drives) != 1 || !drives[0].Removable || drives[0].Size != 8004304896 || len(drives[0].Volumes) != 1 {
		t.Fatalf("got %+v", drives)
	}
	if v := drives[0].Volumes[0]; v.Path != "/dev/sdb" || v.MountPoint != "/media/pi/CAM" {
		t.Errorf("whole-disk volume %+v", v)
	}
}

func TestMountName(t *testing.T) {
	if got := mountName(Volume{Name: "sdb1", Label: "My Files/../x"}); got != "My_Files..x" {
		t.Errorf("got %q", got)
	}
	if got := mountName(Volume{Name: "sdb1", Label: ".."}); got != "sdb1" {
		t.Errorf("got %q", got)
	}
}

func TestParseScheduled(t *testing.T) {
	s, ok := parseScheduled("USEC=1762279200000000\nWARN_WALL=1\nMODE=reboot\n")
	if !ok || s.Action != PowerReboot || !s.At.Equal(time.UnixMicro(1762279200000000)) {
		t.Fatalf("got %+v %v", s, ok)
	}
	if _, ok := parseScheduled(""); ok {
		t.Error("nothing scheduled")
	}
}

// fakeSudo answers like a server where sudo wants the password "pw".
type fakeSudo struct {
	asked    []bool
	answers  []string
	commands []string
}

func (f *fakeSudo) Run(_ context.Context, cmd string) (transport.Output, error) {
	f.commands = append(f.commands, cmd)
	if strings.HasPrefix(cmd, "sudo -n") {
		return transport.Output{ExitCode: 1, Stderr: []byte("sudo: a password is required\n")}, nil
	}
	return transport.Output{ExitCode: 1, Stderr: []byte("shutdown: Permission denied")}, nil
}

func (f *fakeSudo) RunInput(_ context.Context, cmd string, in []byte) (transport.Output, error) {
	f.commands = append(f.commands, cmd)
	if string(in) != "pw\n" {
		return transport.Output{ExitCode: 1, Stderr: []byte("sudo: 1 incorrect password attempt\n")}, nil
	}
	return transport.Output{Stdout: []byte("ok")}, nil
}

func (f *fakeSudo) SudoPassword(retry bool) (string, bool) {
	f.asked = append(f.asked, retry)
	if len(f.answers) == 0 {
		return "", false
	}
	a := f.answers[0]
	f.answers = f.answers[1:]
	return a, true
}

func TestSudoAsksAgainAfterAWrongPassword(t *testing.T) {
	f := &fakeSudo{answers: []string{"nope", "pw"}}
	out, err := runPrivileged(context.Background(), f, "shutdown -c")
	if err != nil || out != "ok" {
		t.Fatalf("got %q, %v", out, err)
	}
	if len(f.asked) != 2 || f.asked[0] || !f.asked[1] {
		t.Errorf("asked %v", f.asked)
	}
	if last := f.commands[len(f.commands)-1]; last != "sudo -S -p '' sh -c 'shutdown -c'" {
		t.Errorf("ran %q", last)
	}
}

func TestSudoDeclined(t *testing.T) {
	f := &fakeSudo{}
	if _, err := runPrivileged(context.Background(), f, "shutdown -c"); err != ErrCanceled {
		t.Fatalf("got %v", err)
	}
}

func TestCheckMountDir(t *testing.T) {
	for dir, ok := range map[string]bool{
		"/mnt/data": true, "/srv/backup": true, "/home/pi/usb": true,
		"/": false, "/etc": false, "/mnt": false, "/home/": false, "data": false, "": false,
	} {
		if err := CheckMountDir(dir); (err == nil) != ok {
			t.Errorf("%q: %v", dir, err)
		}
	}
}

// fakeMounts is a server where /dev/sdb1 is mounted on mounts, udisks
// is missing and umount needs root.
type fakeMounts struct {
	mounts []string
	ran    []string
}

func (f *fakeMounts) Run(_ context.Context, cmd string) (transport.Output, error) {
	f.ran = append(f.ran, cmd)
	switch {
	case strings.HasPrefix(cmd, "udisksctl"):
		return transport.Output{ExitCode: 127, Stderr: []byte("udisksctl: not found")}, nil
	case strings.HasPrefix(cmd, "sh -c") && strings.Contains(cmd, "/proc/mounts"):
		return transport.Output{Stdout: []byte(strings.Join(f.mounts, "\n"))}, nil
	case strings.HasPrefix(cmd, "umount"):
		return transport.Output{ExitCode: 32, Stderr: []byte("umount: must be superuser to use umount.")}, nil
	case strings.HasPrefix(cmd, "sudo -n sh -c 'umount "):
		if len(f.mounts) > 0 {
			f.mounts = f.mounts[:len(f.mounts)-1]
		}
		return transport.Output{}, nil
	}
	return transport.Output{ExitCode: 1, Stderr: []byte("unexpected " + cmd)}, nil
}

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

// stuckMounts never lets go.
type stuckMounts struct{ fakeMounts }

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
