package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"modernc.org/cc/v4"
)

// TestRecordFunctionDuplicateAborts pins the invariant-violation panic in
// recordFunction (duplicate function registration). The parser converts
// unhandled-form panics to skip-with-reason but deliberately keeps invariant
// panics fatal: a duplicate registration means the walk visited the same
// declaration twice, so the bindings are untrustworthy and aborting is correct.
//
// This is the cc/v4 replacement for the deleted libclang test of the same
// invariant. It fabricates the duplicate by pre-seeding the module's function
// map, then walks a real cc/v4 AST carrying a declaration of the same name. The
// recovered panic message must identify the duplicate, so a future refactor that
// turns the abort into a silent skip fails here instead of corrupting generated
// output.
func TestRecordFunctionDuplicateAborts(t *testing.T) {
	const fnName = "ffg_dup_test_fn"

	path := filepath.Join(t.TempDir(), "dup-fn.h")

	if err := os.WriteFile(path, fmt.Appendf(nil, "int %s(int x);\n", fnName), 0o644); err != nil {
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
			functions: map[string]*Function{
				fnName: {Name: fnName}, // pre-seeded duplicate
			},
			structs:   map[string]*Struct{},
			enums:     map[string]*Enum{},
			callbacks: map[string]*Function{},
			constants: map[string]*Constant{},
		},
		abs:   path,
		skips: &SkipCollector{},
		tags:  map[string]bool{},
		table: table,
	}

	// log.Panicln writes the panic message to the default logger before
	// panicking. Suppress that to keep test output clean; restore on exit.
	origOut := log.Writer()
	log.SetOutput(io.Discard)

	defer log.SetOutput(origOut)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("the walk did not panic on duplicate registration of %q", fnName)
		}

		msg := fmt.Sprintf("%v", r)
		if !strings.Contains(msg, "already exists") || !strings.Contains(msg, fnName) {
			t.Fatalf("panic message %q does not identify the duplicate function %q", msg, fnName)
		}
	}()

	w.walkAST(ast)
}

// TestUnnamedStructWithoutSpecifierAborts pins the second invariant panic in
// the walk: unnamedStructType is reached with no tagless specifier on it.
//
// Only a field declaration puts a specifier on the walk, so a tagless struct
// written directly as a function's return type or parameter reaches the type
// walk with none. C allows that shape and cc/v4 accepts it, so the abort is
// reachable in principle; it is unreachable on the pinned headers, because
// every tagless record under include/ is a struct field or a typedef.
//
// The abort is deliberate and this test exists to keep it that way. Recording a
// skip instead would return nil, which is the parser's encoding of C void, and
// classifyReturnShape maps a nil result to returnShapeVoid: the binding would
// compile and silently drop the returned struct. The parameter position happens
// to fail loudly (classifyArgShape sends a nil type to argShapeUnhandledSkip),
// but the return position does not, so a uniform "skip it" is not a safe swap
// for the panic.
//
// The message assertion is part of the pin. The indent breadcrumb has to carry
// the symbol and the slot, or an abort during an FFmpeg upgrade tells a
// maintainer nothing about which declaration to look at.
//
// The tagless union case is included to fix the boundary: parseUnionType
// expands a union in place and never needs a name, so a union return type
// parses rather than aborting.
func TestUnnamedStructWithoutSpecifierAborts(t *testing.T) {
	cfg, err := newCCConfig(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		t.Fatalf("cc/v4 configuration: %v", err)
	}

	dir := t.TempDir()

	tests := []struct {
		name string
		src  string
		// want is the breadcrumb the panic message must carry, or "" where the
		// header must parse without aborting.
		want string
	}{
		{
			name: "tagless struct return type",
			src:  "struct { int x; } ffg_anon_ret(void);\n",
			want: "[ffg_anon_ret][return]",
		},
		{
			name: "tagless struct parameter",
			src:  "void ffg_anon_param(struct { int y; } p);\n",
			want: "[ffg_anon_param][p]",
		},
		{
			name: "pointer to tagless struct return type",
			src:  "struct { int x; } *ffg_anon_pret(void);\n",
			want: "[ffg_anon_pret][return]",
		},
		{
			name: "tagless struct in a callback typedef",
			src:  "typedef int (*ffg_anon_cb)(struct { int y; } a);\n",
			want: "[ffg_anon_cb][a]",
		},
		{
			name: "tagless union return type parses",
			src:  "union { int a; float b; } ffg_anon_uret(void);\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTestHeader(t, dir, strings.ReplaceAll(tt.name, " ", "-"), tt.src)

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
					functions: map[string]*Function{},
					structs:   map[string]*Struct{},
					enums:     map[string]*Enum{},
					callbacks: map[string]*Function{},
					constants: map[string]*Constant{},
				},
				abs:   path,
				skips: &SkipCollector{},
				tags:  map[string]bool{},
				table: table,
			}

			// log.Panicln writes the message to the default logger before
			// panicking. Suppress that to keep test output clean.
			origOut := log.Writer()
			log.SetOutput(io.Discard)

			defer log.SetOutput(origOut)

			defer func() {
				r := recover()

				if tt.want == "" {
					if r != nil {
						t.Fatalf("the walk aborted with %v; this shape needs no synthetic name and must parse", r)
					}

					return
				}

				if r == nil {
					t.Fatalf("the walk did not abort on a tagless record with no specifier; a skip here emits a silently wrong binding")
				}

				msg := fmt.Sprintf("%v", r)
				if !strings.Contains(msg, "tagless record type with no specifier") {
					t.Fatalf("panic message %q is not the tagless-record abort", msg)
				}

				if !strings.Contains(msg, tt.want) {
					t.Fatalf("panic message %q does not carry the breadcrumb %q, so it does not say which declaration aborted the run", msg, tt.want)
				}
			}()

			w.walkAST(ast)
		})
	}
}

// writeTestHeader writes src to a uniquely named header in dir and returns the
// absolute path.
func writeTestHeader(t *testing.T, dir, name, src string) string {
	t.Helper()

	path := filepath.Join(dir, name+".h")
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}

	return path
}

// TestCCTranslateReturnsErrorsInsteadOfPanicking pins the recover in
// ccTranslate. cc.Translate reports a rejected translation unit two different
// ways: a returned error for anything its own diagnostics cover, and a bare
// panic through its todo() helper for a preprocessor branch it never
// implemented. The generator has to see one failure form, because Parse names
// the offending header in its report, so a stack trace out of a dependency is
// not an acceptable outcome.
//
// The two panic cases are real C that any compiler accepts:
//
//   - #ifdef __STDC_MB_MIGHT_NEQ_WC__ reaches cpp.go:1369. The name sits in
//     cc/v4's protectedMacros table but has no arm in the switch that follows,
//     and ccPredefined does not define it, so an ordinary feature test for a C99
//     standard macro falls straight through to todo().
//   - A function-like macro invocation whose closing parenthesis never arrives
//     reaches cpp.go:1497, where parseMacroArgs runs into end of file.
//
// Both are asserted to come back as an error naming a recovered panic, not just
// as some error, so the test cannot pass through the ordinary diagnostic route.
// If a future cc/v4 turns either branch into a diagnostic, this test fails and
// the fix is to find another todo() input, not to drop the case: the recover
// still has to be proven live.
func TestCCTranslateReturnsErrorsInsteadOfPanicking(t *testing.T) {
	cfg, err := newCCConfig(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		t.Fatalf("cc/v4 configuration: %v", err)
	}

	dir := t.TempDir()

	tests := []struct {
		name string
		src  string
		// recovered is true where cc/v4 fails by panicking, so the error must
		// come from the recover rather than from cc.Translate's return.
		recovered bool
	}{
		{
			name: "unparseable declaration",
			src:  "this is not a C declaration ((( ;\n",
		},
		{
			name: "unterminated struct at end of file",
			src:  "struct S { int a;\n",
		},
		{
			name: "missing include",
			src:  "#include <ffg_definitely_missing_header.h>\nint f(void);\n",
		},
		{
			name:      "protected macro with no parser arm",
			src:       "#ifdef __STDC_MB_MIGHT_NEQ_WC__\nint g(void);\n#endif\nint f(void);\n",
			recovered: true,
		},
		{
			name:      "unterminated macro invocation",
			src:       "#define FFG_M(x) x\nFFG_M(1\n",
			recovered: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTestHeader(t, dir, strings.ReplaceAll(tt.name, " ", "-"), tt.src)

			ast, err := translateNoPanic(t, cfg, path)
			if err == nil {
				t.Fatalf("ccTranslate(%s) returned no error; the header must be rejected for this case to test anything", filepath.Base(path))
			}

			if !tt.recovered {
				return
			}

			if ast != nil {
				t.Errorf("ccTranslate returned an AST alongside a recovered panic; callers must not see a partial translation unit")
			}

			if !strings.Contains(err.Error(), "panic translating") {
				t.Fatalf("error %q does not name a recovered panic; cc/v4 no longer panics on this input, so the recover in ccTranslate is no longer covered. Find another todo() input rather than deleting this case", err)
			}
		})
	}
}

// translateNoPanic calls ccTranslate and turns an escaping panic into a test
// failure rather than a crashed test binary. Nothing here recovers on
// ccTranslate's behalf: the assertion is that ccTranslate itself never lets a
// panic out.
func translateNoPanic(t *testing.T, cfg *cc.Config, path string) (ast *cc.AST, err error) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("a panic escaped ccTranslate(%s): %v", filepath.Base(path), r)
		}
	}()

	return ccTranslate(cfg, path)
}

// TestCCConfigTranslatesEveryTarget parses one real header for each of the four
// supported targets, so a configuration that cannot translate FFmpeg at all
// fails here with a named target and cause rather than only in the golden run.
//
// Every target runs on every host. That is new: on the host-header route a
// darwin target on a Linux host preprocessed glibc under Apple's ABI and failed
// its type check on _Float32, and the linux/arm64 target failed to parse glibc's
// aarch64 bits/math-vector.h. Neither host header is in the parse any more.
//
// av_frame_alloc is a long-lived public symbol, so this does not pin the FFmpeg
// version.
func TestCCConfigTranslatesEveryTarget(t *testing.T) {
	root := commentHeaderRoot(t)

	// AVLibPath is resolved at package init against the test's working
	// directory. newCCConfig puts it on both include lists, so it has to point
	// at the real header tree for the parse below to resolve FFmpeg's own
	// cross-includes.
	originalLibPath := AVLibPath
	AVLibPath = root

	t.Cleanup(func() { AVLibPath = originalLibPath })

	header := filepath.Join(root, "libavutil", "frame.h")

	for _, tgt := range crossTargets {
		t.Run(tgt.goos+"/"+tgt.goarch, func(t *testing.T) {
			cfg, err := newCCConfig(tgt.goos, tgt.goarch)
			if err != nil {
				t.Fatalf("newCCConfig: %v", err)
			}

			ast, err := ccTranslate(cfg, header)
			if err != nil {
				t.Fatalf("ccTranslate(libavutil/frame.h) under %s: %v", ccStd, err)
			}

			if !declaresFunction(ast, "av_frame_alloc") {
				t.Errorf("libavutil/frame.h translated but declares no function av_frame_alloc; the parse config reaches fewer symbols than it should")
			}
		})
	}
}

// predefinedMacro looks name up in a `#define` dump and returns its replacement
// text. Each line is "#define NAME REPLACEMENT" or, for a function-like macro,
// "#define NAME(ARGS) REPLACEMENT"; only object-like macros are wanted here.
func predefinedMacro(predefined, name string) (string, bool) {
	prefix := "#define " + name + " "

	for line := range strings.Lines(predefined) {
		line = strings.TrimRight(line, "\r\n")
		if line == "#define "+name {
			return "", true
		}

		if after, ok := strings.CutPrefix(line, prefix); ok {
			return after, true
		}
	}

	return "", false
}

// declaresFunction reports whether the translation unit's file scope holds a
// function declarator named name.
func declaresFunction(ast *cc.AST, name string) bool {
	for _, node := range ast.Scope.Nodes[name] {
		d, ok := node.(*cc.Declarator)
		if !ok {
			continue
		}

		if t := d.Type(); t != nil && t.Kind() == cc.Function {
			return true
		}
	}

	return false
}
