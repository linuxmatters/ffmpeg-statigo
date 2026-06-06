package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/Newbluecake/bootstrap/clang"
)

// TestParseFunctionDuplicateAborts pins the invariant-violation panic in
// parseFunction (duplicate function registration). The parser converts shape
// panics to skip-with-reason but deliberately keeps invariant panics fatal:
// a duplicate registration means libclang misreported or the parser visited
// the same FunctionDecl twice, so the bindings are untrustworthy and aborting
// is correct. This test fabricates the duplicate by pre-seeding the module's
// function map, then runs parseFunction over a real libclang FunctionDecl
// cursor with the same name. The recovered panic message must identify the
// duplicate so a future refactor that turns the abort into a silent skip
// fails this test instead of corrupting generated output.
func TestParseFunctionDuplicateAborts(t *testing.T) {
	tmp, err := os.CreateTemp("", "dup-fn-*.h")
	if err != nil {
		t.Fatalf("create temp header: %v", err)
	}
	defer os.Remove(tmp.Name())

	const fnName = "ffg_dup_test_fn"
	if _, err := fmt.Fprintf(tmp, "int %s(int x);\n", fnName); err != nil {
		t.Fatalf("write temp header: %v", err)
	}
	if err := tmp.Close(); err != nil {
		t.Fatalf("close temp header: %v", err)
	}

	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	// Minimal args: the synthetic header references no system types, so no
	// include discovery (and therefore no gcc dependency) is required.
	tu := idx.ParseTranslationUnit(tmp.Name(), nil, nil, 0)
	defer tu.Dispose()

	if !tu.IsValid() {
		t.Fatalf("invalid translation unit for synthetic header")
	}

	var fnCursor clang.Cursor
	var found bool
	tu.TranslationUnitCursor().Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.Kind() == clang.Cursor_FunctionDecl && cursor.Spelling() == fnName {
			fnCursor = cursor
			found = true
			return clang.ChildVisit_Break
		}
		return clang.ChildVisit_Continue
	})

	if !found {
		t.Fatalf("FunctionDecl %q not found in synthetic translation unit", fnName)
	}

	p := &Parser{
		mod: &Module{
			functions: map[string]*Function{
				fnName: {Name: fnName}, // pre-seeded duplicate
			},
			structs:   map[string]*Struct{},
			enums:     map[string]*Enum{},
			callbacks: map[string]*Function{},
			constants: map[string]*Constant{},
		},
	}

	// log.Panicln writes the panic message to the default logger before
	// panicking. Suppress that to keep test output clean; restore on exit.
	origOut := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(origOut)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("parseFunction did not panic on duplicate registration")
		}
		msg := fmt.Sprintf("%v", r)
		if !strings.Contains(msg, "already exists") || !strings.Contains(msg, fnName) {
			t.Fatalf("panic message %q does not identify the duplicate function %q",
				msg, fnName)
		}
	}()

	p.parseFunction("", fnCursor)
}
