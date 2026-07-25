// Package waves drives the endless assault: the 15-second prep timer between
// waves, spawning bugs in batches, tracking the current wave number, and
// scaling difficulty with the wave number (GDD §7).
//
// Filled in Stage 7. The wave number also drives season changes (see package
// season) and the high-score record (see package save).
package waves
