package theme

import (
	"fmt"
	"io/fs"

	"gioui.org/font"
	"gioui.org/font/opentype"

	"github.com/xmhrnx/giode"
)

// Font families.
const (
	// Sans is the interface face.
	Sans = "Inter"
	// Mono is for terminals, paths, code and numbers that must line up.
	Mono = "JetBrains Mono"
)

// fontFiles maps each embedded file to the face it provides.
var fontFiles = []struct {
	file   string
	family string
	weight font.Weight
	style  font.Style
}{
	{"fonts/Inter-Regular.ttf", Sans, font.Normal, font.Regular},
	{"fonts/Inter-Medium.ttf", Sans, font.Medium, font.Regular},
	{"fonts/Inter-SemiBold.ttf", Sans, font.SemiBold, font.Regular},
	{"fonts/Inter-Bold.ttf", Sans, font.Bold, font.Regular},
	{"fonts/JetBrainsMonoNL-Regular.ttf", Mono, font.Normal, font.Regular},
	{"fonts/JetBrainsMonoNL-Bold.ttf", Mono, font.Bold, font.Regular},
	{"fonts/JetBrainsMonoNL-Italic.ttf", Mono, font.Normal, font.Italic},
}

// LoadFonts registers the embedded faces with giode. Inter comes first,
// so text without a family uses it.
func LoadFonts(fsys fs.FS) error {
	faces := make([]font.FontFace, 0, len(fontFiles))
	for _, f := range fontFiles {
		data, err := fs.ReadFile(fsys, f.file)
		if err != nil {
			return fmt.Errorf("load %s: %w", f.file, err)
		}
		face, err := opentype.Parse(data)
		if err != nil {
			return fmt.Errorf("parse %s: %w", f.file, err)
		}
		faces = append(faces, font.FontFace{
			Font: font.Font{Typeface: font.Typeface(f.family), Weight: f.weight, Style: f.style},
			Face: face,
		})
	}
	giode.RegisterFonts(faces)
	return nil
}
