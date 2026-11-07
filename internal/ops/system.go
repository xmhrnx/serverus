package ops

import (
	"errors"
	"strings"
)

// System reports whether the operating system needs the volume mounted.
func (v Volume) System() bool {
	switch v.MountPoint {
	case "/", "/boot", "/boot/efi", "/boot/firmware", "/usr", "/var", "/home", "[SWAP]":
		return true
	}
	return strings.HasPrefix(v.MountPoint, "/snap/")
}

// ErrNoLsblk is returned when the machine has no lsblk (not Linux, or a
// very small one).
var ErrNoLsblk = errors.New("drive management needs lsblk, part of util-linux")

const (
	lsblkColumns    = "NAME,PATH,SIZE,TYPE,FSTYPE,LABEL,UUID,MOUNTPOINT,RM,HOTPLUG,TRAN,MODEL,FSUSED,FSAVAIL"
	lsblkOldColumns = "NAME,SIZE,TYPE,FSTYPE,LABEL,UUID,MOUNTPOINT,RM,HOTPLUG,TRAN,MODEL"
)
