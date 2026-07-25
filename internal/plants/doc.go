// Package plants holds the companion plants ("pets") the player seeds near a
// Core. They do not conduct water; instead each radiates an aura in a 2-3 cell
// radius (GDD §5):
//
//   - Battery Flower (yellow): +20% energy regen to roots in radius.
//   - Water Moss (blue):       -10% water cost of energy regen in radius.
//   - Winter Thornbush (white): keeps nearby water sources from freezing in winter.
//
// Filled in Stage 12.
package plants
