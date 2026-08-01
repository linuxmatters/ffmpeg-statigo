package main

import "testing"

// Compile-time assertion that the cc/v4 backend satisfies the boundary. The
// build fails here if ccParser drifts from the HeaderParser method set.
var _ HeaderParser = ccParser{}

// TestCCParserImplementsHeaderParser pins the parser boundary: ccParser is
// assignable to HeaderParser, and the backend names itself for the verbose run
// log. An empty Version would leave the log line with no backend.
func TestCCParserImplementsHeaderParser(t *testing.T) {
	var parser HeaderParser = ccParser{}

	if got := parser.Version(); got == "" {
		t.Error("ccParser.Version() = \"\", want a non-empty backend name")
	}
}
