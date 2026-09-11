package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// cell is one name×size intersection in the rendered matrix.
type cell struct {
	State string // "override" | "fallback" | "skip"
	Src   string // svg URL, set for "override" (/raw/…) and "fallback" (/ref/…)
	Path  string // on-disk file path (theme or reference), shown in the hover popup
	Extra bool   // override present in the theme but not in the selected reference
	// For an "override" that also exists in the selected reference, the original
	// it replaced — shown on hover so you can compare. Empty when Extra.
	OrigSrc  string // reference svg URL (/ref/…)
	OrigPath string // reference on-disk path
}

type row struct {
	Name     string
	Cells    []cell // aligned to the category's Sizes
	HasExtra bool   // has at least one "not in reference" cell
}

// category is one tab: a name, its size columns, and the icon matrix.
type category struct {
	Name     string
	Sizes    []string
	Rows     []row
	Fallback int // cells served from the reference (theme has no override)
	Extra    int // theme overrides not present in the selected reference
}

// build merges the theme scan with an optional reference theme into per-category
// view models. The reference is the selectable base theme used as fallback:
//
//   - it defines the universe: rows and size columns are the union of theme +
//     reference, and a cell is only "wanted" when the reference has that exact
//     name/size (a column may exist only because some other icon uses that size);
//   - it renders the gaps: a wanted cell the theme does not override is filled
//     with the reference's own SVG (what a theme Inheriting it would show).
//
// So: theme override wins; else, if the reference has that exact cell, its icon
// is shown as a fallback; otherwise the cell is "not needed". A theme override
// the reference does not contain is flagged (Extra). With no reference selected
// the matrix is theme-only.
func build(theme *iconSet, themeRoot string, ref *iconSet, refName, refRoot string) []category {
	ctxSet := map[string]bool{}
	for c := range theme.m {
		ctxSet[c] = true
	}
	if ref != nil {
		for c := range ref.m {
			ctxSet[c] = true
		}
	}
	ctxs := keys(ctxSet)
	sort.Strings(ctxs)

	var cats []category
	for _, ctx := range ctxs {
		sizes := unionKeys(dispSizes(theme, ctx), dispSizes(ref, ctx))
		sortSizes(sizes)

		names := unionKeys(dispNames(theme, ctx), dispNames(ref, ctx))
		sort.Strings(names)

		cat := category{Name: ctx, Sizes: sizes}
		for _, name := range names {
			r := row{Name: name, Cells: make([]cell, len(sizes))}
			for i, size := range sizes {
				// Theme override wins outright.
				if rs, rn, ok := dispExact(theme, ctx, name, size); ok {
					dir := theme.relDir(ctx, rs)
					c := cell{
						State: "override",
						Src:   "/raw/" + dir + "/" + rn + ".svg",
						Path:  filepath.Join(themeRoot, dir, rn+".svg"),
					}
					if ref != nil {
						if rrs, rrn, want := dispExact(ref, ctx, name, size); want {
							// Overriding an icon the reference also has: keep the
							// original so the viewer can show it on hover.
							rdir := ref.relDir(ctx, rrs)
							c.OrigSrc = "/ref/" + refName + "/" + rdir + "/" + rrn + ".svg"
							c.OrigPath = filepath.Join(refRoot, rdir, rrn+".svg")
						} else {
							c.Extra = true
							cat.Extra++
							r.HasExtra = true
						}
					}
					r.Cells[i] = c
					continue
				}
				// No override: fill from the reference only when it has that
				// exact cell; otherwise the cell is "not needed".
				if ref != nil {
					if rs, rn, ok := dispExact(ref, ctx, name, size); ok {
						dir := ref.relDir(ctx, rs)
						r.Cells[i] = cell{
							State: "fallback",
							Src:   "/ref/" + refName + "/" + dir + "/" + rn + ".svg",
							Path:  filepath.Join(refRoot, dir, rn+".svg"),
						}
						cat.Fallback++
						continue
					}
				}
				r.Cells[i] = cell{State: "skip"}
			}
			cat.Rows = append(cat.Rows, r)
		}
		cats = append(cats, cat)
	}
	return cats
}

// dispNames returns the display row names for a context: base names, with any
// trailing -symbolic variant folded onto its base.
func dispNames(s *iconSet, ctx string) map[string]bool {
	out := map[string]bool{}
	if s == nil || s.m[ctx] == nil {
		return out
	}
	for name := range s.m[ctx] {
		if base, _, ok := symbolicOf(name); ok {
			out[base] = true
		} else {
			out[name] = true
		}
	}
	return out
}

// dispSizes returns the display size columns for a context: numeric buckets plus
// a "symbolic" (and "symbolic-rtl") column folded from -symbolic names.
func dispSizes(s *iconSet, ctx string) map[string]bool {
	out := map[string]bool{}
	if s == nil || s.m[ctx] == nil {
		return out
	}
	for name, sizes := range s.m[ctx] {
		if _, ds, ok := symbolicOf(name); ok {
			out[ds] = true
		} else {
			for sz := range sizes {
				out[sz] = true
			}
		}
	}
	return out
}

// dispExact resolves a display cell to the real (size dir, file name) when the
// set contains it exactly — an exact numeric size, or the -symbolic file for a
// symbolic column.
func dispExact(s *iconSet, ctx, name, size string) (rawSize, rawName string, ok bool) {
	if s == nil || s.m[ctx] == nil {
		return "", "", false
	}
	if strings.HasPrefix(size, "symbolic") {
		rawName = name + "-" + size // symbolic -> -symbolic, symbolic-rtl -> -symbolic-rtl
		sizes := s.m[ctx][rawName]
		if sizes == nil {
			return "", "", false
		}
		return pickSymbolicSource(sizes), rawName, true
	}
	if s.has(ctx, name, size) {
		return size, name, true
	}
	return "", "", false
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func unionKeys(a, b map[string]bool) []string {
	u := map[string]bool{}
	for k := range a {
		u[k] = true
	}
	for k := range b {
		u[k] = true
	}
	return keys(u)
}

// multiFlag collects a repeatable string flag (e.g. -ref a -ref b).
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

// discoverThemes returns name->root for every subdirectory of dir that is a
// freedesktop icon theme (contains an index.theme). Symlinked theme dirs are
// followed.
func discoverThemes(dir string) map[string]string {
	out := map[string]string{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return out
	}
	for _, e := range entries {
		p := filepath.Join(dir, e.Name())
		fi, err := os.Stat(p) // follows symlinks
		if err != nil || !fi.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(p, "index.theme")); err != nil {
			continue
		}
		out[e.Name()] = p
	}
	return out
}

// userIconsDir is the user's local icon-theme directory, honouring
// $XDG_DATA_HOME (default ~/.local/share). Returns "" if the home is unknown.
func userIconsDir() string {
	if x := os.Getenv("XDG_DATA_HOME"); x != "" {
		return filepath.Join(x, "icons")
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".local", "share", "icons")
	}
	return ""
}

// referenceThemes discovers reference icon themes for the selector: those under
// iconsDir plus the user's local themes (userIconsDir). On a name clash the
// user-local theme wins, matching the freedesktop lookup order (user dirs take
// precedence over /usr/share).
func referenceThemes(iconsDir string) map[string]string {
	out := discoverThemes(iconsDir)
	if d := userIconsDir(); d != "" {
		for name, root := range discoverThemes(d) {
			out[name] = root
		}
	}
	return out
}

// setNoStore tells the browser never to cache the response, so edits to the
// theme show up immediately on reload.
func setNoStore(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

// noStore wraps a handler so every response it serves is marked no-store.
func noStore(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setNoStore(w)
		h.ServeHTTP(w, r)
	})
}

// cmdView implements: icons view [-theme D] [-addr A] [-icons-dir D] [-ref P ...]
func cmdView(args []string) {
	fs := flag.NewFlagSet("view", flag.ExitOnError)
	theme := fs.String("theme", defaultTheme(), "path to the icon theme root (<category>/<size>/<name>.svg)")
	addr := fs.String("addr", ":3033", "address to listen on")
	iconsDir := fs.String("icons-dir", "/usr/share/icons",
		"directory scanned for reference icon themes; the user's ~/.local/share/icons (or $XDG_DATA_HOME/icons) is always scanned too")
	var refPaths multiFlag
	fs.Var(&refPaths, "ref",
		"extra reference icon theme path, repeatable (selectable base theme, same as a discovered theme)")
	fs.Parse(args)

	root, err := filepath.Abs(*theme)
	if err != nil {
		log.Fatal(err)
	}

	// Reference themes: everything under -icons-dir that looks like a theme,
	// plus any explicit -ref paths (validated, like the old -breeze flag).
	registry := referenceThemes(*iconsDir)
	for _, p := range refPaths {
		ap, err := filepath.Abs(p)
		if err != nil {
			log.Fatal(err)
		}
		if fi, err := os.Stat(ap); err != nil || !fi.IsDir() {
			log.Fatalf("reference theme not found: %s (pass -ref <path> to a theme root)", ap)
		}
		registry[filepath.Base(ap)] = ap
	}
	var refNames []string
	for name := range registry {
		refNames = append(refNames, name)
	}
	sort.Strings(refNames)

	// A static file server per reference theme, dispatched by name under /ref/.
	refServers := map[string]http.Handler{}
	for name, rroot := range registry {
		refServers[name] = http.StripPrefix("/ref/"+name+"/", http.FileServer(http.Dir(rroot)))
	}
	// References are system themes (static): scan each once, on first use.
	var mu sync.Mutex
	refCache := map[string]*iconSet{}
	getRef := func(name string) (*iconSet, string, error) {
		rroot := registry[name]
		if rroot == "" {
			return nil, "", nil
		}
		mu.Lock()
		defer mu.Unlock()
		if s := refCache[name]; s != nil {
			return s, rroot, nil
		}
		s, err := scanReference(rroot)
		if err != nil {
			return nil, "", err
		}
		refCache[name] = s
		return s, rroot, nil
	}

	// The theme is edited live -> never cache; references are static -> cacheable.
	http.Handle("/raw/", noStore(http.StripPrefix("/raw/", http.FileServer(http.Dir(root)))))
	http.HandleFunc("/ref/", func(w http.ResponseWriter, r *http.Request) {
		rest := strings.TrimPrefix(r.URL.Path, "/ref/")
		i := strings.IndexByte(rest, '/')
		if i < 0 {
			http.NotFound(w, r)
			return
		}
		h := refServers[rest[:i]]
		if h == nil {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		setNoStore(w)
		// Re-scan the theme per request so file changes show on reload.
		theme, err := scanSet(root)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		selected := r.URL.Query().Get("ref")
		if !contains(refNames, selected) {
			selected = ""
		}
		var ref *iconSet
		var refRoot string
		if selected != "" {
			ref, refRoot, err = getRef(selected)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		data := struct {
			Theme       string
			ThemeName   string
			Ref         string
			Categories  []category
			Refs        []string
			SelectedRef string
		}{
			Theme:       root,
			ThemeName:   filepath.Base(root),
			Ref:         refRoot,
			Categories:  build(theme, root, ref, selected, refRoot),
			Refs:        refNames,
			SelectedRef: selected,
		}
		if err := page.Execute(w, data); err != nil {
			log.Println("render:", err)
		}
	})

	fmt.Printf("icons: serving %s on http://localhost%s\n", root, *addr)
	fmt.Printf("references: %d themes (%s + %s + %d passed): %v\n", len(registry), *iconsDir, userIconsDir(), len(refPaths), refNames)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var page = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.ThemeName}} — icons</title>
<style>
  :root {
    --paper:#fcfcfd; --panel:#ffffff;
    --ink:#16181d; --muted:#6b7280; --faint:#9aa1ac;
    --line:#e6e8ec; --line-strong:#d5d8de;
    --accent:#2f86c9; --danger:#d64550;
    /* transparency checker, the vector-editor backdrop */
    --checker:#eef0f3;
    --sans:-apple-system,"Segoe UI",Roboto,Inter,system-ui,sans-serif;
    --mono:"SF Mono","JetBrains Mono","Cascadia Code",ui-monospace,Menlo,Consolas,monospace;
  }
  * { box-sizing:border-box; }
  html { -webkit-text-size-adjust:100%; }
  body {
    margin:0; background:var(--paper); color:var(--ink);
    font-family:var(--sans); font-size:14px; line-height:1.4;
    -webkit-font-smoothing:antialiased;
  }

  /* one centered column, sized to the widest content (the table) */
  .wrap { width:fit-content; max-width:100%; margin:0 auto; }

  header {
    padding:20px 24px 0;
  }
  .theme {
    font-size:16px; font-weight:600; letter-spacing:-.01em;
  }
  .path {
    font-family:var(--mono); font-size:12px; color:var(--muted);
    margin-top:2px; max-width:72ch;
    white-space:nowrap; overflow:hidden; text-overflow:ellipsis;
  }
  .controls {
    display:flex; align-items:center; gap:16px; flex-wrap:wrap;
    margin-top:12px;
  }
  .controls label { font-size:13px; color:var(--muted); }
  .controls select {
    font:inherit; font-size:13px; padding:5px 8px;
    border:1px solid var(--line-strong); border-radius:6px;
    background:var(--panel); color:var(--ink);
  }
  .legend { display:flex; gap:8px; font-size:12px; color:var(--muted); flex-wrap:wrap; }
  .legend .leg {
    appearance:none; border:1px solid transparent; background:none; cursor:pointer;
    font:inherit; font-size:12px; color:var(--muted);
    display:inline-flex; align-items:center; gap:6px; padding:3px 7px; border-radius:5px;
    transition:background .12s ease, color .12s ease, opacity .12s ease;
  }
  .legend .leg:hover { background:#eef1f5; color:var(--ink); }
  .legend .leg:focus-visible { outline:2px solid var(--accent); outline-offset:1px; }
  .legend .leg.off { opacity:.45; text-decoration:line-through; }
  /* toggled-off categories are hidden; content only, so columns stay aligned */
  body.hide-override .c-override,
  body.hide-fallback .c-fallback,
  body.hide-skip .c-skip,
  body.hide-extra .c-extra { display:none; }
  /* "only not-in-reference": collapse to rows that have an extra */
  body.only-extra tbody tr:not(.has-extra) { display:none; }
  .only {
    appearance:none; cursor:pointer; font:inherit; font-size:12px;
    color:var(--danger); background:none;
    border:1px solid var(--line-strong); border-radius:5px; padding:3px 8px;
    transition:background .12s ease, color .12s ease;
  }
  .only:hover { background:#fdecee; }
  .only:focus-visible { outline:2px solid var(--accent); outline-offset:1px; }
  .only.on { background:var(--danger); border-color:var(--danger); color:#fff; }
  .swatch { width:12px; height:12px; border-radius:3px; border:1px solid var(--line-strong); }
  .swatch.ov { background:#dfeaf5; border-color:var(--accent); }
  .swatch.fb { background:var(--panel); border-style:dashed; }
  .swatch.skip { background:var(--panel); }
  .swatch.tri { position:static; background:var(--danger); border-color:var(--danger); clip-path:polygon(100% 0, 0 0, 100% 100%); }
  .swatch.tri-ov { position:static; background:var(--accent); border-color:var(--accent); clip-path:polygon(0 0, 100% 0, 0 100%); }

  .tabs {
    display:flex; flex-wrap:wrap; gap:2px;
    padding:16px 24px 0;
    border-bottom:1px solid var(--line);
    position:sticky; top:0; background:var(--paper); z-index:30;
  }
  .tabs button {
    appearance:none; border:0; background:none; cursor:pointer;
    font:inherit; font-size:13px; color:var(--muted);
    padding:9px 14px; border-radius:6px 6px 0 0;
    border-bottom:2px solid transparent; margin-bottom:-1px;
    transition:color .12s ease, border-color .12s ease;
  }
  .tabs button:hover { color:var(--ink); }
  .tabs button .n {
    font-family:var(--mono); font-size:11px; color:var(--faint);
    margin-left:6px;
  }
  .tabs button .fb {
    font-family:var(--mono); font-size:11px; color:var(--muted);
    margin-left:5px;
  }
  .tabs button .extra {
    font-family:var(--mono); font-size:11px; color:var(--danger);
    margin-left:5px;
  }
  .tabs button.active {
    color:var(--ink); font-weight:600;
    border-bottom-color:var(--accent);
  }
  .tabs button.active .n { color:var(--accent); }
  .tabs button:focus-visible {
    outline:2px solid var(--accent); outline-offset:-2px; border-radius:6px;
  }

  .panel { display:none; padding:16px 24px 48px; }
  .panel.active { display:block; }

  table { border-collapse:separate; border-spacing:0; }
  thead th {
    position:sticky; top:44px; z-index:20;
    background:var(--paper);
    font-size:12px; font-weight:600; color:var(--muted);
    font-variant-numeric:tabular-nums;
    padding:8px 10px; text-align:center;
    border-bottom:1px solid var(--line-strong);
  }
  thead th.name-h {
    left:0; z-index:25; text-align:left;
    font-family:var(--mono); font-weight:600;
  }
  tbody th {
    position:sticky; left:0; z-index:10;
    background:var(--panel);
    font-family:var(--mono); font-size:12.5px; font-weight:400;
    text-align:left; white-space:nowrap;
    padding:0 16px 0 0; color:var(--ink);
    border-bottom:1px solid var(--line);
  }
  tbody td {
    padding:5px 6px; text-align:center;
    border-bottom:1px solid var(--line);
  }
  tbody tr:hover th { background:#f2f6fb; color:var(--accent); }
  tbody tr:hover td { background:#f7fafd; }

  .tile {
    width:56px; height:56px; margin:0 auto; position:relative;
    display:flex; align-items:center; justify-content:center;
    border:1px solid var(--line); border-radius:7px;
    background-color:var(--panel);
    background-image:
      linear-gradient(45deg,var(--checker) 25%,transparent 25%),
      linear-gradient(-45deg,var(--checker) 25%,transparent 25%),
      linear-gradient(45deg,transparent 75%,var(--checker) 75%),
      linear-gradient(-45deg,transparent 75%,var(--checker) 75%);
    background-size:12px 12px;
    background-position:0 0,0 6px,6px -6px,-6px 0;
  }
  .tile img { display:block; max-width:44px; max-height:44px; }
  /* red triangle, upper-right: theme icon not present in the selected reference */
  .tri {
    position:absolute; top:2px; right:2px; width:12px; height:12px;
    background:var(--danger); clip-path:polygon(100% 0, 0 0, 100% 100%);
  }
  /* blue triangle, upper-left: the theme provides this icon (an override) */
  .tri-ov {
    position:absolute; top:2px; left:2px; width:12px; height:12px;
    background:var(--accent); clip-path:polygon(0 0, 100% 0, 0 100%);
  }
  /* fallback: inherited from the reference theme (no theme override) — dimmed,
     dashed, with a small "R" corner tag. */
  .tile.fb { border-style:dashed; border-color:var(--line-strong); }
  .tile.fb img { opacity:.5; }
  .tile.fb::after {
    content:"R"; position:absolute; top:2px; right:3px;
    font:600 8px var(--mono); color:var(--faint);
  }
  /* not-needed: neither the theme nor the reference provides this cell */
  .skip { position:relative; display:inline-block; color:var(--faint); font-family:var(--mono); user-select:none; }

  /* hover popup — pure CSS, shows instantly (no delay) */
  .tip {
    display:none; position:absolute; left:50%; top:calc(100% + 6px);
    transform:translateX(-50%); z-index:100;
    width:max-content; max-width:360px; text-align:left;
    background:#1f2430; color:#f5f6f8;
    border-radius:6px; padding:7px 9px;
    box-shadow:0 6px 20px rgba(0,0,0,.28);
    font-size:11px; line-height:1.4; pointer-events:none;
  }
  .tile:hover .tip, .skip:hover .tip { display:block; }
  .tip-state { display:block; font-weight:600; font-family:var(--sans); margin-bottom:2px; }
  .tip-state.ov   { color:#7fb2e0; }
  .tip-state.fb   { color:#c9cdd6; }
  .tip-state.skip { color:#9aa1ac; }
  .tip-extra { display:block; color:#f2a3a8; font-weight:600; margin-bottom:2px; }
  .tip-path { font-family:var(--mono); color:#cfd3da; word-break:break-all; }
  .tip-orig {
    display:block; margin-top:6px; padding-top:6px;
    border-top:1px solid #3a4150;
  }
  .tip-orig-lbl { display:block; font-weight:600; font-family:var(--sans); color:#c9cdd6; margin-bottom:3px; }
  .tip-orig img {
    display:block; width:48px; height:48px; object-fit:contain;
    background:#fff; border-radius:4px; padding:3px; margin-bottom:3px;
    image-rendering:auto;
  }

  @media (prefers-reduced-motion:reduce) {
    * { transition:none !important; }
  }
</style>
</head>
<body>
<main class="wrap">
<header>
  <div class="theme">{{.ThemeName}}</div>
  <div class="path">{{.Theme}}</div>
  <div class="path">base theme: {{if .Ref}}{{.Ref}}{{else}}— none —{{end}}</div>
  <form class="controls" method="get" id="controls">
    <label for="ref">Base theme:</label>
    <select name="ref" id="ref" onchange="document.getElementById('controls').submit()">
      <option value=""{{if eq .SelectedRef ""}} selected{{end}}>— none —</option>
      {{range .Refs}}
      <option value="{{.}}"{{if eq . $.SelectedRef}} selected{{end}}>{{.}}</option>
      {{end}}
    </select>
    <div class="legend" role="group" aria-label="Toggle cell visibility">
      <button type="button" class="leg" data-cat="override" title="click to hide/show"><span class="swatch tri-ov"></span>theme provides</button>
      <button type="button" class="leg" data-cat="fallback" title="click to hide/show"><span class="swatch fb"></span>base fallback</button>
      <button type="button" class="leg" data-cat="skip" title="click to hide/show"><span class="swatch skip"></span>not needed</button>
      <button type="button" class="leg" data-cat="extra" title="click to hide/show"><span class="swatch tri"></span>not in base</button>
    </div>
    <button type="button" class="only" id="onlyExtra" aria-pressed="false" title="collapse the table to only rows that have a not-in-base icon">▲ only not-in-base</button>
  </form>
</header>

<nav class="tabs" role="tablist" aria-label="Icon categories">
{{range $i, $c := .Categories}}
  <button role="tab" data-tab="{{$c.Name}}" aria-controls="panel-{{$c.Name}}"{{if eq $i 0}} class="active" aria-selected="true"{{else}} aria-selected="false"{{end}}>{{$c.Name}}<span class="n">{{len $c.Rows}}</span>{{if $c.Fallback}}<span class="fb">R{{$c.Fallback}}</span>{{end}}{{if $c.Extra}}<span class="extra">▲{{$c.Extra}}</span>{{end}}</button>
{{end}}
</nav>

{{range $i, $c := .Categories}}
<section class="panel{{if eq $i 0}} active{{end}}" id="panel-{{$c.Name}}" role="tabpanel">
  <table>
    <thead>
      <tr><th class="name-h">name</th>{{range $c.Sizes}}<th>{{.}}</th>{{end}}</tr>
    </thead>
    <tbody>
    {{range $r := $c.Rows}}
      <tr{{if $r.HasExtra}} class="has-extra"{{end}}>
        <th scope="row">{{$r.Name}}</th>
        {{range $cell := $r.Cells}}
          {{if eq $cell.State "override"}}
            <td><span class="tile {{if $cell.Extra}}c-extra{{else}}c-override{{end}}"><img loading="lazy" src="{{$cell.Src}}" alt="{{$r.Name}}"><span class="tri-ov"></span>{{if $cell.Extra}}<span class="tri"></span>{{end}}<span class="tip"><span class="tip-state ov">theme override</span>{{if $cell.Extra}}<span class="tip-extra">not in selected base</span>{{end}}<span class="tip-path">{{$cell.Path}}</span>{{if $cell.OrigSrc}}<span class="tip-orig"><span class="tip-orig-lbl">original (base)</span><img loading="lazy" src="{{$cell.OrigSrc}}" alt="original"><span class="tip-path">{{$cell.OrigPath}}</span></span>{{end}}</span></span></td>
          {{else if eq $cell.State "fallback"}}
            <td><span class="tile fb c-fallback"><img loading="lazy" src="{{$cell.Src}}" alt="{{$r.Name}} (base fallback)"><span class="tip"><span class="tip-state fb">base fallback · not overridden</span><span class="tip-path">{{$cell.Path}}</span></span></span></td>
          {{else}}
            <td><span class="skip c-skip">-<span class="tip"><span class="tip-state skip">not needed</span></span></span></td>
          {{end}}
        {{end}}
      </tr>
    {{end}}
    </tbody>
  </table>
</section>
{{end}}
</main>

<script>
  const tabs = document.querySelectorAll('.tabs button');
  tabs.forEach(function (b) {
    b.addEventListener('click', function () {
      tabs.forEach(function (x) {
        x.classList.remove('active');
        x.setAttribute('aria-selected', 'false');
      });
      document.querySelectorAll('.panel').forEach(function (x) { x.classList.remove('active'); });
      b.classList.add('active');
      b.setAttribute('aria-selected', 'true');
      document.getElementById('panel-' + b.dataset.tab).classList.add('active');
    });
  });

  // Legend items toggle visibility of their cell category.
  document.querySelectorAll('.legend .leg').forEach(function (b) {
    b.addEventListener('click', function () {
      const cls = 'hide-' + b.dataset.cat;
      const hidden = document.body.classList.toggle(cls);
      b.classList.toggle('off', hidden);
      b.setAttribute('aria-pressed', hidden ? 'true' : 'false');
    });
  });

  // "only not-in-reference": collapse the table to rows that have a red triangle.
  const only = document.getElementById('onlyExtra');
  only.addEventListener('click', function () {
    const on = document.body.classList.toggle('only-extra');
    only.classList.toggle('on', on);
    only.setAttribute('aria-pressed', on ? 'true' : 'false');
  });
</script>
</body>
</html>
`))
