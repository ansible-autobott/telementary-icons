package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// parseOverlay must read the set name and links, coercing unquoted numeric
// sizes (16) and the string "symbolic" into a uniform []string.
func TestParseOverlayCoercesNumericSizes(t *testing.T) {
	data := []byte("set: demo\n" +
		"links:\n" +
		"  - src: a.svg\n" +
		"    category: apps\n" +
		"    names: [foo, bar]\n" +
		"    sizes: [16, 24, symbolic]\n")

	cfg, err := parseOverlay(data)
	if err != nil {
		t.Fatalf("parseOverlay: %v", err)
	}
	if cfg.Set != "demo" {
		t.Errorf("set = %q, want demo", cfg.Set)
	}
	if len(cfg.Links) != 1 {
		t.Fatalf("links = %d, want 1", len(cfg.Links))
	}
	l := cfg.Links[0]
	if l.Src != "a.svg" || l.Category != "apps" {
		t.Errorf("link = %+v, want src=a.svg category=apps", l)
	}
	if !reflect.DeepEqual(l.Names, []string{"foo", "bar"}) {
		t.Errorf("names = %v, want [foo bar]", l.Names)
	}
	if !reflect.DeepEqual(l.Sizes, []string{"16", "24", "symbolic"}) {
		t.Errorf("sizes = %v, want [16 24 symbolic]", l.Sizes)
	}
}

// findOverlayConfigs must locate overlay.yaml files at any depth (so design
// sets can be grouped in category subfolders), sorted, ignoring dirs without one.
func TestFindOverlayConfigsRecursive(t *testing.T) {
	root := t.TempDir()
	mk := func(parts ...string) string {
		p := filepath.Join(append([]string{root}, parts...)...)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("set: x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}
	a := mk("apps", "smb4k", "overlay.yaml")
	b := mk("mimetypes", "audio", "overlay.yaml")
	if err := os.MkdirAll(filepath.Join(root, "apps", "noconfig"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := findOverlayConfigs(root)
	if err != nil {
		t.Fatalf("findOverlayConfigs: %v", err)
	}
	want := []string{a, b} // sorted: apps/ before mimetypes/
	if !reflect.DeepEqual(got, want) {
		t.Errorf("configs = %v\nwant %v", got, want)
	}
}

// A missing designs dir yields no configs and no error.
func TestFindOverlayConfigsMissing(t *testing.T) {
	got, err := findOverlayConfigs(filepath.Join(t.TempDir(), "does-not-exist"))
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if len(got) != 0 {
		t.Errorf("configs = %v, want empty", got)
	}
}

// parseOverlay must read optional provenance: a set-level origin (theme + root)
// and a per-link `from` (the source path under that root). These are metadata
// and must not affect the planned jobs.
func TestParseOverlayReadsOrigin(t *testing.T) {
	data := []byte("set: x\n" +
		"origin:\n" +
		"  theme: elementary\n" +
		"links:\n" +
		"  - src: a-48.svg\n" +
		"    from: apps/48/a.svg\n" +
		"    category: apps\n" +
		"    names: [a]\n" +
		"    sizes: [48]\n")

	cfg, err := parseOverlay(data)
	if err != nil {
		t.Fatalf("parseOverlay: %v", err)
	}
	if cfg.Origin == nil {
		t.Fatal("origin = nil, want set")
	}
	if cfg.Origin.Theme != "elementary" {
		t.Errorf("origin = %+v, want theme=elementary", *cfg.Origin)
	}
	if cfg.Links[0].From != "apps/48/a.svg" {
		t.Errorf("from = %q, want apps/48/a.svg", cfg.Links[0].From)
	}
	// Provenance must not leak into the planned jobs.
	jobs := linkJobs("/d", cfg, "/theme")
	want := []job{{src: "/d/a-48.svg", dst: "/theme/apps/48/a.svg"}}
	if !reflect.DeepEqual(jobs, want) {
		t.Errorf("jobs = %+v, want %+v", jobs, want)
	}
}

// linkJobs must copy the source into the first listed size and make the rest of
// the sizes relative symlinks pointing at that first copy — one such group per
// name.
func TestLinkJobsCopiesFirstSizeLinksRest(t *testing.T) {
	cfg := overlayConfig{Links: []overlayLink{{
		Src:      "smb4k_icon.svg",
		Category: "apps",
		Names:    []string{"smb4k", "smb4k-alt"},
		Sizes:    []string{"16", "24", "48"},
	}}}

	jobs := linkJobs("/designs/app_smb4k", cfg, "/theme")

	want := []job{
		// name smb4k: 16 copied, 24 & 48 linked to ../16/smb4k.svg
		{src: "/designs/app_smb4k/smb4k_icon.svg", dst: "/theme/apps/16/smb4k.svg"},
		{dst: "/theme/apps/24/smb4k.svg", link: "../16/smb4k.svg"},
		{dst: "/theme/apps/48/smb4k.svg", link: "../16/smb4k.svg"},
		// name smb4k-alt: its own first copy + links
		{src: "/designs/app_smb4k/smb4k_icon.svg", dst: "/theme/apps/16/smb4k-alt.svg"},
		{dst: "/theme/apps/24/smb4k-alt.svg", link: "../16/smb4k-alt.svg"},
		{dst: "/theme/apps/48/smb4k-alt.svg", link: "../16/smb4k-alt.svg"},
	}
	if !reflect.DeepEqual(jobs, want) {
		t.Errorf("jobs =\n%+v\nwant\n%+v", jobs, want)
	}
}

// A single-size link is just a copy — no symlinks.
func TestLinkJobsSingleSizeIsCopyOnly(t *testing.T) {
	cfg := overlayConfig{Links: []overlayLink{{
		Src: "a.svg", Category: "places", Names: []string{"folder"}, Sizes: []string{"22"},
	}}}
	jobs := linkJobs("/d", cfg, "/theme")
	want := []job{{src: "/d/a.svg", dst: "/theme/places/22/folder.svg"}}
	if !reflect.DeepEqual(jobs, want) {
		t.Errorf("jobs = %+v, want %+v", jobs, want)
	}
}

// applyJobs must write the first copy as a real file and the rest as relative
// symlinks that resolve to that copy's bytes.
func TestApplyJobsWritesCopyThenSymlinks(t *testing.T) {
	root := t.TempDir()
	srcSVG := filepath.Join(root, "src.svg")
	if err := os.WriteFile(srcSVG, []byte("<svg/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	theme := filepath.Join(root, "theme")

	jobs := []job{
		{src: srcSVG, dst: filepath.Join(theme, "apps", "16", "x.svg")},
		{dst: filepath.Join(theme, "apps", "24", "x.svg"), link: "../16/x.svg"},
	}
	if err := applyJobs(jobs); err != nil {
		t.Fatalf("applyJobs: %v", err)
	}

	// 16 is a real file with the source bytes.
	fi, err := os.Lstat(jobs[0].dst)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		t.Errorf("%s is a symlink, want a real file", jobs[0].dst)
	}

	// 24 is a symlink whose target is ../16/x.svg and which resolves to the bytes.
	fi, err = os.Lstat(jobs[1].dst)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Errorf("%s is not a symlink", jobs[1].dst)
	}
	target, err := os.Readlink(jobs[1].dst)
	if err != nil {
		t.Fatal(err)
	}
	if target != "../16/x.svg" {
		t.Errorf("link target = %q, want ../16/x.svg", target)
	}
	got, err := os.ReadFile(jobs[1].dst) // follows the symlink
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "<svg/>" {
		t.Errorf("resolved symlink content = %q, want <svg/>", got)
	}
}

func TestPlanOverlays(t *testing.T) {
	root := t.TempDir()
	set := filepath.Join(root, "apps", "demo")
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

	jobs, sets, err := planOverlays(root, "/theme")
	if err != nil {
		t.Fatalf("planOverlays: %v", err)
	}
	if sets != 1 {
		t.Errorf("sets = %d, want 1", sets)
	}
	// 16 is a copy, 24 is a symlink to it.
	want := []job{
		{src: filepath.Join(set, "a.svg"), dst: "/theme/apps/16/a.svg"},
		{dst: "/theme/apps/24/a.svg", link: "../16/a.svg"},
	}
	if !reflect.DeepEqual(jobs, want) {
		t.Errorf("jobs = %+v\nwant %+v", jobs, want)
	}
}
