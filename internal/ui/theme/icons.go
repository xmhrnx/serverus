package theme

import (
	"io/fs"

	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/icon"
	"github.com/mehran9675/giode/styles/properties"
)

// LoadIcons registers the embedded Lucide icons by file name
// ("icons/server.svg" is "server").
func LoadIcons(fsys fs.FS) error {
	return giode.RegisterIconsFS(fsys, "icons")
}

// Icon draws a registered icon at size dp in color c ("" inherits the
// surrounding text color).
func Icon(name string, size int, c giode.Color) giode.Element {
	return giode.Icon(name, giode.Styles{Width: properties.Width(size), Height: properties.Height(size), Color: c})
}

// IconSource returns a registered icon as an icon.Source, for menus.
func IconSource(name string) icon.Source {
	src, _ := icon.Lookup(name)
	return src
}
