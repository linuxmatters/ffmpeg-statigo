package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"modernc.org/cc/v4"
)

// ctypeProbeHeader is a synthetic header carrying one field of every C type
// form the IR golden holds a ctype for. It is self-contained: cc.Builtin covers
// the compiler builtins and nothing here includes a system header, so the test
// does not depend on the FFmpeg submodule.
//
// The two unnamed records sit at the end because their spelling carries a
// line:col, which ctypeUnnamedAt computes from this same source text.
const ctypeProbeHeader = `typedef unsigned char uint8_t;
typedef signed short int16_t;
typedef struct AVRational { int num; int den; } AVRational;
typedef struct AVStream AVStream;
struct AVCodec;
struct AVCodecContext;
enum AVPixelFormat { AV_PIX_FMT_NONE };

struct Probe {
    int bare;
    unsigned int two_word;
    uint8_t sugared;
    AVRational by_value;
    struct AVRational elaborated;
    enum AVPixelFormat pix_fmt;
    char *ptr;
    struct AVCodec *tag_ptr;
    AVStream **ptr_ptr;
    const char *const *mime_types;
    int fixed[3];
    int16_t nested[3][2];
    uint8_t *plane_ptrs[8];
    const uint8_t (*audio_shuffle)[9];
    char *(*item_name)(void *);
    void (*idct)(int16_t *);
    enum AVPixelFormat (*get_format)(struct AVCodecContext *, const enum AVPixelFormat *);
    int (*execute2)(struct AVCodecContext *, int (*)(struct AVCodecContext *, void *, int, int), void *, int *, int);
    int (*printf_like)(const char *, ...);
    int (*const callback)(int);
    int (*const table)[4];
    int *const *(*levels)[4];
    union { int i; double dbl; } default_val;
    struct { unsigned idx; int horizontal; } *offsets;
};
`

// ctypeUnnamedAt returns the clang spelling of the tagless record whose keyword
// is the nth occurrence of kw in the probe header, with the line and column
// counted from the source text rather than from the parser, so the expectation
// is independent of what the printer produces.
func ctypeUnnamedAt(t *testing.T, path, kw string, nth int) string {
	t.Helper()

	needle := "    " + kw + " {"

	off := -1
	for i := 0; i <= nth; i++ {
		j := strings.Index(ctypeProbeHeader[off+1:], needle)
		if j < 0 {
			t.Fatalf("occurrence %d of %q not found in the probe header", i, needle)
		}

		off += 1 + j
	}

	line := strings.Count(ctypeProbeHeader[:off], "\n") + 1
	col := off - (strings.LastIndexByte(ctypeProbeHeader[:off], '\n') + 1) + 1

	// The keyword sits after the four-space indent the needle matched.
	col += len("    ")

	return fmt.Sprintf("%s (unnamed %s at %s:%d:%d)", kw, kw, path, line, col)
}

// TestCCCTypeName pins one string of every ctype form the golden holds: bare,
// two-word, typedef-sugared, elaborated, pointer, double pointer, pointer-level
// const, fixed and multidimensional arrays, incomplete array, pointer to array,
// function pointers from one to five parameters including a nested one, a const
// written on a '*' inside parentheses, and the unnamed struct and union
// spellings that carry a path and a line:col.
func TestCCCTypeName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "probe.h")

	if err := os.WriteFile(path, []byte(ctypeProbeHeader), 0o600); err != nil {
		t.Fatalf("write probe header: %v", err)
	}

	got := ctypeProbeFields(t, path)

	tests := []struct {
		field string
		want  string
	}{
		{field: "bare", want: "int"},
		{field: "two_word", want: "unsigned int"},
		{field: "sugared", want: "uint8_t"},
		{field: "by_value", want: "AVRational"},
		{field: "elaborated", want: "struct AVRational"},
		{field: "pix_fmt", want: "enum AVPixelFormat"},
		{field: "ptr", want: "char *"},
		{field: "tag_ptr", want: "struct AVCodec *"},
		{field: "ptr_ptr", want: "AVStream **"},
		{field: "mime_types", want: "char *const *"},
		{field: "fixed", want: "int[3]"},
		{field: "nested", want: "int16_t[3][2]"},
		{field: "plane_ptrs", want: "uint8_t *[8]"},
		{field: "audio_shuffle", want: "uint8_t (*)[9]"},
		{field: "item_name", want: "char *(*)(void *)"},
		{field: "idct", want: "void (*)(int16_t *)"},
		{
			field: "get_format",
			want:  "enum AVPixelFormat (*)(struct AVCodecContext *, const enum AVPixelFormat *)",
		},
		{
			field: "execute2",
			want: "int (*)(struct AVCodecContext *, " +
				"int (*)(struct AVCodecContext *, void *, int, int), void *, int *, int)",
		},
		{field: "printf_like", want: "int (*)(const char *, ...)"},
		{field: "callback", want: "int (*const)(int)"},
		{field: "table", want: "int (*const)[4]"},
		{field: "levels", want: "int *const *(*)[4]"},
		{field: "default_val", want: ctypeUnnamedAt(t, path, "union", 0)},
		{field: "offsets", want: ctypeUnnamedAt(t, path, "struct", 0) + " *"},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			if got[tt.field] != tt.want {
				t.Errorf("ccCTypeName(%s) = %q, want %q", tt.field, got[tt.field], tt.want)
			}
		})
	}
}

// TestCCCTypeNameNilDeclarator pins the empty result. The golden's three empty
// ctype values belong to the synthetic UnnamedStruct_avformat_986_5 fields,
// which registerUnnamedStruct creates without a declarator and without a
// spelling, so nothing must invent one for them.
func TestCCCTypeNameNilDeclarator(t *testing.T) {
	if got := ccCTypeName(nil, nil); got != "" {
		t.Errorf("ccCTypeName(nil, nil) = %q, want %q", got, "")
	}
}

// ctypeProbeFields translates the probe header and renders every field of
// struct Probe, keyed by field name.
func ctypeProbeFields(t *testing.T, path string) map[string]string {
	t.Helper()

	cfg, err := newCCConfig(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		t.Fatalf("newCCConfig: %v", err)
	}

	ast, err := ccTranslate(cfg, path)
	if err != nil {
		t.Fatalf("translate probe header: %v", err)
	}

	out := map[string]string{}

	for tu := ast.TranslationUnit; tu != nil; tu = tu.TranslationUnit {
		ed := tu.ExternalDeclaration
		if ed == nil || ed.Case != cc.ExternalDeclarationDecl || ed.Declaration == nil {
			continue
		}

		for ds := ed.Declaration.DeclarationSpecifiers; ds != nil; ds = ds.DeclarationSpecifiers {
			if ds.Case != cc.DeclarationSpecifiersTypeSpec || ds.TypeSpecifier == nil {
				continue
			}

			sus := ds.TypeSpecifier.StructOrUnionSpecifier
			if sus == nil || sus.Case != cc.StructOrUnionSpecifierDef || sus.Token.SrcStr() != "Probe" {
				continue
			}

			for sdl := sus.StructDeclarationList; sdl != nil; sdl = sdl.StructDeclarationList {
				sd := sdl.StructDeclaration
				if sd == nil || sd.Case != cc.StructDeclarationDecl {
					continue
				}

				anon := ccAnonSpecifier(sd)

				for l := sd.StructDeclaratorList; l != nil; l = l.StructDeclaratorList {
					if l.StructDeclarator == nil || l.StructDeclarator.Declarator == nil {
						continue
					}

					dcl := l.StructDeclarator.Declarator
					if name := dcl.Name(); name != "" {
						out[name] = ccCTypeName(dcl, anon)
					}
				}
			}
		}
	}

	if len(out) == 0 {
		t.Fatal("struct Probe produced no fields")
	}

	return out
}
