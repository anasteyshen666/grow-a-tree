// Package plants holds the companion plants ("pets") the player seeds near a
// Core. They do not conduct water; instead each radiates an aura in a 2-3 cell
// radius (GDD §5):
//
//   - Цветок-Батарейка (yellow): +20% energy regen to roots in radius.
//   - Водяной Мох (blue):        -10% water cost of energy regen in radius.
//   - Зимний Терновник (white):  keeps nearby water sources from freezing in winter.
//
// Fills in Этап 12.
package plants
