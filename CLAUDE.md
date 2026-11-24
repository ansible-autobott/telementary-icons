# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
code in this repository.

## What this repo is

This is **not a software project** — it is a KDE/freedesktop **icon theme**
delivered as a single Debian package named `telementary-icons` and published as
a GitHub release.

The theme is built on a vendored copy of the **Tela** icon theme (GPL-3.0) with
hand-designed **elementary**-style overlays applied on top; `Inherits=Breeze`
remains only as a final fallback. It installs under the theme id **`telementary`**
(`/usr/share/icons/telementary`, `Name=Telementary`).

## Build / release

Requires `go` and `nfpm` on PATH.

- `make theme` — regenerate `output/` from the vendored Tela base + `overlays/`.
- `make build` — regenerate the theme, then package it into
  `telementary-icons_<version>_all.deb` via `nfpm`. Version comes from `VERSION`
  (defaults to `git describe`); override with `make build VERSION=1.2.3`.
- `make view` — serve a browsable view of the theme for design work.
- `make review` — the same server; open `/review` to approve/reject the
  elementary overlay candidates (see `review/` below). Needs the elementary
  source theme: `ELEMENTARY=` (default `~/.datos/edit/design/icons/elementary`).
- `make clean` — remove `*.deb` and the generated `output/`.
- `make tag version="v1.2.3"` — create and push an annotated git tag (only from a
  clean `main`), which triggers the release workflow.
- `make help` — list targets.

**Releasing is tag-driven** (`.github/workflows/release.yml`): pushing a
`v*.*.*` tag runs the `Release` workflow, which sets up Go, installs `nfpm`,
runs `make build`, and publishes a GitHub release with the `.deb` attached. The
tag is the single source of truth for the version — `nfpm.yaml` has
`version: "${VERSION}"`, which the workflow sets from the tag name (nfpm strips
the leading `v`). No local publish step and no `GH_TOKEN` (the workflow uses the
built-in token). Tags follow `v0.1.x`.

There is no test or lint step. Validation = the `.deb` builds and installs.

## Architecture / structure

**`nfpm.yaml` `contents:` is the single source of truth for what ships and where
it installs.** Adding a new packaged asset means adding a `src → dst` entry
there. The mapping today:

| Source | Installs to |
|---|---|
| `output/` (generated) | `/usr/share/icons/telementary` |
| `tela-icons/COPYING`, `tela-icons/AUTHORS` | `/usr/share/doc/telementary-icons/` |

### `output/` — the generated theme (gitignored)

The shipped freedesktop theme, laid out as `output/<category>/<size>/<name>.svg`
(e.g. `actions/16`, `apps/symbolic`, `places/64`). `<category>@2x` directories
are symlinks to their base category. `output/index.theme` is the freedesktop
manifest: its `Directories=` line enumerates every `category/size`, and each
`[category/size]` section declares `Size`/`Context`/`MinSize`/`MaxSize`/`Type`.
`Inherits=Breeze`.

**`output/` is generated and gitignored** — never edit or commit files there.
`index.theme`, the `@2x`/`@3x` symlinks, and name aliases are all generated.
Regenerate the whole theme with `make theme` (the Go tool in `tools/`). Adding a
new size bucket or category means teaching that tool's size/context tables, then
rerunning it — not just dropping files in.

### `tela-icons/` — vendored base theme (committed)

The theme is built on top of a **vendored copy of the Tela icon theme**
(`tela-icons/`, GPL-3.0, pinned upstream commit — see `tela-icons/README.md`).
`make theme` lays Tela down first (transposing its size-major `src/`+`links/`
into the theme's category-major layout), then applies `overlays/**` on top. Tela
is the substrate; `Inherits=Breeze` remains only as a final fallback.
`tela-icons/` is committed source; do not edit it — re-vendor to update.

### `overlays/` — hand-designed icon source (committed)

`overlays/` layers on top of the vendored Tela base; anything not overlaid comes
from Tela, and only what neither provides falls through to `Inherits=Breeze`.
The overlay source is hand-designed, committed icon content — alongside
`tela-icons/`, the two together are the committed inputs `output/` is generated
from.

Design sets live under `overlays/<category>/<set>/` — grouped by freedesktop
category (`apps/`, `mimetypes/`, `places/`, `actions/`, …) for navigation, e.g.
`overlays/apps/smb4k/` or `overlays/places/folder/`. Each set folder holds the
source SVGs (plus reference images) and an `overlay.yaml` declaring how they map
into the theme. The grouping is organizational only — the target context is each
link's own `category` field, not the parent folder — and the theme-generation
tool discovers `overlay.yaml` files at any depth. Name aliases live in
`overlays/aliases.txt`.

The workflow:

1. **Design** — draw new/replacement icons in `overlays/<category>/<set>/`. A
   single source SVG is meant to be reused across many sizes (and sometimes
   several icon names). `overlay.yaml` in the set folder captures that mapping:
   per entry, one `src` SVG → a `category`, one or more `names`, and the list of
   `sizes`. For monochrome icons, use the `ColorScheme-Text` pattern (see
   `overlays/README.md`).
2. **Regenerate theme** — `make theme` wipes `output/`, lays down the vendored
   Tela base, applies all overlays, then generates `index.theme` and the
   `@2x`/`@3x`/alias symlinks. After editing a design SVG, re-run `make theme`.
   Then `make build` to package.

`overlay.yaml` schema — see `overlays/apps/smb4k/overlay.yaml` for a worked
example. Only list sizes a category actually has (e.g. `apps` has no `22`; use
`symbolic` for the symbolic bucket).

### `review/` — elementary overlay review (committed)

`review/candidates.yaml` lists proposed overlays: action icons the theme shows
as monochrome glyphs (Tela's, or Breeze's when the theme lacks the name) that
elementary can replace with full-color art. Each item is written as
overlay.yaml-style links, with `from` relative to the elementary theme root.
It is a one-off generated snapshot (its header records how it was built; the
generator is not in the repo). Edit it by hand to add, drop or re-map items:
the review page re-reads it on every request.

The review page (`make review`, then `/review`) walks the candidates one at a
time: for each size, what KDE resolves today next to the proposed art, on light
and dark backgrounds, plus the Tela aliases that would follow. Each approve or
reject (with an optional comment) is written to `review/decisions.md`, a
todo-app file with `## Approved` / `## Rejected` tasks. The page rewrites those
two categories on every decision and regenerates each task's proposal lines;
only the checkbox state, the `Comment:` and any other lines (sub-tasks, notes)
are kept.

**Implementing a decision** = adding the overlay set under `overlays/` that the
task describes. A comment overrides the proposal. Mark the task `[x]` when done.
Implemented and rejected items are then dropped from `candidates.yaml`; their
tasks stay in `decisions.md` verbatim (the page keeps tasks it has no candidate
for).

### `staged/` — overlay sets being reworked (committed)

Same layout and `overlay.yaml` format as `overlays/`, but **not applied**: the
generator only reads `overlays/`, so nothing here ships. It holds rejected
review items (`staged/actions/<group>/`, one set per review group, each block
carrying the review comment) prepared like the shipped sets, for hand editing.
To ship one, move its block and SVGs into the matching `overlays/` set and run
`make theme`. Validate the staged sets with a dry run:
`cd tools && go run . overlay -overlays ../staged`.

### `tools/` — the Go generator

The `tools/` directory holds the Go program (`module icons`) that builds the
theme. Subcommands (run via the make targets):

- `theme` — wipe `output/`, lay down Tela, apply overlays, then `gen-theme`.
- `overlay` — apply just the design overlays (dry-run unless `-apply`).
- `gen-theme` — regenerate `index.theme`, `@2x`/`@3x` scaled dirs, and aliases.
- `view` — serve a browsable name×size matrix (with a base theme as reference),
  plus the overlay review page at `/review`.

Default paths are anchored at the repo root (found by walking up to `.git`), so
the tool works from any cwd: `output/`, `tela-icons/`, `overlays/`.

## Editing icons

Use Inkscape and follow the elementary
[Icon Design Guidelines](https://elementary.io/docs/human-interface-guidelines#iconography).
"Vacuum Defs" in Inkscape to keep SVGs lean before committing.

To turn a flat or symbolic shape into the beveled gray `go-first` icon type
(as done for `go-home-16/24` and the view-modes icons), use the `bevel-icon`
skill (`.claude/skills/bevel-icon/`, script `bevel.py`).
