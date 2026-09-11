package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// symbolicSuffixes are the trailing markers of a symbolic icon variant, longest
// first so "-symbolic-rtl" is matched before "-symbolic".
var symbolicSuffixes = []string{"-symbolic-rtl", "-symbolic-ltr", "-symbolic"}

// symbolicOf reports whether name is a trailing symbolic variant and, if so,
// returns its base name and the display "size" it folds onto ("symbolic",
// "symbolic-rtl", …). Only a trailing marker counts, so names that merely
// contain "-symbolic" mid-string (e.g. emblem-symbolic-link, a real icon) are
// left alone.
func symbolicOf(name string) (base, dispSize string, ok bool) {
	for _, suf := range symbolicSuffixes {
		if strings.HasSuffix(name, suf) {
			return name[:len(name)-len(suf)], "symbolic" + strings.TrimPrefix(suf, "-symbolic"), true
		}
	}
	return "", "", false
}

// isSymbolic reports whether an icon name is a symbolic variant. Symbolic and
// base icons are kept strictly separate: symbolic names go only to the theme's
// symbolic/ bucket, base names only to the numeric buckets.
func isSymbolic(name string) bool {
	_, _, ok := symbolicOf(name)
	return ok
}

// pickSymbolicSource chooses which source size dir to copy a symbolic icon from:
// the source's own symbolic/ dir if it has one (the true symbolic art), else the
// smallest numeric size it exists at, else any.
func pickSymbolicSource(sizes map[string]bool) string {
	if sizes["symbolic"] {
		return "symbolic"
	}
	best, bestVal := "", 1<<30
	for s := range sizes {
		if v, err := strconv.Atoi(s); err == nil && v < bestVal {
			best, bestVal = s, v
		}
	}
	if best != "" {
		return best
	}
	var ss []string
	for s := range sizes {
		ss = append(ss, s)
	}
	sort.Strings(ss)
	if len(ss) > 0 {
		return ss[0]
	}
	return ""
}

// freedesktop Context value for each context directory we recognise.
var genContext = map[string]string{
	"actions":     "Actions",
	"animations":  "Animations",
	"apps":        "Applications",
	"categories":  "Categories",
	"devices":     "Devices",
	"emblems":     "Emblems",
	"emotes":      "Emotes",
	"mimetypes":   "MimeTypes",
	"panel":       "Status",
	"places":      "Places",
	"preferences": "Applications",
	"status":      "Status",
}

// sizeMeta is the per-size-bucket metadata, following current Breeze conventions.
// scaled -> also emit <size>@2x / <size>@3x dirs (listed under ScaledDirectories).
// Fixed -> exact-size match only; Scalable -> covers MinSize..MaxSize.
type sizeMeta struct {
	size    int
	typ     string
	minSize int
	maxSize int
	scaled  bool
}

var genSizes = map[string]sizeMeta{
	"16":       {size: 16, typ: "Fixed", scaled: true},
	"22":       {size: 22, typ: "Fixed", scaled: true},
	"24":       {size: 24, typ: "Fixed", scaled: true}, // GTK compat
	"32":       {size: 32, typ: "Scalable", minSize: 32, maxSize: 256, scaled: true},
	"48":       {size: 48, typ: "Scalable", minSize: 48, maxSize: 256, scaled: false},
	"64":       {size: 64, typ: "Scalable", minSize: 64, maxSize: 256, scaled: false},
	"scalable": {size: 64, typ: "Scalable", minSize: 16, maxSize: 512, scaled: true},
	"symbolic": {size: 16, typ: "Scalable", minSize: 8, maxSize: 512, scaled: false},
}

var genScales = []int{2, 3}

const genHeader = `[Icon Theme]
Name=Telementary
Comment=Custom KDE icon theme based on Tela and elementary
Inherits=Breeze
Example=folder
DisplayDepth=32
KDE-Extensions=.svg
FollowsColorScheme=true
`

// genTheme materializes alias symlinks, writes index.theme, and creates the
// @2x/@3x scaled-dir symlinks for the theme rooted at icons. Returns the count
// of base directories and scaled directories written.
func genTheme(icons string) (int, int, error) {
	applyAliases(icons)

	var directories, scaled, sections []string
	ctxNames, err := realDirs(icons)
	if err != nil {
		return 0, 0, err
	}
	for _, ctx := range ctxNames {
		if _, ok := genContext[ctx]; !ok {
			continue
		}
		ctxDir := filepath.Join(icons, ctx)
		sizeNames, err := realDirs(ctxDir)
		if err != nil {
			return 0, 0, err
		}
		sort.Slice(sizeNames, func(i, j int) bool {
			si := genSizes[sizeNames[i]].sizeOr(999)
			sj := genSizes[sizeNames[j]].sizeOr(999)
			if si != sj {
				return si < sj
			}
			return sizeNames[i] < sizeNames[j]
		})
		for _, size := range sizeNames {
			meta, ok := genSizes[size]
			if !ok {
				fmt.Fprintf(os.Stderr, "warning: unknown size dir %s/%s, skipping\n", ctx, size)
				continue
			}
			if !hasSVG(filepath.Join(ctxDir, size)) {
				continue
			}
			base := ctx + "/" + size
			directories = append(directories, base)
			sections = append(sections, section(base, ctx, meta, 1))
			if meta.scaled {
				for _, s := range genScales {
					makeScaledLink(ctxDir, size, s)
					name := fmt.Sprintf("%s/%s@%dx", ctx, size, s)
					scaled = append(scaled, name)
					sections = append(sections, section(name, ctx, meta, s))
				}
			}
		}
	}

	out := []string{genHeader, "Directories=" + strings.Join(directories, ",")}
	if len(scaled) > 0 {
		out = append(out, "ScaledDirectories="+strings.Join(scaled, ","))
	}
	out = append(out, "")
	out = append(out, sections...)

	dst := filepath.Join(icons, "index.theme")
	if err := os.WriteFile(dst, []byte(strings.Join(out, "\n")+"\n"), 0o644); err != nil {
		return 0, 0, err
	}
	return len(directories), len(scaled), nil
}

// cmdGenTheme implements: icons gen-theme [-theme D]
//
// It mirrors how upstream Breeze builds its theme: the source of truth is only
// the real .svg files under icons/<context>/<size>/ plus the alias map in
// icons/aliases.txt. index.theme, the <size>@2x / <size>@3x directory symlinks,
// and the alias symlinks are all derived here.
func cmdGenTheme(args []string) {
	fs := flag.NewFlagSet("gen-theme", flag.ExitOnError)
	themeDir := fs.String("theme", defaultTheme(), "icon theme root to regenerate")
	fs.Parse(args)

	icons, err := filepath.Abs(*themeDir)
	if err != nil {
		log.Fatal(err)
	}
	dirs, scaled, err := genTheme(icons)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s: %d directories, %d scaled\n", filepath.Join(icons, "index.theme"), dirs, scaled)
}

func (m sizeMeta) sizeOr(def int) int {
	if m.size == 0 {
		return def
	}
	return m.size
}

func section(name, ctx string, meta sizeMeta, scale int) string {
	lines := []string{"[" + name + "]", "Size=" + strconv.Itoa(meta.size)}
	if scale != 1 {
		lines = append(lines, "Scale="+strconv.Itoa(scale))
	}
	lines = append(lines, "Context="+genContext[ctx], "Type="+meta.typ)
	if meta.typ == "Scalable" {
		lines = append(lines, "MinSize="+strconv.Itoa(meta.minSize), "MaxSize="+strconv.Itoa(meta.maxSize))
	}
	lines = append(lines, "")
	return strings.Join(lines, "\n")
}

func hasSVG(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".svg") {
			return true
		}
	}
	return false
}

// makeScaledLink creates (or refreshes) a <size>@<scale>x symlink pointing at
// the base <size> directory.
func makeScaledLink(ctxDir, size string, scale int) {
	link := filepath.Join(ctxDir, fmt.Sprintf("%s@%dx", size, scale))
	if fi, err := os.Lstat(link); err == nil {
		if fi.Mode()&os.ModeSymlink != 0 {
			os.Remove(link)
		} else {
			fmt.Fprintf(os.Stderr, "warning: %s exists and is not a symlink, skipping\n", link)
			return
		}
	}
	if err := os.Symlink(size, link); err != nil {
		fmt.Fprintf(os.Stderr, "warning: symlink %s: %v\n", link, err)
	}
}

// applyAliases reads icons/aliases.txt and creates name-alias symlinks so that
// several icon names resolve to one canonical SVG. Format per line:
//
//	<context>/<size>  <canonical>  <alias> [<alias> ...]
func applyAliases(icons string) {
	path := filepath.Join(icons, "aliases.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		return // no alias map -> nothing to do
	}
	for _, raw := range strings.Split(string(data), "\n") {
		line := raw
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = line[:i]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			fmt.Fprintf(os.Stderr, "warning: bad alias line: %q\n", raw)
			continue
		}
		rel, canonical, aliases := parts[0], parts[1], parts[2:]
		dir := filepath.Join(icons, rel)
		target := canonical + ".svg"
		if _, err := os.Stat(filepath.Join(dir, target)); err != nil {
			fmt.Fprintf(os.Stderr, "warning: alias source missing: %s/%s\n", rel, target)
			continue
		}
		for _, a := range aliases {
			link := filepath.Join(dir, a+".svg")
			if fi, err := os.Lstat(link); err == nil {
				if fi.Mode()&os.ModeSymlink == 0 {
					fmt.Fprintf(os.Stderr, "warning: %s/%s.svg is a real file, not overwriting\n", rel, a)
					continue
				}
				os.Remove(link)
			}
			if err := os.Symlink(target, link); err != nil {
				fmt.Fprintf(os.Stderr, "warning: symlink %s: %v\n", link, err)
			}
		}
	}
}
