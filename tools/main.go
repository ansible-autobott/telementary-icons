// Command icons is a toolbox for this repository's freedesktop icon theme.
//
// Subcommands:
//
//	icons view        serve a browsable name×size matrix of the theme (default)
//	icons overlay     apply design overlays (overlays/*/overlay.yaml) into the theme
//	icons gen-theme   regenerate output/index.theme, the @2x/@3x scaled symlinks, and aliases
//	icons theme       regenerate output/ from overlays/ (wipe + apply overlays + gen-theme)
//
// The viewer can load a base theme (e.g. KDE's Breeze) as a selectable reference
// and fallback: it defines the row and size universe, and any cell the theme does
// not override is filled with the base theme's icon (what a theme Inheriting it
// would show). Design overlay sets live under overlays/<category>/<set>/overlay.yaml.
//
// Run from anywhere inside the repo: default theme paths are anchored at the repo
// root, so cwd does not matter.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]

	// Default to "view" so a bare `icons` (or `icons -addr :x`) still serves.
	sub := "view"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		sub, args = args[0], args[1:]
	}

	switch sub {
	case "view":
		cmdView(args)
	case "overlay":
		cmdOverlay(args)
	case "gen-theme":
		cmdGenTheme(args)
	case "theme":
		cmdTheme(args)
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "icons: unknown command %q\n\n", sub)
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `icons — freedesktop icon theme toolbox

Usage:
  icons [view]                    serve the theme viewer (default)
  icons overlay                   apply design overlays into the theme, driven by overlay.yaml
  icons gen-theme                 regenerate index.theme + @2x/@3x + alias symlinks
  icons theme                     regenerate output/ from overlays/ (wipe + overlay + gen-theme)

Run "icons <command> -h" for that command's flags.
`)
}

// repoRoot walks up from the current directory looking for a .git entry and
// returns that directory, or "" if none is found. It lets the tool resolve its
// default theme paths regardless of where it is launched from — the module lives
// at tools/, so a bare "go run ." has a different cwd than the repo root the
// relative defaults assume.
func repoRoot() string {
	d, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(d, ".git")); err == nil {
			return d
		}
		parent := filepath.Dir(d)
		if parent == d {
			return ""
		}
		d = parent
	}
}

// defaultTheme anchors the common path at the repo root when it can be found,
// so the tool works from any cwd; otherwise falls back to cwd.
func defaultTheme() string {
	if r := repoRoot(); r != "" {
		return filepath.Join(r, "output")
	}
	return "output"
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}
