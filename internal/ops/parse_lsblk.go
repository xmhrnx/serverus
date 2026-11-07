package ops

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

func parseLsblk(out string) ([]Drive, error) {
	var doc struct {
		Devices []blockDevice `json:"blockdevices"`
	}
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		return nil, fmt.Errorf("reading lsblk: %w", err)
	}
	var drives []Drive
	for _, b := range doc.Devices {
		if skippedDisk(b) {
			continue
		}
		d := Drive{
			Name: b.Name, Path: b.devPath(), Size: int64(b.Size), Model: strings.TrimSpace(b.Model),
			Transport: b.Tran, Removable: bool(b.RM) || bool(b.Hotplug),
		}
		if len(b.Children) == 0 && (b.FSType != "" || b.MountPoint != "") {
			d.Volumes = append(d.Volumes, b.volume()) // a file system on the whole disk
		}
		var walk func([]blockDevice)
		walk = func(children []blockDevice) {
			for _, c := range children {
				if c.FSType != "" || c.Type == "part" {
					d.Volumes = append(d.Volumes, c.volume())
				}
				walk(c.Children) // LVM and LUKS volumes inside partitions
			}
		}
		walk(b.Children)
		drives = append(drives, d)
	}
	slices.SortStableFunc(drives, func(a, b Drive) int {
		switch {
		case a.External() && !b.External():
			return -1
		case b.External() && !a.External():
			return 1
		}
		return 0
	})
	return drives, nil
}
