package enemies

import (
	"image/color"
	"math/rand"
)

// stats defines a bug level's behavior and look (GDD §7): levels 1-2 are fast
// and weak, 3-4 are tougher and gnaw faster, 5 is the slow boss. A bug is only
// killed by rot bait, and a level-N bug must eat N rot tiles to die.
type stats struct {
	moveInterval float64 // seconds per cell (higher = slower)
	gnawTime     float64 // seconds to chew through a root
	biteDamage   int     // Core HP lost per bite
	size         float32 // draw size as a fraction of a cell
	col          color.RGBA
}

const maxLevel = 5

// levels is indexed by bug level 1..5 (index 0 is unused).
var levels = [maxLevel + 1]stats{
	1: {moveInterval: 0.26, gnawTime: 1.10, biteDamage: 3, size: 0.55, col: color.RGBA{0xff, 0x9a, 0x8c, 0xff}},
	2: {moveInterval: 0.30, gnawTime: 0.95, biteDamage: 4, size: 0.62, col: color.RGBA{0xff, 0x6e, 0x5e, 0xff}},
	3: {moveInterval: 0.34, gnawTime: 0.80, biteDamage: 6, size: 0.70, col: color.RGBA{0xff, 0x4c, 0x5e, 0xff}},
	4: {moveInterval: 0.40, gnawTime: 0.65, biteDamage: 9, size: 0.80, col: color.RGBA{0xd8, 0x2f, 0x5a, 0xff}},
	5: {moveInterval: 0.52, gnawTime: 0.50, biteDamage: 14, size: 0.96, col: color.RGBA{0x9a, 0x2f, 0x7a, 0xff}},
}

// levelForWave picks a bug level suited to the wave. A new level unlocks every
// 5 waves: level 1 from wave 1, level 2 from wave 5, level 3 from wave 10, ...
func levelForWave(wave int) int {
	top := wave/5 + 1
	if top > maxLevel {
		top = maxLevel
	}
	return 1 + rand.Intn(top)
}
