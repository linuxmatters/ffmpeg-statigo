package main

import "testing"

// Compile-time assertion that the libclang backend satisfies the boundary. The
// build fails here if clangParser drifts from the HeaderParser method set.
var _ HeaderParser = clangParser{}

// TestClangParserImplementsHeaderParser pins the parser boundary: clangParser
// is assignable to HeaderParser, and the backend names itself for the verbose
// run log. An empty Version would leave the log line with no backend.
func TestClangParserImplementsHeaderParser(t *testing.T) {
	var parser HeaderParser = clangParser{}

	if got := parser.Version(); got == "" {
		t.Error("clangParser.Version() = \"\", want a non-empty backend name")
	}
}
