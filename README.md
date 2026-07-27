# Grow a Tree

A meditative-yet-tense grid strategy: you play the mind of an ancient forest,
growing a network of glowing roots, mining water and energy, and fending off
endless waves of bugs. Full concept in [docs/GDD.md](docs/GDD.md).

**Stack:** Go + [Ebiten v2](https://ebitengine.org/). No engine, plain code.

## Run

```sh
go run ./cmd/game
```

## Build

```sh
go build -o bin/growtree.exe ./cmd/game
```

## Controls

- **LMB** — grow a root into an empty (or rotted) cell touching your network. Costs energy.
- **RMB** — cut your own root; refunds 50% of the energy and leaves rot behind.
- **1 / 2 / 3** — plant a companion near the Core (80 seeds): Battery / Moss / Thorn.
- **M** — mature a new Core (100 seeds); link cores with roots to raise the water/energy caps. Each extra Core also drains water faster and spawns one more source.
- **R** — restart after game over.

The game runs fullscreen: the square play field sits centered with a HUD panel
on each side (resource bars, cores/HP, controls, and the wave counter), so the
UI never covers the field.

## Credits

- Font: **Dogica** by Roberto Mocci, licensed under the SIL Open Font License 1.1
  (see [assets/dogica_license.txt](assets/dogica_license.txt)).

## Structure

```
cmd/game/            — entry point (main): opens the window, runs the game loop
internal/
  game/              — root Game: game loop, Update/Draw, wires subsystems together
  world/             — grid, cells, Core, roots, rot, sources, network connectivity
  resources/         — economy: Energy / Water / Seeds
  enemies/           — bugs: 5 levels, pathfinding, AI, combat
  waves/             — wave manager: prep timer, spawning, wave number
  plants/            — companion plants ("pets") and their auras
  season/            — seasons (summer/autumn/winter/spring) and their modifiers
  ui/                — HUD: bars, buttons, counters
  save/              — high-score persistence (save.json)
assets/              — placeholder / future textures (embedded via go:embed)
config/              — JSON balance data (enemies, plants, seasons, general)
docs/                — GDD
bin/                 — built .exe (git-ignored)
```

Each package under `internal/` currently ships a `doc.go` describing its
responsibility and the stage that fills it. Logic is added package by package as
the stages progress.

## Roadmap

- [x] Stage 0 — project skeleton (window, game loop, background)
- [x] Stage 1 — grid and root growth (LMB)
- [x] Stage 2 — resource economy
- [x] Stage 3 — cutting and rot (RMB)
- [x] Stage 4 — network logistics (flood-fill from the Core)
- [x] Stage 5 — water sources
- [x] Stage 6 — enemies: basic AI
- [x] Stage 7 — waves
- [x] Stage 8 — combat
- [x] Stage 9 — 5 bug levels
- [x] Stage 10 — mushrooms and symbiosis
- [x] Stage 11 — rot as bait
- [x] Stage 12 — companion plants
- [x] Stage 13 — multiple Cores
- [x] Stage 14 — seasons
- [ ] Stage 15 — full UI
- [ ] Stage 16 — high-score save
- [ ] Stage 17 — balance and polish
- [ ] Stage 18 — texture integration
