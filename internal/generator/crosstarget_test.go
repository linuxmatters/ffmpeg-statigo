package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// crossTargets are the four platforms WW-141 claims produce identical generator
// output. Only the pair matching the host operating system can be measured on
// any one machine; the test below records why.
var crossTargets = []struct{ goos, goarch string }{
	{"linux", "amd64"},
	{"linux", "arm64"},
	{"darwin", "amd64"},
	{"darwin", "arm64"},
}

// crossTargetsToRun returns the entries of crossTargets this host can measure:
// the host's own operating system, minus the host pair itself, which
// TestIRGoldensMatchFreshRun already runs against the same goldens.
func crossTargetsToRun() []struct{ goos, goarch string } {
	var out []struct{ goos, goarch string }

	for _, tgt := range crossTargets {
		if tgt.goos == runtime.GOOS && tgt.goarch != runtime.GOARCH {
			out = append(out, tgt)
		}
	}

	return out
}

// TestGeneratorIRIsIdenticalAcrossTargets measures the cross-platform
// determinism claim: the generator must produce the same output for
// linux/amd64, linux/arm64, darwin/amd64 and darwin/arm64.
//
// The measurement lives in a test because the binary cannot be pointed at a
// target. ccParser.Parse calls newCCConfig(runtime.GOOS, runtime.GOARCH), both
// compile-time constants, and the flag set carries no target flag, so
// `GOOS=darwin GOARCH=arm64 go run ./internal/generator` only cross-compiles the
// generator and then fails to execute it. newCCConfig does take the target
// explicitly, so the test calls it per target and renders the three IR streams
// dump.go writes.
//
// The comparison is against the committed goldens, which are the host target's
// own output: TestIRGoldensMatchFreshRun pins host output to them, so a
// cross-target run that also matches them puts both targets on the same bytes
// without parsing the host target twice.
//
// What it proves: the parsed Module and the skip set do not vary with the target
// ABI. That is the risk worth guarding. cc.NewConfig(goos, goarch, ...) picks
// cc/v4's ABI table, and the four pairs differ in it - char is unsigned on
// linux/arm64, and long double is 8 bytes rather than 16 on darwin/arm64 - yet
// the generator never asks for a size or an offset. It emits an enum constant
// and a macro alike as `= C.NAME` (generator.go), so cgo resolves every number
// at build time. This test keeps that true.
//
// What it does NOT prove: anything about the other operating system. The include
// paths always come from the host compiler, so a darwin target on a Linux host
// would preprocess glibc headers under Apple's ABI. That is not a macOS parse
// and it does not even translate: cc/v4's darwin ABI tables carry no _Float32,
// _Float64 or _Float64x kind, while glibc's math.h declares prototypes for all
// three, so every header fails its type check with "_Float32 not supported on
// darwin/arm64". The test therefore runs only targets sharing the host's
// operating system. A Linux run leaves both darwin targets untested and a macOS
// run leaves both linux targets untested; covering all four needs a run on each
// operating system.
//
// The comparison is live, not vacuous. The IR is not ABI-proof everywhere:
// FF_PAD_STRUCTURE sizes AVBPrint.reserved_padding from sizeof expressions, so
// pointing the target at linux/386 fails this test on exactly that line,
// char[1004] against char[1000]. All four claimed targets are LP64, which is why
// they agree, and arch_guard.go rejects a 32-bit build anyway.
//
// It is not parallel: it mutates AVLibPath and log output, and t.Chdir is
// process-wide.
func TestGeneratorIRIsIdenticalAcrossTargets(t *testing.T) {
	targets := crossTargetsToRun()
	if len(targets) == 0 {
		t.Skipf("no second target for host %s/%s among the four claimed targets", runtime.GOOS, runtime.GOARCH)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	// go test runs with the package directory as the working directory, so the
	// repo root is two levels up.
	repoRoot, err := filepath.Abs(filepath.Join(cwd, "..", ".."))
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}

	includeDir := filepath.Join(repoRoot, "include")
	if _, err := os.Stat(includeDir); err != nil {
		t.Skipf("FFmpeg headers absent at %s: %v\ninitialise them with: git submodule update --init --recursive", includeDir, err)
	}

	// Read the committed goldens before the chdir, while the relative package
	// paths still resolve.
	golden := targetIR{
		structure: readGolden(t, repoRoot, "structure.txt"),
		comments:  readGolden(t, repoRoot, "comments.txt"),
		skips:     readGolden(t, repoRoot, "skips.txt"),
	}

	// AVLibPath is resolved at package init against the test's working
	// directory, so it points at internal/generator/include, which does not
	// exist. Point it at the real header tree, as TestIRGoldensMatchFreshRun
	// does.
	originalLibPath := AVLibPath
	AVLibPath = includeDir

	t.Cleanup(func() { AVLibPath = originalLibPath })

	// A real parse logs one line per type it visits and would drown go test -v.
	originalLogOutput := log.Writer()

	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(originalLogOutput) })

	// Gen writes the five *.gen.go files relative to the working directory, so
	// they land in a temporary directory and are discarded. A test run cannot
	// touch the committed bindings.
	t.Chdir(t.TempDir())

	for _, tgt := range targets {
		name := tgt.goos + "/" + tgt.goarch
		got := renderTargetIR(t, tgt.goos, tgt.goarch)

		compareStream(t, "structure.txt", name, golden.structure, got.structure)
		compareStream(t, "comments.txt", name, golden.comments, got.comments)
		compareStream(t, "skips.txt", name, golden.skips, got.skips)
	}
}

// readGolden reads one committed golden through its repo-root-relative path.
func readGolden(t *testing.T, repoRoot, name string) string {
	t.Helper()

	path := filepath.Join(repoRoot, irDumpDir, name)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v\nregenerate with: %s", path, err, regenerateCmd)
	}

	return string(data)
}

// targetIR holds the three rendered IR streams of one target, in the form the
// committed goldens carry.
type targetIR struct {
	structure string
	comments  string
	skips     string
}

// renderTargetIR runs the whole generator pass for one target and renders the
// three streams. The order mirrors run(): the Module is rendered at the parser
// boundary, before applyManualFixups, and the skip set after Gen, once the
// codegen layer has added its own markers.
func renderTargetIR(t *testing.T, goos, goarch string) targetIR {
	t.Helper()

	skips := &SkipCollector{}
	m := parseForTarget(t, goos, goarch, skips)

	structure, comments := renderBoth(t, m)

	if err := applyManualFixups(m); err != nil {
		t.Fatalf("%s/%s: applyManualFixups: %v", goos, goarch, err)
	}

	return targetIR{structure: structure, comments: comments, skips: renderSkipsString(t, Gen(m, skips))}
}

// parseForTarget parses every pinned header for one target.
//
// ccParser.Parse cannot be called here: it hardcodes runtime.GOOS and
// runtime.GOARCH, so its loop is repeated with a config built for the requested
// pair. That loop is the only part of Parse this file duplicates, and a drift
// between the two shows up as a failed comparison against the goldens, which
// TestIRGoldensMatchFreshRun holds to Parse's own output.
func parseForTarget(t *testing.T, goos, goarch string, skips *SkipCollector) *Module {
	t.Helper()

	cfg, err := newCCConfig(goos, goarch)
	if err != nil {
		t.Fatalf("newCCConfig(%q, %q): %v", goos, goarch, err)
	}

	w := &ccWalk{
		mod: &Module{
			functions: make(map[string]*Function),
			structs:   make(map[string]*Struct),
			enums:     make(map[string]*Enum),
			callbacks: make(map[string]*Function),
			constants: make(map[string]*Constant),
		},
		skips: skips,
	}

	for _, file := range files {
		w.abs = filepath.Join(AVLibPath, file)
		w.tags = make(map[string]bool)

		w.table, err = BuildCommentTable(w.abs)
		if err != nil {
			t.Fatalf("%s/%s: comment table for %s: %v", goos, goarch, w.abs, err)
		}

		ast, err := ccTranslate(cfg, w.abs)
		if err != nil {
			t.Fatalf("%s/%s: translate %s: %v", goos, goarch, w.abs, err)
		}

		w.walkAST(ast)
	}

	return w.mod
}

// compareStream fails with the first differing line named, in the form the
// golden comparison in dump_test.go uses.
func compareStream(t *testing.T, stream, target, want, got string) {
	t.Helper()

	if want == got {
		return
	}

	diff := diffLines(want, got)
	if len(diff) == 0 {
		t.Errorf("%s for %s differs from the committed golden in trailing bytes only (%d vs %d bytes)",
			stream, target, len(want), len(got))

		return
	}

	d := diff[0]
	t.Errorf("%s for %s differs from the committed golden at line %d (%d differing lines):\n golden %s\n %s %s",
		stream, target, d.num, len(diff), truncate(d.want), target, truncate(d.got))
}
