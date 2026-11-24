---
name: bevel-icon
description: Turn a flat icon shape (a *-symbolic glyph or a flat currentColor icon) into this theme's beveled gray "go-first" icon type — gray gradient body, white bottom glow, two inset borders. Use when asked to make an icon "grayish", "beveled", "like go-first", "Tela-style gray", or to "apply the same transformation" as go-home-16/24 or the view-modes icons.
---

# Bevel an icon into the go-first gray type

The numeric sizes of `overlays/actions/navigation/go-first-*.svg` define the
look: a flat shape rendered as a beveled gray button. `go-home-16/24` (built
from `go-home-symbolic`) and the view-modes icons (`view-list-*-16/22/32`,
built from their flat art) are the same transformation applied to other shapes.
`bevel.py` next to this file does the whole transformation. Use it rather than
rebuilding the layers by hand.

## Run it

```sh
python3 .claude/skills/bevel-icon/bevel.py SRC OUT [--size N] [--look strong]
```

- `SRC` is any flat SVG: its visible shapes are merged into one silhouette, and
  its colors are ignored.
- `--size N` is the output canvas. It defaults to SRC's canvas, and a different
  N scales the shape to fill it (the 16px symbolic → `--size 24` makes the 24).
- `--look strong` is a more prominent bevel (see "The strong look" below). The
  default, `go-first`, is the recipe the committed icons use.
- It needs `inkscape`, `xvfb-run`, and `svgo` (or `npx`). It prints the canvas,
  the gradient span and the layer count, plus warnings to read (below).

Typical run, one call per numeric size the set ships:

```sh
S=.claude/skills/bevel-icon/bevel.py
python3 $S overlays/actions/<set>/<name>-symbolic.svg overlays/actions/<set>/<name>-16.svg
python3 $S overlays/actions/<set>/<name>-symbolic.svg overlays/actions/<set>/<name>-24.svg --size 24
```

Regression check: rebuilding `go-home-16`/`go-home-24` from
`overlays/actions/navigation/go-home-symbolic.svg` reproduces the committed
files byte for byte. Rerun that check if you change the script.

## The recipe (what the script builds)

Four layers stacked over the silhouette, bottom to top:

1. **White bottom glow:** the shape, `transform="translate(0 1)"`,
   `opacity:.35;fill:#fff`.
2. **Body:** the shape, `fill:url(#a)`. `#a` is a vertical linear gradient
   (userSpaceOnUse, `y1`/`y2` = the shape's top/bottom) with go-first's stops
   `0 #737373 / .262 #636363 / .705 #4b4b4b / 1 #3f3f3f`.
3. **Inset border 1:** the outline inset **0.5px**, `opacity:.1;fill:none;
   stroke:#000;stroke-width:1;stroke-linecap:round;stroke-linejoin:round`.
4. **Inset border 2:** the outline inset **1.4px**, same style at `opacity:.03`.

**Rule:** stroke width 1, the 0.5/1.4px insets and the 1px glow shift stay the
same at every canvas size. Only the shape and the gradient span scale.

### The strong look (`--look strong`)

The same four layers, pushed for more contrast (asked for on the `edit-clear`
set, where the go-first bevel read too flat):

- **Body stops:** `0 #8f8f8f / .262 #6c6c6c / .705 #454545 / 1 #2b2b2b`.
- **Glow:** `opacity:.5`.
- **Inset border 1:** 0.5px at `opacity:.35`, a crisper dark edge.
- **Inset border 2:** 1.5px, a white highlight instead of the black line. The
  stroke is `url(#b)`, a second gradient over the same span, from white at
  `.22` (top) to `.04` (bottom). A `.45` highlight was tried first: it read as a
  second outline, especially around holes (the ✕), and was busy at 16/22px.

## Decisions already made (Andrés)

- **Bevel numeric sizes only.** Leave `-symbolic` files flat (`ColorScheme-Text`
  / `currentColor`) so the symbolic bucket keeps following the KDE color scheme.
  Beveled numeric art is fixed gray, so it does not follow the scheme.
- **Skip `*_wip*` files** and anything the user excludes.
- **One svgo pass** (`--pretty --indent=2`). The script does this. Never use
  `--multipass`: it dropped a fill from these icons before.

## Warnings and what to do

- *partly transparent … merged into the silhouette*: the source had a
  semi-opaque part (e.g. a dimmed bell), and it is now solid. Check that the
  shape still reads.
- *strokes were converted to filled outlines*: a stroked source was turned into
  fills first. Check its render.
- *the N px inset vanished*: the shape is too thin for that border, so the layer
  is left out. That is expected for 1px bars (the view-modes text lines); the
  bevel just reads subtler there.
- *source canvas is not square*: it was scaled by its width. Check placement.

## Finish the job

1. **Look at it.** Render the new icon next to `go-first` at the same size and
   view the PNG:
   ```sh
   mkdir -p /tmp/bevel && cp OUT /tmp/bevel/new.svg && cp overlays/actions/navigation/go-first-16.svg /tmp/bevel/ref.svg
   inkscape --export-type=png --export-dpi=384 /tmp/bevel/new.svg /tmp/bevel/ref.svg
   montage /tmp/bevel/ref.png /tmp/bevel/new.png -tile 2x1 -geometry +8+8 -background '#eee' /tmp/bevel/side.png
   ```
   Inkscape writes each PNG next to its SVG, so render copies in `/tmp` and
   never inside `overlays/`.
2. **Update the set's `overlay.yaml`.** Keep `src`/`names`/`sizes`, point `from`
   at the shape the art was derived from, and say in the header comment that the
   numeric art is beveled from that shape (not copied verbatim).
3. `make theme`, then check the icon in the theme viewer (`make view`).

## If it looks too big

Beveled shapes fill the canvas edge to edge, which can look heavy next to other
toolbar icons. The view-modes set was later shrunk to 85% around the center for
Dolphin. The body and borders use `transform="matrix(.85 0 0 .85 c c)"` with
`c = size * .075` (1.65 at 22px), and the glow uses the same matrix with its
y offset raised by 0.85 (the scaled 1px shift). Do this only when asked. It is
not part of the bevel.

## How the hard parts work (for fixing the script)

- **Flatten:** the source's visible shapes (hidden, zero-opacity and unpainted
  shapes are dropped first, e.g. invisible bounding-box rects) are wrapped in
  `<g transform="scale(N/w) translate(-vx -vy)">` and run through headless
  Inkscape: `selection-ungroup` ×3, `object-to-path` (plus
  `object-stroke-to-path` if stroked), then `path-union`. This bakes every
  transform and the scale into one plain path.
- **Insets:** no offset library is installed, and svgo cannot offset paths. So
  Inkscape's dynamic offset is used:
  `<path sodipodi:type="inkscape:offset" inkscape:radius="-0.5"
  inkscape:original="D" d="D"/>`, then `object-to-path` under `xvfb-run` bakes
  the computed geometry (it needs a display). The sodipodi namespace URI must be
  exactly `http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd`. With `-0.0.dtd`
  Inkscape renames the prefix and silently ignores the offset.
- **Gradient span:** `inkscape --query-y --query-height` on the flattened path.
