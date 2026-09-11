# Vendored Tela icon theme

Upstream: https://github.com/vinceliuice/Tela-icon-theme
Pinned commit: a1fffc5bfab716bd022dd228ee96fe3965cdb33d (2026-08-10)
License: GPL-3.0 (see COPYING) — attribution in AUTHORS.
Variant: standard (no recolor; upstream applies color `sed` only to non-standard variants).

## What's here
- `src/`   — upstream `src/` verbatim: size-major real SVGs (`<size>/<category>/<name>.svg`).
- `links/` — upstream `links/` verbatim: name-alias symlinks, size-major.
- `COPYING`, `AUTHORS` — license + attribution.

## What's omitted (not needed — the theme is generated)
`install.sh`, `elementary/`, `colorscheme/`, `release/`, preview PNGs.

## How it's used
`make theme` (the `tools/` program) transposes this tree from Tela's size-major
layout into `output/`'s category-major layout as the theme **base**, then
applies `overlays/**` on top. Do not edit files here; re-vendor from a new
pinned commit to update.
