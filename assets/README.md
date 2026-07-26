# assets

PNG textures, embedded into the binary via `//go:embed` (see `assets.go`) so the
`.exe` stays self-contained. All tiles are 24×24 (one cell).

## Present (wired up)

| File | Used for |
|------|----------|
| `dirt.png`        | terrain / cell background (drawn darkened) |
| `common koren.png`| normal root |
| `koren grib.png`  | mushroom-infected root |
| `core.png`        | Core (drawn on each of the 2×2 core cells) |
| `gnil.png`        | rot (fades out as it ages) |
| `water.png`       | water source (fades as it depletes) |
| `spora.png`       | mushroom spore |

## Still needed (currently drawn as colored squares)

- `bug_1.png` … `bug_5.png` — the 5 bug levels (draw at 24×24; sizing/centering handled in code)
- companion plants — Battery / Moss / Thorn
- a "dead root" look for roots cut off from the Core (optional; now a dim gray square)

Any square, transparent-background PNG works. To keep sprites crisp on larger
windows you can draw them at 2× (48×48) or 4× (96×96) — code scales them to the
24px cell.
