package main

import (
	"os"
	"path/filepath"
	"testing"
)

// mkTheme creates a minimal freedesktop theme dir (root/<name>/index.theme).
func mkTheme(t *testing.T, root, name string) string {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "index.theme"), []byte("[Icon Theme]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// findCell locates a rendered cell by context/name/size.
func findCell(t *testing.T, cats []category, ctx, name, size string) cell {
	t.Helper()
	for _, c := range cats {
		if c.Name != ctx {
			continue
		}
		idx := -1
		for i, s := range c.Sizes {
			if s == size {
				idx = i
			}
		}
		if idx < 0 {
			t.Fatalf("size %s not a column in %s", size, ctx)
		}
		for _, r := range c.Rows {
			if r.Name == name {
				return r.Cells[idx]
			}
		}
	}
	t.Fatalf("cell %s/%s/%s not found", ctx, size, name)
	return cell{}
}

// An override that also exists in the reference must carry the original's URL
// and path so the viewer can show it on hover.
func TestBuildOverrideCarriesOriginal(t *testing.T) {
	theme := newIconSet()
	theme.add("places", "48", "folder")
	ref := newIconSet()
	ref.add("places", "48", "folder")

	cats := build(theme, "/theme", ref, "Breeze", "/refroot")

	c := findCell(t, cats, "places", "folder", "48")
	if c.State != "override" {
		t.Fatalf("state = %q, want override", c.State)
	}
	if c.Extra {
		t.Errorf("Extra = true, want false (reference has this icon)")
	}
	if c.OrigSrc != "/ref/Breeze/places/48/folder.svg" {
		t.Errorf("OrigSrc = %q, want /ref/Breeze/places/48/folder.svg", c.OrigSrc)
	}
	if c.OrigPath != "/refroot/places/48/folder.svg" {
		t.Errorf("OrigPath = %q, want /refroot/places/48/folder.svg", c.OrigPath)
	}
}

// An override the reference lacks is "extra" and has no original to show.
func TestBuildExtraOverrideHasNoOriginal(t *testing.T) {
	theme := newIconSet()
	theme.add("apps", "48", "smb4k")
	ref := newIconSet()
	ref.add("apps", "48", "folder") // ref populates the apps/48 column but lacks smb4k

	cats := build(theme, "/theme", ref, "Breeze", "/refroot")

	c := findCell(t, cats, "apps", "smb4k", "48")
	if c.State != "override" || !c.Extra {
		t.Fatalf("cell = %+v, want override+extra", c)
	}
	if c.OrigSrc != "" {
		t.Errorf("OrigSrc = %q, want empty", c.OrigSrc)
	}
}

// referenceThemes must list themes from the system icons dir AND from the user's
// $XDG_DATA_HOME/icons, with the user copy winning on a name clash.
func TestReferenceThemesIncludesUserAndOverrides(t *testing.T) {
	sys := t.TempDir()
	dataHome := t.TempDir()

	mkTheme(t, sys, "Breeze")
	mkTheme(t, sys, "Shared") // will be overridden by the user copy
	userCustom := mkTheme(t, filepath.Join(dataHome, "icons"), "Custom")
	userShared := mkTheme(t, filepath.Join(dataHome, "icons"), "Shared")

	t.Setenv("XDG_DATA_HOME", dataHome)

	got := referenceThemes(sys)

	if len(got) != 3 {
		t.Fatalf("themes = %v, want 3 (Breeze, Shared, Custom)", got)
	}
	if _, ok := got["Breeze"]; !ok {
		t.Errorf("missing system theme Breeze; got %v", got)
	}
	if got["Custom"] != userCustom {
		t.Errorf("Custom = %q, want user path %q", got["Custom"], userCustom)
	}
	if got["Shared"] != userShared {
		t.Errorf("Shared = %q, want user override %q", got["Shared"], userShared)
	}
}
