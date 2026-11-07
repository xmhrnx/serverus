package ops

import (
	"strconv"
	"strings"
)

// lsblk prints numbers and booleans as JSON values or as strings,
// depending on its version.
type flexInt int64

type flexBool bool

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

func (f *flexBool) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	*f = flexBool(s == "true" || s == "1")
	return nil
}
