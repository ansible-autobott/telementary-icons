#!/usr/bin/env python3
"""Turn a flat icon shape into the theme's beveled gray "go-first" icon type.

    bevel.py SRC OUT [--size N]

SRC is any flat SVG (usually a *-symbolic glyph or a flat currentColor icon);
its visible shapes are merged into one silhouette. OUT gets the 4-layer bevel
(white bottom glow, gray gradient body, two inset borders) on an N x N canvas;
N defaults to SRC's own canvas, and a different N scales the shape to fit.
Needs inkscape, xvfb-run and svgo (or npx) on PATH.
"""

import argparse
import os
import re
import shutil
import subprocess
import sys
import tempfile
import xml.etree.ElementTree as ET

SVG = "http://www.w3.org/2000/svg"
SODIPODI = "http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"  # NOT ...-0.0.dtd
INKSCAPE = "http://www.inkscape.org/namespaces/inkscape"

# Looks: the body gradient's stops (top to bottom), the bottom glow's opacity
# and the inset borders as (inset in px, opacity, stroke). Insets are absolute
# at every canvas size, never scaled. A "highlight" adds gradient #b, a
# top-lit (offset, color, opacity) stroke for an inset to use as url(#b).
LOOKS = {
    # go-first-16's own recipe.
    "go-first": {
        "stops": [("0", "#737373"), (".262", "#636363"), (".705", "#4b4b4b"), ("1", "#3f3f3f")],
        "glow": ".35",
        "insets": [(0.5, ".1", "#000"), (1.4, ".03", "#000")],
    },
    # Higher-contrast body, darker edge and a white inner highlight.
    "strong": {
        "stops": [("0", "#8f8f8f"), (".262", "#6c6c6c"), (".705", "#454545"), ("1", "#2b2b2b")],
        "glow": ".5",
        "insets": [(0.5, ".35", "#000"), (1.5, "1", "url(#b)")],
        "highlight": [("0", "#fff", ".22"), ("1", "#fff", ".04")],
    },
}
BODY = "fill:url(#a);stroke:none"
BORDER = "fill:none;stroke:{};stroke-width:1;stroke-linecap:round;stroke-linejoin:round"

SHAPES = {"path", "rect", "circle", "ellipse", "polygon", "polyline", "line"}
SKIP = {"defs", "title", "desc", "metadata", "style", "namedview"}

ET.register_namespace("", SVG)


def local(tag):
    return tag.rsplit("}", 1)[-1]


def prop(el, name):
    """An element's own presentation value (style= wins over the attribute)."""
    m = re.search(r"(?:^|;)\s*" + re.escape(name) + r"\s*:\s*([^;]+)", el.get("style", ""))
    return (m.group(1) if m else el.get(name, "")).strip()


def opacity(el, name="opacity"):
    try:
        return float(prop(el, name) or 1)
    except ValueError:
        return 1.0


def canvas(root):
    vb = root.get("viewBox")
    if vb:
        x, y, w, h = map(float, vb.replace(",", " ").split())
        return x, y, w, h
    num = lambda v: float(re.match(r"[\d.]+", v).group())
    return 0.0, 0.0, num(root.get("width")), num(root.get("height"))


def prune(el, warnings, parent_hidden=False):
    """Drop invisible elements (hidden, opacity 0, unpainted) so they do not
    grow the silhouette; note partly transparent and stroked ones."""
    stroked = False
    for child in list(el):
        tag = local(child.tag)
        hidden = (parent_hidden or prop(child, "display") == "none"
                  or prop(child, "visibility") == "hidden"
                  or opacity(child) == 0)
        if tag in SKIP or hidden:
            el.remove(child)
            continue
        if tag in SHAPES:
            fill, stroke = prop(child, "fill"), prop(child, "stroke")
            if fill == "none" and stroke in ("", "none"):
                el.remove(child)
                continue
            op = min(opacity(child), opacity(child, "fill-opacity"))
            if op < 1:
                warnings.append(f"partly transparent <{tag}> (opacity {op:g}) merged into the silhouette")
            if stroke not in ("", "none"):
                stroked = True
        stroked |= prune(child, warnings)
    return stroked


def inkscape(src, actions):
    subprocess.run(["xvfb-run", "-a", "inkscape", src, "--actions=" + actions],
                   check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)


def flatten(src, size, tmp, warnings):
    """SRC's visible shapes as one path `d` in an N x N canvas, plus its style."""
    root = ET.parse(src).getroot()
    x, y, w, h = canvas(root)
    if abs(w - h) > 1e-6:
        warnings.append(f"source canvas is {w}x{h}, not square; scaled by its width")
    size = size or int(round(w))
    stroked = prune(root, warnings)
    body = "".join(ET.tostring(c, encoding="unicode") for c in root)
    f = size / w
    wrap = os.path.join(tmp, "wrap.svg")
    with open(wrap, "w") as fh:
        fh.write(f'<svg xmlns="{SVG}" width="{size}" height="{size}" viewBox="0 0 {size} {size}">'
                 f'<g transform="scale({f}) translate({-x} {-y})">{body}</g></svg>')
    flat = os.path.join(tmp, "flat.svg")
    to_path = "object-to-path;object-stroke-to-path" if stroked else "object-to-path"
    if stroked:
        warnings.append("strokes were converted to filled outlines; check the result")
    inkscape(wrap, "select-all:all;selection-ungroup;selection-ungroup;selection-ungroup;"
                   f"select-all:all;{to_path};path-union;"
                   f"export-plain-svg;export-filename:{flat};export-do")
    paths = [p for p in ET.parse(flat).getroot().iter(f"{{{SVG}}}path") if p.get("d")]
    if len(paths) != 1 or paths[0].get("transform"):
        sys.exit(f"flatten: expected one untransformed path in {flat}, got {len(paths)}")
    return size, paths[0].get("d"), prop(paths[0], "fill-rule")


def bbox_y(d, size, tmp):
    """Top and bottom of the silhouette: the gradient's y1/y2."""
    f = os.path.join(tmp, "bbox.svg")
    with open(f, "w") as fh:
        fh.write(f'<svg xmlns="{SVG}" width="{size}" height="{size}"><path id="s" d="{d}"/></svg>')
    out = subprocess.run(["inkscape", "--query-id=s", "--query-y", "--query-height", f],
                         check=True, capture_output=True, text=True).stdout.split()
    top, height = float(out[0]), float(out[1])
    return top, top + height


def insets(d, size, radii, tmp):
    """True geometric insets of `d` via Inkscape's dynamic offset (needs a
    display, hence xvfb). Returns one `d` per radius ("" if it vanished)."""
    f, out = os.path.join(tmp, "offset.svg"), os.path.join(tmp, "offset-out.svg")
    paths = "".join(f'<path id="i{i}" sodipodi:type="inkscape:offset" inkscape:radius="{-r}" '
                    f'inkscape:original="{d}" d="{d}"/>' for i, r in enumerate(radii))
    with open(f, "w") as fh:
        fh.write(f'<svg xmlns="{SVG}" xmlns:sodipodi="{SODIPODI}" xmlns:inkscape="{INKSCAPE}" '
                 f'width="{size}" height="{size}">{paths}</svg>')
    inkscape(f, f"select-all;object-to-path;export-plain-svg;export-filename:{out};export-do")
    got = {p.get("id"): p.get("d", "") for p in ET.parse(out).getroot().iter(f"{{{SVG}}}path")}
    return [got.get(f"i{i}", "").strip() for i in range(len(radii))]


def svgo(path):
    cmd = [shutil.which("svgo")] if shutil.which("svgo") else ["npx", "--yes", "svgo"]
    # A single pass: --multipass has dropped fills from these icons before.
    subprocess.run(cmd + ["--pretty", "--indent=2", "-q", "-i", path, "-o", path], check=True)


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("src")
    ap.add_argument("out")
    ap.add_argument("--size", type=int, help="output canvas (default: the source's)")
    ap.add_argument("--look", choices=LOOKS, default="go-first",
                    help="bevel strength (default: go-first, the committed icons' recipe)")
    a = ap.parse_args()
    look = LOOKS[a.look]

    warnings = []
    with tempfile.TemporaryDirectory() as tmp:
        size, d, rule = flatten(a.src, a.size, tmp, warnings)
        y1, y2 = bbox_y(d, size, tmp)
        borders = insets(d, size, [r for r, _, _ in look["insets"]], tmp)
    fill_rule = f";fill-rule:{rule}" if rule else ""
    c = f"{size / 2:g}"
    span = f'x1="{c}" x2="{c}" y1="{y1:g}" y2="{y2:g}" gradientUnits="userSpaceOnUse"'
    stops = "".join(f'<stop offset="{o}" style="stop-color:{col};stop-opacity:1"/>' for o, col in look["stops"])
    defs = f'<linearGradient id="a" {span}>{stops}</linearGradient>'
    if "highlight" in look:
        stops = "".join(f'<stop offset="{o}" style="stop-color:{col};stop-opacity:{op}"/>'
                        for o, col, op in look["highlight"])
        defs += f'<linearGradient id="b" {span}>{stops}</linearGradient>'
    layers = [f'<path d="{d}" style="opacity:{look["glow"]};fill:#fff;stroke:none{fill_rule}" '
              'transform="translate(0 1)"/>',
              f'<path d="{d}" style="{BODY}{fill_rule}"/>']
    for (r, op, stroke), bd in zip(look["insets"], borders):
        if bd:
            layers.append(f'<path d="{bd}" style="opacity:{op};{BORDER.format(stroke)}"/>')
        else:
            warnings.append(f"the {r}px inset vanished (shape too thin); that border layer is left out")
    with open(a.out, "w") as fh:
        fh.write(f'<svg xmlns="{SVG}" width="{size}" height="{size}"><defs>{defs}</defs>'
                 + "".join(layers) + "</svg>\n")
    svgo(a.out)
    for w in warnings:
        print(f"warning: {w}", file=sys.stderr)
    print(f"{a.out}: {size}x{size}, gradient y {y1:g}..{y2:g}, {len(layers)} layers")


if __name__ == "__main__":
    main()
