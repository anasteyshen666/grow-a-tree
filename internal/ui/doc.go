// Package ui draws the HUD and reads UI input (GDD §9):
//
//   - Top-left:   Energy / Water / Seeds bars with numbers.
//   - Top-right:  current Wave, Season icon, and the high-score Record.
//   - Bottom-center: the "Созреть (Seed)" button and the countdown to the
//     next bug attack.
//
// Controls are mouse-first (LMB grow, RMB cut/rot) plus a couple of hotkeys
// (Space pause, R restart after death). Fills in Этап 15.
package ui
