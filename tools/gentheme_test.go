package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenThemeWritesIndex(t *testing.T) {
	theme := t.TempDir()
	dir := filepath.Join(theme, "apps", "48")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "foo.svg"), []byte("<svg/>"), 0o644); err != nil {
		t.Fatal(err)
	}

	dirs, _, err := genTheme(theme)
	if err != nil {
		t.Fatalf("genTheme: %v", err)
	}
	if dirs < 1 {
		t.Errorf("dirs = %d, want >= 1", dirs)
	}
	data, err := os.ReadFile(filepath.Join(theme, "index.theme"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "Inherits=Breeze") {
		t.Errorf("index.theme missing Inherits=Breeze")
	}
	if !strings.Contains(s, "apps/48") || !strings.Contains(s, "[apps/48]") {
		t.Errorf("index.theme missing apps/48 directory/section:\n%s", s)
	}
}

func TestGenThemeScalableAndPanel(t *testing.T) {
	theme := t.TempDir()
	writeSVG := func(rel string) {
		p := filepath.Join(theme, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("<svg/>"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeSVG("apps/scalable/firefox.svg")
	writeSVG("panel/16/battery.svg")

	if _, _, err := genTheme(theme); err != nil {
		t.Fatalf("genTheme: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(theme, "index.theme"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	for _, want := range []string{
		"[apps/scalable]",
		"MinSize=16",
		"MaxSize=512",
		"apps/scalable@2x", // scaled dir emitted
		"[panel/16]",
		"Context=Status", // panel maps to Status
	} {
		if !strings.Contains(s, want) {
			t.Errorf("index.theme missing %q:\n%s", want, s)
		}
	}
}
