package ops

import (
	"context"
	"slices"
	"strings"
)

// DriveUser is a program keeping a mounted drive busy.
type DriveUser struct {
	PID  int
	Name string
	Path string // the open file or working folder on the drive
}

// BusyError is returned by Unmount when programs still use the volume.
type BusyError struct {
	Volume  Volume
	Targets []string    // where it's still mounted
	Users   []DriveUser // the programs found; may be empty
}

func (e *BusyError) Error() string {
	if len(e.Users) == 0 {
		return e.Volume.Title() + " is in use"
	}
	names := []string{}
	for _, u := range e.Users {
		if !slices.Contains(names, u.Name) {
			names = append(names, u.Name)
		}
	}
	return e.Volume.Title() + " is in use by " + strings.Join(names, ", ")
}

func isBusy(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "busy")
}

func (d Drives) busyError(ctx context.Context, v Volume, targets []string) error {
	return &BusyError{Volume: v, Targets: targets, Users: d.users(ctx, targets)}
}
