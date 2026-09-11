package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// defaultTela anchors the vendored Tela source dir at the repo root, so the tool
// works from any cwd.
func defaultTela() string {
	if r := repoRoot(); r != "" {
		return filepath.Join(r, "tela-icons")
	}
	return "tela-icons"
}

// transposeLinkTarget rewrites a Tela alias symlink target from Tela's
// size-major layout (<size>/<category>/<name>.svg) into this theme's
// category-major layout (<category>/<size>/<name>.svg).
//
// A Tela link at links/<size>/<cat>/<name>.svg carries a target relative to that
// directory. Resolving it against "<size>/<cat>" yields the target's Tela coords
// "<size2>/<cat2>/<name2>.svg"; the rewritten target is that file expressed
// relative to the link's new home "<cat>/<size>". Targets that don't resolve to
// exactly three components are rejected (we vendor a pinned commit — surprises
// should surface, not silently produce a broken link).
func transposeLinkTarget(size, cat, target string) (string, error) {
	coord := filepath.Clean(filepath.Join(size, cat, target))
	parts := strings.Split(coord, string(os.PathSeparator))
	if len(parts) != 3 {
		return "", fmt.Errorf("tela link %s/%s -> %q resolves to %q, want <size>/<cat>/<name>", size, cat, target, coord)
	}
	// Reject if any part is . or .. (escaping the tree)
	for _, p := range parts {
		if p == "." || p == ".." {
			return "", fmt.Errorf("tela link %s/%s -> %q resolves to %q, want <size>/<cat>/<name>", size, cat, target, coord)
		}
	}
	size2, cat2, name2 := parts[0], parts[1], parts[2]
	return filepath.Rel(filepath.Join(cat, size), filepath.Join(cat2, size2, name2))
}

// mergeTela lays the vendored Tela theme down as the base of themeDir: it
// transposes Tela's size-major src/ (real SVGs, copied) and links/ (alias
// symlinks, recreated with rewritten targets) into the category-major layout.
// src is copied first and links are laid on top, so if a name exists as both a
// real src file and a links/ alias to a different icon, the alias intentionally
// wins over the colliding real src file — matching Tela's upstream install
// order (cp src, then cp links).
// Returns how many real files and symlinks were written. A missing telaDir is
// not an error (yields 0, 0) so the pipeline degrades to overlays-only.
func mergeTela(telaDir, themeDir string) (files, links int, err error) {
	files, err = copyTelaTree(filepath.Join(telaDir, "src"), themeDir)
	if err != nil {
		return files, 0, err
	}
	links, err = linkTelaTree(filepath.Join(telaDir, "links"), themeDir)
	return files, links, err
}

// copyTelaTree copies every real file at srcRoot/<size>/<cat>/<name> to
// themeDir/<cat>/<size>/<name>. Non-3-level entries (e.g. src/index.theme) and
// symlinks are skipped.
func copyTelaTree(srcRoot, themeDir string) (int, error) {
	n := 0
	err := filepath.WalkDir(srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil // missing tree: nothing to lay down
			}
			return err
		}
		if d.IsDir() || d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}
		parts := strings.Split(rel, string(os.PathSeparator))
		if len(parts) != 3 {
			return nil
		}
		size, cat, name := parts[0], parts[1], parts[2]
		dst := filepath.Join(themeDir, cat, size, name)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := copyFile(path, dst); err != nil {
			return fmt.Errorf("copy %s: %w", path, err)
		}
		n++
		return nil
	})
	return n, err
}

// linkTelaTree recreates every symlink at linksRoot/<size>/<cat>/<name> as a
// symlink at themeDir/<cat>/<size>/<name> with a transposed relative target.
func linkTelaTree(linksRoot, themeDir string) (int, error) {
	n := 0
	err := filepath.WalkDir(linksRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if d.Type()&os.ModeSymlink == 0 {
			return nil // only symlinks are aliases
		}
		rel, err := filepath.Rel(linksRoot, path)
		if err != nil {
			return err
		}
		parts := strings.Split(rel, string(os.PathSeparator))
		if len(parts) != 3 {
			return nil
		}
		size, cat, name := parts[0], parts[1], parts[2]
		target, err := os.Readlink(path)
		if err != nil {
			return err
		}
		newTarget, err := transposeLinkTarget(size, cat, target)
		if err != nil {
			return err
		}
		dst := filepath.Join(themeDir, cat, size, name)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
			return err
		}
		if err := os.Symlink(newTarget, dst); err != nil {
			return err
		}
		n++
		return nil
	})
	return n, err
}
