// Package format turns numbers, sizes and times into the short strings
// the interface shows.
package format

import (
	"fmt"
	"io/fs"
	"math"
	"strings"
	"time"
	"unicode"
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

// Duration is an uptime-style duration: "3d 4h", "5h 12m", "48s".
func Duration(d time.Duration) string {
	d = d.Round(time.Second)
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh", days, hours)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, mins)
	case mins > 0:
		return fmt.Sprintf("%dm %ds", mins, int(d.Seconds())%60)
	}
	return fmt.Sprintf("%ds", int(d.Seconds()))
}

// Latency is a round trip: "23 ms".
func Latency(d time.Duration) string {
	if d <= 0 {
		return "—"
	}
	if d < time.Millisecond {
		return "<1 ms"
	}
	return fmt.Sprintf("%d ms", d.Milliseconds())
}

// Ago is a relative time: "just now", "5 min ago", "3 days ago", or the
// date for anything older than a month.
func Ago(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d min ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d h ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	}
	return t.Format("Jan 2, 2006")
}

// Timestamp is a full local date and time.
func Timestamp(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.Local().Format("2006-01-02 15:04")
}

// Mode is a file mode as ls prints it: "drwxr-xr-x".
func Mode(m fs.FileMode) string {
	s := []byte("----------")
	switch {
	case m&fs.ModeDir != 0:
		s[0] = 'd'
	case m&fs.ModeSymlink != 0:
		s[0] = 'l'
	}
	const rwx = "rwxrwxrwx"
	for i := 0; i < 9; i++ {
		if m&(1<<uint(8-i)) != 0 {
			s[i+1] = rwx[i]
		}
	}
	return string(s)
}

// Celsius is a temperature: "54°C".
func Celsius(c float64) string { return fmt.Sprintf("%.0f°C", c) }

// Sentence capitalizes the first letter of an error message and ends it
// with a period, to show it as a sentence.
func Sentence(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	s = string(r)
	if !strings.HasSuffix(s, ".") && !strings.HasSuffix(s, "!") && !strings.HasSuffix(s, "?") {
		s += "."
	}
	return s
}
