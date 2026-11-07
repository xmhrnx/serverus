package ops

import (
	"slices"
)

// System reports whether the operating system runs from the drive.
func (d Drive) System() bool {
	return slices.ContainsFunc(d.Volumes, func(v Volume) bool { return v.System() })
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
