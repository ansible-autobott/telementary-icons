package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// writeFiles creates files (path -> content) under root; content "->target"
// makes a symlink instead.
func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for p, content := range files {
		full := filepath.Join(root, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if target, ok := strings.CutPrefix(content, "->"); ok {
			if err := os.Symlink(target, full); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

const (
	monoSVG  = `<svg><style>.ColorScheme-Text { color:#565656; }</style><path class="ColorScheme-Text" style="fill:currentColor"/></svg>`
	graySVG  = `<svg><path fill="url(#gpa:foreground) #555761"/></svg>`
	colorSVG = `<svg><linearGradient><stop stop-color="#9bdb4d"/><stop stop-color="#68b723"/></linearGradient><path fill="#206b00"/></svg>`
)

func TestSvgStyle(t *testing.T) {
	for svg, want := range map[string]string{monoSVG: "mono", graySVG: "gray", colorSVG: "color"} {
		if got := svgStyle([]byte(svg)); got != want {
			t.Errorf("svgStyle(%.40q…) = %q, want %q", svg, got, want)
		}
	}
}

// recolorSVG paints the scheme's text color, GNOME's gpa:foreground paint, and
// currentColor (via the root) in the requested color.
func TestRecolorSVG(t *testing.T) {
	got := string(recolorSVG([]byte(monoSVG), "#fcfcfc"))
	if !strings.Contains(got, ".ColorScheme-Text { color:#fcfcfc; }") || !strings.HasPrefix(got, `<svg color="#fcfcfc">`) {
		t.Errorf("mono recolor = %s", got)
	}
	got = string(recolorSVG([]byte(graySVG), "#fcfcfc"))
	if !strings.Contains(got, `fill="#fcfcfc"`) {
		t.Errorf("gpa recolor = %s", got)
	}
}

// nearestBucket prefers the closest numeric size and, on a tie, the larger.
func TestNearestBucket(t *testing.T) {
	s := newIconSet()
	s.add("actions", "16", "x")
	s.add("actions", "24", "x")
	s.add("actions", "symbolic", "x")
	for _, tc := range []struct{ from, want string }{{"16", "16"}, {"22", "24"}, {"20", "24"}, {"48", "24"}} {
		b, ctx, ok := nearestBucket(s, "x", tc.from)
		if !ok || b != tc.want || ctx != "actions" {
			t.Errorf("nearestBucket(%s) = %s %s %v, want %s", tc.from, b, ctx, ok, tc.want)
		}
	}
	if _, _, ok := nearestBucket(s, "missing", "22"); ok {
		t.Errorf("nearestBucket found a missing name")
	}
}

// The alias index lists every symlink whose chain passes through a file,
// directly or through another alias.
func TestBuildAliasIndexFollowsChains(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		"actions/22/a.svg": monoSVG,
		"actions/22/b.svg": "->a.svg",
		"actions/22/c.svg": "->b.svg",
	})
	idx := buildAliasIndex(root, []string{"actions"})
	a := filepath.Join(root, "actions", "22", "a.svg")
	b := filepath.Join(root, "actions", "22", "b.svg")
	if got := idx[a]; len(got) != 2 {
		t.Errorf("aliases of a = %v, want b and c", got)
	}
	if got := idx[b]; !reflect.DeepEqual(got, []aliasRef{{"c", "22"}}) {
		t.Errorf("aliases of b = %v, want [c]", got)
	}
}

// current resolves like the icon lookup: the exact bucket, then the nearest
// size the theme has, and only for names the theme lacks the base theme.
func TestLookupCurrent(t *testing.T) {
	root, baseRoot := t.TempDir(), t.TempDir()
	writeFiles(t, root, map[string]string{
		"actions/22/doc.svg": monoSVG,
		"actions/22/alt.svg": "->doc.svg",
	})
	writeFiles(t, baseRoot, map[string]string{"actions/16/only-base.svg": monoSVG})
	theme, err := scanSet(root)
	if err != nil {
		t.Fatal(err)
	}
	base := newIconSet()
	base.add("actions", "16", "only-base")
	lk := &lookup{root: root, theme: theme, base: base, baseName: "breeze", baseRoot: baseRoot}

	if ref := lk.current("actions", "doc", "22"); ref.Path != "actions/22/doc.svg" || ref.Style != "mono" || ref.Note != "" {
		t.Errorf("exact = %+v", ref)
	}
	if ref := lk.current("actions", "alt", "22"); ref.Note != "alias of doc" {
		t.Errorf("alias note = %q", ref.Note)
	}
	if ref := lk.current("actions", "doc", "16"); ref.Path != "actions/22/doc.svg" || !strings.Contains(ref.Note, "KDE scales actions/22") {
		t.Errorf("nearest = %+v", ref)
	}
	if ref := lk.current("actions", "only-base", "22"); ref.Source != "base" || ref.URL != "/review/img/base/actions/16/only-base.svg" {
		t.Errorf("base fallback = %+v", ref)
	}
	if ref := lk.current("actions", "nowhere", "22"); ref.URL != "" {
		t.Errorf("missing name = %+v, want no image", ref)
	}
}

// An item's rows describe each proposed bucket, and its decisions-file task
// names the elementary source per link (resolving symlinks) and the aliases
// that follow the replaced files.
func TestReviewerItem(t *testing.T) {
	theme, elem := t.TempDir(), t.TempDir()
	writeFiles(t, theme, map[string]string{
		"actions/16/doc.svg":                monoSVG,
		"actions/22/doc.svg":                monoSVG,
		"actions/22/doc-alias.svg":          "->doc.svg",
		"actions/symbolic/doc-symbolic.svg": monoSVG,
	})
	writeFiles(t, elem, map[string]string{
		"actions/16/doc.svg":                colorSVG,
		"places/24/folder.svg":              colorSVG,
		"actions/24/doc.svg":                "->../../places/24/folder.svg",
		"actions/symbolic/doc-symbolic.svg": graySVG,
	})
	set, err := scanSet(theme)
	if err != nil {
		t.Fatal(err)
	}
	rv := &reviewer{theme: theme, elementary: elem}
	lk := &lookup{root: theme, theme: set, aliases: buildAliasIndex(theme, []string{"actions"}), ids: map[string]bool{"doc": true}}
	it := rv.item(candidate{
		ID: "doc", Kind: "same-name", Source: "doc", Category: "actions",
		Links: []candidateLink{
			{Names: []string{"doc"}, From: "actions/16/doc.svg", Sizes: []string{"16"}},
			{Names: []string{"doc"}, From: "actions/24/doc.svg", Sizes: []string{"24", "22"}},
			{Names: []string{"doc-symbolic"}, From: "actions/symbolic/doc-symbolic.svg", Sizes: []string{"symbolic"}},
		},
	}, lk)

	var buckets []string
	for _, r := range it.Rows {
		buckets = append(buckets, r.Name+"@"+r.Bucket)
	}
	if want := []string{"doc@16", "doc@22", "doc@24", "doc-symbolic@symbolic"}; !reflect.DeepEqual(buckets, want) {
		t.Errorf("rows = %v, want %v", buckets, want)
	}
	if p := it.Rows[1].Proposed; p.Note != "→ places/24/folder.svg · 24px art, scaled down" || p.Style != "color" {
		t.Errorf("22 proposed = %+v", p)
	}
	if p := it.Rows[3].Proposed; p.Style != "mono" || p.Note != "recolor to ColorScheme-Text" {
		t.Errorf("symbolic proposed = %+v", p)
	}
	if want := []aliasUse{{Name: "doc-alias", Buckets: []string{"22"}}}; !reflect.DeepEqual(it.Also, want) {
		t.Errorf("also = %+v, want %+v", it.Also, want)
	}
	wantPreview := []string{
		"names: doc, doc-symbolic · category: actions",
		"16 ← actions/16/doc.svg",
		"24, 22 ← actions/24/doc.svg (→ places/24/folder.svg)",
		"doc-symbolic: symbolic ← actions/symbolic/doc-symbolic.svg (recolor to ColorScheme-Text)",
		"also changes: doc-alias (22)",
	}
	if it.Title != "`doc` — elementary overlay" || !reflect.DeepEqual(it.Preview, wantPreview) {
		t.Errorf("task = %q\n%s\nwant\n%s", it.Title, strings.Join(it.Preview, "\n"), strings.Join(wantPreview, "\n"))
	}
}

// The decide endpoint only accepts JSON (a cross-site form post is refused),
// validates the candidate, and writes the decisions file.
func TestHandleDecide(t *testing.T) {
	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"theme/actions/16/doc.svg": monoSVG,
		"candidates.yaml": "items:\n  - id: doc\n    group: g\n    kind: same-name\n    source: doc\n    category: actions\n" +
			"    links:\n      - {names: [doc], from: actions/16/doc.svg, sizes: [16]}\n",
	})
	rv := &reviewer{
		theme:      filepath.Join(dir, "theme"),
		candidates: filepath.Join(dir, "candidates.yaml"),
		decisions:  filepath.Join(dir, "review", "decisions.md"),
	}
	mux := http.NewServeMux()
	rv.register(mux)
	post := func(ctype, body string) int {
		req := httptest.NewRequest(http.MethodPost, "/review/api/decide", strings.NewReader(body))
		req.Header.Set("Content-Type", ctype)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		return rec.Code
	}

	if code := post("application/x-www-form-urlencoded", "id=doc&verdict=approved"); code != http.StatusUnsupportedMediaType {
		t.Errorf("form post = %d, want 415", code)
	}
	if code := post("application/json", `{"id":"nope","verdict":"approved"}`); code != http.StatusNotFound {
		t.Errorf("unknown id = %d, want 404", code)
	}
	if code := post("application/json", `{"id":"doc","verdict":"maybe"}`); code != http.StatusBadRequest {
		t.Errorf("bad verdict = %d, want 400", code)
	}
	if code := post("application/json", `{"id":"doc","verdict":"approved","comment":"ship it"}`); code != http.StatusOK {
		t.Fatalf("decide = %d, want 200", code)
	}
	data, err := os.ReadFile(rv.decisions)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "## Approved\n\n- [ ] `doc` — elementary overlay\n") ||
		!strings.Contains(string(data), "  Comment: ship it\n") {
		t.Errorf("decisions file:\n%s", data)
	}
}
