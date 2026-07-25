// Package world holds the playing field: the grid of cells, the Core (2x2),
// roots the player grows, rot tiles, and the water sources.
//
// It also owns the network connectivity logic — the flood-fill from each Core
// that decides which roots are "alive" (conducting water/energy) and which have
// been cut off. This is the strategic heart of the game (GDD §3, §4).
//
// Stages that fill this package:
//   - Stage 1: Grid, Cell, Root, Core, growth on left-click.
//   - Stage 3: cutting roots (right-click) and Rot tiles.
//   - Stage 4: flood-fill connectivity from the Core.
//   - Stage 5: water sources (spawn, depletion).
package world
