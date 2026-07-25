# assets

Textures/sprites will live here once they are drawn (Stage 18).

For now the game is rendered with Ebiten primitives using the color coding from
the GDD (green — energy/roots, blue — water, yellow — seeds, brown — rot). When
sprites arrive, cell drawing switches to `ebiten.Image` while all game logic
stays untouched.

Planned assets:
- `core.png`      — the Core (2×2)
- `root.png`      — a root segment (+ a mushroom variant)
- `rot.png`       — rot
- `water.png`     — water source
- `bug_1..5.png`  — the 5 bug levels
- `plant_*.png`   — companion plants (battery flower, water moss, thornbush)
- `spore.png`     — mushroom spore

Assets will be embedded into the binary via `//go:embed` so the `.exe` stays
self-contained (a single file).
