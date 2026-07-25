# config

JSON files with balance data, so the game can be tuned without recompiling logic:

- `enemies.json` — bug stats per level 1–5 (speed, HP, damage, armor).
- `plants.json`  — companion plants (seed cost, radius, aura strength).
- `seasons.json` — season modifiers (water drying, spore spawn, freezing,
  bug armor, spawn multiplier).
- `balance.json` — general numbers: root growth cost, energy regen rate,
  water cost, seed accrual, prep phase length, etc.

These files appear as the matching stages are built (9, 12, 14, and general
balance in Stage 17). Default values are mirrored in code so the game runs even
without external configs.
