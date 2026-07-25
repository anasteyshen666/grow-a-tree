// Package resources models the game economy: the three balances the whole
// game is built around (GDD §2).
//
//   - Energy — spent to grow roots, plant companions, activate defense.
//     Regenerates over time, but only while there is Water.
//   - Water  — spent to regenerate Energy. Mined from sources that deplete.
//   - Seeds  — accumulate over time. Spent on companion plants and new Cores.
//
// Fills in Этап 2 (energy/water/seeds tick) and is extended by later stages
// (sources in Этап 5, plant auras in Этап 12, seasons in Этап 14).
package resources
