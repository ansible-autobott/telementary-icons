package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// buildTheme regenerates themeDir end to end: wipe the output, lay the vendored
// Tela theme down as the base (mergeTela), apply every overlay.yaml on top, copy
// aliases.txt, then run the index/symlink generator. Returns counts to report.
func buildTheme(telaDir, overlaysDir, themeDir string) (baseFiles, baseLinks, nJobs, sets, dirs, scaled int, err error) {
	if err = os.RemoveAll(themeDir); err != nil {
		return
	}
	baseFiles, baseLinks, err = mergeTela(telaDir, themeDir)
	if err != nil {
		err = fmt.Errorf("merge tela: %w", err)
		return
	}
	jobs, s, e := planOverlays(overlaysDir, themeDir)
	if e != nil {
		err = e
		return
	}
	if e := applyJobs(jobs); e != nil {
		err = e
		return
	}
	aliasSrc := filepath.Join(overlaysDir, "aliases.txt")
	if _, statErr := os.Stat(aliasSrc); statErr == nil {
		if e := copyFile(aliasSrc, filepath.Join(themeDir, "aliases.txt")); e != nil {
			err = fmt.Errorf("copy aliases.txt: %w", e)
			return
		}
	}
	d, sc, e := genTheme(themeDir)
	if e != nil {
		err = e
		return
	}
	return baseFiles, baseLinks, len(jobs), s, d, sc, nil
}

// cmdTheme implements: icons theme [-tela D] [-overlays D] [-theme D]
//
// It regenerates the shipped icon theme from the overlay sets. The output dir is
// wiped first, so the result contains only the current overlays plus the
// generated index.theme / symlinks.
func cmdTheme(args []string) {
	fs := flag.NewFlagSet("theme", flag.ExitOnError)
	telaDir := fs.String("tela", defaultTela(), "vendored Tela source dir (src/ + links/)")
	overlaysDir := fs.String("overlays", defaultOverlays(), "root of overlay sets")
	themeDir := fs.String("theme", defaultTheme(), "generated theme output dir")
	fs.Parse(args)

	baseFiles, baseLinks, nJobs, sets, dirs, scaled, err := buildTheme(*telaDir, *overlaysDir, *themeDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("theme generated at %s\n", *themeDir)
	fmt.Printf("  tela base: %s (%d file(s), %d alias link(s))\n", *telaDir, baseFiles, baseLinks)
	fmt.Printf("  overlays: %s (%d set(s))\n", *overlaysDir, sets)
	fmt.Printf("  %d overlay file(s) written; index.theme: %d directories, %d scaled\n", nJobs, dirs, scaled)
}
