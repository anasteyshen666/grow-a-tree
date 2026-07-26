// Package assets embeds the game's image files so the built binary stays
// self-contained (a single .exe).
package assets

import "embed"

//go:embed *.png
var FS embed.FS
