package theme

import (
	"fmt"
	"io/fs"

	"gioui.org/font"
	"gioui.org/font/opentype"

	"github.com/xmhrnx/giode"
)

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
