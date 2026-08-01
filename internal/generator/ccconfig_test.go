package main

import (
	"io/fs"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"

	"modernc.org/cc/v4"
)

// stubHeaders are the C library headers sysinclude/ replaces. The list is every
// name the FFmpeg headers in files reach for with #include <...>, and a name
// missing from it is a parse failure rather than a silent difference in output.
var stubHeaders = []string{
	"errno.h",
	"inttypes.h",
	"limits.h",
	"math.h",
	"stdarg.h",
	"stddef.h",
	"stdint.h",
	"stdio.h",
	"stdlib.h",
	"string.h",
	"time.h",
}

// TestCCConfigReadsNoHostHeaders is the hermeticity assertion. cc/v4 opens only
// the file handed to Translate and whatever it finds on these two search lists,
// so pinning the lists pins the whole set of files a parse can read.
//
// This is the property three CI legs lost when the config took its include roots
// from the host compiler: the darwin legs picked up the Apple SDK's math.h and
// the linux/arm64 leg picked up glibc's aarch64 bits/math-vector.h, which cc/v4
// cannot parse.
func TestCCConfigReadsNoHostHeaders(t *testing.T) {
	cfg, err := newCCConfig(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		t.Fatalf("newCCConfig: %v", err)
	}

	wantSys := []string{sysIncludeRoot, AVLibPath}
	if !slices.Equal(cfg.SysIncludePaths, wantSys) {
		t.Errorf("SysIncludePaths = %v, want %v; a path outside this list lets host headers into the parse", cfg.SysIncludePaths, wantSys)
	}

	// The leading "" is cc/v4's marker for the including file's own directory,
	// which is what a quoted include resolves against first.
	wantInc := []string{"", sysIncludeRoot, AVLibPath}
	if !slices.Equal(cfg.IncludePaths, wantInc) {
		t.Errorf("IncludePaths = %v, want %v", cfg.IncludePaths, wantInc)
	}

	if len(cfg.HostIncludePaths) != 0 || len(cfg.HostSysIncludePaths) != 0 {
		t.Errorf("host include discovery ran: HostIncludePaths=%v HostSysIncludePaths=%v", cfg.HostIncludePaths, cfg.HostSysIncludePaths)
	}

	if cfg.CC != "" {
		t.Errorf("cfg.CC = %q; the host compiler was executed", cfg.CC)
	}

	if cfg.Predefined != ccPredefined {
		t.Error("cfg.Predefined is not the committed literal; the macro set no longer comes from this repository")
	}

	if !cfg.Header {
		t.Error("cfg.Header is false; function bodies would be type checked")
	}
}

// TestStubHeadersResolve asserts every stub is reachable under the sentinel
// search path and that nothing else is. The second half is what lets the real
// FFmpeg headers through: cc/v4 offers each path to Config.FS and falls back to
// os.Open only when the FS refuses it.
func TestStubHeadersResolve(t *testing.T) {
	stubs := sysIncludeStubs{}

	for _, name := range stubHeaders {
		f, err := stubs.Open(filepath.Join(sysIncludeRoot, name))
		if err != nil {
			t.Errorf("stub <%s> is missing: %v", name, err)

			continue
		}

		f.Close()
	}

	for _, name := range []string{
		"/usr/include/math.h",
		filepath.Join(AVLibPath, "libavutil", "frame.h"),
		filepath.Join(sysIncludeRoot, "vulkan", "vulkan.h"),
	} {
		if f, err := stubs.Open(name); err == nil {
			f.Close()
			t.Errorf("the stub FS answered for %q; only sysinclude/ names may resolve through it", name)
		}
	}
}

// TestCCPredefinedPinsDialect pins the dialect and the compiler identity the
// FFmpeg headers branch on. It replaces a test that read the host compiler's
// `-dM -E -` dump, which is the thing this configuration no longer runs.
//
// Only these four names change what the pinned headers expand to.
// libavutil/attributes.h picks its whole av_* attribute set from __GNUC__,
// __GNUC_MINOR__ and __GNUC_STDC_INLINE__; libavutil/avassert.h and
// libavutil/attributes.h both test __STDC_VERSION__ against 202311L, and the
// 201112L below is what -std=gnu11 used to report. A drift in any of them
// changes which declarations the generator sees, so it fails here with a name
// rather than only as a wall of golden diff.
func TestCCPredefinedPinsDialect(t *testing.T) {
	want := map[string]string{
		"__STDC_VERSION__":     "201112L",
		"__GNUC__":             "15",
		"__GNUC_MINOR__":       "3",
		"__GNUC_STDC_INLINE__": "1",
	}

	for name, wantVal := range want {
		got, ok := predefinedMacro(ccPredefined, name)
		if !ok || got != wantVal {
			t.Errorf("%s = %q (defined: %v), want %q", name, got, ok, wantVal)
		}
	}

	// __STRICT_ANSI__ is what -std=c11 would have added. It hides declarations
	// the FFmpeg headers reach through, so gnu11 is named in ccStd and must stay
	// the dialect described here.
	if _, ok := predefinedMacro(ccPredefined, "__STRICT_ANSI__"); ok {
		t.Errorf("__STRICT_ANSI__ is defined; the parse describes a strict ISO dialect instead of %s", ccStd)
	}
}

// mathConstantsOf parses one header and returns the M_* constants it contributes
// to the constant table, sorted.
func mathConstantsOf(t *testing.T, cfg *cc.Config, header string) []string {
	t.Helper()

	ast, err := ccTranslate(cfg, header)
	if err != nil {
		t.Fatalf("translate %s: %v", header, err)
	}

	w := &ccWalk{
		mod:   &Module{constants: make(map[string]*Constant)},
		abs:   header,
		skips: &SkipCollector{},
	}
	w.macros(ast)

	var out []string

	for _, name := range w.mod.constantOrder {
		if strings.HasPrefix(name, "M_") {
			out = append(out, name)
		}
	}

	sort.Strings(out)

	return out
}

// TestStubMathHeaderPinsTheMConstants is the regression test for the macOS
// failure. libavutil/mathematics.h wraps all 26 M_* constants in an #ifndef, so
// each one the C library leaves undefined becomes a macro of an FFmpeg header,
// passes the location filter in ccWalk.macros and enters the bindings.
//
// glibc defines all 26 under _GNU_SOURCE. The Apple SDK defines only the 13
// without the f suffix, so the darwin CI legs emitted 13 constants, M_Ef first,
// that the Linux legs never produced, and the IR goldens failed on both.
//
// The second half reproduces that condition on any host, without a Mac: it
// swaps the stub math.h for the Apple SDK's macro set and asserts the
// f-suffixed names come back. Without it the first half would still pass if
// sysinclude/math.h were emptied, because the four survivors are the ones no C
// library defines.
func TestStubMathHeaderPinsTheMConstants(t *testing.T) {
	root := commentHeaderRoot(t)

	originalLibPath := AVLibPath
	AVLibPath = root

	t.Cleanup(func() { AVLibPath = originalLibPath })

	cfg, err := newCCConfig(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		t.Fatalf("newCCConfig: %v", err)
	}

	header := filepath.Join(root, "libavutil", "mathematics.h")

	// M_LOG2_10, M_PHI and their f variants are the four no C library defines, so
	// libavutil/mathematics.h owns them everywhere and they are the four the
	// bindings carry.
	want := []string{"M_LOG2_10", "M_LOG2_10f", "M_PHI", "M_PHIf"}

	got := mathConstantsOf(t, cfg, header)
	if !slices.Equal(got, want) {
		t.Errorf("M_* constants = %v, want %v; a C library macro set is reaching the constant table", got, want)
	}

	// Now the Apple SDK's shape: math.h present, f variants absent.
	cfg.FS = appleShapedMath{}

	partial := mathConstantsOf(t, cfg, header)
	if slices.Equal(partial, want) {
		t.Fatal("the Apple-shaped math.h changed nothing, so the first assertion proves nothing about the stub")
	}

	for _, name := range []string{"M_Ef", "M_PIf", "M_SQRT2f"} {
		if !slices.Contains(partial, name) {
			t.Errorf("under the Apple-shaped math.h the M_* set is %v, which does not carry %s; this is the state the macOS legs failed in", partial, name)
		}
	}
}

// appleMathH is <math.h> with the Apple SDK's M_* set: the 13 names without an
// f suffix and no others.
const appleMathH = `#define M_E 2.71828182845904523536028747135266250
#define M_LOG2E 1.44269504088896340735992468100189214
#define M_LOG10E 0.434294481903251827651128918916605082
#define M_LN2 0.693147180559945309417232121458176568
#define M_LN10 2.30258509299404568401799145468436421
#define M_PI 3.14159265358979323846264338327950288
#define M_PI_2 1.57079632679489661923132169163975144
#define M_PI_4 0.785398163397448309615660845819875721
#define M_1_PI 0.318309886183790671537767526745028724
#define M_2_PI 0.636619772367581343075535053490057448
#define M_2_SQRTPI 1.12837916709551257389615890312154517
#define M_SQRT2 1.41421356237309504880168872420969808
#define M_SQRT1_2 0.707106781186547524400844362104849039
`

// appleShapedMath serves the ordinary stubs with math.h swapped for the Apple
// SDK's macro set, so a Linux host can put the parse in the exact state the
// macOS CI legs failed in.
type appleShapedMath struct{}

func (appleShapedMath) Open(name string) (fs.File, error) {
	if filepath.Base(name) == "math.h" && strings.HasPrefix(filepath.ToSlash(name), sysIncludeRoot+"/") {
		return stringFile{Reader: strings.NewReader(appleMathH), name: "math.h", size: int64(len(appleMathH))}, nil
	}

	return sysIncludeStubs{}.Open(name)
}

// stringFile is the minimal fs.File over a string that cc/v4's scanner needs: it
// reads the bytes and stats the size.
type stringFile struct {
	*strings.Reader

	name string
	size int64
}

func (f stringFile) Stat() (fs.FileInfo, error) { return stringFileInfo(f), nil }
func (stringFile) Close() error                 { return nil }

type stringFileInfo stringFile

func (i stringFileInfo) Name() string     { return i.name }
func (i stringFileInfo) Size() int64      { return i.size }
func (stringFileInfo) Mode() fs.FileMode  { return 0o444 }
func (stringFileInfo) ModTime() time.Time { return time.Time{} }
func (stringFileInfo) IsDir() bool        { return false }
func (stringFileInfo) Sys() any           { return nil }
