package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// job is a single planned overlay action at dst. A copy job has src set (the
// source file to read) and link empty; a symlink job has link set (the relative
// symlink target) and src empty.
type job struct {
	src  string
	dst  string
	link string
}

// overlayConfig is one design set's overlay.yaml: how its source SVGs map into
// the icon theme. Origin is optional provenance and does not affect the plan.
type overlayConfig struct {
	Set    string         `yaml:"set"`
	Origin *overlayOrigin `yaml:"origin,omitempty"`
	Links  []overlayLink  `yaml:"links"`
}

// overlayOrigin records which upstream theme a design set's artwork came from.
// Documentation only. No absolute path is stored — that is machine-specific and
// may change; each link's `from` gives the theme-relative source path.
type overlayOrigin struct {
	Theme string `yaml:"theme"`
}

// overlayLink maps one source SVG to a freedesktop icon: a category, one or more
// names (aliases), and the size buckets to publish it at. Sizes are strings so
// numeric buckets ("16") and the symbolic bucket ("symbolic") share one field.
// From is optional provenance: the source path under the set's origin.root.
type overlayLink struct {
	Src      string   `yaml:"src"`
	From     string   `yaml:"from,omitempty"`
	Category string   `yaml:"category"`
	Names    []string `yaml:"names"`
	Sizes    []string `yaml:"sizes"`
}

// parseOverlay decodes an overlay.yaml. Unquoted numeric sizes decode straight
// into the []string field (yaml.v3 coerces scalar values), so 16 becomes "16".
func parseOverlay(data []byte) (overlayConfig, error) {
	var cfg overlayConfig
	err := yaml.Unmarshal(data, &cfg)
	return cfg, err
}

// linkJobs expands a config into overlay actions under dstRoot. For each
// (link, name), the source SVG is copied into its first listed size
// (dstRoot/<category>/<size0>/<name>.svg) and every further size becomes a
// relative symlink to that first copy (../<size0>/<name>.svg) — so one SVG is
// stored once per name and shared across sizes. setDir anchors the relative
// src paths; dstRoot is the target icon theme root.
func linkJobs(setDir string, cfg overlayConfig, dstRoot string) []job {
	var jobs []job
	for _, l := range cfg.Links {
		src := filepath.Join(setDir, l.Src)
		for _, name := range l.Names {
			for i, size := range l.Sizes {
				dst := filepath.Join(dstRoot, l.Category, size, name+".svg")
				if i == 0 {
					jobs = append(jobs, job{src: src, dst: dst})
					continue
				}
				link := filepath.Join("..", l.Sizes[0], name+".svg")
				jobs = append(jobs, job{dst: dst, link: link})
			}
		}
	}
	return jobs
}

// planOverlays reads every overlay.yaml under overlaysDir and expands them into
// the jobs that place each design SVG into themeDir. Returns the jobs and the
// number of sets (overlay.yaml files) found.
func planOverlays(overlaysDir, themeDir string) ([]job, int, error) {
	cfgPaths, err := findOverlayConfigs(overlaysDir)
	if err != nil {
		return nil, 0, err
	}
	var jobs []job
	sets := 0
	for _, cfgPath := range cfgPaths {
		setDir := filepath.Dir(cfgPath)
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			return nil, 0, fmt.Errorf("read %s: %w", cfgPath, err)
		}
		cfg, err := parseOverlay(data)
		if err != nil {
			return nil, 0, fmt.Errorf("parse %s: %w", cfgPath, err)
		}
		sets++
		jobs = append(jobs, linkJobs(setDir, cfg, themeDir)...)
	}
	return jobs, sets, nil
}

// findOverlayConfigs returns every overlay.yaml under root, at any depth, sorted.
// Design sets are grouped in category subfolders (overlays/<category>/<set>/),
// so discovery walks the whole tree rather than a single level. A missing root
// yields no configs and no error.
func findOverlayConfigs(root string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil // missing root (or vanished entry): nothing to collect
			}
			return err
		}
		if !d.IsDir() && d.Name() == "overlay.yaml" {
			out = append(out, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}

func defaultOverlays() string {
	if r := repoRoot(); r != "" {
		return filepath.Join(r, "overlays")
	}
	return "overlays"
}

// applyJobs writes each planned job to disk: copies for src jobs, relative
// symlinks (replacing anything already at dst) for link jobs.
func applyJobs(jobs []job) error {
	for _, p := range jobs {
		if err := os.MkdirAll(filepath.Dir(p.dst), 0o755); err != nil {
			return err
		}
		if p.link != "" {
			// Overlay overwrites: drop whatever is at dst, then link.
			if err := os.Remove(p.dst); err != nil && !os.IsNotExist(err) {
				return err
			}
			if err := os.Symlink(p.link, p.dst); err != nil {
				return err
			}
			continue
		}
		// Copy job: remove any existing dst first — it may be a Tela alias
		// symlink from the base, and copyFile writes THROUGH a symlink onto its
		// target. Removing it makes the overlay replace the alias, not clobber
		// the file the alias points at.
		if err := os.Remove(p.dst); err != nil && !os.IsNotExist(err) {
			return err
		}
		if err := copyFile(p.src, p.dst); err != nil {
			return fmt.Errorf("copy %s: %w", p.src, err)
		}
	}
	return nil
}

// cmdOverlay implements: icons overlay [flags]
//
// It reads every overlays/<category>/<set>/overlay.yaml (at any depth) and
// copies each design's source SVG straight into
// output/<category>/<size>/<name>.svg, fanned across the names and sizes it
// lists — overwriting any icon of the same name. Design sets are grouped in
// category subfolders for navigation; the target context is the link's own
// `category` field, not the folder. Dry-run unless -apply. Then run gen-theme.
func cmdOverlay(args []string) {
	fs := flag.NewFlagSet("overlay", flag.ExitOnError)
	overlaysDir := fs.String("overlays", defaultOverlays(),
		"root of overlay sets, grouped in <category>/<set>/overlay.yaml")
	themeDir := fs.String("theme", defaultTheme(), "target icon theme root")
	apply := fs.Bool("apply", false, "actually copy files (default: dry run)")
	fs.Parse(args)

	jobs, sets, err := planOverlays(*overlaysDir, *themeDir)
	if err != nil {
		log.Fatalf("overlays: %v", err)
	}

	// Warn about design sources that do not exist on disk (copy jobs only).
	var missing []string
	copies, links := 0, 0
	for _, j := range jobs {
		if j.link != "" {
			links++
			continue
		}
		copies++
		if _, err := os.Stat(j.src); err != nil {
			missing = append(missing, j.src)
		}
	}

	if *apply {
		if err := applyJobs(jobs); err != nil {
			log.Fatal(err)
		}
	}

	mode := "DRY RUN (use -apply to write)"
	if *apply {
		mode = "APPLIED"
	}
	fmt.Printf("overlay %s\n", mode)
	fmt.Printf("  overlays: %s (%d set(s) with overlay.yaml)\n", *overlaysDir, sets)
	fmt.Printf("  theme:    %s\n\n", *themeDir)
	fmt.Printf("  %d file(s) into the theme (%d copied, %d symlinked)\n", len(jobs), copies, links)

	if len(missing) > 0 {
		fmt.Printf("\nWARNING: %d design source(s) missing:\n", len(missing))
		for _, m := range uniqStrings(missing) {
			fmt.Printf("  %s\n", m)
		}
	}

	if *apply {
		fmt.Printf("\nDone. Next: run `icons gen-theme` to refresh index.theme and symlinks.\n")
	} else {
		fmt.Printf("\nRe-run with -apply to write them.\n")
	}
}

// uniqStrings returns the sorted, de-duplicated elements of ss.
func uniqStrings(ss []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}
