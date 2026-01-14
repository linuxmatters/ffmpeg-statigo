package main

import (
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

// Library represents a third-party library to build
type Library struct {
	Name          string
	URL           string
	Enabled       *bool    // Defaults to true (nil = enabled) - set to false to disable without removing code
	Platform      []string // empty = all platforms, ["linux"], ["darwin"], etc.
	FFmpegEnables []string // Optional FFmpeg --enable-* flags (e.g. ["libx264"], ["nvenc", "nvdec"])
	BuildSystem   BuildSystem
	ConfigureArgs func(os string) []string
	PostExtract   func(srcPath string) error // optional patches
	SkipAutoFlags bool                       // Skip automatic CFLAGS/LDFLAGS (for non-standard configure scripts like zlib)
	LinkLibs      []string                   // Libraries to link in final static lib (nil for header-only)
	Dependencies  []*Library                 // Build-time dependencies; platform-specific deps are auto-filtered via ShouldBuild()
}

// BuildSystem defines the interface for different build systems
type BuildSystem interface {
	Configure(lib *Library, srcPath, buildDir, installDir string) error
	Build(lib *Library, srcPath, buildDir string) error
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
func (lib *Library) Build(buildRoot, installDir string, logger io.Writer) error {
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
		if err := lib.PostExtract(srcPath); err != nil {
			return fmt.Errorf("post-extract failed: %w", err)
		}
	}

	// Build
	buildDir := filepath.Join(buildRoot, "build", lib.Name)
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build dir: %w", err)
	}

	if err := lib.BuildSystem.Configure(lib, srcPath, buildDir, installDir); err != nil {
		return fmt.Errorf("configure failed: %w", err)
	}

	if err := lib.BuildSystem.Build(lib, srcPath, buildDir); err != nil {
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

// runCommand executes a command and streams output to logger
func runCommand(dir string, logger io.Writer, installDir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = logger
	cmd.Stderr = logger

	// Set environment with PKG_CONFIG_PATH
	cmd.Env = buildEnv(installDir)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", name, err)
	}
	return nil
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

// buildEnv returns environment variables for building
func buildEnv(installDir string) []string {
	env := os.Environ()
	pkgConfigPath := pkgConfigPath(installDir)

	// Update or add PKG_CONFIG_PATH
	updatedPkg := false
	for i, e := range env {
		if strings.HasPrefix(e, "PKG_CONFIG_PATH=") {
			existing := strings.TrimPrefix(e, "PKG_CONFIG_PATH=")
			env[i] = "PKG_CONFIG_PATH=" + pkgConfigPath + ":" + existing
			updatedPkg = true
			break
		}
	}
	if !updatedPkg {
		env = append(env, "PKG_CONFIG_PATH="+pkgConfigPath)
	}

	// Update or add PATH to include staging/bin for tools like glslang, spirv-*
	binPath := filepath.Join(installDir, "bin")
	updatedPath := false
	for i, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			existing := strings.TrimPrefix(e, "PATH=")
			env[i] = "PATH=" + binPath + ":" + existing
			updatedPath = true
			break
		}
	}
	if !updatedPath {
		env = append(env, "PATH="+binPath)
	}

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
			updatedCflags := false
			for i, e := range env {
				if strings.HasPrefix(e, "CFLAGS=") {
					existing := strings.TrimPrefix(e, "CFLAGS=")
					env[i] = "CFLAGS=" + existing + " " + cgoCflags
					updatedCflags = true
					break
				}
			}
			if !updatedCflags {
				env = append(env, "CFLAGS="+cgoCflags)
			}

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
			updatedCxxflags := false
			for i, e := range env {
				if strings.HasPrefix(e, "CXXFLAGS=") {
					existing := strings.TrimPrefix(e, "CXXFLAGS=")
					env[i] = "CXXFLAGS=" + existing + " " + cxxExtra
					updatedCxxflags = true
					break
				}
			}
			if !updatedCxxflags {
				env = append(env, "CXXFLAGS="+cxxExtra)
			}

			// Extract SDK path from CGO_CFLAGS (-isysroot <path>) for LDFLAGS
			// This ensures cargo/rustc can find SDK libraries like libiconv
			var sdkPath string
			parts := strings.Fields(cgoCflags)
			for i, p := range parts {
				if p == "-isysroot" && i+1 < len(parts) {
					sdkPath = parts[i+1]
					break
				}
			}
			if sdkPath != "" {
				ldExtra := "-L" + filepath.Join(sdkPath, "usr", "lib")
				updatedLdflags := false
				for i, e := range env {
					if strings.HasPrefix(e, "LDFLAGS=") {
						existing := strings.TrimPrefix(e, "LDFLAGS=")
						env[i] = "LDFLAGS=" + existing + " " + ldExtra
						updatedLdflags = true
						break
					}
				}
				if !updatedLdflags {
					env = append(env, "LDFLAGS="+ldExtra)
				}

				// Also set LIBRARY_PATH for cargo/rustc which may not use LDFLAGS
				libraryPath := filepath.Join(sdkPath, "usr", "lib")
				updatedLibraryPath := false
				for i, e := range env {
					if strings.HasPrefix(e, "LIBRARY_PATH=") {
						existing := strings.TrimPrefix(e, "LIBRARY_PATH=")
						env[i] = "LIBRARY_PATH=" + libraryPath + ":" + existing
						updatedLibraryPath = true
						break
					}
				}
				if !updatedLibraryPath {
					env = append(env, "LIBRARY_PATH="+libraryPath)
				}
			}
		}
	}

	return env
}
