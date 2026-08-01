package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// irDumpDir is where the committed IR goldens live. It is relative to the
// process working directory, matching Gen, which writes the *.gen.go files
// relative to cwd; the generator is always run from the repo root. Every writer
// below takes the directory as a parameter so a test can point it at a temp
// directory instead.
const irDumpDir = "internal/generator/testdata/ir"

const (
	structureFileHeader = "# generator IR structural dump (WW-139); test artefact, format is unstable by design"
	commentsFileHeader  = "# generator IR comment dump (WW-139); one record per symbol carrying a comment"
	skipsFileHeader     = "# generator skip dump (WW-139)"
)

// includeSegment is the checkout-relative form the dump substitutes for the
// absolute header directory, and the anchor normalizeCTypeName cuts back to.
const includeSegment = "include" + string(filepath.Separator)

// normalizeCTypeName rewrites the absolute header path libclang embeds in a C
// type spelling into a checkout-relative one.
//
// libclang spells an unnamed union or struct as "union (unnamed union at
// <absolute-header-path>:<line>:<col>)", so Field.CTypeName carries the path of
// the checkout that parsed the headers. Written to the golden unchanged, that
// path would make the file match only on the machine and at the directory that
// generated it, and fail everywhere else, CI included.
//
// AVLibPath is the absolute include directory and is exactly the prefix
// libclang reports, so replacing it covers the normal case. The fallback covers
// a libclang that reports a resolved or symlinked path that does not literally
// match AVLibPath: an absolute path left in the spelling still names a header
// below an include directory, so it is cut back to the last "include/" segment,
// which is stable wherever the checkout lives.
func normalizeCTypeName(s string) string {
	s = strings.ReplaceAll(s, AVLibPath+string(filepath.Separator), includeSegment)

	i := strings.LastIndex(s, includeSegment)
	if i < 0 {
		return s
	}

	// The spelling puts the path after a space ("... at /abs/path/x.h:1:1)"), so
	// the path token starts at the first byte after the preceding space.
	start := strings.LastIndexByte(s[:i], ' ') + 1
	if s[start] != filepath.Separator {
		return s
	}

	return s[:start] + s[i:]
}

// renderType expands a parsed Type into a single-line textual form for the IR
// golden dump. It reuses (*IdentType).String() and reproduces the pointer and
// array bracket forms their String() methods produce, but it does not build on
// String() throughout, because three implementations are lossy:
//
//   - FuncType.String() returns the literal "func", dropping Args and Result.
//   - UnionType.String() returns the literal "union", dropping Fields.
//   - IdentType.String() returns Name, dropping IsAnonEnum.
//
// A dump built on those three would hide a parser change in a callback argument
// type, a union member, or the anonymous-enum classification, which is most of
// what the dump exists to catch. renderType expands those three cases and
// recurses on Inner, so a lossy node nested inside a pointer or an array is
// still expanded.
//
// A nil Type is the parser's encoding of C void, so it renders "void" rather
// than an empty string.
func renderType(t Type) string {
	switch t := t.(type) {
	case nil:
		return "void"

	case *IdentType:
		if t.IsAnonEnum {
			return t.String() + "/anonenum"
		}
		return t.String()

	case *PointerType:
		return "*" + renderType(t.Inner)

	case *ConstArray:
		return fmt.Sprintf("%s[%d]", renderType(t.Inner), t.Size)

	case *Array:
		return renderType(t.Inner) + "[?]"

	case *FuncType:
		args := make([]string, len(t.Args))
		for i, a := range t.Args {
			args[i] = renderType(a)
		}
		return fmt.Sprintf("func(%s)%s", strings.Join(args, ","), renderType(t.Result))

	case *UnionType:
		var b strings.Builder
		b.WriteString("union{")
		for _, f := range t.Fields {
			fmt.Fprintf(&b, "%s %s;", f.Name, renderType(f.Type))
		}
		b.WriteString("}")
		return b.String()

	default:
		// An unhandled implementation must be loud in the diff rather than
		// silently rendering as something plausible, so name the Go type.
		return fmt.Sprintf("unknown(%T)", t)
	}
}

// commentHash fingerprints a doc comment for the structural dump, so a comment
// change shows as a one-line diff without the comment body bloating that
// stream. Twelve hex characters of SHA-256 are wide enough that a collision is
// not a practical concern at this corpus size, and fixed width keeps the
// structural lines aligned.
//
// An empty comment renders "-" rather than the hash of the empty string, so
// "no comment" and "some comment" can never collide and the distinction stays
// readable.
func commentHash(s string) string {
	if s == "" {
		return "-"
	}

	sum := sha256.Sum256([]byte(s))

	return hex.EncodeToString(sum[:])[:12]
}

// renderBool exists so every boolean field in the dump goes through one
// spelling ("true"/"false"), rather than depending on a format verb chosen per
// call site.
func renderBool(b bool) string {
	return strconv.FormatBool(b)
}

// orDash substitutes "-" for an empty string field. Every dump line is
// whitespace-separated, so an empty field would otherwise collapse two columns
// into one and shift the type expansion that always ends the line.
func orDash(s string) string {
	if s == "" {
		return "-"
	}

	return s
}

// renderModuleStructure writes the structural stream: the five order slices
// verbatim, then one record per symbol. Every line names its symbol, so
// grepping a symbol name returns its whole record, and no line carries an index
// into the order slices, so inserting a symbol is a one-line diff rather than a
// shifted tail.
//
// Determinism is the whole point of the artefact, so nothing here ranges over a
// map. Every walk goes through an order slice and looks the symbol up in the
// matching map; an order entry with no map entry emits a MISSING marker, which
// catches a parser that populates one side only.
//
// A variable-width type expansion always ends its line, because a union or
// function expansion contains spaces. The C type spelling has the same problem
// but sits mid-line, so it is quoted with %q instead, which also renders an
// empty spelling as "" rather than needing orDash.
func renderModuleStructure(w io.Writer, m *Module) error {
	bw := bufio.NewWriter(w)

	fmt.Fprintln(bw, structureFileHeader)

	for _, name := range m.functionOrder {
		fmt.Fprintf(bw, "order function %s\n", name)
	}

	for _, name := range m.structOrder {
		fmt.Fprintf(bw, "order struct %s\n", name)
	}

	for _, name := range m.enumOrder {
		fmt.Fprintf(bw, "order enum %s\n", name)
	}

	for _, name := range m.callbackOrder {
		fmt.Fprintf(bw, "order callback %s\n", name)
	}

	for _, name := range m.constantOrder {
		fmt.Fprintf(bw, "order constant %s\n", name)
	}

	for _, name := range m.functionOrder {
		f, ok := m.functions[name]
		if !ok {
			fmt.Fprintf(bw, "function %s MISSING\n", name)
			continue
		}

		renderFunctionRecord(bw, "function", name, f)
	}

	for _, name := range m.structOrder {
		s, ok := m.structs[name]
		if !ok {
			fmt.Fprintf(bw, "struct %s MISSING\n", name)
			continue
		}

		fmt.Fprintf(bw, "struct %s typedefd=%s byvalue=%s comment=%s\n",
			name, renderBool(s.Typedefd), renderBool(s.ByValue), commentHash(s.Comment))

		for i, f := range s.Fields {
			fmt.Fprintf(bw, "struct %s field %d %s ctype=%q bits=%d comment=%s %s\n",
				name, i, orDash(f.Name), normalizeCTypeName(f.CTypeName), f.BitWidth, commentHash(f.Comment), renderType(f.Type))
		}
	}

	for _, name := range m.enumOrder {
		e, ok := m.enums[name]
		if !ok {
			fmt.Fprintf(bw, "enum %s MISSING\n", name)
			continue
		}

		fmt.Fprintf(bw, "enum %s typedefd=%s comment=%s\n", name, renderBool(e.Typedefd), commentHash(e.Comment))

		for i, c := range e.Constants {
			fmt.Fprintf(bw, "enum %s constant %d %s fromenum=%s comment=%s\n",
				name, i, orDash(c.Name), renderBool(c.FromEnum), commentHash(c.Comment))
		}
	}

	for _, name := range m.callbackOrder {
		f, ok := m.callbacks[name]
		if !ok {
			fmt.Fprintf(bw, "callback %s MISSING\n", name)
			continue
		}

		renderFunctionRecord(bw, "callback", name, f)
	}

	for _, name := range m.constantOrder {
		c, ok := m.constants[name]
		if !ok {
			fmt.Fprintf(bw, "constant %s MISSING\n", name)
			continue
		}

		fmt.Fprintf(bw, "constant %s fromenum=%s comment=%s\n", name, renderBool(c.FromEnum), commentHash(c.Comment))
	}

	return bw.Flush()
}

// renderFunctionRecord writes the shared function record form. Functions and
// callbacks are both *Function and differ only in the kind word that opens each
// line, so one writer covers both and the two streams cannot drift apart.
func renderFunctionRecord(w io.Writer, kind, name string, f *Function) {
	fmt.Fprintf(w, "%s %s variadic=%s ptr=%s comment=%s\n",
		kind, name, renderBool(f.Variadic), renderBool(f.Ptr), commentHash(f.Comment))
	fmt.Fprintf(w, "%s %s result %s\n", kind, name, renderType(f.Result))

	for i, p := range f.Args {
		fmt.Fprintf(w, "%s %s param %d %s %s\n", kind, name, i, orDash(p.Name), renderType(p.Type))
	}
}

// renderModuleComments writes the comment stream: one record per symbol that
// carries a comment, in the same kind sequence and order-slice order as the
// structural stream, with each sub-symbol immediately after its parent. A symbol
// with an empty comment produces no record here, but still carries comment=- in
// the structural stream, so the absence is visible there.
//
// Keys match the structural stream's identity, so a different parser backend can
// be compared symbol by symbol.
//
// Determinism: as in renderModuleStructure, no map is ranged over.
func renderModuleComments(w io.Writer, m *Module) error {
	bw := bufio.NewWriter(w)

	fmt.Fprintln(bw, commentsFileHeader)

	for _, name := range m.functionOrder {
		if f, ok := m.functions[name]; ok {
			renderCommentRecord(bw, "function "+name, f.Comment)
		}
	}

	for _, name := range m.structOrder {
		s, ok := m.structs[name]
		if !ok {
			continue
		}

		renderCommentRecord(bw, "struct "+name, s.Comment)

		for _, f := range s.Fields {
			renderCommentRecord(bw, "struct "+name+"."+f.Name, f.Comment)
		}
	}

	for _, name := range m.enumOrder {
		e, ok := m.enums[name]
		if !ok {
			continue
		}

		renderCommentRecord(bw, "enum "+name, e.Comment)

		for _, c := range e.Constants {
			renderCommentRecord(bw, "enum "+name+"."+c.Name, c.Comment)
		}
	}

	for _, name := range m.callbackOrder {
		if f, ok := m.callbacks[name]; ok {
			renderCommentRecord(bw, "callback "+name, f.Comment)
		}
	}

	for _, name := range m.constantOrder {
		if c, ok := m.constants[name]; ok {
			renderCommentRecord(bw, "constant "+name, c.Comment)
		}
	}

	return bw.Flush()
}

// renderCommentRecord writes one comment record. The header carries an "@ "
// prefix and an explicit line count, so a comment body that itself begins "@ "
// cannot be misread as a header. The body is the post-processComment text
// verbatim, which is what a parser change has to reproduce.
func renderCommentRecord(w io.Writer, key, comment string) {
	if comment == "" {
		return
	}

	lines := strings.Split(comment, "\n")

	fmt.Fprintf(w, "@ %s lines=%d hash=%s\n", key, len(lines), commentHash(comment))

	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
}

// writeIRDump writes the structural and comment goldens into dir. It is the only
// function in this file that touches the filesystem; the renderers take an
// io.Writer so a test can render to a buffer instead.
func writeIRDump(dir string, m *Module) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create IR dump directory %v: %w", dir, err)
	}

	if err := writeDumpFile(filepath.Join(dir, "structure.txt"), func(w io.Writer) error {
		return renderModuleStructure(w, m)
	}); err != nil {
		return err
	}

	return writeDumpFile(filepath.Join(dir, "comments.txt"), func(w io.Writer) error {
		return renderModuleComments(w, m)
	})
}

// renderSkips writes the skip stream: a counts line, then one line per unique
// (Symbol, Reason) pair, sorted by symbol then reason. Deduplicating on the pair
// keeps the body readable while the counts line still moves when a marker is
// added or removed, so a count-only change is a visible one-line diff.
//
// The three counts measure different things and all three matter: markers is
// every recorded skip, symbols collapses the markers a single struct with
// several skipped fields produces, and pairs is the number of body lines.
//
// SkipEntry.Manual is excluded; sortedSkipPairs documents why.
//
// A nil collector yields the two header lines with zero counts and no body,
// which keeps the golden a valid file rather than an empty one.
func renderSkips(w io.Writer, c *SkipCollector) error {
	bw := bufio.NewWriter(w)

	pairs := sortedSkipPairs(c)

	fmt.Fprintln(bw, skipsFileHeader)
	fmt.Fprintf(bw, "# markers=%d symbols=%d pairs=%d\n", c.Total(), len(sortedUniqueSymbols(c)), len(pairs))

	for _, e := range pairs {
		fmt.Fprintf(bw, "%s: %s\n", e.Symbol, e.Reason)
	}

	return bw.Flush()
}

// writeSkipDump writes the skip golden into dir. It is separate from
// writeIRDump because the two are snapshotted at different points in a run: the
// Module at the parser boundary, the skip set after the codegen pass has added
// its own markers.
func writeSkipDump(dir string, c *SkipCollector) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create IR dump directory %v: %w", dir, err)
	}

	return writeDumpFile(filepath.Join(dir, "skips.txt"), func(w io.Writer) error {
		return renderSkips(w, c)
	})
}

// writeDumpFile writes one golden atomically, mirroring saveFormatted: render to
// a temporary sibling then rename over the destination, so a crash mid-write
// never leaves a half-written golden that a later diff would report as a parser
// change. saveFormatted itself is not reused because it is jennifer-specific.
func writeDumpFile(path string, render func(io.Writer) error) error {
	tmp := path + ".tmp"

	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	if err := render(f); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)

		return err
	}

	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)

		return err
	}

	return os.Rename(tmp, path)
}
