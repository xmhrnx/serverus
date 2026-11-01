// Package assets embeds the static files the app ships with: fonts,
// icons and the window icon.
package assets

import "embed"

// Files holds every embedded asset, rooted at this directory.
//
//go:embed fonts/*.ttf icons/*.svg icon.png
var Files embed.FS

// Icon is the window icon as PNG bytes.
//
//go:embed icon.png
var Icon []byte
