package main

// Host C preprocessor discovery for the modernc.org/cc/v4 parser backend.
//
// The generator has always taken its include roots from the host compiler. The
// libclang backend ran `gcc -E -x c -v /dev/null` itself and scraped the
// "#include <...> search starts here:" block into -I flags; cc.NewConfig runs
// the same probe, so nothing here has to.
//
// WW-141 phase 1.3 planned to lift hostCppConfig from xlab/c-for-go
// (parser/parser.go, MIT). It is deliberately not lifted. cc.NewConfig already
// does that work and more: cc/v4's own newConfig (cc.go) runs `<cc> -dM -E -`
// for the predefined macros and `<cc> -v -E -` for the search list, both under
// LC_ALL=C, then parses the "..." block into HostIncludePaths and the <...>
// block into HostSysIncludePaths. The parsing loop is the same code, line for
// line. c-for-go calls hostCppConfig only so a separately configured CPP can be
// prepended to the lists cc.NewConfig produces anyway. Lifting it here would
// add a second copy of a probe cc/v4 runs on every NewConfig call.
//
// The host compiler route is kept on purpose. A hermetic header set is an
// explicit non-goal of this port.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"modernc.org/cc/v4"
)

// ccStd is the dialect the parse requires. The default config fails on
// stddef.h "undefined: nullptr", and -std=c11 or -std=c17 reach several hundred
// fewer functions than -std=gnu11.
const ccStd = "-std=gnu11"

// ccDefines are the feature-test macros the parse requires.
// Without them glibc's math.h leaves M_Ef and its siblings undefined, FFmpeg's
// own #ifndef fallbacks in libavutil/mathematics.h fire, and 11 extra M_*f
// macros land in the constant table.
var ccDefines = []string{"-D_GNU_SOURCE", "-D_DEFAULT_SOURCE", "-D__STDC_CONSTANT_MACROS"}

// newCCConfig builds the cc/v4 config for one target. goos and goarch select
// the ABI; the include paths always come from the host compiler, so a
// cross-target run preprocesses against host headers.
func newCCConfig(goos, goarch string) (*cc.Config, error) {
	cfg, err := cc.NewConfig(goos, goarch, append([]string{ccStd}, ccDefines...)...)
	if err != nil {
		return nil, fmt.Errorf("cc/v4 host compiler configuration: %w", err)
	}

	if err := checkHostIncludes(hostIncludePaths(cfg), os.Getenv("NIX_CC")); err != nil {
		return nil, err
	}

	if sdk := darwinSDKIncludePath(); sdk != "" && !slices.Contains(cfg.SysIncludePaths, sdk) {
		cfg.IncludePaths = append(cfg.IncludePaths, sdk)
		cfg.SysIncludePaths = append(cfg.SysIncludePaths, sdk)
	}

	// FFmpeg headers include each other in both forms, "libavutil/x.h" and
	// <libavutil/x.h>, so the header root goes on both lists.
	cfg.IncludePaths = append(cfg.IncludePaths, AVLibPath)
	cfg.SysIncludePaths = append(cfg.SysIncludePaths, AVLibPath)

	// Header mode drops type checking of function bodies. FFmpeg headers carry
	// static inline definitions the generator never emits.
	cfg.Header = true

	return cfg, nil
}

// hostIncludePaths returns the include roots cc.NewConfig discovered by running
// the host compiler, quoted-include roots first. cc/v4 takes bare paths, so
// there are no "-I" prefixes on them.
func hostIncludePaths(cfg *cc.Config) []string {
	paths := make([]string, 0, len(cfg.HostIncludePaths)+len(cfg.HostSysIncludePaths))
	paths = append(paths, cfg.HostIncludePaths...)
	paths = append(paths, cfg.HostSysIncludePaths...)

	return paths
}

// checkHostIncludes enforces the dev-shell contract: inside the Nix shell, host
// include discovery must return something.
// Outside it an empty list passes, because a non-Nix host is not held to that
// contract.
func checkHostIncludes(paths []string, nixCC string) error {
	if len(paths) > 0 || nixCC == "" {
		return nil
	}

	return errors.New("host compiler returned no system include paths under NIX_CC; dev-shell contract requires successful include discovery")
}

// darwinSDKIncludePath resolves the macOS SDK include root through xcrun, or ""
// when the host is not macOS or xcrun is unavailable. On macOS cc.NewConfig
// runs the same clang driver xcrun configures, so the SDK root normally already
// sits in HostSysIncludePaths and newCCConfig skips this; it stays as a fallback
// because that could not be measured from a Linux host. Bare macOS without the Command Line Tools is
// outside the Nix contract, so a missing xcrun is not an error.
func darwinSDKIncludePath() string {
	if runtime.GOOS != "darwin" {
		return ""
	}

	out, err := exec.CommandContext(context.Background(), "xcrun", "--show-sdk-path").CombinedOutput()
	if err != nil {
		return ""
	}

	sdkPath := strings.TrimSpace(string(out))
	if sdkPath == "" {
		return ""
	}

	return filepath.Join(sdkPath, "usr", "include")
}
