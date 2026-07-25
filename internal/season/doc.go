// Package season runs the meta-game cycle, tied to the wave number
// (e.g. a new season every 5 waves) (GDD §6):
//
//   - Лето:   standard rules, best time to expand and hoard Seeds.
//   - Осень:  water sources dry up faster; mushroom spores spawn actively.
//   - Зима:   sources freeze (give no water); bugs get more armored.
//   - Весна:  ice melts, rich new sources appear, but bug spawn count doubles.
//
// Fills in Этап 14. It applies modifiers to resources, world, and enemies.
package season
