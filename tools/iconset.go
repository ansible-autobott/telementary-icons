package main

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// iconSet is a theme's icon presence: context -> name -> set of sizes. It is the
// common shape for a theme scan (a live overlay theme or a reference/base theme).
type iconSet struct {
	// m[context][name][size] == true when that svg exists.
	m map[string]map[string]map[string]bool
	// sizes[context] is the set of size buckets seen under that context.
	sizes map[string]map[string]bool
	// dirs[context][size] is the theme-relative directory the files live in.
	// Usually "<context>/<size>", but size-major themes (hicolor: "16x16/apps")
	// differ, so the real dir is recorded to build correct URLs/paths.
	dirs map[string]map[string]string
}

func newIconSet() *iconSet {
	return &iconSet{
		m:     map[string]map[string]map[string]bool{},
		sizes: map[string]map[string]bool{},
		dirs:  map[string]map[string]string{},
	}
}

// add records an icon whose files live in the conventional "<context>/<size>" dir.
func (s *iconSet) add(ctx, size, name string) {
	s.addDir(ctx, size, name, ctx+"/"+size)
}

// addDir records an icon that lives in relDir (relative to the theme root).
func (s *iconSet) addDir(ctx, size, name, relDir string) {
	if s.m[ctx] == nil {
		s.m[ctx] = map[string]map[string]bool{}
		s.sizes[ctx] = map[string]bool{}
		s.dirs[ctx] = map[string]string{}
	}
	if s.m[ctx][name] == nil {
		s.m[ctx][name] = map[string]bool{}
	}
	s.m[ctx][name][size] = true
	s.sizes[ctx][size] = true
	s.dirs[ctx][size] = relDir
}

// relDir returns the theme-relative directory for (context, size), defaulting to
// the conventional "<context>/<size>".
func (s *iconSet) relDir(ctx, size string) string {
	if d := s.dirs[ctx]; d != nil {
		if r := d[size]; r != "" {
			return r
		}
	}
	return ctx + "/" + size
}

func (s *iconSet) has(ctx, name, size string) bool {
	return s.m[ctx] != nil && s.m[ctx][name] != nil && s.m[ctx][name][size]
}

// scanSet walks a freedesktop theme root (<context>/<size>/<name>.svg) into an
// iconSet. Symlinked size directories (the @2x / @3x aliases) are skipped so
// they do not appear as duplicate size buckets, but symlinked *svg files*
// (name aliases such as edit-undo -> undo) are kept — a desktop can request an
// icon by any of those names, so they are part of the requestable set.
func scanSet(root string) (*iconSet, error) {
	set := newIconSet()
	ctxDirs, err := realDirs(root)
	if err != nil {
		return nil, err
	}
	for _, ctx := range ctxDirs {
		ctxPath := filepath.Join(root, ctx)
		sizes, err := realDirs(ctxPath)
		if err != nil {
			return nil, err
		}
		for _, size := range sizes {
			entries, err := os.ReadDir(filepath.Join(ctxPath, size))
			if err != nil {
				return nil, err
			}
			for _, e := range entries {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".svg") {
					continue
				}
				set.add(ctx, size, strings.TrimSuffix(e.Name(), ".svg"))
			}
		}
	}
	return set, nil
}

// realDirs returns the names of the immediate subdirectories of dir, excluding
// symlinks (freedesktop themes symlink @2x/@3x buckets to their base).
func realDirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.Type()&os.ModeSymlink != 0 || !e.IsDir() {
			continue
		}
		out = append(out, e.Name())
	}
	sort.Strings(out)
	return out, nil
}

// sortSizes orders size buckets numerically ascending, with non-numeric
// buckets (e.g. "symbolic", "scalable") after the numeric ones, alphabetical.
func sortSizes(sizes []string) {
	sort.Slice(sizes, func(i, j int) bool {
		ni, ei := strconv.Atoi(sizes[i])
		nj, ej := strconv.Atoi(sizes[j])
		switch {
		case ei == nil && ej == nil:
			return ni < nj
		case ei == nil:
			return true // numeric before non-numeric
		case ej == nil:
			return false
		default:
			return sizes[i] < sizes[j]
		}
	})
}
