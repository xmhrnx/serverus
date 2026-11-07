package ops

import (
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

// Drives lists and mounts the disks of a Linux machine.
type Drives struct {
	Exec transport.Executor
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
