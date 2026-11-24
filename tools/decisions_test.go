package main

import (
	"path/filepath"
	"strings"
	"testing"
)

var testProps = []proposal{
	{ID: "document-open", Title: "`document-open` — elementary overlay", Lines: []string{
		"names: document-open · category: actions",
		"16 ← actions/16/document-open.svg",
		"24, 22 ← actions/24/document-open.svg",
	}},
	{ID: "edit-cut", Title: "`edit-cut` — elementary overlay", Lines: []string{
		"names: edit-cut · category: actions",
		"16 ← actions/16/edit-cut.svg",
	}},
}

// A missing decisions file starts from the template: preamble, title and both
// managed categories, with the decided task under its verdict.
func TestDecisionsStartFromTemplate(t *testing.T) {
	doc, err := loadDecisions(filepath.Join(t.TempDir(), "decisions.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc.set("edit-cut", verdictRejected, "keep Tela's scissors")
	got := doc.render(testProps)

	for _, want := range []string{
		"<!--\nreview/decisions.md",
		"# Elementary overlay review\n\n## Approved\n\n## Rejected\n\n" +
			"- [ ] `edit-cut` — elementary overlay\n" +
			"  names: edit-cut · category: actions\n" +
			"  16 ← actions/16/edit-cut.svg\n" +
			"  Comment: keep Tela's scissors\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("render missing %q\n--- got ---\n%s", want, got)
		}
	}
	if !strings.HasSuffix(got, "Comment: keep Tela's scissors\n") {
		t.Errorf("render should end with the last task and one newline:\n%s", got)
	}
}

// A rewrite keeps everything the review page does not own: the preamble, other
// categories, a task's checkbox state and sub-tasks, and tasks it does not
// know. The generated proposal lines are replaced by the current proposal.
func TestDecisionsRewritePreservesForeignContent(t *testing.T) {
	src := `<!-- todo:guide — managed by todo -->

# Elementary overlay review

## Approved

- [x] ` + "`document-open`" + ` — elementary overlay
  names: stale · category: actions
  16 ← actions/16/stale.svg
  Comment: use the 32 art for 22
  done in overlays/actions/document/
  - [ ] verify on the dark panel
    check both toolbars
- [ ] a task of my own
  with a description

## Rejected

# Notes

- [ ] ask about symbolic recolor
`
	doc := parseDecisions(src)
	doc.set("edit-cut", verdictRejected, "")
	got := doc.render(testProps)

	want := `<!-- todo:guide — managed by todo -->

# Elementary overlay review

## Approved

- [x] ` + "`document-open`" + ` — elementary overlay
  names: document-open · category: actions
  16 ← actions/16/document-open.svg
  24, 22 ← actions/24/document-open.svg
  Comment: use the 32 art for 22
  done in overlays/actions/document/
  - [ ] verify on the dark panel
    check both toolbars
- [ ] a task of my own
  with a description

## Rejected

- [ ] ` + "`edit-cut`" + ` — elementary overlay
  names: edit-cut · category: actions
  16 ← actions/16/edit-cut.svg

# Notes

- [ ] ask about symbolic recolor
`
	if got != want {
		t.Errorf("render mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	if again := parseDecisions(got).render(testProps); again != got {
		t.Errorf("render is not stable\n--- first ---\n%s\n--- second ---\n%s", got, again)
	}
}

// Changing a verdict moves the task between categories with its checkbox
// state; clearing it removes the task.
func TestDecisionsMoveAndClear(t *testing.T) {
	doc := parseDecisions(decisionsTemplate)
	doc.set("document-open", verdictApproved, "first")
	doc = parseDecisions(doc.render(testProps))
	if d := doc.decided["document-open"]; d == nil || d.Verdict != verdictApproved || d.Comment != "first" {
		t.Fatalf("after approve: %+v", d)
	}

	doc.set("document-open", verdictRejected, "changed my mind")
	out := doc.render(testProps)
	approved := out[strings.Index(out, "## Approved"):strings.Index(out, "## Rejected")]
	if strings.Contains(approved, "document-open") {
		t.Errorf("task still under Approved:\n%s", out)
	}
	if !strings.Contains(out[strings.Index(out, "## Rejected"):], "`document-open`") {
		t.Errorf("task not under Rejected:\n%s", out)
	}

	doc = parseDecisions(out)
	doc.set("document-open", "", "")
	if out := doc.render(testProps); strings.Contains(out, "`document-open`") {
		t.Errorf("cleared task still present:\n%s", out)
	}
}

// Multi-line comments survive a round trip, including blank lines and lines
// that would otherwise read as a sub-task or a generated proposal line.
func TestDecisionsMultiLineCommentRoundTrip(t *testing.T) {
	comment := "use the 32 art for 22\n\n- [ ] not a sub-task\n22 ← actions/32/x.svg\nComment: literal"
	doc := parseDecisions(decisionsTemplate)
	doc.set("edit-cut", verdictApproved, comment)
	out := doc.render(testProps)
	if !strings.Contains(out, "  Comment:\n  > use the 32 art for 22\n  >\n  > - [ ] not a sub-task\n") {
		t.Errorf("comment not rendered as a quote block:\n%s", out)
	}
	back := parseDecisions(out).decided["edit-cut"]
	if back == nil || back.Comment != comment {
		t.Errorf("comment = %q, want %q", back.Comment, comment)
	}
	if len(back.extra) != 0 {
		t.Errorf("comment leaked into extra lines: %q", back.extra)
	}
}

// Only the Approved/Rejected categories are decisions: a backticked task
// elsewhere is not picked up.
func TestDecisionsIgnoreOtherCategories(t *testing.T) {
	doc := parseDecisions("# Notes\n\n- [ ] `edit-cut` — a note\n\n## Approved\n\n- [ ] `document-open` — x\n")
	if doc.decided["edit-cut"] != nil {
		t.Errorf("a task outside the managed categories was taken as a decision")
	}
	if d := doc.decided["document-open"]; d == nil || d.Verdict != verdictApproved {
		t.Errorf("document-open = %+v, want approved", d)
	}
}
