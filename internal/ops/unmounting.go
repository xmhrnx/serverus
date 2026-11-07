package ops

import (
	"context"
	"fmt"
	"strings"
)

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
