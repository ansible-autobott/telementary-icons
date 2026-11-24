package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// The review page records its approve/reject decisions in a markdown file in
// the todo app's format (https://github.com/candy-tools/todo): headings are
// categories, "- [ ]" lines are tasks, indented lines are a task's description.
// Two categories are managed here — "Approved" and "Rejected" — holding one task
// per decided candidate:
//
//	## Approved
//
//	- [ ] `document-open` — elementary overlay
//	  names: document-open, document-open-symbolic · category: actions
//	  16 ← actions/16/document-open.svg (→ places/16/folder-open.svg)
//	  24, 22 ← actions/24/document-open.svg
//	  Comment: use the 32 art for 22
//
// Everything outside those two categories (a preamble, todo's guide block,
// other categories) is kept verbatim. Within a managed task only its checkbox
// state, its comment and the lines this tool does not generate (sub-tasks,
// notes an implementer added) survive a rewrite: the proposal lines are
// regenerated from the candidates file every time.

const (
	verdictApproved = "approved"
	verdictRejected = "rejected"
)

// managedCategories are the categories this file owns, in the order they are
// appended to a file that lacks them.
var managedCategories = []struct{ verdict, heading string }{
	{verdictApproved, "Approved"},
	{verdictRejected, "Rejected"},
}

const decisionsTitle = "Elementary overlay review"

// decisionsTemplate seeds a decisions file that does not exist yet.
const decisionsTemplate = `<!--
review/decisions.md — written by the telementary overlay review page
(make review, then http://localhost:3033/review), which rewrites the Approved
and Rejected categories below on every decision. It is also a todo-app list:
todo review/decisions.md

For the implementer: each task proposes publishing elementary art (paths are
relative to the elementary icon theme root) at the listed theme buckets. A line
"24, 22 ← actions/24/x.svg" is one overlay.yaml link: copy that file into an
overlays/<category>/<set>/ folder and publish it at sizes [24, 22] (the first
is the real copy). "(→ …)" marks an elementary symlink: copy the file it
resolves to. Symbolic art must be recolored to the ColorScheme-Text pattern
(overlays/README.md). "also changes" lists Tela aliases that follow the
replaced file. A Comment overrides the proposal. Approved: implement it.
Rejected: do not overlay it, but still act on its comment if it asks for
something else. Mark a task [x] once it is done.
-->

# ` + decisionsTitle + `

## Approved

## Rejected
`

var (
	mdHeadingRe  = regexp.MustCompile(`^(#{1,6})[ \t]+(.*?)[ \t]*$`)
	mdTaskRe     = regexp.MustCompile(`^[-*+][ \t]+\[([ xX/>])\][ \t]?(.*)$`)
	mdSubtaskRe  = regexp.MustCompile(`^[ \t]+[-*+][ \t]+\[[ xX/>]\]`)
	mdBacktickRe = regexp.MustCompile("`([^`]+)`")
)

// decision is one decided candidate as recorded in the decisions file.
type decision struct {
	Verdict string `json:"verdict"` // verdictApproved | verdictRejected
	Status  string `json:"status"`  // the task's checkbox marker: " ", "x", "/" or ">"
	Comment string `json:"comment"`
	// extra holds the task's description lines this tool does not generate
	// (sub-tasks, implementer notes), re-emitted verbatim after the comment.
	// raw is the whole task, re-emitted untouched when its id is not a
	// candidate (anymore) — the file never loses what someone wrote.
	extra []string
	raw   []string
}

// proposal is a candidate as the decisions file shows it: the task title and
// the generated description lines (the proposed overlay), without the comment.
type proposal struct {
	ID    string
	Title string
	Lines []string
}

// decisionDoc is a parsed decisions file.
type decisionDoc struct {
	lines    []string
	sections map[string][2]int // verdict -> [heading line, end line) in lines
	decided  map[string]*decision
	order    map[string][]string // verdict -> task keys in file order
}

// parseDecisions reads a decisions file. The first category titled
// "Approved" (and "Rejected") at any heading level is the managed one; its
// tasks are keyed by the first backtick-quoted token in their title (the
// candidate id). A task without one, or repeating an id already seen, is kept
// verbatim under a key no candidate can match.
func parseDecisions(src string) *decisionDoc {
	doc := &decisionDoc{
		sections: map[string][2]int{},
		decided:  map[string]*decision{},
		order:    map[string][]string{},
	}
	src = strings.TrimRight(strings.ReplaceAll(src, "\r\n", "\n"), "\n")
	if src == "" {
		return doc
	}
	doc.lines = strings.Split(src, "\n")

	var heads []int
	for i, l := range doc.lines {
		if mdHeadingRe.MatchString(l) {
			heads = append(heads, i)
		}
	}
	for k, h := range heads {
		end := len(doc.lines)
		if k+1 < len(heads) {
			end = heads[k+1]
		}
		v := verdictOfHeading(mdHeadingRe.FindStringSubmatch(doc.lines[h])[2])
		if v == "" {
			continue
		}
		if _, seen := doc.sections[v]; seen {
			continue
		}
		doc.sections[v] = [2]int{h, end}
		doc.parseTasks(v, doc.lines[h+1:end])
	}
	return doc
}

func verdictOfHeading(title string) string {
	for _, c := range managedCategories {
		if strings.EqualFold(strings.TrimSpace(title), c.heading) {
			return c.verdict
		}
	}
	return ""
}

// parseTasks splits a managed category's body into tasks: a top-level
// checkbox line plus the indented (or blank) lines under it. Stray
// non-indented prose ends a task and is dropped, as the todo app does.
func (doc *decisionDoc) parseTasks(verdict string, body []string) {
	var cur []string
	flush := func() {
		for len(cur) > 1 && strings.TrimSpace(cur[len(cur)-1]) == "" {
			cur = cur[:len(cur)-1]
		}
		if len(cur) > 0 {
			doc.addTask(verdict, cur)
		}
		cur = nil
	}
	for _, l := range body {
		switch {
		case mdTaskRe.MatchString(l):
			flush()
			cur = []string{l}
		case cur != nil && (strings.TrimSpace(l) == "" || l[0] == ' ' || l[0] == '\t'):
			cur = append(cur, l)
		default:
			flush()
		}
	}
	flush()
}

func (doc *decisionDoc) addTask(verdict string, lines []string) {
	m := mdTaskRe.FindStringSubmatch(lines[0])
	key := ""
	if b := mdBacktickRe.FindStringSubmatch(m[2]); b != nil {
		key = b[1]
	}
	if _, dup := doc.decided[key]; key == "" || dup {
		// Unkeyed or duplicate: keep it verbatim where it is.
		key = "\x00" + verdict + strconv.Itoa(len(doc.order[verdict]))
	}
	status := m[1]
	if status == "X" {
		status = "x"
	}
	comment, extra := splitTaskDescription(lines[1:])
	doc.decided[key] = &decision{Verdict: verdict, Status: status, Comment: comment, extra: extra, raw: lines}
	doc.order[verdict] = append(doc.order[verdict], key)
}

// splitTaskDescription separates a task's description into its comment — a
// "Comment: …" line, or "Comment:" followed by "> …" lines — and the lines to
// keep verbatim. Generated proposal lines are dropped (they are regenerated),
// and everything from the first sub-task on is kept as is.
func splitTaskDescription(desc []string) (string, []string) {
	var comment, extra []string
	inComment, inSubtasks := false, false
	for _, l := range desc {
		if inSubtasks || mdSubtaskRe.MatchString(l) {
			inSubtasks = true
			extra = append(extra, l)
			continue
		}
		d := strings.TrimSpace(l)
		switch {
		case strings.HasPrefix(d, "Comment:"):
			inComment = true
			if first := strings.TrimSpace(strings.TrimPrefix(d, "Comment:")); first != "" {
				comment = append(comment, first)
			}
		case inComment && strings.HasPrefix(d, ">"):
			comment = append(comment, strings.TrimPrefix(strings.TrimPrefix(d, ">"), " "))
		case d == "":
		case isGeneratedLine(d):
			inComment = false
		default:
			inComment = false
			extra = append(extra, l)
		}
	}
	return strings.TrimSpace(strings.Join(comment, "\n")), extra
}

// isGeneratedLine reports whether a description line is one of the proposal
// lines this tool writes (see reviewer.proposal).
func isGeneratedLine(d string) bool {
	for _, p := range []string{"names:", "note:", "also changes:"} {
		if strings.HasPrefix(d, p) {
			return true
		}
	}
	return strings.Contains(d, " ← ")
}

// set records a verdict and comment for id, keeping the task's checkbox state
// and extra lines. An empty verdict clears the decision.
func (doc *decisionDoc) set(id, verdict, comment string) {
	if verdict == "" {
		delete(doc.decided, id)
		return
	}
	d := doc.decided[id]
	if d == nil {
		d = &decision{Status: " "}
		doc.decided[id] = d
	}
	d.Verdict, d.Comment = verdict, comment
}

// decisionsFor returns the recorded decisions for the given candidate ids.
func (doc *decisionDoc) decisionsFor(props []proposal) map[string]decision {
	out := map[string]decision{}
	for _, p := range props {
		if d := doc.decided[p.ID]; d != nil {
			out[p.ID] = *d
		}
	}
	return out
}

// render returns the file with both managed categories regenerated: decided
// candidates in candidate order, then the tasks it does not know, verbatim.
// Everything outside the managed categories is left as it was; missing ones
// are appended (under the review title heading, if the file lacks it too).
// Blank lines are folded like the todo app's canonical form, so saving the
// file from todo does not churn it.
func (doc *decisionDoc) render(props []proposal) string {
	known := map[string]bool{}
	for _, p := range props {
		known[p.ID] = true
	}
	body := func(verdict string) []string {
		var out []string
		for _, p := range props {
			if d := doc.decided[p.ID]; d != nil && d.Verdict == verdict {
				out = append(out, d.task(p)...)
			}
		}
		for _, key := range doc.order[verdict] {
			if d := doc.decided[key]; d != nil && !known[key] {
				out = append(out, d.raw...)
			}
		}
		return out
	}
	var out []string
	emit := func(heading string, tasks []string) {
		out = append(out, "", heading, "")
		out = append(out, tasks...)
		out = append(out, "")
	}

	headingAt := map[int]string{}
	for v, sec := range doc.sections {
		headingAt[sec[0]] = v
	}
	for i := 0; i < len(doc.lines); {
		if v, ok := headingAt[i]; ok {
			emit(doc.lines[i], body(v))
			i = doc.sections[v][1]
			continue
		}
		out = append(out, doc.lines[i])
		i++
	}
	hasTitle := false
	for _, l := range doc.lines {
		if m := mdHeadingRe.FindStringSubmatch(l); m != nil && m[2] == decisionsTitle {
			hasTitle = true
		}
	}
	for _, c := range managedCategories {
		if _, ok := doc.sections[c.verdict]; ok {
			continue
		}
		if !hasTitle {
			out = append(out, "", "# "+decisionsTitle)
			hasTitle = true
		}
		emit("## "+c.heading, body(c.verdict))
	}
	return strings.Join(foldBlankLines(out), "\n") + "\n"
}

// task renders a decided candidate's task: checkbox and title, the generated
// proposal, the comment, then the lines kept from the file.
func (d *decision) task(p proposal) []string {
	status := d.Status
	if status == "" {
		status = " "
	}
	lines := []string{"- [" + status + "] " + p.Title}
	for _, l := range p.Lines {
		lines = append(lines, "  "+l)
	}
	lines = append(lines, commentLines(d.Comment)...)
	return append(lines, d.extra...)
}

// commentLines renders a comment as task description lines: one line inline,
// several as a "Comment:" line followed by quoted lines.
func commentLines(c string) []string {
	c = strings.TrimSpace(c)
	if c == "" {
		return nil
	}
	ls := strings.Split(c, "\n")
	if len(ls) == 1 {
		return []string{"  Comment: " + ls[0]}
	}
	out := []string{"  Comment:"}
	for _, l := range ls {
		if l = strings.TrimRight(l, " \t"); l == "" {
			out = append(out, "  >")
		} else {
			out = append(out, "  > "+l)
		}
	}
	return out
}

// foldBlankLines drops leading and trailing blank lines and folds runs of
// blank lines into one — the todo app's canonical spacing.
func foldBlankLines(lines []string) []string {
	var out []string
	prevBlank := true
	for _, l := range lines {
		blank := strings.TrimSpace(l) == ""
		if blank && prevBlank {
			continue
		}
		if blank {
			l = ""
		}
		out = append(out, l)
		prevBlank = blank
	}
	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return out
}

// loadDecisions reads the decisions file; a missing one starts from the
// template.
func loadDecisions(path string) (*decisionDoc, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return parseDecisions(decisionsTemplate), nil
	}
	if err != nil {
		return nil, err
	}
	return parseDecisions(string(data)), nil
}

// writeFileAtomic writes data to a temp file beside path and renames it over
// path, so a concurrent reader (the todo app polls the file) never sees a
// half-written file.
func writeFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".decisions-*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Chmod(tmp, 0o644); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}
