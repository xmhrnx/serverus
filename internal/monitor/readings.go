package monitor

// HottestSensor returns the warmest sensor, ok false when there is none.
func (s Snapshot) HottestSensor() (Sensor, bool) {
	if len(s.Sensors) == 0 {
		return Sensor{}, false
	}
	best := s.Sensors[0]
	for _, x := range s.Sensors[1:] {
		if x.Celsius > best.Celsius {
			best = x
		}
	}
	return best, true
}

// RootDisk returns the file system mounted at "/" (or the largest one
// when there's no root, as on Windows).
func (s Snapshot) RootDisk() (DiskUsage, bool) {
	var best DiskUsage
	found := false
	for _, d := range s.Disks {
		if d.Mount == "/" {
			return d, true
		}
		if !found || d.Total > best.Total {
			best, found = d, true
		}
	}
	return best, found
}
