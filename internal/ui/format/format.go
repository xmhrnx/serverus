// Package format turns numbers, sizes and times into the short strings
// the interface shows.
package format

import (
	"fmt"
	"math"
)

// Bytes is a size in binary units: "512 B", "1.4 GiB".
func Bytes(n uint64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := uint64(unit), 0
	for m := n / unit; m >= unit && exp < 5; m /= unit {
		div *= unit
		exp++
	}
	v := float64(n) / float64(div)
	if v >= 100 {
		return fmt.Sprintf("%.0f %ciB", v, "KMGTPE"[exp])
	}
	return fmt.Sprintf("%.1f %ciB", v, "KMGTPE"[exp])
}

// Size is Bytes for a signed size; negative sizes are unknown.
func Size(n int64) string {
	if n < 0 {
		return "—"
	}
	return Bytes(uint64(n))
}

// Rate is a transfer speed: "3.2 MiB/s".
func Rate(bytesPerSec float64) string {
	if bytesPerSec < 0 || math.IsNaN(bytesPerSec) {
		bytesPerSec = 0
	}
	return Bytes(uint64(bytesPerSec)) + "/s"
}

// Percent is a percentage without decimals above 10: "7.5%", "42%".
func Percent(p float64) string {
	if p < 10 && p > 0 {
		return fmt.Sprintf("%.1f%%", p)
	}
	return fmt.Sprintf("%.0f%%", p)
}
