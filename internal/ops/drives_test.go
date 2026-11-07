package ops

import (
	"testing"
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
