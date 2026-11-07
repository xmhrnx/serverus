package ops

import (
	"context"
	"fmt"
)

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
