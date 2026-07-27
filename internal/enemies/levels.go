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
	size         float32 // draw size in field pixels
	col          color.RGBA
}

const maxLevel = 5

// levels is indexed by bug level 1..5 (index 0 is unused). Each level is a
// distinct color and a distinct size (20, 24, 28, 32, 36 px).
var levels = [maxLevel + 1]stats{
	1: {moveInterval: 0.26, gnawTime: 1.10, biteDamage: 3, size: 20, col: color.RGBA{0xff, 0xc8, 0x3c, 0xff}},
	2: {moveInterval: 0.30, gnawTime: 0.95, biteDamage: 4, size: 24, col: color.RGBA{0xff, 0x8c, 0x3c, 0xff}},
	3: {moveInterval: 0.34, gnawTime: 0.80, biteDamage: 6, size: 28, col: color.RGBA{0xff, 0x4c, 0x5e, 0xff}},
	4: {moveInterval: 0.40, gnawTime: 0.65, biteDamage: 9, size: 32, col: color.RGBA{0xd8, 0x3c, 0xb0, 0xff}},
	5: {moveInterval: 0.52, gnawTime: 0.50, biteDamage: 14, size: 36, col: color.RGBA{0x9a, 0x3c, 0xff, 0xff}},
}

// levelForWave picks a random bug level from those unlocked at this wave. A new
// level unlocks every 2 waves — waves 1-2 send level 1, waves 3-4 mix levels
// 1-2, ... waves 9-10 mix 1-5 — then it resets to level 1 and the 10-wave cycle
// repeats. Levels are always mixed within the unlocked range.
func levelForWave(wave int) int {
	if wave < 1 {
		wave = 1
	}
	top := ((wave-1)%(maxLevel*2))/2 + 1
	return 1 + rand.Intn(top)
}
