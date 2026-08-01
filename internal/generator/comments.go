package main

// Position-indexed comment table for one header, built by lexing the raw bytes.
//
// The cc/v4 token stream cannot supply absolute comment positions. A token's
// Sep() accumulates everything skipped since the previous token, which spans
// #include boundaries, preprocessor directive lines are stripped from it, and
// cc/v4 exposes no absolute sep position. Measured over the pinned headers,
// deriving a position from a sep is byte-exact only 93.5% of the time. Lexing
// the raw bytes gives exact spans, so it is the table's only source.

import (
	"os"
	"sort"
	"strings"
)

// Comment is one comment, or one merged run of comments, in a header.
//
// Lines are 1-based. StartCol is the 1-based byte column of the first byte of
// the comment. EndLine and EndCol point one byte past the comment's last byte,
// matching the end-column convention used for declaration tokens. StartOff and
// EndOff are the same two points as byte offsets into the header, so the
// matcher can read the raw text between a comment and a declaration.
type Comment struct {
	StartLine, StartCol int
	EndLine, EndCol     int
	StartOff, EndOff    int
	Text                string // raw text including markers, run members joined with "\n"
	StartsLine          bool   // only whitespace precedes the comment on its start line
	Members             int    // comments folded into this run; 1 when nothing merged

	// Trailing marks a comment that documents what precedes it: either the
	// `///<` marker family, or a plain comment that follows code on its line.
	// A merged run keeps its first member's value.
	Trailing bool
}

// CommentTable holds every comment in one header, with runs merged, indexed by
// line for O(1) lookup.
//
// A line can carry more than one comment, so both line indexes hold a slice.
// `/* a */ int x; /* b */` is two runs starting and ending on the same line, and
// the trailing pass has to see both to pick the leftmost one past the name.
type CommentTable struct {
	Path     string
	src      []byte
	starts   []int
	comments []*Comment
	byEnd    map[int][]*Comment
	byStart  map[int][]*Comment
}

// BuildCommentTable lexes a header file and returns its comment table.
func BuildCommentTable(path string) (*CommentTable, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return BuildCommentTableFromSource(path, src), nil
}

// BuildCommentTableFromSource builds a table from bytes already in memory. The
// tests use it for small unambiguous sources.
func BuildCommentTableFromSource(path string, src []byte) *CommentTable {
	starts := lineStarts(src)
	raw := lexComments(src, starts)
	merged := mergeCommentRuns(raw, src, starts)

	t := &CommentTable{
		Path:     path,
		src:      src,
		starts:   starts,
		comments: merged,
		byEnd:    make(map[int][]*Comment, len(merged)),
		byStart:  make(map[int][]*Comment, len(merged)),
	}

	for _, c := range merged {
		t.byEnd[c.EndLine] = append(t.byEnd[c.EndLine], c)
		t.byStart[c.StartLine] = append(t.byStart[c.StartLine], c)
	}

	return t
}

// EndingAt returns the first merged run whose EndLine is line.
func (t *CommentTable) EndingAt(line int) (Comment, bool) {
	cs := t.byEnd[line]
	if len(cs) == 0 {
		return Comment{}, false
	}

	return *cs[0], true
}

// StartingAt returns the first merged run whose StartLine is line.
func (t *CommentTable) StartingAt(line int) (Comment, bool) {
	cs := t.byStart[line]
	if len(cs) == 0 {
		return Comment{}, false
	}

	return *cs[0], true
}

// All returns the merged runs in source order, which is also ascending StartOff
// order, so a matching rule can binary-search it.
func (t *CommentTable) All() []*Comment { return t.comments }

// Src returns the header's raw bytes. A matching rule reads the text between a
// comment and a declaration out of it, so the rules work on the written source
// rather than on the token stream.
func (t *CommentTable) Src() []byte { return t.src }

// Offset converts a 1-based line and byte column to a byte offset into Src. An
// out-of-range position clamps to the end of the source, so a caller never
// slices out of bounds.
func (t *CommentTable) Offset(line, col int) int {
	if line < 1 || line >= len(t.starts) {
		return len(t.src)
	}

	off := t.starts[line] + col - 1
	if off > len(t.src) {
		return len(t.src)
	}

	if off < 0 {
		return 0
	}

	return off
}

// lineStarts returns the byte offset of each line start, indexed 1-based. Index
// 0 is unused so a line number indexes directly.
func lineStarts(src []byte) []int {
	starts := []int{0, 0}

	for i, b := range src {
		if b == '\n' {
			starts = append(starts, i+1)
		}
	}

	return starts
}

// lexComments scans C source and returns every comment in source order.
//
// It skips string and character literals, honours backslash-newline line
// continuations inside line comments, and does not treat comment markers found
// inside a literal as comments. Trigraphs are not handled; the FFmpeg headers
// do not need them.
func lexComments(src []byte, starts []int) []*Comment {
	var out []*Comment

	n := len(src)
	line, col := 1, 1

	// step consumes the byte at i and returns the next index.
	step := func(i int) int {
		if src[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}

		return i + 1
	}

	// continuation reports the length of a backslash-newline at i, or 0.
	continuation := func(i int) int {
		if i >= n || src[i] != '\\' {
			return 0
		}

		j := i + 1
		if j < n && src[j] == '\r' {
			j++
		}

		if j < n && src[j] == '\n' {
			return j + 1 - i
		}

		return 0
	}

	for i := 0; i < n; {
		switch {
		case src[i] == '/' && i+1 < n && src[i+1] == '*':
			startOff, startLine, startCol := i, line, col
			i = step(step(i))

			for i < n {
				if src[i] == '*' && i+1 < n && src[i+1] == '/' {
					i = step(step(i))

					break
				}

				i = step(i)
			}

			out = append(out, newComment(src, starts, startOff, i, startLine, startCol, line, col))

		case src[i] == '/' && i+1 < n && src[i+1] == '/':
			startOff, startLine, startCol := i, line, col
			i = step(step(i))

			for i < n {
				if k := continuation(i); k > 0 {
					for end := i + k; i < end; {
						i = step(i)
					}

					continue
				}

				if src[i] == '\n' {
					break
				}

				i = step(i)
			}

			out = append(out, newComment(src, starts, startOff, i, startLine, startCol, line, col))

		case src[i] == '"' || src[i] == '\'':
			quote := src[i]
			i = step(i)

			for i < n && src[i] != quote {
				if src[i] == '\\' && i+1 < n {
					i = step(step(i))

					continue
				}

				if src[i] == '\n' {
					break // unterminated literal; do not swallow the rest of the file
				}

				i = step(i)
			}

			if i < n && src[i] == quote {
				i = step(i)
			}

		case src[i] == '\\':
			if k := continuation(i); k > 0 {
				for end := i + k; i < end; {
					i = step(i)
				}

				continue
			}

			i = step(i)

		default:
			i = step(i)
		}
	}

	return out
}

// newComment builds a Comment from a lexed span. endOff, endLine and endCol are
// exclusive: they point one byte past the comment.
func newComment(src []byte, starts []int, startOff, endOff, startLine, startCol, endLine, endCol int) *Comment {
	text := string(src[startOff:endOff])
	startsLine := strings.TrimSpace(string(src[starts[startLine]:startOff])) == ""

	return &Comment{
		StartLine:  startLine,
		StartCol:   startCol,
		EndLine:    endLine,
		EndCol:     endCol,
		StartOff:   startOff,
		EndOff:     endOff,
		Text:       text,
		StartsLine: startsLine,
		Trailing:   isTrailingMarker(text) || (!isDocMarker(text) && !startsLine),
	}
}

// isDocMarker reports whether the comment opens with a doxygen marker: ///,
// //!, /** or /*!. Everything else is an ordinary comment.
func isDocMarker(text string) bool {
	for _, p := range []string{"///", "//!", "/**", "/*!"} {
		if strings.HasPrefix(text, p) {
			return true
		}
	}

	return false
}

// isTrailingMarker reports whether the comment opens with the member form of a
// doxygen marker: ///<, //!<, /**< or /*!<.
func isTrailingMarker(text string) bool {
	for _, p := range []string{"///<", "//!<", "/**<", "/*!<"} {
		if strings.HasPrefix(text, p) {
			return true
		}
	}

	return false
}

// mergeCommentRuns joins comments that sit on consecutive lines with nothing but
// whitespace between them, because such a run documents one declaration and the
// comment golden pins it as one entry.
//
// Comment b joins the run ending with a when b starts on the line after a ends,
// b is the first thing on its own line, nothing but whitespace follows a on a's
// end line, and the two agree about trailing-ness. A blank line between two
// comments always separates them.
//
// The run's Text is the members joined with "\n", so each member keeps its own
// line text but loses its indentation. processComment trims every line before
// it strips markers, so that costs nothing.
//
// Marker family alone does not decide a merge: a /* */ block still merges with
// a following /// run.
func mergeCommentRuns(raw []*Comment, src []byte, starts []int) []*Comment {
	var (
		out  []*Comment
		last *Comment // last member folded into out's final entry
	)

	for _, c := range raw {
		if run := lastComment(out); run != nil && joinsRun(last, c, src, starts) {
			run.EndLine, run.EndCol, run.EndOff = c.EndLine, c.EndCol, c.EndOff
			run.Text += "\n" + c.Text
			run.Members++
			last = c

			continue
		}

		clone := *c
		clone.Members = 1
		out = append(out, &clone)
		last = c
	}

	return out
}

func lastComment(cs []*Comment) *Comment {
	if len(cs) == 0 {
		return nil
	}

	return cs[len(cs)-1]
}

// joinsRun reports whether b continues the run whose last member is a.
func joinsRun(a, b *Comment, src []byte, starts []int) bool {
	if a == nil || !b.StartsLine {
		return false
	}

	if b.StartLine != a.EndLine+1 {
		return false
	}

	if a.Trailing != b.Trailing && !wrapsTrailing(a, b) {
		return false
	}

	return restOfLineIsBlank(src, starts, a.EndLine, a.EndCol)
}

// wrapsTrailing reports whether b continues the trailing comment a rather than
// documenting the declaration below it. A trailing comment absorbs a following
// comment only when that comment is an ordinary one starting at the same
// column, which is how FFmpeg writes a wrapped trailing comment. A doxygen
// block below a `///<`, or an ordinary comment back at the indent of the code,
// documents the next declaration instead.
//
// The alignment test is empirical. It is taken from the comment golden in
// testdata/ir, which answers the question in both directions: a wrapped
// trailing comment in libavcodec/dv_profile.h merges, while an ordinary comment
// back at the code's indent in libavcodec/avcodec.h and in
// libavformat/avformat.h does not. It is not a mechanism read out of clang's
// source.
func wrapsTrailing(a, b *Comment) bool {
	return a.Trailing && !isDocMarker(b.Text) && b.StartCol == a.StartCol
}

// restOfLineIsBlank reports whether only whitespace follows column col on line.
func restOfLineIsBlank(src []byte, starts []int, line, col int) bool {
	if line < 1 || line >= len(starts) {
		return true
	}

	from := starts[line] + col - 1
	if from >= len(src) {
		return true
	}

	end := len(src)
	if line+1 < len(starts) {
		end = starts[line+1]
	}

	return strings.TrimSpace(string(src[from:end])) == ""
}

// Matching a comment to a declaration (WW-141 phase 4.2).
//
// Two passes per declaration, in this order.
//
//  1. Trailing. The first comment starting after the declared name, on the
//     name's own line. The anchor is the identifier, not the last token of the
//     declaration, which is why "void (*idct)(int16_t *block /* align 16 */);"
//     (libavcodec/avdct.h:32) gives "align 16" to idct rather than to nothing.
//  2. Leading, only where pass 1 claimed nothing. The nearest comment beginning
//     before the declared name, rejected when it documents what precedes it or
//     when the raw source between it and the name holds a rejectSet character.
//
// The order is the design: a trailing comment is searched for first, and only
// then the comment in front. What stops the previous
// enumerator's trailing comment reaching the next enumerator is pass 2's
// trailing test, not adjacency.
//
// Neither pass is exclusive and neither tests adjacency, so one comment can
// document several declarations. That is deliberate: a comment above a run of
// enumerators reaches every one of them, because commas separate enumerators
// and nothing in rejectSet intervenes.

// rejectSet is the character set clang refuses to carry a leading comment
// across, from ASTContext::getRawCommentForDeclNoCache. A ";" or "}" means
// another declaration ended, "{" means the declaration sits inside a body the
// comment is outside of, "#" is a preprocessor directive and "@" an Objective-C
// keyword.
const rejectSet = ";{}#@"

// CommentAt returns the processed doc comment claimed for a declaration whose
// name ends at nameLine and nameCol, or "" when nothing is claimed. nameCol is
// one byte past the name's last byte.
//
// The text goes through processComment, which turns a raw C comment into the
// text the emitter writes.
func (t *CommentTable) CommentAt(nameLine, nameCol int) string {
	c, _ := t.claim(nameLine, nameCol)
	if c == nil {
		return ""
	}

	return processComment(c.Text)
}

// claim runs both passes and returns the comment claimed and the rule that
// claimed it: "trailing", "leading", or "" when nothing was claimed.
func (t *CommentTable) claim(nameLine, nameCol int) (*Comment, string) {
	if c := t.trailingCandidate(nameLine, nameCol); c != nil {
		return c, "trailing"
	}

	if c := t.leadingCandidate(nameLine, nameCol); c != nil {
		return c, "leading"
	}

	return nil, ""
}

// trailingCandidate returns the leftmost trailing-form comment starting on
// nameLine past column nameCol.
//
// Adjacency is not required: an Enumerator's token span excludes the separating
// comma, so a comment after "A = 1," starts past the name but not next to it.
// The comment may also sit inside the declaration's own token span, because
// nameCol is one past the name rather than past the declaration.
//
// Only a Trailing comment is claimed, and the golden pins the distinction:
// "AV_PIX_FMT_OHCODEC, /// hardware decoding through openharmony"
// (libavutil/pixfmt.h:500) opens with a doxygen marker that carries no "<", so
// it documents what follows rather than what precedes, and the golden holds no
// comment for that enumerator.
func (t *CommentTable) trailingCandidate(nameLine, nameCol int) *Comment {
	var best *Comment

	for _, c := range t.byStart[nameLine] {
		if !c.Trailing || c.StartCol <= nameCol {
			continue
		}

		if best == nil || c.StartCol < best.StartCol {
			best = c
		}
	}

	return best
}

// Macro-collapsed positions (WW-141 phase 5.1).
//
// cc/v4 rewrites the position of every token a macro expands to. cpp.go sets
// subst[i].off = T.off, where T is the invocation token, so a declaration
// written inside a macro invocation reports the invocation's own line and
// column instead of its own. A macro argument has to be resolved back to the
// position the argument is written at, or it claims the wrong comment.
//
// The written positions are still in the raw bytes this table already lexes, so
// they are recovered from there: read the identifier at the reported position,
// and when it is not the name being looked for, the position is an invocation
// site and the name is resolved back to where it comes from, which for an
// argument is the invocation's own argument text.

// IdentifierAt returns the identifier written at a 1-based line and byte
// column, or "" when no identifier starts there.
func (t *CommentTable) IdentifierAt(line, col int) string {
	return identifierAt(t.src, t.Offset(line, col))
}

// MacroArgumentSite returns the position at which name is written inside the
// argument list of the macro invocation at line and col.
//
// This is clang's macro-argument branch of SourceManager::getFileLoc: a token
// that reaches the expansion through an argument keeps the argument's own
// spelling position. The first whole-identifier match wins, because an argument
// is written once and both structs FF_PAD_STRUCTURE expands to share it.
//
// It reports false when the position is not an invocation, when the argument
// list is unterminated, or when the name is not written in it, which is the
// case for a name no source position spells because "##" pasted it together.
func (t *CommentTable) MacroArgumentSite(name string, line, col int) (int, int, bool) {
	off := t.Offset(line, col)

	invoked := identifierAt(t.src, off)
	if invoked == "" {
		return 0, 0, false
	}

	open := skipTrivia(t.src, off+len(invoked))
	if open >= len(t.src) || t.src[open] != '(' {
		return 0, 0, false
	}

	closing := matchingParen(t.src, open)
	if closing < 0 {
		return 0, 0, false
	}

	at := identifierOffset(t.src, open+1, closing, name)
	if at < 0 {
		return 0, 0, false
	}

	l, c := t.position(at)

	return l, c, true
}

// position converts a byte offset into a 1-based line and byte column.
func (t *CommentTable) position(off int) (line, col int) {
	line = max(sort.Search(len(t.starts), func(i int) bool { return t.starts[i] > off })-1, 1)

	return line, off - t.starts[line] + 1
}

// identifierAt returns the identifier starting at off, or "".
func identifierAt(src []byte, off int) string {
	if off < 0 || off >= len(src) || !isIdentStart(src[off]) {
		return ""
	}

	end := off
	for end < len(src) && isIdentByte(src[end]) {
		end++
	}

	return string(src[off:end])
}

func isIdentStart(b byte) bool {
	return b == '_' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

func isIdentByte(b byte) bool {
	return isIdentStart(b) || b >= '0' && b <= '9'
}

// skipNonCode returns the index just past a comment or a string or character
// literal starting at i, or i itself when the byte at i is ordinary code.
func skipNonCode(src []byte, i int) int {
	n := len(src)

	switch {
	case src[i] == '/' && i+1 < n && src[i+1] == '*':
		for j := i + 2; j < n; j++ {
			if src[j] == '*' && j+1 < n && src[j+1] == '/' {
				return j + 2
			}
		}

		return n

	case src[i] == '/' && i+1 < n && src[i+1] == '/':
		for j := i + 2; j < n; j++ {
			if src[j] == '\\' && j+1 < n && src[j+1] == '\n' {
				j++

				continue
			}

			if src[j] == '\n' {
				return j
			}
		}

		return n

	case src[i] == '"' || src[i] == '\'':
		quote := src[i]

		for j := i + 1; j < n; j++ {
			switch {
			case src[j] == '\\' && j+1 < n:
				j++
			case src[j] == quote:
				return j + 1
			case src[j] == '\n':
				return j // unterminated literal; do not swallow the rest of the file
			}
		}

		return n
	}

	return i
}

// skipTrivia advances past whitespace, line continuations and comments.
func skipTrivia(src []byte, i int) int {
	for i < len(src) {
		if j := skipNonCode(src, i); j != i {
			i = j

			continue
		}

		switch {
		case src[i] == ' ', src[i] == '\t', src[i] == '\r', src[i] == '\n':
			i++
		case src[i] == '\\' && i+1 < len(src) && (src[i+1] == '\n' || src[i+1] == '\r'):
			i++
		default:
			return i
		}
	}

	return i
}

// matchingParen returns the index of the ")" closing the "(" at open, or -1.
func matchingParen(src []byte, open int) int {
	depth := 0

	for i := open; i < len(src); {
		if j := skipNonCode(src, i); j != i {
			i = j

			continue
		}

		switch src[i] {
		case '(':
			depth++
		case ')':
			depth--

			if depth == 0 {
				return i
			}
		}

		i++
	}

	return -1
}

// identifierOffset returns the offset of the first whole identifier equal to
// name in src[from:to), skipping comments and literals, or -1.
func identifierOffset(src []byte, from, to int, name string) int {
	for i := from; i < to; {
		if j := skipNonCode(src, i); j != i {
			i = j

			continue
		}

		if !isIdentStart(src[i]) {
			i++

			continue
		}

		end := i
		for end < len(src) && isIdentByte(src[end]) {
			end++
		}

		if end <= to && string(src[i:end]) == name {
			return i
		}

		i = end
	}

	return -1
}

// leadingCandidate returns the comment clang would attach in front of a name:
// the nearest one beginning before it, unless it documents what precedes it, or
// unless the raw source strictly between its end and the name holds a rejectSet
// character.
//
// The reject set is what carries the struct boundary: a struct's own "{" lies
// between its comment and its first field, so the comment cannot leak down.
func (t *CommentTable) leadingCandidate(nameLine, nameCol int) *Comment {
	nameEnd := t.Offset(nameLine, nameCol)

	all := t.All()

	i := sort.Search(len(all), func(i int) bool { return all[i].StartOff >= nameEnd })
	if i == 0 {
		return nil
	}

	c := all[i-1]

	// A comment that documents what precedes it never documents what follows it.
	// The golden confirms it in both marker forms: libavcodec/codec_id.h:54
	// documents AV_CODEC_ID_MPEG2VIDEO with a "///<" and the golden holds no
	// record for AV_CODEC_ID_H261 on the line after, and
	// libavcodec/avcodec.h:447 closes with a plain "/* see AVMEDIA_TYPE_xxx */"
	// and the golden holds no struct AVCodecContext.codec. So the rejection is
	// not about the "<" marker alone, which is what Comment.Trailing encodes.
	if c.Trailing {
		return nil
	}

	if c.EndOff > nameEnd {
		return nil
	}

	if strings.ContainsAny(string(t.Src()[c.EndOff:nameEnd]), rejectSet) {
		return nil
	}

	return c
}
