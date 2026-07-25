// Package season runs the meta-game cycle, tied to the wave number
// (e.g. a new season every 5 waves) (GDD §6):
//
//   - Summer: standard rules, best time to expand and hoard Seeds.
//   - Autumn: water sources dry up faster; mushroom spores spawn actively.
//   - Winter: sources freeze (give no water); bugs get more armored.
//   - Spring: ice melts, rich new sources appear, but bug spawn count doubles.
//
// Filled in Stage 14. It applies modifiers to resources, world, and enemies.
package season
