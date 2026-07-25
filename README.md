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
- [ ] Stage 4 — network logistics (flood-fill from the Core)
- [ ] Stage 5 — water sources
- [ ] Stage 6 — enemies: basic AI
- [ ] Stage 7 — waves
- [ ] Stage 8 — combat
- [ ] Stage 9 — 5 bug levels
- [ ] Stage 10 — mushrooms and symbiosis
- [ ] Stage 11 — rot as bait
- [ ] Stage 12 — companion plants
- [ ] Stage 13 — multiple Cores
- [ ] Stage 14 — seasons
- [ ] Stage 15 — full UI
- [ ] Stage 16 — high-score save
- [ ] Stage 17 — balance and polish
- [ ] Stage 18 — texture integration
