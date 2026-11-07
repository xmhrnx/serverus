package ops

import (
	"context"
	"strconv"
	"strings"
)

// mountPointsIn reads the mount table at table; tests pass their own.
func (d Drives) mountPointsIn(ctx context.Context, dev, table string) ([]string, error) {
	// The loop's status is its last test's, false whenever the table's
	// last line is another device, so the script ends with exit 0.
	script := `while read -r src dir rest; do
  case "$src" in /dev/*) if [ "$(readlink -f "$src" 2>/dev/null)" = ` + Quote(dev) + ` ]; then printf '%s\n' "$dir"; fi;; esac
done < ` + Quote(table) + `
exit 0`
	out, err := run(ctx, d.Exec, "sh -c "+Quote(script))
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, line := range strings.Split(out, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			dirs = append(dirs, unescapeMount(line))
		}
	}
	return dirs, nil
}

// unescapeMount decodes the octal escapes /proc/mounts uses for spaces,
// tabs and backslashes in folder names.
func unescapeMount(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+3 < len(s) {
			if n, err := strconv.ParseUint(s[i+1:i+4], 8, 8); err == nil {
				b.WriteByte(byte(n))
				i += 3
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
