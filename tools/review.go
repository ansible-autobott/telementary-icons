package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// The overlay review page (/review on the `view` server) walks
// review/candidates.yaml one proposal at a time: for every size bucket it shows
// what KDE resolves today next to the elementary art that would replace it, and
// records an approve/reject verdict plus a free-form comment in
// review/decisions.md (see decisions.go) for an implementer to act on.

//go:embed review.html
var reviewPage []byte

// candidateFile is review/candidates.yaml.
type candidateFile struct {
	Origin string      `yaml:"origin"`
	Items  []candidate `yaml:"items"`
}

// candidate is one proposed overlay: elementary art (links) published under a
// theme icon name (ID) that today renders as a monochrome glyph.
type candidate struct {
	ID       string          `yaml:"id" json:"id"`
	Group    string          `yaml:"group" json:"group"`
	Kind     string          `yaml:"kind" json:"kind"`     // same-name | substitute
	Source   string          `yaml:"source" json:"source"` // elementary icon name the art comes from
	Category string          `yaml:"category" json:"category"`
	Usage    int             `yaml:"usage" json:"usage"`
	UsedBy   []string        `yaml:"used_by" json:"usedBy"`
	Links    []candidateLink `yaml:"links" json:"links"`
	Notes    []string        `yaml:"notes" json:"notes"`
}

// candidateLink is an overlay.yaml link without the local src: From is the
// elementary source relative to the elementary theme root.
type candidateLink struct {
	Names []string `yaml:"names" json:"names"`
	From  string   `yaml:"from" json:"from"`
	Sizes []string `yaml:"sizes" json:"sizes"`
}

func loadCandidates(file string) (*candidateFile, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cf candidateFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("parse %s: %w", file, err)
	}
	seen := map[string]bool{}
	for _, c := range cf.Items {
		if c.ID == "" || seen[c.ID] {
			return nil, fmt.Errorf("%s: missing or duplicate id %q", file, c.ID)
		}
		seen[c.ID] = true
	}
	return &cf, nil
}

// iconRef is one rendered icon in a review row.
type iconRef struct {
	URL    string `json:"url,omitempty"`  // image URL; empty when there is nothing to show
	Path   string `json:"path,omitempty"` // theme-, base- or elementary-relative path
	Note   string `json:"note,omitempty"`
	Style  string `json:"style,omitempty"`  // mono | gray | color (see svgStyle)
	Source string `json:"source,omitempty"` // theme | base | elem
}

// reviewRow is one name × bucket of a candidate: the icon KDE resolves today
// and the elementary art proposed for it.
type reviewRow struct {
	Name     string  `json:"name"`
	Bucket   string  `json:"bucket"`
	Current  iconRef `json:"current"`
	Proposed iconRef `json:"proposed"`
}

// aliasUse is a theme name whose alias symlink resolves through a file the
// candidate replaces, so it would switch to the elementary art too.
type aliasUse struct {
	Name    string   `json:"name"`
	Buckets []string `json:"buckets"`
	OwnItem bool     `json:"ownItem"` // the name is a candidate of its own
}

type reviewItem struct {
	candidate
	Rows []reviewRow `json:"rows"`
	Also []aliasUse  `json:"also"`
	// Title and Preview are the decisions-file task this item writes (checkbox
	// and comment aside).
	Title   string   `json:"title"`
	Preview []string `json:"preview"`
}

// reviewer serves the review page and its API.
type reviewer struct {
	theme      string // generated theme root (output/)
	elementary string // elementary theme root the candidates' `from` paths are relative to
	candidates string // review/candidates.yaml
	decisions  string // review/decisions.md
	// base returns the inherited base theme (Breeze) used for names the theme
	// lacks; nil set when unavailable.
	base func() (set *iconSet, name, root string)

	mu sync.Mutex // serializes the decisions file's read-modify-write
}

func (rv *reviewer) register(mux *http.ServeMux) {
	mux.HandleFunc("GET /review", func(w http.ResponseWriter, r *http.Request) {
		setNoStore(w)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(reviewPage)
	})
	mux.HandleFunc("GET /review/api/state", rv.handleState)
	mux.HandleFunc("POST /review/api/decide", rv.handleDecide)
	mux.HandleFunc("GET /review/img/{src}/{path...}", rv.handleImg)
}

type reviewState struct {
	Candidates string              `json:"candidates"`
	Decisions  string              `json:"decisions"`
	Elementary string              `json:"elementary"`
	Base       string              `json:"base"`
	Warnings   []string            `json:"warnings"`
	Items      []*reviewItem       `json:"items"`
	Decided    map[string]decision `json:"decided"`
}

// load reads the candidates and decisions and builds every item against the
// current theme. The theme is re-scanned each time, so a `make theme` run shows
// up on reload.
func (rv *reviewer) load() (*reviewState, *decisionDoc, []proposal, error) {
	cf, err := loadCandidates(rv.candidates)
	if err != nil {
		return nil, nil, nil, err
	}
	theme, err := scanSet(rv.theme)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("scan theme: %w", err)
	}
	st := &reviewState{Candidates: rv.candidates, Decisions: rv.decisions, Elementary: rv.elementary}
	lk := &lookup{root: rv.theme, theme: theme, ids: map[string]bool{}}
	if rv.base != nil {
		lk.base, lk.baseName, lk.baseRoot = rv.base()
	}
	st.Base = lk.baseName
	if lk.base == nil {
		st.Warnings = append(st.Warnings, "no Breeze base theme found: names missing from the theme show as empty")
	}
	switch fi, err := os.Stat(rv.elementary); {
	case rv.elementary == "":
		st.Warnings = append(st.Warnings, "elementary theme root not set (pass -elementary, or ELEMENTARY= to make): proposed art cannot be shown")
	case err != nil || !fi.IsDir():
		st.Warnings = append(st.Warnings, "elementary theme root not found: "+rv.elementary)
	}

	cats := map[string]bool{}
	for _, c := range cf.Items {
		lk.ids[c.ID] = true
		cats[c.Category] = true
	}
	lk.aliases = buildAliasIndex(rv.theme, keys(cats))

	var props []proposal
	for _, c := range cf.Items {
		it := rv.item(c, lk)
		st.Items = append(st.Items, it)
		props = append(props, proposal{ID: it.ID, Title: it.Title, Lines: it.Preview})
	}
	doc, err := loadDecisions(rv.decisions)
	if err != nil {
		return nil, nil, nil, err
	}
	st.Decided = doc.decisionsFor(props)
	return st, doc, props, nil
}

func (rv *reviewer) handleState(w http.ResponseWriter, r *http.Request) {
	rv.mu.Lock()
	st, _, _, err := rv.load()
	rv.mu.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, st)
}

type decideRequest struct {
	ID      string `json:"id"`
	Verdict string `json:"verdict"` // verdictApproved | verdictRejected | "" to clear
	Comment string `json:"comment"`
}

func (rv *reviewer) handleDecide(w http.ResponseWriter, r *http.Request) {
	// JSON only: a cross-site form post cannot send application/json without a
	// CORS preflight, which this server never grants.
	if mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type")); mt != "application/json" {
		http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	var req decideRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Verdict != "" && req.Verdict != verdictApproved && req.Verdict != verdictRejected {
		http.Error(w, "verdict must be approved, rejected or empty", http.StatusBadRequest)
		return
	}

	rv.mu.Lock()
	defer rv.mu.Unlock()
	st, doc, props, err := rv.load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	known := false
	for _, p := range props {
		known = known || p.ID == req.ID
	}
	if !known {
		http.Error(w, "unknown candidate "+strconv.Quote(req.ID), http.StatusNotFound)
		return
	}
	comment := strings.TrimSpace(strings.ReplaceAll(req.Comment, "\r\n", "\n"))
	doc.set(req.ID, req.Verdict, comment)
	if err := writeFileAtomic(rv.decisions, []byte(doc.render(props))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("review: %s %s -> %s", req.ID, orDefault(req.Verdict, "cleared"), rv.decisions)
	st.Decided = doc.decisionsFor(props)
	writeJSON(w, map[string]any{"decided": st.Decided})
}

var fgRe = regexp.MustCompile(`^[0-9a-fA-F]{6}$`)

// handleImg serves an SVG from the theme, the base theme or the elementary
// theme. ?fg=rrggbb paints monochrome icons in that color (see recolorSVG), so
// the page can preview them on light and dark backgrounds.
func (rv *reviewer) handleImg(w http.ResponseWriter, r *http.Request) {
	var root string
	switch r.PathValue("src") {
	case "theme":
		root = rv.theme
	case "elem":
		root = rv.elementary
	case "base":
		if rv.base != nil {
			_, _, root = rv.base()
		}
	}
	p := r.PathValue("path")
	if root == "" || !fs.ValidPath(p) || path.Ext(p) != ".svg" {
		http.NotFound(w, r)
		return
	}
	data, err := fs.ReadFile(os.DirFS(root), p)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if fg := r.URL.Query().Get("fg"); fgRe.MatchString(fg) {
		data = recolorSVG(data, "#"+fg)
	}
	setNoStore(w)
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(data)
}

func writeJSON(w http.ResponseWriter, v any) {
	setNoStore(w)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("review:", err)
	}
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// item builds a candidate's review rows, the aliases it would carry along, and
// the decisions-file task it writes.
func (rv *reviewer) item(c candidate, lk *lookup) *reviewItem {
	it := &reviewItem{candidate: c}
	for _, l := range c.Links {
		for _, name := range l.Names {
			for _, size := range l.Sizes {
				it.Rows = append(it.Rows, reviewRow{
					Name:     name,
					Bucket:   size,
					Current:  lk.current(c.Category, name, size),
					Proposed: rv.proposed(l.From, size),
				})
			}
		}
	}
	// Group rows by name (first-seen order), buckets ascending within a name.
	nameRank := map[string]int{}
	for _, r := range it.Rows {
		if _, ok := nameRank[r.Name]; !ok {
			nameRank[r.Name] = len(nameRank)
		}
	}
	sort.SliceStable(it.Rows, func(i, j int) bool {
		a, b := it.Rows[i], it.Rows[j]
		if a.Name != b.Name {
			return nameRank[a.Name] < nameRank[b.Name]
		}
		return sizeLess(a.Bucket, b.Bucket)
	})
	it.Also = lk.aliasesOf(c.Category, it.Rows)
	it.Title, it.Preview = rv.proposal(it)
	return it
}

// proposal renders the decisions-file task for an item: its title and the
// generated description lines (names, one line per overlay link, notes, and
// the aliases that follow).
func (rv *reviewer) proposal(it *reviewItem) (string, []string) {
	title := "`" + it.ID + "` — elementary overlay"
	if it.Source != "" && it.Source != it.ID {
		title += " from `" + it.Source + "`"
	}
	var names []string
	seen := map[string]bool{}
	for _, r := range it.Rows {
		if !seen[r.Name] {
			seen[r.Name] = true
			names = append(names, r.Name)
		}
	}
	lines := []string{"names: " + strings.Join(names, ", ") + " · category: " + it.Category}
	for _, l := range it.Links {
		line := strings.Join(l.Sizes, ", ") + " ← " + l.From
		if len(l.Names) != 1 || l.Names[0] != it.ID {
			line = strings.Join(l.Names, ", ") + ": " + line
		}
		var extra []string
		if real := rv.resolveElementary(l.From); real != "" && real != l.From {
			extra = append(extra, "→ "+real)
		}
		if contains(l.Sizes, "symbolic") {
			extra = append(extra, "recolor to ColorScheme-Text")
		}
		if len(extra) > 0 {
			line += " (" + strings.Join(extra, "; ") + ")"
		}
		lines = append(lines, line)
	}
	for _, n := range it.Notes {
		lines = append(lines, "note: "+n)
	}
	if len(it.Also) > 0 {
		var also []string
		for _, a := range it.Also {
			also = append(also, a.Name+" ("+strings.Join(a.Buckets, ", ")+")")
		}
		lines = append(lines, "also changes: "+strings.Join(also, ", "))
	}
	return title, lines
}

// resolveElementary returns the elementary-relative path an elementary file
// resolves to (following symlinks), or "" when it cannot be resolved.
func (rv *reviewer) resolveElementary(from string) string {
	if rv.elementary == "" {
		return ""
	}
	real, err := filepath.EvalSymlinks(filepath.Join(rv.elementary, filepath.FromSlash(from)))
	if err != nil {
		return ""
	}
	root, err := filepath.EvalSymlinks(rv.elementary)
	if err != nil {
		return ""
	}
	rel, err := filepath.Rel(root, real)
	if err != nil {
		return ""
	}
	return filepath.ToSlash(rel)
}

// proposed describes the elementary art proposed for a bucket.
func (rv *reviewer) proposed(from, bucket string) iconRef {
	ref := iconRef{Path: from, Source: "elem"}
	real := rv.resolveElementary(from)
	if real == "" {
		if rv.elementary == "" {
			ref.Note = "elementary theme root not set"
		} else {
			ref.Note = "missing from the elementary theme"
		}
		return ref
	}
	ref.URL = "/review/img/elem/" + from
	var notes []string
	if real != from {
		notes = append(notes, "→ "+real)
	}
	src := "" // the elementary size dir: <context>/<size>/<name>.svg
	if parts := strings.Split(from, "/"); len(parts) == 3 {
		src = parts[1]
	}
	switch {
	case bucket == "symbolic":
		notes = append(notes, "recolor to ColorScheme-Text")
		ref.Style = "mono"
	case src != "" && src != bucket:
		dir := "up"
		if atoiOr(src, 0) > atoiOr(bucket, 0) {
			dir = "down"
		}
		notes = append(notes, src+"px art, scaled "+dir)
	}
	if ref.Style == "" {
		ref.Style = svgStyleOf(filepath.Join(rv.elementary, filepath.FromSlash(real)))
	}
	ref.Note = strings.Join(notes, " · ")
	return ref
}

// lookup resolves what KDE shows today for a name at a size bucket.
type lookup struct {
	root     string
	theme    *iconSet
	base     *iconSet
	baseName string
	baseRoot string
	aliases  map[string][]aliasRef
	ids      map[string]bool // candidate ids
}

// current emulates the freedesktop lookup closely enough for a side-by-side:
// the exact bucket in the candidate's category, then in any other context, then
// the nearest bucket the theme has the name at (KDE scales it), and only when
// the theme lacks the name entirely the base theme it inherits (Breeze).
func (lk *lookup) current(cat, name, bucket string) iconRef {
	if ref, ok := lk.themeAt(cat, name, bucket); ok {
		return ref
	}
	for _, ctx := range sortedKeys(lk.theme.m) {
		if ctx == cat {
			continue
		}
		if ref, ok := lk.themeAt(ctx, name, bucket); ok {
			ref.Note = joinNotes("from "+ctx+"/"+bucket, ref.Note)
			return ref
		}
	}
	if b, ctx, ok := nearestBucket(lk.theme, name, bucket); ok {
		ref, _ := lk.themeAt(ctx, name, b)
		ref.Note = joinNotes("no "+bucket+": KDE scales "+ctx+"/"+b, ref.Note)
		return ref
	}
	for _, ctx := range sortedKeys(lk.theme.m) {
		if ref, ok := lk.themeAt(ctx, name, "scalable"); ok {
			ref.Note = joinNotes("from "+ctx+"/scalable", ref.Note)
			return ref
		}
	}
	if lk.theme.hasName(name) {
		return iconRef{Note: "not at this size in the theme"}
	}
	if lk.base != nil {
		if ref, ok := lk.baseAt(name, bucket); ok {
			return ref
		}
	}
	return iconRef{Note: "no icon in the theme or its base"}
}

func (lk *lookup) themeAt(ctx, name, bucket string) (iconRef, bool) {
	if !lk.theme.has(ctx, name, bucket) {
		return iconRef{}, false
	}
	rel := ctx + "/" + bucket + "/" + name + ".svg"
	abs := filepath.Join(lk.root, filepath.FromSlash(rel))
	ref := iconRef{URL: "/review/img/theme/" + rel, Path: rel, Source: "theme", Style: svgStyleOf(abs)}
	if target, err := os.Readlink(abs); err == nil {
		ref.Note = "alias of " + strings.TrimSuffix(target, ".svg")
	}
	return ref, true
}

// baseAt finds name in the base theme: the exact bucket in any context, else
// the nearest size (for a symbolic name, the smallest size it has).
func (lk *lookup) baseAt(name, bucket string) (iconRef, bool) {
	var ctx, size string
	for _, c := range sortedKeys(lk.base.m) {
		if lk.base.has(c, name, bucket) {
			ctx, size = c, bucket
			break
		}
	}
	if ctx == "" {
		want := bucket
		if bucket == "symbolic" {
			want = "0"
		}
		b, c, ok := nearestBucket(lk.base, name, want)
		if !ok {
			return iconRef{}, false
		}
		ctx, size = c, b
	}
	dir := lk.base.relDir(ctx, size)
	return iconRef{
		URL:    "/review/img/base/" + dir + "/" + name + ".svg",
		Path:   dir + "/" + name + ".svg",
		Note:   "not in the theme: " + lk.baseName + " fallback",
		Source: "base",
		Style:  svgStyleOf(filepath.Join(lk.baseRoot, dir, name+".svg")),
	}, true
}

// nearestBucket returns the numeric bucket (and a context holding it) closest
// to want that s has name at, preferring the larger on a tie (downscaling
// looks better than upscaling).
func nearestBucket(s *iconSet, name, want string) (bucket, ctx string, ok bool) {
	w, err := strconv.Atoi(want)
	if err != nil {
		return "", "", false
	}
	best, bestCtx, bestDist := -1, "", 1<<30
	for _, c := range sortedKeys(s.m) {
		for size := range s.m[c][name] {
			n, err := strconv.Atoi(size)
			if err != nil {
				continue
			}
			d := n - w
			if d < 0 {
				d = -d
			}
			if d < bestDist || (d == bestDist && n > best) {
				best, bestCtx, bestDist = n, c, d
			}
		}
	}
	if best < 0 {
		return "", "", false
	}
	return strconv.Itoa(best), bestCtx, true
}

func (s *iconSet) hasName(name string) bool {
	for _, names := range s.m {
		if names[name] != nil {
			return true
		}
	}
	return false
}

// aliasesOf lists the theme names that follow the files this candidate
// replaces: overlaying a file swaps it in place, so every alias symlink whose
// chain passes through it resolves to the new art too. The candidate's own
// names are left out.
func (lk *lookup) aliasesOf(cat string, rows []reviewRow) []aliasUse {
	own := map[string]bool{}
	for _, r := range rows {
		own[r.Name] = true
	}
	byName := map[string]map[string]bool{}
	for _, r := range rows {
		target := filepath.Join(lk.root, cat, r.Bucket, r.Name+".svg")
		for _, a := range lk.aliases[target] {
			if own[a.name] {
				continue
			}
			if byName[a.name] == nil {
				byName[a.name] = map[string]bool{}
			}
			byName[a.name][a.bucket] = true
		}
	}
	var out []aliasUse
	for _, name := range sortedKeys(byName) {
		b := keys(byName[name])
		sort.Slice(b, func(i, j int) bool { return sizeLess(b[i], b[j]) })
		out = append(out, aliasUse{Name: name, Buckets: b, OwnItem: lk.ids[name]})
	}
	return out
}

type aliasRef struct{ name, bucket string }

// buildAliasIndex maps every file in the given contexts to the alias symlinks
// whose chain passes through it (directly or via other aliases).
func buildAliasIndex(root string, cats []string) map[string][]aliasRef {
	idx := map[string][]aliasRef{}
	for _, cat := range cats {
		buckets, err := realDirs(filepath.Join(root, cat))
		if err != nil {
			continue
		}
		for _, b := range buckets {
			dir := filepath.Join(root, cat, b)
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, e := range entries {
				if e.Type()&os.ModeSymlink == 0 || !strings.HasSuffix(e.Name(), ".svg") {
					continue
				}
				ref := aliasRef{strings.TrimSuffix(e.Name(), ".svg"), b}
				cur := filepath.Join(dir, e.Name())
				for hop := 0; hop < 8; hop++ {
					t, err := os.Readlink(cur)
					if err != nil {
						break
					}
					if !filepath.IsAbs(t) {
						t = filepath.Join(filepath.Dir(cur), t)
					}
					t = filepath.Clean(t)
					idx[t] = append(idx[t], ref)
					if fi, err := os.Lstat(t); err != nil || fi.Mode()&os.ModeSymlink == 0 {
						break
					}
					cur = t
				}
			}
		}
	}
	return idx
}

var (
	svgHexRe     = regexp.MustCompile(`#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})\b`)
	svgStyleElRe = regexp.MustCompile(`(?s)<style[^>]*>.*?</style>`)
)

// svgStyle classifies an icon's coloring: "mono" when it paints with the color
// scheme's currentColor (optionally with a colored accent), "gray" when it only
// uses fixed grays, "color" otherwise. The <style> block is ignored: the
// ColorScheme-* classes defined there are just defaults.
func svgStyle(svg []byte) string {
	s := string(svg)
	chromatic := 0
	for _, m := range svgHexRe.FindAllStringSubmatch(svgStyleElRe.ReplaceAllString(s, ""), -1) {
		h := m[1]
		if len(h) == 3 {
			h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
		}
		r, _ := strconv.ParseUint(h[0:2], 16, 8)
		g, _ := strconv.ParseUint(h[2:4], 16, 8)
		b, _ := strconv.ParseUint(h[4:6], 16, 8)
		if max(r, g, b)-min(r, g, b) > 30 {
			chromatic++
		}
	}
	lower := strings.ToLower(s)
	recolors := strings.Contains(lower, "currentcolor")
	switch {
	case chromatic > 0 && !(recolors && chromatic <= 2 && !strings.Contains(lower, "gradient")):
		return "color"
	case recolors:
		return "mono"
	case chromatic > 0:
		return "color"
	default:
		return "gray"
	}
}

func svgStyleOf(file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		return ""
	}
	return svgStyle(data)
}

var (
	schemeTextRe = regexp.MustCompile(`(\.ColorScheme-Text\s*\{[^}]*?\bcolor\s*:\s*)#[0-9a-fA-F]{3,6}`)
	gpaPaintRe   = regexp.MustCompile(`url\(#gpa:[\w-]+\)(\s*#[0-9a-fA-F]{3,6})?`)
	svgOpenRe    = regexp.MustCompile(`<svg\b[^>]*>`)
)

// recolorSVG paints a monochrome icon in fg the way a desktop does for its
// color scheme: KDE's ColorScheme-Text stylesheet color, GNOME's gpa:foreground
// paint (elementary's symbolic icons), and plain currentColor, which inherits
// from the root element. Full-color art is left as is.
func recolorSVG(svg []byte, fg string) []byte {
	s := schemeTextRe.ReplaceAllString(string(svg), "${1}"+fg)
	s = gpaPaintRe.ReplaceAllString(s, fg)
	if loc := svgOpenRe.FindStringIndex(s); loc != nil && !strings.Contains(s[loc[0]:loc[1]], " color=") {
		s = s[:loc[0]+len("<svg")] + ` color="` + fg + `"` + s[loc[0]+len("<svg"):]
	}
	return []byte(s)
}

func atoiOr(s string, def int) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}

func joinNotes(a, b string) string {
	if a == "" || b == "" {
		return a + b
	}
	return a + " · " + b
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// defaultReviewFile anchors a review/ file at the repo root, like defaultTheme.
func defaultReviewFile(name string) string {
	if r := repoRoot(); r != "" {
		return filepath.Join(r, "review", name)
	}
	return filepath.Join("review", name)
}
