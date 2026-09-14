package zaplint_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"

	"github.com/rleungx/zaplint"
	"golang.org/x/tools/go/analysis/analysistest"
)

// Apply the actual emitted edits, compile both programs and compare the encoded
// log fields. Diagnostic text alone cannot prove that a replacement is safe.
func TestReplacementSemantics(t *testing.T) {
	t.Parallel()
	a := zaplint.New(&zaplint.Options{ReplaceAny: true})
	results := analysistest.Run(t, analysistest.TestData(), a, "replacement_semantics")
	if len(results) != 1 {
		t.Fatalf("got %d analysis results, want 1", len(results))
	}
	result := results[0]
	filename := filepath.Join(analysistest.TestData(), "src", "replacement_semantics", "main.go")
	source, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	type edit struct {
		start, end int
		text       []byte
	}
	var edits []edit
	for _, diagnostic := range result.Diagnostics {
		if len(diagnostic.SuggestedFixes) != 1 {
			t.Fatalf("expected one applicable fix for %q", diagnostic.Message)
		}
		for _, change := range diagnostic.SuggestedFixes[0].TextEdits {
			start := result.Pass.Fset.PositionFor(change.Pos, false)
			end := result.Pass.Fset.PositionFor(change.End, false)
			if start.Filename != filename || end.Filename != filename {
				t.Fatalf("edit targets unexpected file")
			}
			edits = append(edits, edit{start.Offset, end.Offset, change.NewText})
		}
	}
	if len(edits) == 0 {
		t.Fatal("no replacement opportunities exercised")
	}
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	modified := bytes.Clone(source)
	boundary := len(source)
	for _, edit := range edits {
		if edit.start < 0 || edit.end < edit.start || edit.end > boundary {
			t.Fatal("invalid or overlapping replacement edits")
		}
		modified = append(append(append([]byte{}, modified[:edit.start]...), edit.text...), modified[edit.end:]...)
		boundary = edit.start
	}
	fixed := filepath.Join(t.TempDir(), "main.go")
	if err := os.WriteFile(fixed, modified, 0600); err != nil {
		t.Fatal(err)
	}
	run := func(file string) []byte {
		cmd := exec.Command("go", "run", "-mod=vendor", file)
		cmd.Dir = filepath.Join(analysistest.TestData(), "src")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("replacement program did not compile/run: %v\n%s", err, output)
		}
		return output
	}
	before, after := run(filename), run(fixed)
	if !bytes.Equal(before, after) {
		t.Fatalf("replacement changed encoded logs:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}
