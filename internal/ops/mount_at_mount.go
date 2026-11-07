package ops

import (
	"context"
	"fmt"
	"strings"
)

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
