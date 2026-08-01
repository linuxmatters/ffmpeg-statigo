package main

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// commentHeaderRoot resolves the committed FFmpeg header tree, or skips the test
// when the submodule is not initialised. AVLibPath is resolved at package init
// against the working directory, which under go test is the package directory,
// so it cannot be used as-is.
func commentHeaderRoot(t *testing.T) string {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	repoRoot, err := filepath.Abs(filepath.Join(cwd, "..", ".."))
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}

	includeDir := filepath.Join(repoRoot, "include")
	if _, err := os.Stat(includeDir); err != nil {
		t.Skipf("FFmpeg headers absent at %s: %v\ninitialise them with: git submodule update --init --recursive", includeDir, err)
	}

	return includeDir
}

func buildCommentTableFor(t *testing.T, root, rel string) *CommentTable {
	t.Helper()

	table, err := BuildCommentTable(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("BuildCommentTable(%s): %v", rel, err)
	}

	return table
}

// TestCommentSpans checks hand-verified spans across the block, /// run and
// ///< trailing forms.
func TestCommentSpans(t *testing.T) {
	root := commentHeaderRoot(t)

	cases := []struct {
		header     string
		startLine  int
		startCol   int
		endLine    int
		endCol     int
		startsLine bool
		wantPrefix string
	}{
		// Block comments, first thing on their line.
		{"libavutil/rational.h", 1, 1, 20, 4, true, "/*\n * rational numbers"},
		{"libavutil/rational.h", 22, 1, 27, 4, true, "/**\n * @file"},
		{"libavutil/rational.h", 36, 1, 53, 4, true, "/**\n * @defgroup lavu_math_rational"},
		{"libavutil/rational.h", 55, 1, 57, 4, true, "/**\n * Rational number"},
		{"libavutil/rational.h", 63, 1, 70, 4, true, "/**\n * Create an AVRational."},

		// ///< trailing comments on the AVRational fields.
		{"libavutil/rational.h", 59, 14, 59, 28, false, "///< Numerator"},
		{"libavutil/rational.h", 60, 14, 60, 30, false, "///< Denominator"},

		// Indented block comments inside AVFrame.
		{"libavutil/frame.h", 398, 1, 426, 4, true, "/**\n * This structure describes decoded"},
		{"libavutil/frame.h", 429, 5, 447, 8, true, "/**\n     * pointer to the picture/channel planes."},
		{"libavutil/frame.h", 450, 5, 471, 8, true, "/**\n     * For video, a positive or negative value"},

		// A two-member // run, merged into one entry spanning both lines.
		{"libavutil/channel_layout.h", 109, 5, 110, 60, true, "// leave space for 1024 ids"},
	}

	tables := map[string]*CommentTable{}

	for _, c := range cases {
		if tables[c.header] == nil {
			tables[c.header] = buildCommentTableFor(t, root, c.header)
		}

		got, ok := tables[c.header].StartingAt(c.startLine)
		if !ok {
			t.Errorf("%s: no comment starting at line %d", c.header, c.startLine)

			continue
		}

		if got.StartCol != c.startCol || got.EndLine != c.endLine || got.EndCol != c.endCol {
			t.Errorf("%s:%d span = (%d,%d)-(%d,%d), want (%d,%d)-(%d,%d)",
				c.header, c.startLine, got.StartLine, got.StartCol, got.EndLine, got.EndCol,
				c.startLine, c.startCol, c.endLine, c.endCol)
		}

		if got.StartsLine != c.startsLine {
			t.Errorf("%s:%d StartsLine = %v, want %v", c.header, c.startLine, got.StartsLine, c.startsLine)
		}

		if !strings.HasPrefix(got.Text, c.wantPrefix) {
			t.Errorf("%s:%d text = %.60q, want prefix %q", c.header, c.startLine, got.Text, c.wantPrefix)
		}

		if end, ok := tables[c.header].EndingAt(c.endLine); !ok || end.StartLine != c.startLine {
			t.Errorf("%s: EndingAt(%d) did not return the run starting at %d", c.header, c.endLine, c.startLine)
		}
	}
}

// TestChannelLayoutRunMerged names the merged run explicitly, because the two
// members must arrive as one entry with both line texts kept.
func TestChannelLayoutRunMerged(t *testing.T) {
	root := commentHeaderRoot(t)
	table := buildCommentTableFor(t, root, "libavutil/channel_layout.h")

	run, ok := table.StartingAt(109)
	if !ok {
		t.Fatal("no comment starting at libavutil/channel_layout.h:109")
	}

	// Each member keeps its own line text; the second member's indentation is
	// dropped by the join, which is harmless because processComment trims every
	// line before stripping markers.
	want := "// leave space for 1024 ids, which correspond to maximum order-32 harmonics,\n" +
		"// which should be enough for the foreseeable use cases"

	if run.Text != want {
		t.Errorf("merged run text =\n%q\nwant\n%q", run.Text, want)
	}

	if run.Members != 2 {
		t.Errorf("merged run Members = %d, want 2", run.Members)
	}

	if _, ok := table.StartingAt(110); ok {
		t.Error("line 110 still starts its own entry, so the run did not merge")
	}
}

func TestCommentMergeBehaviour(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want []Comment
	}{
		{
			name: "consecutive run merges",
			src: "/// one\n" +
				"/// two\n" +
				"/// three\n" +
				"int x;\n",
			want: []Comment{{StartLine: 1, EndLine: 3, Text: "/// one\n/// two\n/// three", StartsLine: true}},
		},
		{
			name: "blank line separates",
			src: "/// one\n" +
				"\n" +
				"/// two\n" +
				"int x;\n",
			want: []Comment{
				{StartLine: 1, EndLine: 1, Text: "/// one", StartsLine: true},
				{StartLine: 3, EndLine: 3, Text: "/// two", StartsLine: true},
			},
		},
		{
			// libavcodec/dv_profile.h holds this form, and the golden holds the
			// merged run for audio_min_samples, so a trailing comment starts a run
			// when an aligned ordinary comment follows it.
			name: "trailing comment absorbs an aligned ordinary continuation",
			src: "int x;  /* trailing */\n" +
				"        /* continued */\n" +
				"int y;\n",
			want: []Comment{
				{StartLine: 1, EndLine: 2, Text: "/* trailing */\n/* continued */", StartsLine: false},
			},
		},
		{
			// libavcodec/avcodec.h holds this form. The second comment is back at
			// the code's indent, so it documents the declaration below it.
			name: "trailing comment does not absorb a comment at the code indent",
			src: "int x;  /* trailing */\n" +
				"/* leading */\n" +
				"int y;\n",
			want: []Comment{
				{StartLine: 1, EndLine: 1, Text: "/* trailing */", StartsLine: false},
				{StartLine: 2, EndLine: 2, Text: "/* leading */", StartsLine: true},
			},
		},
		{
			// libavutil/pixfmt.h holds this form. A doxygen block below a `///<`
			// is the next declaration's comment, never a continuation.
			name: "trailing marker does not absorb a doxygen block",
			src: "int x; ///< trailing\n" +
				"/** leading */\n" +
				"int y;\n",
			want: []Comment{
				{StartLine: 1, EndLine: 1, Text: "///< trailing", StartsLine: false},
				{StartLine: 2, EndLine: 2, Text: "/** leading */", StartsLine: true},
			},
		},
		{
			name: "trailing comment does not join the preceding run",
			src: "/// leading\n" +
				"int x; ///< trailing\n",
			want: []Comment{
				{StartLine: 1, EndLine: 1, Text: "/// leading", StartsLine: true},
				{StartLine: 2, EndLine: 2, Text: "///< trailing", StartsLine: false},
			},
		},
		{
			name: "block then line comment merges",
			src: "/* one */\n" +
				"/// two\n" +
				"int x;\n",
			want: []Comment{{StartLine: 1, EndLine: 2, Text: "/* one */\n/// two", StartsLine: true}},
		},
		{
			name: "code after a block comment blocks the merge",
			src: "/* one */ int x;\n" +
				"/// two\n",
			want: []Comment{
				{StartLine: 1, EndLine: 1, Text: "/* one */", StartsLine: true},
				{StartLine: 2, EndLine: 2, Text: "/// two", StartsLine: true},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			table := BuildCommentTableFromSource("mem.h", []byte(c.src))
			got := table.All()

			if len(got) != len(c.want) {
				t.Fatalf("got %d entries, want %d: %+v", len(got), len(c.want), got)
			}

			for i, w := range c.want {
				g := got[i]

				if g.StartLine != w.StartLine || g.EndLine != w.EndLine {
					t.Errorf("entry %d lines = %d-%d, want %d-%d", i, g.StartLine, g.EndLine, w.StartLine, w.EndLine)
				}

				if g.Text != w.Text {
					t.Errorf("entry %d text = %q, want %q", i, g.Text, w.Text)
				}

				if g.StartsLine != w.StartsLine {
					t.Errorf("entry %d StartsLine = %v, want %v", i, g.StartsLine, w.StartsLine)
				}
			}
		})
	}
}

// TestCommentOffsets pins the byte offsets a matching rule reads the source
// with, because an off-by-one there would scan the comment's own last byte.
func TestCommentOffsets(t *testing.T) {
	src := "/* one */\nint x;\n"

	table := BuildCommentTableFromSource("mem.h", []byte(src))

	c := table.All()[0]
	if c.StartOff != 0 || c.EndOff != 9 {
		t.Fatalf("offsets = %d..%d, want 0..9", c.StartOff, c.EndOff)
	}

	if got := string(table.Src()[c.EndOff:table.Offset(2, 7)]); got != "\nint x;" {
		t.Errorf("text between comment and end of `x;` = %q, want %q", got, "\nint x;")
	}
}

func TestCommentLexerEdgeCases(t *testing.T) {
	cases := []struct {
		name  string
		src   string
		want  []string
		spans [][2]int // start line, end line
	}{
		{
			name: "block marker inside a string is not a comment",
			src:  "const char *s = \"/* not a comment */\";\n/* real */\n",
			want: []string{"/* real */"},
		},
		{
			name: "line marker inside a string is not a comment",
			src:  "const char *s = \"http://example.com\";\n// real\n",
			want: []string{"// real"},
		},
		{
			name: "quote inside a block comment does not open a string",
			src:  "/* he said \"hi\" and 'bye' */\nint x;\n// real\n",
			want: []string{"/* he said \"hi\" and 'bye' */", "// real"},
		},
		{
			name: "line comment inside a block comment stays one comment",
			src:  "/* outer // inner\n * still outer\n */\n",
			want: []string{"/* outer // inner\n * still outer\n */"},
		},
		{
			name: "block marker inside a line comment stays one comment",
			src:  "// outer /* inner\nint x;\n",
			want: []string{"// outer /* inner"},
		},
		{
			name:  "backslash continues a line comment onto the next line",
			src:   "#define M 1 // first \\\nsecond line still comment\nint x;\n",
			want:  []string{"// first \\\nsecond line still comment"},
			spans: [][2]int{{1, 2}},
		},
		{
			name: "escaped quote does not end a string",
			src:  "const char *s = \"a\\\"/* still string */b\";\n// real\n",
			want: []string{"// real"},
		},
		{
			name: "escaped quote in a character literal",
			src:  "char c = '\\'';\n/* real */\n",
			want: []string{"/* real */"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			src := []byte(c.src)
			got := lexComments(src, lineStarts(src))

			if len(got) != len(c.want) {
				t.Fatalf("got %d comments, want %d: %q", len(got), len(c.want), commentTexts(got))
			}

			for i, w := range c.want {
				if got[i].Text != w {
					t.Errorf("comment %d = %q, want %q", i, got[i].Text, w)
				}
			}

			for i, s := range c.spans {
				if got[i].StartLine != s[0] || got[i].EndLine != s[1] {
					t.Errorf("comment %d span = %d-%d, want %d-%d", i, got[i].StartLine, got[i].EndLine, s[0], s[1])
				}
			}
		})
	}
}

func commentTexts(cs []*Comment) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Text
	}

	return out
}

// TestCommentTableTotals pins the table's totals across every header in the
// files slice. The numbers are the measured corpus, so a lexer or merge-rule
// change that moves a single comment fails here rather than silently moving a
// doc comment in the generated output.
//
// Every merged run is two members: adjacent comment lines are rare in the
// FFmpeg headers, and the 1020 entries that do not start their line carry the
// volume of the trailing form.
func TestCommentTableTotals(t *testing.T) {
	root := commentHeaderRoot(t)

	const (
		wantRaw      = 3883
		wantEntries  = 3863
		wantRuns     = 20
		wantAbsorbed = 20
		wantTrailing = 1020
	)

	var raw, entries, runs, absorbed, trailing int

	for _, rel := range files {
		path := filepath.Join(root, rel)

		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}

		raw += len(lexComments(src, lineStarts(src)))

		table := buildCommentTableFor(t, root, rel)

		for _, c := range table.All() {
			entries++

			if c.Members > 1 {
				runs++
				absorbed += c.Members - 1

				if c.Members != 2 {
					t.Errorf("%s:%d merged run has %d members, want 2", rel, c.StartLine, c.Members)
				}
			}

			if !c.StartsLine {
				trailing++
			}
		}
	}

	t.Logf("headers=%d raw=%d entries=%d runs=%d absorbed=%d trailing=%d",
		len(files), raw, entries, runs, absorbed, trailing)

	for _, c := range []struct {
		name      string
		got, want int
	}{
		{"raw comments", raw, wantRaw},
		{"table entries", entries, wantEntries},
		{"merged runs", runs, wantRuns},
		{"members absorbed", absorbed, wantAbsorbed},
		{"entries that do not start their line", trailing, wantTrailing},
	} {
		if c.got != c.want {
			t.Errorf("%s = %d, want %d", c.name, c.got, c.want)
		}
	}
}

// matchSource writes src to a temp header, translates it with cc/v4, walks it
// and returns the doc comment of every declaration, keyed the way the comment
// golden keys it. It goes through the production walk rather than calling the
// claim passes directly, because the name positions the passes anchor on are
// half of what is under test.
func matchSource(t *testing.T, name, src string) map[string]string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)

	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}

	cfg, err := newCCConfig(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		t.Fatalf("cc/v4 configuration: %v", err)
	}

	ast, err := ccTranslate(cfg, path)
	if err != nil {
		t.Fatalf("translate: %v", err)
	}

	table, err := BuildCommentTable(path)
	if err != nil {
		t.Fatalf("comment table: %v", err)
	}

	w := &ccWalk{
		mod: &Module{
			functions: make(map[string]*Function),
			structs:   make(map[string]*Struct),
			enums:     make(map[string]*Enum),
			callbacks: make(map[string]*Function),
			constants: make(map[string]*Constant),
		},
		abs:   path,
		skips: &SkipCollector{},
		tags:  make(map[string]bool),
		table: table,
	}
	w.walkAST(ast)

	out := map[string]string{}

	for n, f := range w.mod.functions {
		out["function "+n] = f.Comment
	}

	for n, s := range w.mod.structs {
		out["struct "+n] = s.Comment

		for _, f := range s.Fields {
			out["struct "+n+"."+f.Name] = f.Comment
		}
	}

	for n, e := range w.mod.enums {
		out["enum "+n] = e.Comment

		for _, c := range e.Constants {
			out["enum "+n+"."+c.Name] = c.Comment
		}
	}

	return out
}

// wantComment asserts the comment claimed for one declaration key.
func wantComment(t *testing.T, matches map[string]string, key, want string) {
	t.Helper()

	got, ok := matches[key]
	if !ok {
		keys := make([]string, 0, len(matches))
		for k := range matches {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		t.Fatalf("%s: no declaration walked; keys=%v", key, keys)
	}

	if got != want {
		t.Errorf("%s: comment = %q, want %q", key, got, want)
	}
}

// TestMatchEnumTrailingComments pins the enum-member case a "block ending on
// the line above" rule gets wrong: it hands A's trailing comment to B as well.
func TestMatchEnumTrailingComments(t *testing.T) {
	matches := matchSource(t, "enum.h", `enum E {
    A = 1, ///< doc for A
    B = 2, ///< doc for B
};
`)

	wantComment(t, matches, "enum E.A", "   doc for A")
	wantComment(t, matches, "enum E.B", "   doc for B")
}

// TestMatchTrailingDoesNotCascade is the other half of that case: a member with
// no comment of its own gets nothing, because a comment documenting what
// precedes it never documents what follows it. libavcodec/codec_id.h:54-55,
// where the golden holds no AV_CODEC_ID_H261.
func TestMatchTrailingDoesNotCascade(t *testing.T) {
	matches := matchSource(t, "nocascade.h", `enum E {
    A, ///< doc for A
    B,
};
`)

	wantComment(t, matches, "enum E.A", "   doc for A")
	wantComment(t, matches, "enum E.B", "")
}

// TestMatchEnumRunCascade is libavcodec/codec_id.h:376-382. Enumerators are
// separated by commas, so nothing in the reject set stops one comment reaching
// every member of the run. The run ends at the next comment, not at the first
// member, and a blank line does not end it.
func TestMatchEnumRunCascade(t *testing.T) {
	matches := matchSource(t, "cascade.h", `enum E {
    /* various ADPCM codecs */
    A = 0x11000,
    B,
    C,

    /* next group */
    D,
};
`)

	for _, key := range []string{"enum E.A", "enum E.B", "enum E.C"} {
		wantComment(t, matches, key, "  various ADPCM codecs")
	}

	wantComment(t, matches, "enum E.D", "  next group")
}

// TestMatchTrailingNeedsTrailingForm covers libavutil/pixfmt.h:500. The comment
// follows the enumerator on its line but opens "///" with no "<", so it
// documents what follows rather than what precedes. The trailing pass does not
// claim it and the leading pass hands it to the next member down.
func TestMatchTrailingNeedsTrailingForm(t *testing.T) {
	matches := matchSource(t, "form.h", `enum E {
    A, /// not a trailing comment
    B,
};
`)

	wantComment(t, matches, "enum E.A", "")
	wantComment(t, matches, "enum E.B", "  not a trailing comment")
}

// TestMatchTrailingAbutsName pins the exclusive end of the name column. A
// comment cannot begin inside the name, so a comment starting exactly at that
// column sits immediately after the name and is that name's trailing comment.
// Rejecting the column dropped such a comment while the identical comment one
// space further along was claimed, which is a difference the writer never
// intended. No FFmpeg header writes the abutting form today, so only a test
// holds the rule.
func TestMatchTrailingAbutsName(t *testing.T) {
	matches := matchSource(t, "abut.h", `struct S {
    int abutting/**< doc for abutting */;
    int spaced /**< doc for spaced */;
};
`)

	wantComment(t, matches, "struct S.abutting", "   doc for abutting")
	wantComment(t, matches, "struct S.spaced", "   doc for spaced")
}

// TestMatchStructFieldLeadingBlock checks a leading block goes to the field
// below it and not to the next field down.
func TestMatchStructFieldLeadingBlock(t *testing.T) {
	matches := matchSource(t, "struct.h", `struct S {
    /**
     * doc for first
     */
    int first;

    int second;
};
`)

	wantComment(t, matches, "struct S.first", "  doc for first")
	wantComment(t, matches, "struct S.second", "")
}

// TestMatchMultiDeclaratorShared is "int width, height;",
// libavcodec/avcodec.h:600. The text between the comment and height is
// "int width,", which holds nothing in the reject set, so both names share it.
func TestMatchMultiDeclaratorShared(t *testing.T) {
	matches := matchSource(t, "multi.h", `struct S {
    /**
     * picture size
     */
    int width, height;
};
`)

	wantComment(t, matches, "struct S.width", "  picture size")
	wantComment(t, matches, "struct S.height", "  picture size")
}

// TestMatchInDeclarationTrailing is libavcodec/avdct.h:32. The comment sits
// inside the declaration's token span but after the name, so it is the name's
// trailing comment. Anchoring on the declaration's last token misses it.
func TestMatchInDeclarationTrailing(t *testing.T) {
	matches := matchSource(t, "indecl.h", `struct S {
    void (*idct)(short *block /* align 16 */);
};
`)

	wantComment(t, matches, "struct S.idct", "  align 16")
}

// TestMatchBlankLineGap checks that a blank line between a comment and the
// declaration does not block attachment.
func TestMatchBlankLineGap(t *testing.T) {
	matches := matchSource(t, "gap.h", `/**
 * doc for S
 */

struct S {
    int x;
};
`)

	wantComment(t, matches, "struct S", "  doc for S")
}

// TestMatchStructCommentDoesNotLeakToField checks the boundary the reject set
// carries, rather than assuming it: the struct's own "{" lies between its
// comment and its first field, so the comment cannot reach the field.
func TestMatchStructCommentDoesNotLeakToField(t *testing.T) {
	matches := matchSource(t, "leak.h", `/**
 * doc for S
 */
struct S {
    int x;
};
`)

	wantComment(t, matches, "struct S", "  doc for S")
	wantComment(t, matches, "struct S.x", "")
}

// TestMatchDirectiveBlocksLeading covers the rest of the reject set: a
// preprocessor directive between a comment and a declaration blocks
// attachment, because "#" is in the set.
func TestMatchDirectiveBlocksLeading(t *testing.T) {
	matches := matchSource(t, "directive.h", `/**
 * doc for S
 */
#define SOMETHING 1
struct S {
    int x;
};
`)

	wantComment(t, matches, "struct S", "")
}

// TestMatchEmbeddedTagClaimsNothing pins clang's getDeclLocForCommentSearch
// rule: a tag written inside another declaration's decl-specifier-seq, and not
// completed there, has no search location and so claims no comment. The block
// belongs to the typedef or the function instead.
func TestMatchEmbeddedTagClaimsNothing(t *testing.T) {
	matches := matchSource(t, "embedded.h", `/**
 * doc for the typedef
 */
typedef struct Opaque Opaque;

/**
 * doc for the function
 */
struct Ret *make(void);

/**
 * doc for the bare tag
 */
struct Bare;
`)

	wantComment(t, matches, "struct Opaque", "")
	wantComment(t, matches, "struct Ret", "")
	wantComment(t, matches, "function make", "  doc for the function")
	wantComment(t, matches, "struct Bare", "  doc for the bare tag")
}

// TestMatchMacroCollapsedPositions pins the three routes spellingSite takes when
// cc/v4 collapses a declaration's position onto the macro invocation that wrote
// it. The macro is written out rather than taken from libavutil/bprint.h,
// because the rule is keyed on "the reported position does not spell the name"
// and on nothing about FF_PAD_STRUCTURE or AVBPrint.
func TestMatchMacroCollapsedPositions(t *testing.T) {
	matches := matchSource(t, "pad.h", `/**
 * Define a padded structure
 */
#define PAD_STRUCT(name, size, ...) \
struct helper_##name { __VA_ARGS__ }; \
typedef struct name { \
    __VA_ARGS__ \
    char padding[size - sizeof(struct helper_##name)]; \
} name;

/**
 * A buffer
 */
PAD_STRUCT(Buf, 64,
    char *str;    /**< the string */
    unsigned len; /**< the length */
    char tail[1];
)
`)

	// A macro argument keeps the position it is written at, so each field claims
	// its own trailing comment and both expansions of the argument claim the
	// same one.
	for _, owner := range []string{"struct helper_Buf", "struct Buf"} {
		wantComment(t, matches, owner+".str", "   the string")
		wantComment(t, matches, owner+".len", "   the length")
		wantComment(t, matches, owner+".tail", "")
	}

	// The tag written as an argument keeps its position too, and the tag built by
	// "##" pasting has no written position at all, so it keeps the invocation's.
	// Both sit under the same block comment and both claim it.
	wantComment(t, matches, "struct Buf", "  A buffer")
	wantComment(t, matches, "struct helper_Buf", "  A buffer")

	// A name the macro body declares is positioned at its own token in the
	// replacement list, where the "#define" blocks any leading comment.
	wantComment(t, matches, "struct Buf.padding", "")
}

// TestMacroArgumentSite pins the raw-source scan MacroArgumentSite runs: it
// crosses nested parentheses and comments, matches whole identifiers only, and
// reports the first written occurrence.
func TestMacroArgumentSite(t *testing.T) {
	src := "M(A, f(x, y), /* size, len */\n  int size_max;\n  int size;\n)\n"

	table := BuildCommentTableFromSource("m.h", []byte(src))

	for _, c := range []struct {
		name      string
		line, col int
		wantLine  int
		wantCol   int
		wantFound bool
	}{
		{name: "y", line: 1, col: 1, wantLine: 1, wantCol: 11, wantFound: true},
		{name: "size", line: 1, col: 1, wantLine: 3, wantCol: 7, wantFound: true},
		{name: "size_max", line: 1, col: 1, wantLine: 2, wantCol: 7, wantFound: true},
		{name: "absent", line: 1, col: 1},
		{name: "size", line: 2, col: 3},
	} {
		line, col, ok := table.MacroArgumentSite(c.name, c.line, c.col)
		if ok != c.wantFound || line != c.wantLine || col != c.wantCol {
			t.Errorf("MacroArgumentSite(%q, %d, %d) = %d, %d, %v; want %d, %d, %v",
				c.name, c.line, c.col, line, col, ok, c.wantLine, c.wantCol, c.wantFound)
		}
	}
}
