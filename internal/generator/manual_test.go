package main

import (
	"os"
	"path/filepath"
	"testing"
)

// writeFixture writes name with the given Go source into dir, failing the test
// on error.
func writeFixture(t *testing.T, dir, name, src string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(src), 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", name, err)
	}
}

// TestScanManualBindings verifies the detector collects exported top-level
// functions and methods from the root-package files while ignoring `*.gen.go`,
// `*_test.go`, unexported declarations, and subdirectories.
func TestScanManualBindings(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "display.go", `package ffmpeg

func AVDisplayRotationGet(m *int32) float64 { return 0 }
func (s *AVCodecContext) IntraMatrix() *int { return nil }
func unexportedHelper() {}
`)
	// Generated and test files must be ignored.
	writeFixture(t, dir, "functions.gen.go", `package ffmpeg

func AVGeneratedThing() {}
`)
	writeFixture(t, dir, "display_test.go", `package ffmpeg_test

func TestSomething() {}
`)
	// A subdirectory must not be visited.
	sub := filepath.Join(dir, "av")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeFixture(t, sub, "sub.go", `package av

func ShouldBeIgnored() {}
`)

	mb, err := scanManualBindings(dir)
	if err != nil {
		t.Fatalf("scanManualBindings: %v", err)
	}

	if !mb.hasFunc("AVDisplayRotationGet") {
		t.Error("expected AVDisplayRotationGet among detected functions")
	}
	if mb.hasFunc("unexportedHelper") {
		t.Error("unexported function should not be detected")
	}
	if mb.hasFunc("AVGeneratedThing") {
		t.Error("generated-file function should be ignored")
	}
	if mb.hasFunc("ShouldBeIgnored") {
		t.Error("subpackage function should be ignored")
	}
	if !mb.hasMethod("AVCodecContext", "IntraMatrix") {
		t.Error("expected IntraMatrix method on AVCodecContext")
	}
	if mb.hasMethod("AVCodecContext", "Missing") {
		t.Error("unexpected method reported present")
	}
}

// newTestGenerator builds a Generator backed by an empty struct table so
// convCamel runs without a libclang parse. The struct-collision suffix path is
// inert with no structs registered, matching the symbols under test.
func newTestGenerator() *Generator {
	return &Generator{
		input: &Module{structs: map[string]*Struct{}},
		skips: &SkipCollector{},
	}
}

// TestEnrichManualBindingsFunctions covers the function-skip path: a C symbol
// whose convCamel binding exists is annotated, one without a binding is left
// bare.
func TestEnrichManualBindingsFunctions(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, dir, "bindings.go", `package ffmpeg

func AVDisplayRotationGet(m *int32) float64 { return 0 }
func AVSizeMult(a, b uint) (uint, error) { return 0, nil }
`)

	g := newTestGenerator()
	g.skips.Record("av_display_rotation_get", "matrix array")
	g.skips.Record("av_size_mult", "out pointer")
	g.skips.Record("av_never_bound", "out pointer")

	g.enrichManualBindings(dir)

	got := map[string]string{}
	for _, e := range g.skips.Entries {
		got[e.Symbol] = e.Manual
	}

	if got["av_display_rotation_get"] != "AVDisplayRotationGet" {
		t.Errorf("av_display_rotation_get Manual = %q, want AVDisplayRotationGet", got["av_display_rotation_get"])
	}
	if got["av_size_mult"] != "AVSizeMult" {
		t.Errorf("av_size_mult Manual = %q, want AVSizeMult", got["av_size_mult"])
	}
	if got["av_never_bound"] != "" {
		t.Errorf("av_never_bound Manual = %q, want empty", got["av_never_bound"])
	}
}

// TestEnrichManualBindingsStructFields covers the "Type.field" path: a struct
// field whose convCamel accessor exists on the receiver is annotated, one
// without is left bare, and a matching method name on the wrong receiver does
// not match.
func TestEnrichManualBindingsStructFields(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, dir, "fields.go", `package ffmpeg

func (s *AVCodecContext) IntraMatrix() *int { return nil }
func (s *AVFrame) ExtendedData() *int { return nil }
`)

	g := newTestGenerator()
	g.skips.Record("AVCodecContext.intra_matrix", "quant matrix")
	g.skips.Record("AVFrame.extended_data", "extended data")
	g.skips.Record("AVCodecContext.unbound_field", "union")
	// IntraMatrix exists, but on AVCodecContext, not AVFrame.
	g.skips.Record("AVFrame.intra_matrix", "wrong receiver")

	g.enrichManualBindings(dir)

	got := map[string]string{}
	for _, e := range g.skips.Entries {
		got[e.Symbol] = e.Manual
	}

	if got["AVCodecContext.intra_matrix"] != "IntraMatrix" {
		t.Errorf("intra_matrix Manual = %q, want IntraMatrix", got["AVCodecContext.intra_matrix"])
	}
	if got["AVFrame.extended_data"] != "ExtendedData" {
		t.Errorf("extended_data Manual = %q, want ExtendedData", got["AVFrame.extended_data"])
	}
	if got["AVCodecContext.unbound_field"] != "" {
		t.Errorf("unbound_field Manual = %q, want empty", got["AVCodecContext.unbound_field"])
	}
	if got["AVFrame.intra_matrix"] != "" {
		t.Errorf("AVFrame.intra_matrix Manual = %q, want empty (wrong receiver)", got["AVFrame.intra_matrix"])
	}
}

// TestEnrichManualBindingsMissingDir confirms a scan failure degrades to leaving
// every Manual field empty rather than panicking.
func TestEnrichManualBindingsMissingDir(t *testing.T) {
	g := newTestGenerator()
	g.skips.Record("av_size_mult", "out pointer")

	g.enrichManualBindings(filepath.Join(t.TempDir(), "does-not-exist"))

	if g.skips.Entries[0].Manual != "" {
		t.Errorf("Manual = %q, want empty when scan fails", g.skips.Entries[0].Manual)
	}
}

// TestManualBindingBySymbol confirms only annotated entries appear in the
// printSkipSummary lookup map.
func TestManualBindingBySymbol(t *testing.T) {
	c := &SkipCollector{}
	c.Entries = []SkipEntry{
		{Symbol: "av_size_mult", Reason: "out pointer", Manual: "AVSizeMult"},
		{Symbol: "av_never_bound", Reason: "out pointer"},
	}

	m := manualBindingBySymbol(c)
	if m["av_size_mult"] != "AVSizeMult" {
		t.Errorf("map[av_size_mult] = %q, want AVSizeMult", m["av_size_mult"])
	}
	if _, ok := m["av_never_bound"]; ok {
		t.Error("unannotated symbol should not appear in the map")
	}
}
