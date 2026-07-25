// Package enemies contains the bugs: their state, the 5 difficulty levels,
// pathfinding across the grid, and the AI that chooses targets — roots, the
// Core, or (preferentially) fresh Rot used as bait (GDD §4, §7).
//
// Stages that fill this package:
//   - Stage 6: a single bug type — spawn, path to nearest root/Core, gnaw.
//   - Stage 8: combat (roots damage bugs standing on them), death.
//   - Stage 9: 5 levels balanced from config (speed/HP/damage/armor).
//   - Stage 11: Rot as bait — retarget toward rot, poison on eating it.
package enemies
