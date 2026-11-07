package format

import (
	"fmt"
	"io/fs"
	"strings"
	"unicode"
)

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
