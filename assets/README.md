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
| `fly1.png`…`fly5.png` | the 5 bug levels (20/24/28/32/36 px), drawn facing up, rotated toward movement |
| `pet1.png` / `pet2.png` / `pet3.png` | companion plants Battery / Moss / Thorn (24×24) |

A root cut off from the Core is drawn as the normal root texture, darkened
(no separate asset needed). All game art is now in place.

Any square, transparent-background PNG works. To keep sprites crisp on larger
windows you can draw them at 2× (48×48) or 4× (96×96) — code scales them to the
24px cell.
