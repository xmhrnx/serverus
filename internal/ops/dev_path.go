package ops

func (b blockDevice) devPath() string {
	if b.Path != "" {
		return b.Path
	}
	return "/dev/" + b.Name
}

func (b blockDevice) volume() Volume {
	return Volume{
		Name: b.Name, Path: b.devPath(), Size: int64(b.Size), FSType: b.FSType, Label: b.Label, UUID: b.UUID,
		MountPoint: b.MountPoint, Used: int64(b.FSUsed), Avail: int64(b.FSAvail),
	}
}
