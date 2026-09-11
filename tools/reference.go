package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// scanReference scans a reference icon theme into an iconSet keyed by the
// freedesktop (context, size) of each directory, taken from the theme's
// index.theme. This makes the scan layout-agnostic: it handles the conventional
// "<context>/<size>" (Breeze, elementary) as well as size-major themes such as
// hicolor / Adwaita ("16x16/apps"), so a reference lines up with our theme
// regardless of how it lays its directories out. Scaled (@2x/@3x) directories
// are skipped. If there is no usable index.theme it falls back to the plain
// directory scan.
func scanReference(root string) (*iconSet, error) {
	data, err := os.ReadFile(filepath.Join(root, "index.theme"))
	if err != nil {
		return scanSet(root)
	}
	set := newIconSet()
	any := false
	for dir, kv := range parseIniSections(string(data)) {
		if dir == "Icon Theme" {
			continue
		}
		if sc := kv["Scale"]; sc != "" && sc != "1" {
			continue // skip @2x/@3x scaled directories
		}
		size := kv["Size"]
		if size == "" {
			continue
		}
		ctx := contextFromDir(dir)
		if ctx == "" {
			continue
		}
		entries, err := os.ReadDir(filepath.Join(root, dir))
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".svg") {
				continue
			}
			set.addDir(ctx, size, strings.TrimSuffix(e.Name(), ".svg"), dir)
			any = true
		}
	}
	if !any {
		return scanSet(root)
	}
	return set, nil
}

var sizeTokenRe = regexp.MustCompile(`^\d+(x\d+)?$`)

// isSizeToken reports whether a directory path component denotes a size rather
// than a context (16, 16x16, scalable, symbolic, 22@2x, 32@2.25x, …).
func isSizeToken(p string) bool {
	if p == "scalable" || p == "symbolic" {
		return true
	}
	if strings.ContainsAny(p, "@") {
		return true
	}
	return sizeTokenRe.MatchString(p)
}

// contextFromDir returns the context component of a theme directory path — the
// part that is not a size token ("apps" from "apps/16" or "16x16/apps").
func contextFromDir(dir string) string {
	ctx := ""
	for _, p := range strings.Split(dir, "/") {
		if p == "" || isSizeToken(p) {
			continue
		}
		ctx = p
	}
	return ctx
}

// parseIniSections parses an INI-style file (index.theme) into
// section -> key -> value.
func parseIniSections(s string) map[string]map[string]string {
	out := map[string]map[string]string{}
	var cur map[string]string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			cur = map[string]string{}
			out[line[1:len(line)-1]] = cur
			continue
		}
		if cur == nil {
			continue
		}
		if i := strings.IndexByte(line, '='); i >= 0 {
			cur[strings.TrimSpace(line[:i])] = strings.TrimSpace(line[i+1:])
		}
	}
	return out
}
