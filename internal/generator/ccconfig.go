package main

// Hermetic parse configuration for the modernc.org/cc/v4 backend.
//
// The generator does not read the host's system headers and does not run the
// host compiler. Every macro and every declaration the parse sees comes either
// from the committed FFmpeg headers under include/ or from the stub headers in
// sysinclude/, which are embedded into the binary.
//
// It started out on the host route: newCCConfig called cc.NewConfig, which runs
// `<cc> -dM -E -` and `<cc> -v -E -` and hands back whatever the host's libc
// ships. That made the emitted bindings a function of the build machine, and it
// broke two ways at once:
//
//   - libavutil/mathematics.h wraps every M_* constant in an #ifndef. glibc
//     defines all 26 of them under _GNU_SOURCE and the Apple SDK defines only
//     the 13 without the f suffix, so on macOS the 13 f-suffixed fallbacks
//     became macros of an FFmpeg header, passed the location filter, and added
//     13 constants the Linux run never produced.
//   - glibc's aarch64 bits/math-vector.h declares its vector-math prototypes
//     with a construct cc/v4 cannot parse, so every header failed to translate
//     on linux/arm64. That header has no amd64 counterpart, which is why an
//     amd64 dev host never met it.
//
// Both are the same defect: the output depended on headers this repository does
// not own. Stubbing the eleven C library headers the FFmpeg headers include
// removes the dependency rather than papering over either symptom.
//
// What this does not solve: the FFmpeg headers still have to translate, so a
// construct cc/v4 cannot parse in an FFmpeg header would still fail, and the
// stubs have to keep up with any new C library name a future FFmpeg release
// reaches for. A missing one is a loud parse failure naming the header, not a
// silent difference in output.

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"modernc.org/cc/v4"
)

// sysIncludeFS holds the stub C library headers. They are embedded rather than
// read from the checkout so the parse does not depend on the working directory,
// which under `go test` is the package directory and not the repo root.
//
//go:embed sysinclude
var sysIncludeFS embed.FS

// sysIncludeRoot is the search path the stub headers answer on. It is a
// sentinel, not a directory: cc/v4 joins it with the included name and offers
// the result to Config.FS before falling back to os.Open, so nothing on disk
// has to exist at this path, and a stub position that ever leaked into output
// would be obviously synthetic.
const sysIncludeRoot = "ffmpeg-statigo-sysinclude"

// ccStd names the dialect the predefined macros below describe. It is no longer
// a compiler flag; ccPredefined states the same thing directly.
const ccStd = "gnu11"

// ccPredefined is the predefined macro set the parse runs under, in the form
// cc/v4 expects: the text a compiler's `-dM -E -` would print.
//
// It is a fixed literal and does not vary with the target. Only four of these
// names change what the FFmpeg headers expand to. libavutil/attributes.h picks
// its av_* attribute macros from __GNUC__, __GNUC_MINOR__ and
// __GNUC_STDC_INLINE__, and libavutil/avassert.h and libavutil/attributes.h
// both test __STDC_VERSION__ against 202311L. None of the four are target
// specific, and the FFmpeg headers this generator parses branch on no operating
// system or architecture macro at all: __APPLE__ appears only in
// libavutil/hwcontext_opencl.h and _WIN32 only in
// libavutil/hwcontext_vulkan.h, neither of which is in files.
//
// The GCC version claims 15.3, which is the dev shell's compiler and the one
// that produced the committed IR goldens.
//
// The __*_TYPE__ names are required by cc.Builtin, which ccTranslate passes to
// every translation unit, and by sysinclude/stddef.h. The widths are LP64;
// arch_guard.go rejects every other build.
const ccPredefined = `#define __STDC__ 1
#define __STDC_HOSTED__ 1
#define __STDC_VERSION__ 201112L
#define __GNUC__ 15
#define __GNUC_MINOR__ 3
#define __GNUC_PATCHLEVEL__ 0
#define __GNUC_STDC_INLINE__ 1
#define __CHAR_BIT__ 8
#define __SIZEOF_SHORT__ 2
#define __SIZEOF_INT__ 4
#define __SIZEOF_LONG__ 8
#define __SIZEOF_LONG_LONG__ 8
#define __SIZEOF_POINTER__ 8
#define __SIZEOF_FLOAT__ 4
#define __SIZEOF_DOUBLE__ 8
#define __SIZEOF_LONG_DOUBLE__ 16
#define __SIZEOF_SIZE_T__ 8
#define __SIZEOF_PTRDIFF_T__ 8
#define __SIZEOF_WCHAR_T__ 4
#define __SIZE_TYPE__ long unsigned int
#define __PTRDIFF_TYPE__ long int
#define __WCHAR_TYPE__ int
#define __INTPTR_TYPE__ long int
#define __UINTPTR_TYPE__ long unsigned int
#define __INT8_TYPE__ signed char
#define __INT16_TYPE__ short int
#define __INT32_TYPE__ int
#define __INT64_TYPE__ long int
#define __UINT8_TYPE__ unsigned char
#define __UINT16_TYPE__ short unsigned int
#define __UINT32_TYPE__ unsigned int
#define __UINT64_TYPE__ long unsigned int
#define __ORDER_LITTLE_ENDIAN__ 1234
#define __ORDER_BIG_ENDIAN__ 4321
#define __ORDER_PDP_ENDIAN__ 3412
#define __BYTE_ORDER__ __ORDER_LITTLE_ENDIAN__
#define __LP64__ 1
#define _LP64 1
`

// newCCConfig builds the parse configuration for one target.
//
// goos and goarch select the cc/v4 ABI table and nothing else. The four targets
// this module supports differ in it - char is unsigned on linux/arm64, and long
// double is 8 bytes rather than 16 on darwin/arm64 - but the generator asks for
// no size and no offset, so TestGeneratorIRIsIdenticalAcrossTargets holds all
// four to the same bytes.
func newCCConfig(goos, goarch string) (*cc.Config, error) {
	abi, err := cc.NewABI(goos, goarch)
	if err != nil {
		return nil, fmt.Errorf("cc/v4 ABI for %s/%s: %w", goos, goarch, err)
	}

	// The FFmpeg headers include each other in both forms, "libavutil/x.h" and
	// <libavutil/x.h>, so the header root goes on both lists. The leading ""
	// on IncludePaths is cc/v4's marker for "the directory of the including
	// file", which is what a quoted include resolves against first.
	paths := []string{sysIncludeRoot, AVLibPath}

	return &cc.Config{
		ABI:             abi,
		FS:              sysIncludeStubs{},
		IncludePaths:    append([]string{""}, paths...),
		SysIncludePaths: paths,
		Predefined:      ccPredefined,

		// Header mode drops type checking of function bodies. The FFmpeg headers
		// carry static inline definitions the generator never emits.
		Header: true,
	}, nil
}

// sysIncludeStubs serves the embedded stub headers and refuses everything else.
//
// cc/v4 offers every path it resolves to Config.FS first and falls back to
// os.Open on an error, so refusing a path is how the real FFmpeg headers get
// read from disk. It is also how a search path that is not sysIncludeRoot moves
// on to the next candidate.
type sysIncludeStubs struct{}

// Open returns the stub header at name, which cc/v4 has already joined onto one
// of the configured search paths.
func (sysIncludeStubs) Open(name string) (fs.File, error) {
	rel, ok := strings.CutPrefix(filepath.ToSlash(name), sysIncludeRoot+"/")
	if !ok {
		return nil, fs.ErrNotExist
	}

	return sysIncludeFS.Open("sysinclude/" + rel)
}
