package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// manualBindings indexes the hand-written `package ffmpeg` bindings that sit
// alongside the generated `*.gen.go` files. The generator skips some C symbols
// it cannot emit (fixed-size array params, single out-params, variadics,
// function-pointer params, unions); many of those still gain a hand-written
// wrapper in a topic file in the repo root. This index lets the skip summary
// distinguish "skipped and genuinely missing" from "skipped but manually
// covered".
//
// funcs holds exported top-level function names. methods maps a receiver type
// name to the set of exported method names declared on it, so a struct-field
// skip recorded as "Type.field" can be matched against the type's accessors.
type manualBindings struct {
	funcs   map[string]bool
	methods map[string]map[string]bool
}

// scanManualBindings parses every root-package `.go` file in dir, excluding the
// generated `*.gen.go` files and `*_test.go` files, and collects the exported
// top-level functions and methods they declare. Subpackages are not visited:
// only files directly in dir are read, so the generator's siblings (the
// hand-written root-package bindings) are indexed and the `av/` and `internal/`
// trees are ignored.
//
// Parse errors on individual files are skipped rather than fatal: a single
// unparsable sibling must not abort the skip-summary enrichment.
func scanManualBindings(dir string) (*manualBindings, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	mb := &manualBindings{
		funcs:   make(map[string]bool),
		methods: make(map[string]map[string]bool),
	}

	fset := token.NewFileSet()

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".go") ||
			strings.HasSuffix(name, ".gen.go") ||
			strings.HasSuffix(name, "_test.go") {
			continue
		}

		file, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, 0)
		if err != nil {
			continue
		}

		mb.collect(file)
	}

	return mb, nil
}

// collect records the exported top-level functions and methods declared in file.
func (mb *manualBindings) collect(file *ast.File) {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil || !fn.Name.IsExported() {
			continue
		}

		if fn.Recv == nil || len(fn.Recv.List) == 0 {
			mb.funcs[fn.Name.Name] = true
			continue
		}

		recv := receiverTypeName(fn.Recv.List[0].Type)
		if recv == "" {
			continue
		}

		set := mb.methods[recv]
		if set == nil {
			set = make(map[string]bool)
			mb.methods[recv] = set
		}
		set[fn.Name.Name] = true
	}
}

// receiverTypeName extracts the bare type name from a method receiver, peeling
// a leading pointer (*T -> T). Generic receivers and other shapes return "".
func receiverTypeName(expr ast.Expr) string {
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// hasFunc reports whether an exported top-level function with the given name was
// found among the hand-written bindings.
func (mb *manualBindings) hasFunc(name string) bool {
	return mb != nil && mb.funcs[name]
}

// hasMethod reports whether the given exported method was found on receiver.
func (mb *manualBindings) hasMethod(receiver, name string) bool {
	if mb == nil {
		return false
	}
	set := mb.methods[receiver]
	return set != nil && set[name]
}
