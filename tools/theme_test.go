package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildThemeEndToEnd(t *testing.T) {
	root := t.TempDir()
	overlays := filepath.Join(root, "overlays")
	set := filepath.Join(overlays, "apps", "demo")
	if err := os.MkdirAll(set, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(set, "a.svg"), []byte("<svg/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	yaml := "set: demo\nlinks:\n  - src: a.svg\n    category: apps\n    names: [a]\n    sizes: [16, 24]\n"
	if err := os.WriteFile(filepath.Join(set, "overlay.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	// aliases.txt: alias b -> a at apps/16
	if err := os.WriteFile(filepath.Join(overlays, "aliases.txt"),
		[]byte("apps/16  a  b\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	theme := filepath.Join(root, "theme")
	// Pre-seed a stale file to prove the wipe happens.
	if err := os.MkdirAll(filepath.Join(theme, "stale"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(theme, "stale", "old.svg"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	tela := filepath.Join(root, "tela") // absent -> base is a no-op
	baseFiles, baseLinks, nJobs, sets, dirs, _, err := buildTheme(tela, overlays, theme)
	if err != nil {
		t.Fatalf("buildTheme: %v", err)
	}
	if baseFiles != 0 || baseLinks != 0 {
		t.Errorf("absent tela should lay no base: files=%d links=%d", baseFiles, baseLinks)
	}
	if sets != 1 || nJobs != 2 || dirs < 1 {
		t.Errorf("counts: jobs=%d sets=%d dirs=%d", nJobs, sets, dirs)
	}
	// Stale content is gone (wiped).
	if _, err := os.Stat(filepath.Join(theme, "stale")); !os.IsNotExist(err) {
		t.Errorf("stale dir survived the rebuild")
	}
	// The overlay SVG landed as a real copy at 16.
	if b, err := os.ReadFile(filepath.Join(theme, "apps", "16", "a.svg")); err != nil || string(b) != "<svg/>" {
		t.Errorf("apps/16/a.svg = %q, err %v", b, err)
	}
	// aliases.txt was copied in and the alias symlink materialized.
	if _, err := os.Stat(filepath.Join(theme, "aliases.txt")); err != nil {
		t.Errorf("aliases.txt not copied: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(theme, "apps", "16", "b.svg")); err != nil {
		t.Errorf("alias b.svg not created: %v", err)
	}
	// index.theme written with the inherit line.
	data, err := os.ReadFile(filepath.Join(theme, "index.theme"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Inherits=Breeze") {
		t.Errorf("index.theme missing Inherits=Breeze")
	}
}
