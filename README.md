# telementary-icons

A custom KDE / freedesktop **icon theme**, built on the vendored **Tela** base
with hand-designed **elementary**-style overlays on top. It installs under the
theme id `telementary` and `Inherits=Breeze` as a final fallback. Packaged as a
single Debian package (`telementary-icons`) and published as a GitHub release.

## Build

Requires `go` and `nfpm` on PATH.

- `make theme` — regenerate `output/` from the vendored Tela base + `overlays/`
  (wipe → lay Tela → apply overlays → generate `index.theme`, `@2x`/`@3x`, aliases).
- `make build` — regenerate the theme, then package it into
  `telementary-icons_<version>_all.deb` via `nfpm`. Version comes from
  `git describe`; override with `make build VERSION=1.2.3`.
- `make view` — serve a browsable name×size matrix of the theme (dev viewer).
- `make clean` — remove `*.deb` and the generated `output/`.
- `make help` — list targets.

## Release

Tag-driven. `make tag version="vX.Y.Z"` (only from a clean `main`) pushes an
annotated tag, which triggers `.github/workflows/release.yml`: it sets up Go,
installs `nfpm`, runs `make build`, and publishes a GitHub release with the
`.deb` attached. The tag is the single source of truth for the version.

## Layout

| Path | What |
|---|---|
| `overlays/`   | hand-designed icon sources + `overlay.yaml` mappings (committed) |
| `tela-icons/` | vendored Tela base theme, GPL-3.0, pinned upstream (committed) |
| `output/`     | generated freedesktop theme — **gitignored**, never edit or commit |
| `tools/`      | Go generator that builds `output/` from the two sources above |
| `nfpm.yaml`   | single source of truth for what the `.deb` ships and where |

`output/` is what installs to `/usr/share/icons/telementary`. Regenerate it with
`make theme` after editing any source; never hand-edit `output/`.

See `CLAUDE.md` for the design/overlay workflow in detail.
