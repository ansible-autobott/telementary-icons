# Tela icon conventions (for `-tela` overlay sets)

The base icon theme for the `*-tela` sets is **Tela**
([vinceliuice/Tela-icon-theme](https://github.com/vinceliuice/Tela-icon-theme),
GPL-3.0). Tela has **no published design guidelines** — its README only covers
installation and color variants. The rules below are **derived by inspecting the
upstream `src/` SVGs** so that new/edited icons stay visually uniform with Tela.

When in doubt, copy the closest upstream Tela SVG and swap the glyph — don't draw
from a blank canvas.

## Source layout upstream

Tela ships real per-size sets, not one scalable set:

```
src/16/<category>/     # pixel-tuned, viewBox "0 0 16 16"
src/22/<category>/     # viewBox "0 0 22 22"
src/24/<category>/     # viewBox "0 0 24 24"
src/32/<category>/     # sparse (mostly status/emblems)
src/scalable/<cat>/    # large icons: apps, big mimetypes
src/symbolic/<cat>/    # monochrome, currentColor
```

- Small folders/places are drawn at **16**; decorated folder variants
  (`folder-documents`, `folder-download`, …) live at **22/24**.
- App and full-color mimetype icons live in **scalable**.

## Folders — the signature shape

Two stacked shapes, **same hue**, back tab at 50 % opacity (this is how the
two-tone look is achieved — one color, not two):

```xml
<!-- 16px, viewBox 0 0 16 16 -->
<rect id="bottom_layer" x="1" y="1" width="14" height="5" rx="1" ry="1" fill="#5294e2" opacity="0.5"/>
<path id="top_layer" d="m8 3-1 2h-6c-0.554 0-1 0.446-1 1v8c0 0.554 0.446 1 1 1h14c0.554 0 1-0.446 1-1v-10c0-0.5-0.5-1-1-1z" fill="#5294e2"/>
```

- **1-unit outer margin**, **1-unit corner radius** (at 16px).
- Front face has a small **tab notch** at top-left (`m8 3 -1 2`).
- The **default** folder uses the ColorScheme highlight `#5294e2` via
  `class="ColorScheme-Highlight" fill="currentColor"`; named color variants
  hardcode a `fill=`.

### Folder color palette (from `src/16/places/folder-<color>.svg`)

| name    | hex       | | name    | hex       |
|---------|-----------|-|---------|-----------|
| (default/highlight) | `#5294e2` | | orange  | `#f67400` |
| blue    | `#3daee9` | | brown   | `#d35400` |
| green   | `#2ecc71` | | grey    | `#b2b2b2` |
| red     | `#da4453` | | cyan    | `#08a19d` |
| yellow  | `#fdbc4b` | | magenta | `#a10865` |
| violet  | `#9b59b6` | | black   | `#333333` |

Recolor **only** the hue — never change the folder shape, margin, radius, or the
0.5 back-tab opacity.

## Decorated folders (22/24) & symbolic

Overlays/glyphs on folders and all `symbolic/` icons are **monochrome** via the
ColorScheme pattern — same mechanism this repo already uses (see `README.md`):

```xml
<style id="current-color-scheme" type="text/css">
  .ColorScheme-Text { color:#727272; } .ColorScheme-Highlight { color:#5294e2; }
</style>
<path class="ColorScheme-Text" fill="currentColor" d="…"/>
```

- Tela's glyph grey is `#727272`; highlight is `#5294e2`.
- Folder body at 22/24 uses a **2-unit corner radius** and **no stroke**.

## Full-color mimetypes (scalable)

Base = a rounded **paper sheet** with a folded top-right corner:

```xml
<!-- viewBox "0 0 16.933 16.933", width/height 64; content in a scale(.26458) group -->
<rect x="8" y="4" width="48" height="56" ry="5" fill="#f4f4f4"/>           <!-- paper -->
<path d="m56 46-14 14h9c2.77 0 5-2.23 5-5z" fill="url(#a)" opacity="0.1"/> <!-- corner-fold shadow -->
<!-- glyph / text lines: fill="#b3b3b3" -->
```

- Paper fill `#f4f4f4`, corner radius `ry=5` (in the 64-unit group), **no stroke**.
- The only "shadow" is a **10 %-opacity black→transparent linear gradient** on the
  folded corner (`<stop offset="0"/>` black → `<stop offset="1" stop-opacity="0"/>`).
  No drop shadows, no bevels anywhere else.
- Glyph/content grey `#b3b3b3`.

## Global style rules

- **Flat**, with at most **one subtle linear gradient**; no drop shadows/bevels.
- **No strokes** on shapes (`fill` only; `paint-order:stroke fill markers` when set).
- **Rounded corners everywhere**, constant per size (≈1 @16, ≈2 @22/24, ≈5 @64).
- Consistent **outer padding / keyline** — content never touches the canvas edge.
- Match the target `viewBox` to the size bucket (`0 0 16 16`, `0 0 22 22`,
  `0 0 24 24`, or `0 0 64 64` for scalable).

## Sourcing reference SVGs

Fetch archetypes to start from:

```
https://raw.githubusercontent.com/vinceliuice/Tela-icon-theme/master/src/<size>/<category>/<name>.svg
```

Keep a representative folder / app / mimetype SVG beside the sources in each
`-tela` set folder as a design starting point.
