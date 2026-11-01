package theme

import (
	"bytes"
	"image"
	_ "image/png" // the app icon is a PNG

	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"
)

var logo image.Image

// LoadLogo decodes the app icon for use inside the interface.
func LoadLogo(png []byte) error {
	img, _, err := image.Decode(bytes.NewReader(png))
	logo = img
	return err
}

// Logo draws the app icon at size dp.
func Logo(size int) giode.Element {
	if logo == nil {
		return giode.Box(giode.Styles{Width: properties.Width(size), Height: properties.Height(size)})
	}
	return giode.Image(logo, giode.Styles{Width: properties.Width(size), Height: properties.Height(size), Fit: giode.FitContain})
}
