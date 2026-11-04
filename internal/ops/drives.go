package ops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
)

// Drive is a disk attached to a Linux machine, with its volumes.
type Drive struct {
	Name      string // "sdb"
	Path      string // "/dev/sdb"
	Size      int64
	Model     string
	Transport string // "usb", "sata", "nvme", "mmc"…
	// Removable is set for USB sticks, card readers and other drives
	// that can be unplugged.
	Removable bool
	Volumes   []Volume
}

// Title is the drive's display name: its model, or its device name.
func (d Drive) Title() string {
	if d.Model != "" {
		return d.Model
	}
	return d.Name
}

// External reports whether the drive is plugged in rather than built in.
func (d Drive) External() bool { return d.Removable || d.Transport == "usb" }

// System reports whether the operating system runs from the drive.
func (d Drive) System() bool {
	return slices.ContainsFunc(d.Volumes, func(v Volume) bool { return v.System() })
}

// Volume is a file system on a drive: a partition, or the whole disk.
type Volume struct {
	Name       string // "sdb1"
	Path       string // "/dev/sdb1"
	Size       int64
	FSType     string
	Label      string
	UUID       string
	MountPoint string
	Used       int64 // bytes, when mounted
	Avail      int64
}

// Title is the volume's label, or its device name.
func (v Volume) Title() string {
	if v.Label != "" {
		return v.Label
	}
	return v.Name
}

// Mounted reports whether the volume is in use.
func (v Volume) Mounted() bool { return v.MountPoint != "" }

// Mountable reports whether the volume holds a file system that can be
// mounted.
func (v Volume) Mountable() bool {
	switch v.FSType {
	case "", "swap", "LVM2_member", "crypto_LUKS", "linux_raid_member", "zfs_member", "bcache":
		return false
	}
	return true
}

// System reports whether the operating system needs the volume mounted.
func (v Volume) System() bool {
	switch v.MountPoint {
	case "/", "/boot", "/boot/efi", "/boot/firmware", "/usr", "/var", "/home", "[SWAP]":
		return true
	}
	return strings.HasPrefix(v.MountPoint, "/snap/")
}

// Drives lists and mounts the disks of a Linux machine.
type Drives struct {
	Exec transport.Executor
}

// ErrNoLsblk is returned when the machine has no lsblk (not Linux, or a
// very small one).
var ErrNoLsblk = errors.New("drive management needs lsblk, part of util-linux")

const (
	lsblkColumns    = "NAME,PATH,SIZE,TYPE,FSTYPE,LABEL,UUID,MOUNTPOINT,RM,HOTPLUG,TRAN,MODEL,FSUSED,FSAVAIL"
	lsblkOldColumns = "NAME,SIZE,TYPE,FSTYPE,LABEL,UUID,MOUNTPOINT,RM,HOTPLUG,TRAN,MODEL"
)

// List returns the machine's drives, external ones first.
func (d Drives) List(ctx context.Context) ([]Drive, error) {
	out, err := run(ctx, d.Exec, "lsblk -J -b -o "+lsblkColumns)
	var ce *CommandError
	if errors.As(err, &ce) && strings.Contains(ce.Message, "unknown column") {
		// util-linux before 2.33 has no PATH, FSUSED or FSAVAIL.
		out, err = run(ctx, d.Exec, "lsblk -J -b -o "+lsblkOldColumns)
	}
	if errors.As(err, &ce) && (ce.ExitCode == 127 || strings.Contains(ce.Message, "not found")) {
		return nil, ErrNoLsblk
	}
	if err != nil {
		return nil, err
	}
	return parseLsblk(out)
}

// Mount mounts v and returns where. It asks udisks first, which mounts
// under /media like a desktop does and needs no root on desktops; when
// udisks is missing or refuses, it mounts under /mnt as root, owned by
// the signed-in user for file systems without owners (FAT, exFAT, NTFS).
func (d Drives) Mount(ctx context.Context, v Volume) (string, error) {
	out, err := run(ctx, d.Exec, "udisksctl mount --no-user-interaction -b "+Quote(v.Path))
	if err == nil {
		if m := regexp.MustCompile(` at (.+?)\.?$`).FindStringSubmatch(strings.TrimSpace(out)); m != nil {
			return m[1], nil
		}
		return "", nil
	}
	if !udisksUnavailable(err) {
		return "", err
	}
	dir := autoMountDir(v)
	return dir, d.mount(ctx, v, dir)
}

// autoMountDir is where Mount puts v without udisks.
func autoMountDir(v Volume) string { return path.Join("/mnt", mountName(v)) }

// ErrFolderNotEmpty is returned by MountAt when the folder has files,
// which the mounted drive would hide until it's unmounted.
var ErrFolderNotEmpty = errors.New("the folder isn't empty; its files would be hidden while the drive is mounted there")

// MountAt mounts v on dir, creating the folder when it doesn't exist.
// A folder with files in it is refused unless overlay is set.
func (d Drives) MountAt(ctx context.Context, v Volume, dir string, overlay bool) error {
	dir = path.Clean(strings.TrimSpace(dir))
	if err := CheckMountDir(dir); err != nil {
		return err
	}
	if !overlay {
		out, err := run(ctx, d.Exec, "ls -A "+Quote(dir)+" 2>/dev/null | head -n 1")
		if err == nil && strings.TrimSpace(out) != "" {
			return ErrFolderNotEmpty
		}
	}
	return d.mount(ctx, v, dir)
}

// CheckMountDir rejects folders that can't or shouldn't hold a mount:
// relative paths and the system's own folders.
func CheckMountDir(dir string) error {
	if !path.IsAbs(dir) {
		return errors.New("enter a full path, starting with /")
	}
	switch path.Clean(dir) {
	case "/", "/bin", "/boot", "/dev", "/etc", "/lib", "/lib64", "/proc", "/root", "/run", "/sbin", "/sys", "/usr", "/var", "/home", "/tmp", "/media", "/mnt":
		return fmt.Errorf("%s is a system folder; choose a folder inside it, like /mnt/data", path.Clean(dir))
	}
	return nil
}

// mount mounts v on dir as root, owned by the signed-in user for file
// systems without owners (FAT, exFAT, NTFS).
func (d Drives) mount(ctx context.Context, v Volume, dir string) error {
	opts := ""
	switch v.FSType {
	case "vfat", "exfat", "ntfs", "ntfs3", "msdos":
		if ids, err := run(ctx, d.Exec, "printf '%s,%s' \"$(id -u)\" \"$(id -g)\""); err == nil {
			if uid, gid, ok := strings.Cut(strings.TrimSpace(ids), ","); ok {
				opts = fmt.Sprintf(" -o uid=%s,gid=%s", uid, gid)
			}
		}
	}
	_, err := runPrivileged(ctx, d.Exec, fmt.Sprintf("mkdir -p %s && mount%s %s %s", Quote(dir), opts, Quote(v.Path), Quote(dir)))
	return err
}

// Unmount unmounts v from every folder it's mounted on, removing the
// folder when Mount made it under /mnt; folders the user chose are kept.
// It succeeds only once the volume is mounted nowhere.
func (d Drives) Unmount(ctx context.Context, v Volume) error {
	// udisks unmounts what it mounted without root. Whatever it can't,
	// or doesn't, is unmounted below as root, which also reports the
	// real reason (busy) when that fails too.
	_, _ = run(ctx, d.Exec, "udisksctl unmount --no-user-interaction -b "+Quote(v.Path))
	// A volume can be mounted in several places (by the desktop and by
	// fstab, or twice by hand).
	for range 8 {
		targets, err := d.mountPoints(ctx, v.Path)
		if err != nil {
			return err
		}
		if len(targets) == 0 {
			return nil
		}
		target := targets[len(targets)-1] // the newest mount first
		cmd := "umount " + Quote(target)
		if target == autoMountDir(v) {
			cmd += " && { rmdir " + Quote(target) + " 2>/dev/null; true; }"
		}
		if _, err := runPrivileged(ctx, d.Exec, cmd); err != nil {
			if isBusy(err) {
				return d.busyError(ctx, v, targets)
			}
			return err
		}
	}
	targets, _ := d.mountPoints(ctx, v.Path)
	return fmt.Errorf("%s is still mounted at %s", v.Title(), strings.Join(targets, ", "))
}

// mountPoints lists the folders the device at dev is mounted on, oldest
// first, from /proc/mounts, resolving device links (/dev/disk/by-uuid…).
func (d Drives) mountPoints(ctx context.Context, dev string) ([]string, error) {
	return d.mountPointsIn(ctx, dev, "/proc/mounts")
}

// mountPointsIn reads the mount table at table; tests pass their own.
func (d Drives) mountPointsIn(ctx context.Context, dev, table string) ([]string, error) {
	// The loop's status is its last test's, false whenever the table's
	// last line is another device, so the script ends with exit 0.
	script := `while read -r src dir rest; do
  case "$src" in /dev/*) if [ "$(readlink -f "$src" 2>/dev/null)" = ` + Quote(dev) + ` ]; then printf '%s\n' "$dir"; fi;; esac
done < ` + Quote(table) + `
exit 0`
	out, err := run(ctx, d.Exec, "sh -c "+Quote(script))
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, line := range strings.Split(out, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			dirs = append(dirs, unescapeMount(line))
		}
	}
	return dirs, nil
}

// unescapeMount decodes the octal escapes /proc/mounts uses for spaces,
// tabs and backslashes in folder names.
func unescapeMount(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+3 < len(s) {
			if n, err := strconv.ParseUint(s[i+1:i+4], 8, 8); err == nil {
				b.WriteByte(byte(n))
				i += 3
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// PowerOff unmounts every volume of drv and powers it down, so it can be
// unplugged safely.
func (d Drives) PowerOff(ctx context.Context, drv Drive) error {
	for _, v := range drv.Volumes {
		if v.Mounted() {
			if v.System() {
				return fmt.Errorf("%s holds the system and can't be ejected", drv.Title())
			}
			if err := d.Unmount(ctx, v); err != nil {
				return err
			}
		}
	}
	_, err := run(ctx, d.Exec, "udisksctl power-off --no-user-interaction -b "+Quote(drv.Path))
	if err == nil || !udisksUnavailable(err) {
		return err
	}
	_, err = runPrivileged(ctx, d.Exec, "sync && echo 1 > "+Quote("/sys/block/"+drv.Name+"/device/delete"))
	return err
}

// udisksUnavailable reports whether udisksctl is missing or not allowed
// to act for this user, so the root fallback should run.
func udisksUnavailable(err error) bool {
	var ce *CommandError
	if !errors.As(err, &ce) {
		return false
	}
	m := strings.ToLower(ce.Message)
	return ce.ExitCode == 127 || strings.Contains(m, "not found") || strings.Contains(m, "not authorized") ||
		strings.Contains(m, "notauthorized") || strings.Contains(m, "serviceunknown") || strings.Contains(m, "cannot autolaunch") ||
		strings.Contains(m, "no such interface") || strings.Contains(m, "error connecting")
}

// mountName is a folder name for v: its label, made safe, or its device
// name.
func mountName(v Volume) string {
	name := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_', r == '.':
			return r
		case r == ' ':
			return '_'
		}
		return -1
	}, v.Label)
	if name == "" || name == "." || name == ".." {
		return v.Name
	}
	return name
}

// lsblk prints numbers and booleans as JSON values or as strings,
// depending on its version.
type flexInt int64

func (f *flexInt) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		*f = 0
		return nil
	}
	n, err := strconv.ParseInt(s, 10, 64)
	*f = flexInt(n)
	return err
}

type flexBool bool

func (f *flexBool) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	*f = flexBool(s == "true" || s == "1")
	return nil
}

type blockDevice struct {
	Name       string        `json:"name"`
	Path       string        `json:"path"`
	Size       flexInt       `json:"size"`
	Type       string        `json:"type"`
	FSType     string        `json:"fstype"`
	Label      string        `json:"label"`
	UUID       string        `json:"uuid"`
	MountPoint string        `json:"mountpoint"`
	RM         flexBool      `json:"rm"`
	Hotplug    flexBool      `json:"hotplug"`
	Tran       string        `json:"tran"`
	Model      string        `json:"model"`
	FSUsed     flexInt       `json:"fsused"`
	FSAvail    flexInt       `json:"fsavail"`
	Children   []blockDevice `json:"children"`
}

func (b blockDevice) devPath() string {
	if b.Path != "" {
		return b.Path
	}
	return "/dev/" + b.Name
}

func (b blockDevice) volume() Volume {
	return Volume{
		Name: b.Name, Path: b.devPath(), Size: int64(b.Size), FSType: b.FSType, Label: b.Label, UUID: b.UUID,
		MountPoint: b.MountPoint, Used: int64(b.FSUsed), Avail: int64(b.FSAvail),
	}
}

// skippedDisk reports devices that aren't drives a person manages:
// loop devices, RAM disks, zram swap and eMMC boot areas.
func skippedDisk(b blockDevice) bool {
	if b.Type != "disk" {
		return true
	}
	for _, p := range []string{"zram", "ram", "mtdblock"} {
		if strings.HasPrefix(b.Name, p) {
			return true
		}
	}
	return strings.HasPrefix(b.Name, "mmcblk") && (strings.Contains(b.Name, "boot") || strings.Contains(b.Name, "rpmb"))
}

func parseLsblk(out string) ([]Drive, error) {
	var doc struct {
		Devices []blockDevice `json:"blockdevices"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		return nil, fmt.Errorf("reading lsblk: %w", err)
	}
	var drives []Drive
	for _, b := range doc.Devices {
		if skippedDisk(b) {
			continue
		}
		d := Drive{
			Name: b.Name, Path: b.devPath(), Size: int64(b.Size), Model: strings.TrimSpace(b.Model),
			Transport: b.Tran, Removable: bool(b.RM) || bool(b.Hotplug),
		}
		if len(b.Children) == 0 && (b.FSType != "" || b.MountPoint != "") {
			d.Volumes = append(d.Volumes, b.volume()) // a file system on the whole disk
		}
		var walk func([]blockDevice)
		walk = func(children []blockDevice) {
			for _, c := range children {
				if c.FSType != "" || c.Type == "part" {
					d.Volumes = append(d.Volumes, c.volume())
				}
				walk(c.Children) // LVM and LUKS volumes inside partitions
			}
		}
		walk(b.Children)
		drives = append(drives, d)
	}
	slices.SortStableFunc(drives, func(a, b Drive) int {
		switch {
		case a.External() && !b.External():
			return -1
		case b.External() && !a.External():
			return 1
		}
		return 0
	})
	return drives, nil
}
