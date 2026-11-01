package monitor

import "strings"

// splitSections splits the probe's output into its "@@name" sections,
// each a list of lines.
func splitSections(out string) map[string][]string {
	sections := map[string][]string{}
	current := ""
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimRight(line, "\r")
		if name, ok := strings.CutPrefix(line, "@@"); ok {
			current = strings.TrimSpace(name)
			sections[current] = nil
			continue
		}
		if current != "" && line != "" {
			sections[current] = append(sections[current], line)
		}
	}
	return sections
}
