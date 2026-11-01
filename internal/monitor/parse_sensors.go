package monitor

import (
	"strconv"
	"strings"
)

// parseSensors reads the hwmon and thermal sections. hwmon is
// preferred: it names chips and reports thresholds. Thermal zones fill
// in on machines without hwmon temperatures (many ARM boards).
func parseSensors(hwmon, thermal []string) []Sensor {
	var out []Sensor
	for _, line := range hwmon {
		f := strings.Split(line, "|")
		if len(f) < 3 {
			continue
		}
		c, ok := milliCelsius(f[2])
		if !ok {
			continue
		}
		name := strings.TrimSpace(f[1])
		chip := strings.TrimSpace(f[0])
		switch {
		case name == "":
			name = chip
		case chip != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(chip)):
			name = chip + " " + name
		}
		s := Sensor{Name: name, Celsius: c}
		if len(f) > 3 {
			s.High, _ = milliCelsius(f[3])
		}
		if len(f) > 4 {
			s.Critical, _ = milliCelsius(f[4])
		}
		out = append(out, s)
	}
	if len(out) > 0 {
		return out
	}
	for _, line := range thermal {
		name, val, ok := strings.Cut(line, "|")
		if !ok {
			continue
		}
		if c, ok := milliCelsius(val); ok {
			out = append(out, Sensor{Name: strings.TrimSpace(name), Celsius: c})
		}
	}
	return out
}

// milliCelsius parses a sysfs temperature in thousandths of a degree,
// rejecting the bogus values some drivers report.
func milliCelsius(s string) (float64, bool) {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, false
	}
	c := v / 1000
	if c <= -40 || c >= 200 {
		return 0, false
	}
	return c, true
}
