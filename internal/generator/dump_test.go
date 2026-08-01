package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// regenerateCmd is the one command that refreshes every golden under
// irDumpDir. Failure messages quote it verbatim so a developer staring at a
// diff does not have to go looking for it.
const regenerateCmd = "nix develop -c go run ./internal/generator -dump-ir"

// TestRenderTypeExpandsLossyForms pins the D6 rendering table across every Type
// implementation in type.go, plus nil and the nested forms that matter.
//
// The three lossy String() methods are the point of the test: FuncType.String()
// returns the literal "func", UnionType.String() returns the literal "union",
// and IdentType.String() drops IsAnonEnum. If renderType ever regressed to
// calling String() throughout, a parser change in a callback argument type, a
// union member, or the anonymous-enum classification would stop showing in the
// golden. The nested cases (a union inside an array, a func behind a pointer)
// guard the recursion, because a renderer that expanded only at the top level
// would still pass the flat cases.
func TestRenderTypeExpandsLossyForms(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   Type
		want string
	}{
		{
			name: "nil is C void",
			in:   nil,
			want: "void",
		},
		{
			name: "ident reuses String",
			in:   &IdentType{Name: "AVCodecContext"},
			want: "AVCodecContext",
		},
		{
			name: "ident anon enum keeps the classification String drops",
			in:   &IdentType{Name: "uint32_t", IsAnonEnum: true},
			want: "uint32_t/anonenum",
		},
		{
			name: "pointer to ident",
			in:   &PointerType{Inner: &IdentType{Name: "AVFrame"}},
			want: "*AVFrame",
		},
		{
			name: "pointer to void",
			in:   &PointerType{Inner: nil},
			want: "*void",
		},
		{
			name: "pointer to pointer",
			in:   &PointerType{Inner: &PointerType{Inner: &IdentType{Name: "AVFrame"}}},
			want: "**AVFrame",
		},
		{
			name: "const array carries its size",
			in:   &ConstArray{Inner: &IdentType{Name: "int32_t"}, Size: 9},
			want: "int32_t[9]",
		},
		{
			name: "incomplete array",
			in:   &Array{Inner: &IdentType{Name: "uint8_t"}},
			want: "uint8_t[?]",
		},
		{
			name: "func keeps args and result String drops",
			in: &FuncType{
				Args:   []Type{&PointerType{Inner: &IdentType{Name: "AVCodecContext"}}, &IdentType{Name: "int"}},
				Result: &IdentType{Name: "int"},
			},
			want: "func(*AVCodecContext,int)int",
		},
		{
			name: "func with no args and void result",
			in:   &FuncType{},
			want: "func()void",
		},
		{
			name: "func with a void arg",
			in:   &FuncType{Args: []Type{nil}, Result: nil},
			want: "func(void)void",
		},
		{
			name: "pointer to func expands through the pointer",
			in: &PointerType{Inner: &FuncType{
				Args:   []Type{&PointerType{Inner: &IdentType{Name: "AVIOContext"}}},
				Result: &IdentType{Name: "int64_t"},
			}},
			want: "*func(*AVIOContext)int64_t",
		},
		{
			name: "union keeps members String drops",
			in: &UnionType{Fields: []*Field{
				{Name: "i64", Type: &IdentType{Name: "int64_t"}},
				{Name: "dbl", Type: &IdentType{Name: "double"}},
			}},
			want: "union{i64 int64_t;dbl double;}",
		},
		{
			name: "empty union",
			in:   &UnionType{},
			want: "union{}",
		},
		{
			name: "const array of union expands both",
			in: &ConstArray{
				Size: 2,
				Inner: &UnionType{Fields: []*Field{
					{Name: "num", Type: &IdentType{Name: "int"}},
					{Name: "ptr", Type: &PointerType{Inner: &IdentType{Name: "AVFrame"}}},
				}},
			},
			want: "union{num int;ptr *AVFrame;}[2]",
		},
		{
			name: "union member carrying a func expands through both",
			in: &UnionType{Fields: []*Field{
				{Name: "cb", Type: &PointerType{Inner: &FuncType{Args: []Type{&IdentType{Name: "int"}}}}},
			}},
			want: "union{cb *func(int)void;}",
		},
		{
			name: "array of anon enum keeps the classification",
			in:   &Array{Inner: &IdentType{Name: "uint32_t", IsAnonEnum: true}},
			want: "uint32_t/anonenum[?]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := renderType(tt.in); got != tt.want {
				t.Errorf("renderType() = %q, want %q", got, tt.want)
			}
		})
	}
}

// fabricatedModule builds a small Module by hand, covering every kind the
// structural stream walks and every field the record lines carry. It is used
// instead of a real parse so the dump tests are fast and their assertions are
// exact: the test knows every symbol name it expects to see.
//
// A fresh value is returned on each call so a test that mutates one copy cannot
// affect another.
func fabricatedModule() *Module {
	return &Module{
		functionOrder: []string{"av_frame_alloc", "av_log_set_callback"},
		functions: map[string]*Function{
			"av_frame_alloc": {
				Name:    "av_frame_alloc",
				Result:  &PointerType{Inner: &IdentType{Name: "AVFrame"}},
				Comment: "Allocate an AVFrame and set its fields to default values.",
			},
			"av_log_set_callback": {
				Name:     "av_log_set_callback",
				Args:     []*Param{{Name: "callback", Type: &PointerType{Inner: &FuncType{}}}},
				Variadic: true,
				Comment:  "Set the logging callback.",
			},
		},

		structOrder: []string{"AVRational", "AVFrame"},
		structs: map[string]*Struct{
			"AVRational": {
				Name:     "AVRational",
				Typedefd: true,
				Fields: []*Field{
					{Name: "num", Type: &IdentType{Name: "int"}, CTypeName: "int", Comment: "Numerator"},
					{Name: "den", Type: &IdentType{Name: "int"}, CTypeName: "int", Comment: "Denominator"},
				},
				Comment: "Rational number (pair of numerator and denominator).",
			},
			"AVFrame": {
				Name:     "AVFrame",
				Typedefd: true,
				Fields: []*Field{
					{Name: "data", Type: &ConstArray{Inner: &PointerType{Inner: &IdentType{Name: "uint8_t"}}, Size: 8}},
					{Name: "pts", Type: &IdentType{Name: "int64_t"}, CTypeName: "int64_t", Comment: "Presentation timestamp."},
					{Name: "flags", Type: &IdentType{Name: "int"}, BitWidth: 3},
				},
				Comment: "This structure describes decoded (raw) audio or video data.",
			},
		},

		enumOrder: []string{"AVMediaType"},
		enums: map[string]*Enum{
			"AVMediaType": {
				Name:     "AVMediaType",
				Typedefd: true,
				Constants: []*Constant{
					{Name: "AVMEDIA_TYPE_VIDEO", FromEnum: true, Comment: "Video stream."},
					{Name: "AVMEDIA_TYPE_AUDIO", FromEnum: true},
				},
				Comment: "Media type of a stream.",
			},
		},

		callbackOrder: []string{"av_log_callback"},
		callbacks: map[string]*Function{
			"av_log_callback": {
				Name: "av_log_callback",
				Args: []*Param{
					{Name: "avcl", Type: &PointerType{Inner: nil}},
					{Name: "level", Type: &IdentType{Name: "int"}},
				},
				Ptr:     true,
				Comment: "Logging callback typedef.",
			},
		},

		constantOrder: []string{"AV_TIME_BASE"},
		constants: map[string]*Constant{
			"AV_TIME_BASE": {Name: "AV_TIME_BASE", Comment: "Internal time base."},
		},
	}
}

// renderBoth renders the structural and comment streams of m into strings,
// failing the test on any writer error. Both dump tests below need the pair, and
// keeping the buffer plumbing in one place stops the two from drifting.
func renderBoth(t *testing.T, m *Module) (structure, comments string) {
	t.Helper()

	var sb, cb bytes.Buffer

	if err := renderModuleStructure(&sb, m); err != nil {
		t.Fatalf("renderModuleStructure: %v", err)
	}

	if err := renderModuleComments(&cb, m); err != nil {
		t.Fatalf("renderModuleComments: %v", err)
	}

	return sb.String(), cb.String()
}

// lineDiff is one differing line between two renderings, carrying the 1-indexed
// line number so a failure message points at the offending line rather than at
// the whole file.
type lineDiff struct {
	num  int
	want string
	got  string
}

// diffLines compares two dumps line by line and returns every position where
// they differ. A short file is padded with a missing marker rather than
// truncating the comparison, so a dump that loses its tail is still reported.
//
// This is a positional comparison, not a longest-common-subsequence diff: an
// inserted line shifts everything after it. That is adequate here because the
// tests mutate a field in place rather than adding or removing symbols, and the
// golden comparison only needs the first divergence.
func diffLines(want, got string) []lineDiff {
	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(got, "\n")

	n := max(len(wantLines), len(gotLines))

	var out []lineDiff

	for i := range n {
		w, g := "<missing>", "<missing>"

		if i < len(wantLines) {
			w = wantLines[i]
		}

		if i < len(gotLines) {
			g = gotLines[i]
		}

		if w != g {
			out = append(out, lineDiff{num: i + 1, want: w, got: g})
		}
	}

	return out
}

// truncate caps a line for a failure message. A comment body line or a wide
// struct record would otherwise bury the useful part of the message.
func truncate(s string) string {
	const limit = 120

	if len(s) <= limit {
		return s
	}

	return s[:limit] + "..."
}

// TestIRDumpDeterministic proves the writers produce byte-identical output for
// the same Module across repeated renders. Determinism is the whole premise of
// the artefact: if any walk in the writers ranged over one of Module's maps
// instead of its order slices, the goldens would churn between runs and every
// diff would become noise. A Module with several entries in each map gives Go's
// randomised map iteration a real chance to expose such a walk.
func TestIRDumpDeterministic(t *testing.T) {
	t.Parallel()

	m := fabricatedModule()

	firstStructure, firstComments := renderBoth(t, m)

	// Several rounds, because a single repeat can coincidentally reproduce a
	// map order.
	for round := range 8 {
		structure, comments := renderBoth(t, m)

		if structure != firstStructure {
			d := diffLines(firstStructure, structure)
			t.Fatalf("round %d: structure dump not deterministic, first diff at line %d:\n want %s\n  got %s",
				round, d[0].num, truncate(d[0].want), truncate(d[0].got))
		}

		if comments != firstComments {
			d := diffLines(firstComments, comments)
			t.Fatalf("round %d: comment dump not deterministic, first diff at line %d:\n want %s\n  got %s",
				round, d[0].num, truncate(d[0].want), truncate(d[0].got))
		}
	}
}

// TestIRDumpNamesEditedSymbol is the D7 test: a hand-edited field on one symbol
// must produce a diff that names that symbol, so a reviewer reading a golden
// diff can tell which binding moved without cross-referencing anything.
//
// The structural stream carries the strict form of the claim, because D4 makes
// every structural line name its symbol: every differing line must contain the
// symbol name. The comment stream can only offer the weaker form, because a
// record's body is verbatim comment text that need not mention the symbol; there
// the assertion is that a differing line names it, which is the record header.
//
// The mutations are deliberately of different kinds (a type, a comment, a
// boolean, a name) so the test covers more than one path into the record lines.
func TestIRDumpNamesEditedSymbol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// symbol is the name every structural diff line must carry.
		symbol string
		mutate func(m *Module)
		// commentsMustDiffer marks the mutations that also move the comment
		// stream, so the test can assert the comment diff rather than tolerate
		// an empty one.
		commentsMustDiffer bool
	}{
		{
			name:   "struct field type",
			symbol: "AVFrame",
			mutate: func(m *Module) {
				m.structs["AVFrame"].Fields[1].Type = &IdentType{Name: "int32_t"}
			},
		},
		{
			name:   "struct field bit width",
			symbol: "AVFrame",
			mutate: func(m *Module) {
				m.structs["AVFrame"].Fields[2].BitWidth = 4
			},
		},
		{
			name:   "struct byvalue flag",
			symbol: "AVRational",
			mutate: func(m *Module) {
				m.structs["AVRational"].ByValue = true
			},
		},
		{
			name:   "function comment",
			symbol: "av_frame_alloc",
			mutate: func(m *Module) {
				m.functions["av_frame_alloc"].Comment = "Allocate an AVFrame; edited by hand."
			},
			commentsMustDiffer: true,
		},
		{
			name:   "function result type",
			symbol: "av_frame_alloc",
			mutate: func(m *Module) {
				m.functions["av_frame_alloc"].Result = nil
			},
		},
		{
			name:   "callback param name",
			symbol: "av_log_callback",
			mutate: func(m *Module) {
				m.callbacks["av_log_callback"].Args[1].Name = "loglevel"
			},
		},
		{
			name:   "enum constant comment",
			symbol: "AVMediaType",
			mutate: func(m *Module) {
				m.enums["AVMediaType"].Constants[0].Comment = "Video stream; edited by hand."
			},
			commentsMustDiffer: true,
		},
		{
			name:   "top-level constant fromenum flag",
			symbol: "AV_TIME_BASE",
			mutate: func(m *Module) {
				m.constants["AV_TIME_BASE"].FromEnum = true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := fabricatedModule()
			beforeStructure, beforeComments := renderBoth(t, m)

			tt.mutate(m)

			afterStructure, afterComments := renderBoth(t, m)

			structureDiff := diffLines(beforeStructure, afterStructure)
			if len(structureDiff) == 0 {
				t.Fatalf("mutating %s produced no structural diff; the edit is invisible in the golden", tt.symbol)
			}

			for _, d := range structureDiff {
				if !strings.Contains(d.want, tt.symbol) || !strings.Contains(d.got, tt.symbol) {
					t.Errorf("structural diff line %d does not name %q:\n want %s\n  got %s",
						d.num, tt.symbol, truncate(d.want), truncate(d.got))
				}
			}

			commentDiff := diffLines(beforeComments, afterComments)

			if !tt.commentsMustDiffer {
				return
			}

			if len(commentDiff) == 0 {
				t.Fatalf("mutating the comment on %s produced no comment-stream diff", tt.symbol)
			}

			named := false

			for _, d := range commentDiff {
				if strings.Contains(d.want, tt.symbol) || strings.Contains(d.got, tt.symbol) {
					named = true
					break
				}
			}

			if !named {
				t.Errorf("no comment-stream diff line names %q; first diff at line %d:\n want %s\n  got %s",
					tt.symbol, commentDiff[0].num, truncate(commentDiff[0].want), truncate(commentDiff[0].got))
			}
		})
	}
}

// renderSkipsString renders a collector to a string, failing on writer error.
func renderSkipsString(t *testing.T, c *SkipCollector) string {
	t.Helper()

	var b bytes.Buffer

	if err := renderSkips(&b, c); err != nil {
		t.Fatalf("renderSkips: %v", err)
	}

	return b.String()
}

// TestSkipDumpCatchesReasonChange is the load-bearing test for the acceptance
// criterion that the skip golden catches a changed reason string, not only a
// changed count.
//
// The two collectors are built to defeat a count check on purpose: same symbols,
// same number of markers, same number of unique pairs. The test asserts that
// equality first, so it demonstrably shows a count-only comparison would pass
// while the reason text has silently changed. Only then does it assert the
// renderings differ and that the differing line carries both the symbol and the
// new reason text.
func TestSkipDumpCatchesReasonChange(t *testing.T) {
	t.Parallel()

	const (
		symbol    = "av_display_rotation_set"
		oldReason = "const array param matrix"
		newReason = "const array param rotation_matrix"
	)

	a := &SkipCollector{}
	a.Record("AVBlowfish.s", "multi dim const array")
	a.Record(symbol, oldReason)
	a.Record("AVChannelLayout.u", "union type")

	b := &SkipCollector{}
	b.Record("AVBlowfish.s", "multi dim const array")
	b.Record(symbol, newReason)
	b.Record("AVChannelLayout.u", "union type")

	// The premise: every count a cheaper check could compare is identical, so a
	// count check would not have caught this change.
	if a.Total() != b.Total() {
		t.Fatalf("marker totals differ (%d vs %d); the fixture no longer proves a count check is insufficient",
			a.Total(), b.Total())
	}

	if got, want := len(sortedUniqueSymbols(a)), len(sortedUniqueSymbols(b)); got != want {
		t.Fatalf("unique symbol counts differ (%d vs %d); the fixture no longer proves a count check is insufficient", got, want)
	}

	if got, want := len(sortedSkipPairs(a)), len(sortedSkipPairs(b)); got != want {
		t.Fatalf("unique pair counts differ (%d vs %d); the fixture no longer proves a count check is insufficient", got, want)
	}

	before := renderSkipsString(t, a)
	after := renderSkipsString(t, b)

	if before == after {
		t.Fatal("skip dumps are identical after a reason change; the golden cannot catch a reason-only regression")
	}

	diff := diffLines(before, after)

	if len(diff) != 1 {
		t.Errorf("want exactly 1 differing line for a single changed reason, got %d", len(diff))
	}

	d := diff[0]

	if !strings.Contains(d.got, symbol) {
		t.Errorf("differing line %d does not name the symbol %q: %s", d.num, symbol, truncate(d.got))
	}

	if !strings.Contains(d.got, newReason) {
		t.Errorf("differing line %d does not carry the new reason %q: %s", d.num, newReason, truncate(d.got))
	}

	if !strings.Contains(d.want, oldReason) {
		t.Errorf("differing line %d does not carry the old reason %q: %s", d.num, oldReason, truncate(d.want))
	}
}

// TestSkipDumpSortedAndDeduplicated pins the skip stream's body contract:
// sorted by symbol then reason, one line per unique (Symbol, Reason) pair, and
// header counts that keep the three different measurements distinct.
//
// The input is deliberately out of insertion order and contains both an exact
// duplicate pair (which must collapse to one line) and one symbol with two
// different reasons (which must not collapse). Those two cases pull in opposite
// directions, so a dedup keyed on the symbol alone and a dedup that does not
// dedup at all both fail here.
//
// Manual is set on one entry to confirm it is excluded: it is derived by
// scanning the repo root, so including it would make the golden depend on the
// working directory.
func TestSkipDumpSortedAndDeduplicated(t *testing.T) {
	t.Parallel()

	c := &SkipCollector{}
	c.Entries = []SkipEntry{
		{Symbol: "av_display_rotation_set", Reason: "const array param matrix"},
		{Symbol: "AVFrame.data", Reason: "multi dim const array"},
		{Symbol: "AVBlowfish.s", Reason: "multi dim const array", Manual: "AVBlowfishInit"},
		// Exact duplicate of the first entry: must collapse to one body line
		// while still counting as a marker.
		{Symbol: "av_display_rotation_set", Reason: "const array param matrix"},
		// Same symbol, different reason: must not collapse.
		{Symbol: "av_display_rotation_set", Reason: "unhandled arg type *main.ConstArray (matrix)"},
	}

	got := renderSkipsString(t, c)

	want := strings.Join([]string{
		skipsFileHeader,
		"# markers=5 symbols=3 pairs=4",
		"AVBlowfish.s: multi dim const array",
		"AVFrame.data: multi dim const array",
		"av_display_rotation_set: const array param matrix",
		"av_display_rotation_set: unhandled arg type *main.ConstArray (matrix)",
		"",
	}, "\n")

	if got != want {
		for _, d := range diffLines(want, got) {
			t.Errorf("line %d:\n want %s\n  got %s", d.num, truncate(d.want), truncate(d.got))
		}
	}

	if strings.Contains(got, "AVBlowfishInit") {
		t.Error("skip dump leaked SkipEntry.Manual, which is derived from the working directory")
	}

	// The nil collector must still yield a valid file: both headers with zero
	// counts and no body.
	nilWant := skipsFileHeader + "\n# markers=0 symbols=0 pairs=0\n"

	if nilGot := renderSkipsString(t, nil); nilGot != nilWant {
		t.Errorf("nil collector dump = %q, want %q", nilGot, nilWant)
	}
}

// TestCommentGoldenLineCount pins the IR comment-line count recorded in
// irCommentLines against the committed golden, so a change to processComment
// fails a test instead of quietly moving the golden.
//
// The count honours each record header's declared lines=N rather than
// subtracting header lines from the total. A comment body line that itself
// begins "@ " is real FFmpeg doc-comment text and would be miscounted as a
// header by the subtraction method; walking the declared counts cannot make
// that mistake. Walking the records this way also validates the file's
// structure, because a header where a body line was expected, or a body that
// runs off the end, fails the walk.
// It reads the golden through a working-directory-relative path, so it is not
// marked parallel: TestIRGoldensMatchFreshRun changes the working directory for
// the whole process.
func TestCommentGoldenLineCount(t *testing.T) {
	// The golden is committed at a repo-root-relative path; the test's working
	// directory is the package directory, two levels down.
	path := filepath.Join("testdata", "ir", "comments.txt")

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read comment golden %s: %v\nregenerate with: %s", path, err, regenerateCmd)
	}

	lines := strings.Split(string(raw), "\n")

	// A trailing newline splits into a final empty element that is not a line.
	if n := len(lines); n > 0 && lines[n-1] == "" {
		lines = lines[:n-1]
	}

	if len(lines) == 0 || lines[0] != commentsFileHeader {
		t.Fatalf("comment golden does not open with the expected header line\nregenerate with: %s", regenerateCmd)
	}

	dataLines, records := 0, 0

	for i := 1; i < len(lines); {
		n, err := recordBodyLines(lines[i])
		if err != nil {
			t.Fatalf("comment golden line %d: %v\nregenerate with: %s", i+1, err, regenerateCmd)
		}

		if i+1+n > len(lines) {
			t.Fatalf("comment golden line %d declares lines=%d but the file ends after %d more lines",
				i+1, n, len(lines)-i-1)
		}

		records++
		dataLines += n
		i += 1 + n
	}

	if dataLines != irCommentLines {
		t.Errorf("comment golden holds %d data lines, irCommentLines = %d\n"+
			"if processComment changed on purpose, regenerate with: %s\nand update irCommentLines to match",
			dataLines, irCommentLines, regenerateCmd)
	}

	// Recorded alongside the line count at generation time; a change here
	// without a line-count change means symbols moved in or out of the dump.
	const wantRecords = 2907

	if records != wantRecords {
		t.Errorf("comment golden holds %d records, want %d\nregenerate with: %s", records, wantRecords, regenerateCmd)
	}
}

// recordBodyLines reads the declared body-line count from a comment-record
// header. The key can contain a space ("struct AVFrame.pts"), so the fixed
// trailing fields are read from the end rather than by position from the start.
func recordBodyLines(line string) (int, error) {
	if !strings.HasPrefix(line, "@ ") {
		return 0, fmt.Errorf("expected a record header starting %q, got %q", "@ ", truncate(line))
	}

	fields := strings.Fields(line)
	if len(fields) < 4 {
		return 0, fmt.Errorf("malformed record header %q", truncate(line))
	}

	// Trailing fields are "lines=<n> hash=<hash>".
	raw, ok := strings.CutPrefix(fields[len(fields)-2], "lines=")
	if !ok {
		return 0, fmt.Errorf("record header %q has no lines= field where one is expected", truncate(line))
	}

	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("record header %q has a non-numeric line count: %w", truncate(line), err)
	}

	if n < 1 {
		return 0, fmt.Errorf("record header %q declares lines=%d; a record with no body should not exist", truncate(line), n)
	}

	return n, nil
}

// TestGoldensCarryNoMachinePaths fails if any committed golden holds an
// absolute filesystem path.
//
// A golden that names the directory it was generated in matches only on that
// machine, so TestIRGoldensMatchFreshRun would fail for everyone else and in
// CI. The C type spellings are the source of such paths: an unnamed union or
// struct is spelled with the absolute path of the header that declares it,
// which reaches the structural stream through Field.CTypeName and is why
// normalizeCTypeName exists.
//
// The assertion is on the generic property, not on one machine's string: the
// home and temporary directory prefixes a checkout normally lives under, plus
// this checkout's own root and its symlink-resolved form. Hard-coding a literal
// path would pass everywhere it was not generated.
//
// It reads the goldens through working-directory-relative paths, so it is not
// marked parallel: TestIRGoldensMatchFreshRun changes the working directory for
// the whole process.
func TestGoldensCarryNoMachinePaths(t *testing.T) {
	repoRoot := testRepoRoot(t)

	needles := []string{"/home/", "/Users/", "/tmp/", repoRoot}

	// A header may be reported through its resolved path, which differs from the
	// repo root when the checkout is reached through a symlink.
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil && resolved != repoRoot {
		needles = append(needles, resolved)
	}

	for _, name := range []string{"structure.txt", "comments.txt", "skips.txt"} {
		path := filepath.Join("testdata", "ir", name)

		raw, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read golden %s: %v\nregenerate with: %s", path, err, regenerateCmd)
			continue
		}

		for i, line := range strings.Split(string(raw), "\n") {
			for _, needle := range needles {
				if strings.Contains(line, needle) {
					t.Errorf("%s line %d carries the machine-local path %q, so the golden cannot match on another machine:\n %s",
						name, i+1, needle, truncate(line))

					break
				}
			}
		}
	}
}

// TestIRGoldensMatchFreshRun regenerates all three goldens from a real cc/v4
// parse and compares them byte for byte with the committed files. This is the
// test that makes the goldens a parser fingerprint rather than a snapshot
// someone has to remember to refresh: any parser change that moves the IR fails
// here with the offending line named.
//
// The run mirrors run(): dump the Module at the HeaderParser boundary, before
// applyManualFixups, and the skip set after Gen, once the codegen layer has
// added its own markers.
//
// It is not parallel and must never become parallel: t.Chdir is process-wide.
// The chdir puts the five *.gen.go files Gen writes into a temporary directory,
// where they are discarded, so a test run cannot touch the committed bindings.
func TestIRGoldensMatchFreshRun(t *testing.T) {
	repoRoot := setupHeaderParse(t)

	// Read the committed goldens before the chdir, while the relative package
	// paths still resolve.
	goldenNames := []string{"structure.txt", "comments.txt", "skips.txt"}
	golden := make(map[string][]byte, len(goldenNames))

	for _, name := range goldenNames {
		path := filepath.Join(repoRoot, irDumpDir, name)

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read golden %s: %v\nregenerate with: %s", path, err, regenerateCmd)
		}

		golden[name] = data
	}

	tmp := t.TempDir()
	t.Chdir(tmp)

	skips := &SkipCollector{}

	m := ccParser{}.Parse(skips)

	if err := writeIRDump(tmp, m); err != nil {
		t.Fatalf("writeIRDump: %v", err)
	}

	if err := applyManualFixups(m); err != nil {
		t.Fatalf("applyManualFixups: %v", err)
	}

	skips = Gen(m, skips)

	if err := writeSkipDump(tmp, skips); err != nil {
		t.Fatalf("writeSkipDump: %v", err)
	}

	for _, name := range goldenNames {
		fresh, err := os.ReadFile(filepath.Join(tmp, name))
		if err != nil {
			t.Errorf("read freshly generated %s: %v", name, err)
			continue
		}

		if bytes.Equal(fresh, golden[name]) {
			continue
		}

		diff := diffLines(string(golden[name]), string(fresh))
		if len(diff) == 0 {
			t.Errorf("%s differs from the committed golden in trailing bytes only (%d vs %d bytes)\nregenerate with: %s",
				name, len(golden[name]), len(fresh), regenerateCmd)

			continue
		}

		d := diff[0]
		t.Errorf("%s differs from the committed golden at line %d (%d differing lines):\n golden %s\n  fresh %s\n"+
			"if the parser changed on purpose, regenerate with: %s",
			name, d.num, len(diff), truncate(d.want), truncate(d.got), regenerateCmd)
	}
}

// testRepoRoot returns the absolute path of the checkout root.
//
// go test runs with the package directory as the working directory, so the
// repo root is two levels up.
func testRepoRoot(t *testing.T) string {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	repoRoot, err := filepath.Abs(filepath.Join(cwd, "..", ".."))
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}

	return repoRoot
}

// setupHeaderParse prepares a test that runs a real parse over the committed
// FFmpeg headers and returns the repo root. It skips the test when those
// headers are absent.
//
// AVLibPath is resolved at package init against the test's working directory,
// so it points at internal/generator/include, which does not exist. It is
// repointed at the real header tree for the duration of the test.
//
// A real parse logs one line per type it visits and would drown go test -v, so
// the log is discarded too. Both values are restored through t.Cleanup, and
// both are process-wide, so no caller may be marked parallel.
func setupHeaderParse(t *testing.T) string {
	t.Helper()

	repoRoot := testRepoRoot(t)

	includeDir := filepath.Join(repoRoot, "include")
	if _, err := os.Stat(includeDir); err != nil {
		t.Skipf("FFmpeg headers absent at %s: %v\ninitialise them with: git submodule update --init --recursive", includeDir, err)
	}

	originalLibPath := AVLibPath
	AVLibPath = includeDir

	t.Cleanup(func() { AVLibPath = originalLibPath })

	originalLogOutput := log.Writer()

	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(originalLogOutput) })

	return repoRoot
}
