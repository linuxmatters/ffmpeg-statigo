package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

// extractSDKPath returns the path following -isysroot in a flags string,
// or empty string if not found.
func extractSDKPath(flags string) string {
	parts := strings.Fields(flags)
	for i, p := range parts {
		if p == "-isysroot" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// Library represents a third-party library to build
type Library struct {
	Name          string
	URL           string
	Enabled       *bool    // Defaults to true (nil = enabled) - set to false to disable without removing code
	Platform      []string // empty = all platforms, ["linux"], ["darwin"], etc.
	FFmpegEnables []string // Optional FFmpeg --enable-* flags (e.g. ["libx264"], ["nvenc", "nvdec"])
	BuildSystem   BuildSystem
	ConfigureArgs func(os string) []string
	PostExtract   func(ctx context.Context, srcPath string) error // optional patches
	BuildEnv      func() []string                                 // optional extra KEY=value env applied to all build commands
	SkipAutoFlags bool                                            // Skip automatic CFLAGS/LDFLAGS (for non-standard configure scripts like zlib)
	LinkLibs      []string                                        // Libraries to link in final static lib (nil for header-only)
	Dependencies  []*Library                                      // Build-time dependencies; platform-specific deps are auto-filtered via ShouldBuild()
}

// BuildSystem defines the interface for different build systems
type BuildSystem interface {
	Configure(ctx context.Context, lib *Library, srcPath, buildDir, installDir string) error
	Build(ctx context.Context, lib *Library, srcPath, buildDir string) error
}

// ShouldBuild checks if this library should be built on the current platform
func (lib *Library) ShouldBuild() bool {
	// Check if library is enabled (nil = true, explicitly set to false = disabled)
	if lib.Enabled != nil && !*lib.Enabled {
		return false
	}

	// Check platform restrictions
	if len(lib.Platform) == 0 {
		return true // no platform restriction = build everywhere
	}
	return slices.Contains(lib.Platform, runtime.GOOS)
}

// extraEnv returns the library's extra build env, or nil if none is set.
func (lib *Library) extraEnv() []string {
	if lib.BuildEnv == nil {
		return nil
	}
	return lib.BuildEnv()
}

// archiveExtensions maps file suffixes to canonical archive type names.
var archiveExtensions = map[string]string{
	".tar.gz":  "tar.gz",
	".tgz":     "tar.gz",
	".tar.bz2": "tar.bz2",
	".tbz2":    "tar.bz2",
	".tar.xz":  "tar.xz",
	".txz":     "tar.xz",
	".zip":     "zip",
}

// ArchiveType derives the archive type from the URL
func (lib *Library) ArchiveType() string {
	url := strings.ToLower(lib.URL)
	for ext, archiveType := range archiveExtensions {
		if strings.HasSuffix(url, ext) {
			return archiveType
		}
	}
	return ""
}

// Build performs the complete build process for this library
func (lib *Library) Build(ctx context.Context, buildRoot, installDir string, logger io.Writer) error {
	if !lib.ShouldBuild() {
		fmt.Fprintf(logger, "Skipping %s (platform: %v, current: %s)\n", lib.Name, lib.Platform, runtime.GOOS)
		return nil
	}

	// Check if we can skip this build
	state := NewBuildState(lib, buildRoot)
	if state.CanSkip(installDir) {
		fmt.Fprintf(logger, "Skipping %s (already built)\n", lib.Name)
		return nil
	}

	fmt.Fprintf(logger, "Building %s...\n", lib.Name)

	// Download and extract
	archivePath := filepath.Join(buildRoot, "downloads", filepath.Base(lib.URL))
	if err := DownloadFile(lib.URL, archivePath, logger); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	srcPath := filepath.Join(buildRoot, "src", lib.Name)
	if err := os.RemoveAll(srcPath); err != nil {
		return fmt.Errorf("failed to clean source dir: %w", err)
	}

	if err := ExtractArchive(archivePath, srcPath, lib.ArchiveType(), logger); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	// Post-extract hook (for patches, etc.)
	if lib.PostExtract != nil {
		if err := lib.PostExtract(ctx, srcPath); err != nil {
			return fmt.Errorf("post-extract failed: %w", err)
		}
	}

	// Build
	buildDir := filepath.Join(buildRoot, "build", lib.Name)
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return fmt.Errorf("failed to create build dir: %w", err)
	}

	if err := lib.BuildSystem.Configure(ctx, lib, srcPath, buildDir, installDir); err != nil {
		return fmt.Errorf("configure failed: %w", err)
	}

	if err := lib.BuildSystem.Build(ctx, lib, srcPath, buildDir); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Save state
	if err := state.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Fprintf(logger, "✓ %s complete\n", lib.Name)
	return nil
}

// ConfigHash computes a hash of the library's configuration
func (lib *Library) ConfigHash() string {
	h := sha256.New()
	h.Write([]byte(lib.URL))
	h.Write([]byte(lib.Name))
	for _, p := range lib.Platform {
		h.Write([]byte(p))
	}
	if lib.ConfigureArgs != nil {
		for _, arg := range lib.ConfigureArgs(runtime.GOOS) {
			h.Write([]byte(arg))
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// runCommandEnv runs a command with extra KEY=value entries appended to the build env.
func runCommandEnv(ctx context.Context, dir string, logger io.Writer, installDir string, extraEnv []string, name string, args ...string) error {
	// name and args come from internal build definitions, not user input.
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec // G702: build commands are project-defined, not external
	cmd.Dir = dir
	cmd.Stdout = logger
	cmd.Stderr = logger
	cmd.Env = append(buildEnv(installDir), extraEnv...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", name, err)
	}
	return nil
}

// runCommand executes a command and streams output to logger
func runCommand(ctx context.Context, dir string, logger io.Writer, installDir string, name string, args ...string) error {
	return runCommandEnv(ctx, dir, logger, installDir, nil, name, args...)
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// pkgConfigPath returns the PKG_CONFIG_PATH environment variable
func pkgConfigPath(installDir string) string {
	return filepath.Join(installDir, "lib", "pkgconfig")
}

// upsertEnv sets key in env. If an existing "key=..." entry is present, value is
// combined with the existing value using sep — prepended when prepend is true,
// appended otherwise. If key is absent, "key=value" is appended.
func upsertEnv(env []string, key, value, sep string, prepend bool) []string {
	prefix := key + "="
	for i, e := range env {
		if existing, ok := strings.CutPrefix(e, prefix); ok {
			if prepend {
				env[i] = prefix + value + sep + existing
			} else {
				env[i] = prefix + existing + sep + value
			}
			return env
		}
	}
	return append(env, prefix+value)
}

// buildEnv returns environment variables for building
func buildEnv(installDir string) []string {
	env := os.Environ()
	pkgConfigPath := pkgConfigPath(installDir)

	// Update or add PKG_CONFIG_PATH
	env = upsertEnv(env, "PKG_CONFIG_PATH", pkgConfigPath, ":", true)

	// Update or add PATH to include staging/bin for tools like glslang, spirv-*
	binPath := filepath.Join(installDir, "bin")
	env = upsertEnv(env, "PATH", binPath, ":", true)

	// On macOS, remove NIX_CFLAGS_COMPILE which interferes with C++ header search order
	// The Nix clang wrapper injects -isystem paths that cause libc++ headers to be
	// searched after C standard library headers, breaking <cstddef> and similar wrappers
	if runtime.GOOS == "darwin" {
		filtered := env[:0]
		for _, e := range env {
			if !strings.HasPrefix(e, "NIX_CFLAGS_COMPILE=") {
				filtered = append(filtered, e)
			}
		}
		env = filtered

		// Force clang as the compiler on macOS
		// The Nix dev shell includes both gcc and clang, but our CFLAGS include
		// paths to Clang's builtin headers (stddef.h, stdarg.h) which use Clang-specific
		// features like __has_feature and __building_module that gcc doesn't understand
		env = append(env, "CC=clang", "CXX=clang++")
	}

	// On macOS, ensure CFLAGS/CXXFLAGS include SDK path and clang builtin headers
	// This is required for both C (stdarg.h, stddef.h) and C++ (<algorithm>, <cstring>) compilation
	if runtime.GOOS == "darwin" {
		cgoCflags := os.Getenv("CGO_CFLAGS")
		if cgoCflags != "" {
			// Set CFLAGS with full CGO_CFLAGS (includes -isysroot and -I.../clang/18/include)
			env = upsertEnv(env, "CFLAGS", cgoCflags, " ", false)

			// Build CXXFLAGS with -nostdinc++ and explicit libcxx include path
			// Use -nostdinc++ to disable built-in C++ paths, preventing NIX_CFLAGS_COMPILE
			// from interfering with header search order
			// Then add libc++ headers before clang builtins
			var cxxExtra string
			libcxxInclude := os.Getenv("LIBCXX_INCLUDE")
			if libcxxInclude != "" {
				cxxExtra = "-nostdinc++ -I" + libcxxInclude + " " + cgoCflags
			} else {
				cxxExtra = cgoCflags
			}

			// Set CXXFLAGS with same flags for C++ builds
			env = upsertEnv(env, "CXXFLAGS", cxxExtra, " ", false)

			// Extract SDK path from CGO_CFLAGS (-isysroot <path>) for LDFLAGS
			// This ensures cargo/rustc can find SDK libraries like libiconv
			sdkPath := extractSDKPath(cgoCflags)
			if sdkPath != "" {
				ldExtra := "-L" + filepath.Join(sdkPath, "usr", "lib")
				env = upsertEnv(env, "LDFLAGS", ldExtra, " ", false)

				// Also set LIBRARY_PATH for cargo/rustc which may not use LDFLAGS
				libraryPath := filepath.Join(sdkPath, "usr", "lib")
				env = upsertEnv(env, "LIBRARY_PATH", libraryPath, ":", true)
			}
		}
	}

	return env
}
