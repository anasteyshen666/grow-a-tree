// Package assets embeds the game's image and font files so the built binary
// stays self-contained (a single .exe).
//
// The Dogica font (dogica*.ttf) by Roberto Mocci is licensed under the SIL Open
// Font License 1.1 — see dogica_license.txt.
package assets

import "embed"

//go:embed *.png *.ttf dogica_license.txt
var FS embed.FS
