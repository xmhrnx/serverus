package drives

import (
	"strings"
)

// transportName is how the drive is attached, for people.
func transportName(t string) string {
	switch t {
	case "usb":
		return "USB"
	case "sata", "ata":
		return "SATA"
	case "nvme":
		return "NVMe"
	case "mmc":
		return "Memory card"
	case "sas", "scsi":
		return strings.ToUpper(t)
	case "virtio":
		return "Virtual disk"
	}
	return ""
}
