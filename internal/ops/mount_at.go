package ops

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
)

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
