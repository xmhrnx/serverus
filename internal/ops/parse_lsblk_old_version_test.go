package ops

import (
	"testing"
)

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
