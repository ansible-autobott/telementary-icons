# overlays/ — icon theme source

This directory contains the hand-designed **overlay** source for the
telementary theme, applied on top of the vendored Tela base
(`tela-icons/`, see its `README.md`). The theme itself (`output/`) is
**generated** and gitignored — never edit or commit files there.

## Layout

```
overlays/
  aliases.txt                     # name-alias map (edit this)
  <category>/<set>/
    *.svg                         # source SVGs
    overlay.yaml                  # mapping: which names/sizes each SVG becomes
```

- **Categories**: `actions`, `apps`, `categories`, `devices`, `emblems`, `emotes`,
  `mimetypes`, `places`, `status`, etc. — the freedesktop icon-theme contexts.
  Each set directory sits under a category for organizational purposes; the actual
  target context comes from the `category` field in each overlay entry.
- **Sets**: a logical grouping of related icons (e.g., `apps/smb4k/`, `places/folder/`).
  Each set contains source SVGs plus an `overlay.yaml` that declares how they map
  into the theme.
- **aliases.txt**: defines name aliases (one line per alias → canonical name). An
  alias whose canonical exists in neither the overlays nor the Tela base makes
  `make theme` print an "alias source missing" warning (the name then falls
  through Tela → Breeze). If the Tela base already ships a real icon at an
  alias's name, `make theme` instead prints "… is a real file, not overwriting"
  and keeps the real icon. Both warnings are harmless.

## Recoloring (why icons are monochrome)

Small icons are monochrome and follow the active color scheme. Each SVG embeds:

```xml
<style id="current-color-scheme" type="text/css">.ColorScheme-Text { color:#232629; }</style>
...
<path class="ColorScheme-Text" style="fill:currentColor" d="..."/>
```

`FollowsColorScheme=true` in the generated `index.theme` lets KDE swap that color,
so **one icon works on light and dark panels** — there is no separate dark theme.

## Workflow

1. **Design**: create or edit an SVG in `overlays/<category>/<set>/`. Use
   the `ColorScheme-Text` pattern for monochrome icons. Match the `viewBox` to
   the largest target size you'll use (`0 0 64 64` for a 64px icon scaled down,
   or `0 0 16 16` for a 16px-only icon).

2. **Map**: edit `overlay.yaml` in the set directory. Each entry declares one
   source SVG (`src`), the theme `category`, one or more `names`, and the list
   of `sizes` it should appear at. Example:

   ```yaml
   overlays:
     - src: smb4k.svg
       category: apps
       names:
         - smb4k
       sizes: [16, 24, 32, 48, 64]
   ```

   A single source can map to multiple names and sizes. The tool writes the SVG
   once (at the first size) and symlinks the rest.

3. **Regenerate**: run `make theme` (the `tools/` icon generator). This wipes
   `output/`, lays down the vendored Tela base, applies all overlays on top,
   generates `index.theme` and the HiDPI/alias symlinks, and materializes the
   theme. Validate with `gtk-update-icon-cache -t -f output`.

4. **Build**: run `make build` to package the theme into the `.deb`.

## Source of truth

- **Committed source**: `overlays/` (SVGs + overlay.yaml + aliases.txt),
  layered on top of the vendored Tela base in `tela-icons/` (also committed —
  do not edit it, re-vendor to update; see its `README.md` for provenance)
- **Generated, never committed**: `output/` (the freedesktop icon theme)

`make theme` applies overlays **on top of the vendored Tela base**, not onto an
empty tree, so an overlay only needs to differ from Tela where it wants to
override it. Anything not overlaid comes from Tela, then finally
`Inherits=Breeze` as a last-resort fallback.
