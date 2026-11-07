package ops

import (
	"context"
	"errors"
	"path"
	"regexp"
	"strings"
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
