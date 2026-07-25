package enemies

import (
	"image/color"
	"math/rand"
)

// stats defines a bug level's behavior and look (GDD §7): levels 1-2 are fast
// and weak, 3-4 are tougher and gnaw faster, 5 is the slow, armored boss.
type stats struct {
	hp           float64
	moveInterval float64 // seconds per cell (higher = slower)
	gnawTime     float64 // seconds to chew through a root
	biteDamage   int     // Core HP lost per bite
	armor        float64 // fraction of network damage ignored
	size         float32 // draw size as a fraction of a cell
	col          color.RGBA
}

const maxLevel = 5

// levels is indexed by bug level 1..5 (index 0 is unused).
var levels = [maxLevel + 1]stats{
	1: {hp: 18, moveInterval: 0.26, gnawTime: 1.10, biteDamage: 3, armor: 0.00, size: 0.55, col: color.RGBA{0xff, 0x9a, 0x8c, 0xff}},
	2: {hp: 30, moveInterval: 0.30, gnawTime: 0.95, biteDamage: 4, armor: 0.05, size: 0.62, col: color.RGBA{0xff, 0x6e, 0x5e, 0xff}},
	3: {hp: 48, moveInterval: 0.34, gnawTime: 0.80, biteDamage: 6, armor: 0.15, size: 0.70, col: color.RGBA{0xff, 0x4c, 0x5e, 0xff}},
	4: {hp: 72, moveInterval: 0.40, gnawTime: 0.65, biteDamage: 9, armor: 0.28, size: 0.80, col: color.RGBA{0xd8, 0x2f, 0x5a, 0xff}},
	5: {hp: 150, moveInterval: 0.52, gnawTime: 0.50, biteDamage: 14, armor: 0.42, size: 0.96, col: color.RGBA{0x9a, 0x2f, 0x7a, 0xff}},
}

// levelForWave picks a bug level suited to the wave: early waves send weak bugs,
// later waves unlock tougher ones up to the boss.
func levelForWave(wave int) int {
	top := wave/2 + 1
	if top > maxLevel {
		top = maxLevel
	}
	return 1 + rand.Intn(top)
}
