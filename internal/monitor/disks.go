package monitor

// DiskUsage is one mounted file system, in bytes.
type DiskUsage struct {
	Device, FSType, Mount string
	Total, Used, Free     uint64
}

// DiskCounter is the bytes a block device has read and written.
type DiskCounter struct {
	Device                string
	ReadBytes, WriteBytes uint64
}

// NetCounter is the bytes an interface has received and sent.
type NetCounter struct {
	Interface        string
	RxBytes, TxBytes uint64
}

// Percent is the share of the file system in use, 0-100.
func (d DiskUsage) Percent() float64 { return percent(float64(d.Used), float64(d.Used+d.Free)) }
