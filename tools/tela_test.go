package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTransposeLinkTarget(t *testing.T) {
	cases := []struct {
		size, cat, target, want string
	}{
		{"16", "places", "user-trash.svg", "user-trash.svg"},
		{"24", "places", "./bookmark-new.svg", "bookmark-new.svg"},
		{"scalable", "mimetypes", "../apps/application-default-icon.svg", "../../apps/scalable/application-default-icon.svg"},
		{"22", "actions", "../../16/actions/zoom-out.svg", "../16/zoom-out.svg"},
		{"scalable", "apps", "Zoom.svg", "Zoom.svg"},
	}
	for _, c := range cases {
		got, err := transposeLinkTarget(c.size, c.cat, c.target)
		if err != nil {
			t.Errorf("transposeLinkTarget(%q,%q,%q) error: %v", c.size, c.cat, c.target, err)
			continue
		}
		if got != c.want {
			t.Errorf("transposeLinkTarget(%q,%q,%q) = %q, want %q", c.size, c.cat, c.target, got, c.want)
		}
	}
}

func TestTransposeLinkTargetBadTargetErrors(t *testing.T) {
	if _, err := transposeLinkTarget("16", "places", "../../../etc/passwd"); err == nil {
		t.Errorf("expected error for target escaping the <size>/<cat>/<name> shape")
	}
}

// buildTela writes a tiny size-major Tela tree (src + links) under dir.
func buildTela(t *testing.T, dir string) {
	t.Helper()
	writeFile(t, filepath.Join(dir, "src", "16", "places", "user-trash.svg"), "<svg id=trash/>")
	writeFile(t, filepath.Join(dir, "src", "scalable", "apps", "application-default-icon.svg"), "<svg id=appdefault/>")
	writeFile(t, filepath.Join(dir, "src", "index.theme"), "[Icon Theme]\n") // must be skipped
	symlink(t, "user-trash.svg", filepath.Join(dir, "links", "16", "places", "xfce-trash_empty.svg"))
	symlink(t, "../apps/application-default-icon.svg", filepath.Join(dir, "links", "scalable", "mimetypes", "unknown.svg"))
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func symlink(t *testing.T, target, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
}

func TestMergeTela(t *testing.T) {
	root := t.TempDir()
	tela := filepath.Join(root, "tela")
	theme := filepath.Join(root, "theme")
	buildTela(t, tela)

	files, links, err := mergeTela(tela, theme)
	if err != nil {
		t.Fatalf("mergeTela: %v", err)
	}
	if files != 2 { // index.theme skipped
		t.Errorf("files = %d, want 2", files)
	}
	if links != 2 {
		t.Errorf("links = %d, want 2", links)
	}

	// Real file transposed size-major -> category-major.
	if b, err := os.ReadFile(filepath.Join(theme, "places", "16", "user-trash.svg")); err != nil || string(b) != "<svg id=trash/>" {
		t.Errorf("places/16/user-trash.svg = %q err %v", b, err)
	}
	// index.theme not copied into a bucket.
	if _, err := os.Stat(filepath.Join(theme, "index.theme")); !os.IsNotExist(err) {
		t.Errorf("src/index.theme should not be merged")
	}
	// Same-dir alias symlink preserved.
	got, err := os.Readlink(filepath.Join(theme, "places", "16", "xfce-trash_empty.svg"))
	if err != nil || got != "user-trash.svg" {
		t.Errorf("xfce-trash_empty.svg -> %q err %v, want user-trash.svg", got, err)
	}
	// Cross-category alias rewritten and resolvable.
	got, err = os.Readlink(filepath.Join(theme, "mimetypes", "scalable", "unknown.svg"))
	if err != nil || got != "../../apps/scalable/application-default-icon.svg" {
		t.Errorf("unknown.svg -> %q err %v, want ../../apps/scalable/application-default-icon.svg", got, err)
	}
	if _, err := os.ReadFile(filepath.Join(theme, "mimetypes", "scalable", "unknown.svg")); err != nil {
		t.Errorf("cross-category alias does not resolve: %v", err)
	}
}

func TestMergeTelaLinkWinsOverCollidingSrc(t *testing.T) {
	root := t.TempDir()
	tela := filepath.Join(root, "tela")
	theme := filepath.Join(root, "theme")
	// A name present BOTH as a real src svg and as a links/ alias to a different
	// icon: Tela copies src then links, so the alias must win.
	writeFile(t, filepath.Join(tela, "src", "16", "devices", "foo.svg"), "SRCART")
	writeFile(t, filepath.Join(tela, "src", "16", "devices", "bar.svg"), "BARART")
	symlink(t, "bar.svg", filepath.Join(tela, "links", "16", "devices", "foo.svg"))
	if _, _, err := mergeTela(tela, theme); err != nil {
		t.Fatalf("mergeTela: %v", err)
	}
	dst := filepath.Join(theme, "devices", "16", "foo.svg")
	got, err := os.Readlink(dst)
	if err != nil || got != "bar.svg" {
		t.Fatalf("foo.svg should be the winning alias symlink -> bar.svg; readlink=%q err=%v", got, err)
	}
	if b, err := os.ReadFile(dst); err != nil || string(b) != "BARART" {
		t.Errorf("resolved foo.svg = %q err %v, want BARART (link wins over real src)", b, err)
	}
}

func TestOverlayReplacesTelaAliasSymlink(t *testing.T) {
	root := t.TempDir()
	tela := filepath.Join(root, "tela")
	overlays := filepath.Join(root, "overlays")
	theme := filepath.Join(root, "theme")

	// Tela base: real canonical + an alias symlink pointing at it.
	writeFile(t, filepath.Join(tela, "src", "16", "places", "user-trash.svg"), "CANON")
	symlink(t, "user-trash.svg", filepath.Join(tela, "links", "16", "places", "trash.svg"))

	// Overlay that overwrites the alias name `trash` at places/16 with new art.
	writeFile(t, filepath.Join(overlays, "places", "trash", "trash.svg"), "OVERLAY")
	writeFile(t, filepath.Join(overlays, "places", "trash", "overlay.yaml"),
		"set: trash\nlinks:\n  - src: trash.svg\n    category: places\n    names: [trash]\n    sizes: [16]\n")

	if _, _, _, _, _, _, err := buildTheme(tela, overlays, theme); err != nil {
		t.Fatalf("buildTheme: %v", err)
	}
	// The overlay replaced the symlink with a real file...
	if b, err := os.ReadFile(filepath.Join(theme, "places", "16", "trash.svg")); err != nil || string(b) != "OVERLAY" {
		t.Errorf("places/16/trash.svg = %q err %v, want OVERLAY", b, err)
	}
	// ...and the alias's original target was NOT clobbered.
	if b, err := os.ReadFile(filepath.Join(theme, "places", "16", "user-trash.svg")); err != nil || string(b) != "CANON" {
		t.Errorf("places/16/user-trash.svg = %q err %v, want CANON (alias target must be intact)", b, err)
	}
}
