// Package resources models the game economy: the three balances the whole
// game is built around (GDD §2).
//
//   - Energy — spent to grow roots, plant companions, activate defense.
//     Regenerates over time, but only while there is Water.
//   - Water  — spent to regenerate Energy. Mined from sources that deplete.
//   - Seeds  — accumulate over time. Spent on companion plants and new Cores.
//
// Filled in Stage 2 (energy/water/seeds tick) and extended by later stages
// (sources in Stage 5, plant auras in Stage 12, seasons in Stage 14).
package resources
