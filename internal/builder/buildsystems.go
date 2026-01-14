package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// stagingDir derives the staging/install directory from a build directory path.
// Build directories follow the pattern: .build/<lib>/build, so staging is at .build/staging.
func stagingDir(buildDir string) string {
	return filepath.Join(filepath.Dir(filepath.Dir(buildDir)), "staging")
}

// openBuildLog opens a build log file and returns a multiwriter that writes to both
// the log file and stdout. The returned cleanup function must be called when done.
// If append is true, the log file is opened in append mode; otherwise it is truncated.
func openBuildLog(buildDir string, append bool) (io.Writer, func(), error) {
	logFile := filepath.Join(buildDir, "build.log")
	var logger *os.File
	var err error
	if append {
		logger, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		logger, err = os.Create(logFile)
	}
	if err != nil {
		return nil, nil, err
	}
	return io.MultiWriter(logger, os.Stdout), func() { logger.Close() }, nil
}

// withBuildLog opens a build log and executes the provided function with the logger.
// This helper consolidates the repeated log setup/cleanup pattern across build systems.
func withBuildLog(buildDir string, append bool, fn func(output io.Writer) error) error {
	output, cleanup, err := openBuildLog(buildDir, append)
	if err != nil {
		return err
	}
	defer cleanup()
	return fn(output)
}

// AutoconfBuild implements the BuildSystem interface for autoconf-based builds
type AutoconfBuild struct{}

func (a *AutoconfBuild) Configure(lib *Library, srcPath, buildDir, installDir string) error {
	args := []string{
		fmt.Sprintf("--prefix=%s", installDir),
	}

	if lib.ConfigureArgs != nil {
		args = append(args, lib.ConfigureArgs(runtime.GOOS)...)
	}

	// Add standard compiler and linker flags (unless library opts out)
	// Some libraries like zlib have non-standard configure scripts that reject these
	if !lib.SkipAutoFlags {
		incDir := filepath.Join(installDir, "include")
		libDir := filepath.Join(installDir, "lib")

		cflags := fmt.Sprintf("-O3 -I%s", incDir)
		cppflags := fmt.Sprintf("-I%s", incDir)
		cxxflags := "-O3" // Start with optimization flag
		ldflags := fmt.Sprintf("-L%s", libDir)

		// On macOS, configure proper SDK and header paths
		// This is required for CGO compilation and C++ builds
		if runtime.GOOS == "darwin" {
			// CGO_CFLAGS format: -isysroot /path/to/sdk -I/nix/store/.../lib/clang/18/include
			cgoCflags := os.Getenv("CGO_CFLAGS")
			libcxxInclude := os.Getenv("LIBCXX_INCLUDE")

			// For C: add SDK path and clang builtins
			if cgoCflags != "" {
				cflags = fmt.Sprintf("%s %s", cflags, cgoCflags)
			}

			// For CPPFLAGS: only add -isysroot, NOT the clang builtin include path
			// This prevents C++ from finding clang's stddef.h before libc++'s wrapper
			if cgoCflags != "" {
				// Extract just the -isysroot part from CGO_CFLAGS
				// CGO_CFLAGS = "-isysroot /path/to/sdk -I/path/to/clang/include"
				// We only want the -isysroot portion for CPPFLAGS
				parts := strings.Fields(cgoCflags)
				for i, part := range parts {
					if part == "-isysroot" && i+1 < len(parts) {
						cppflags = fmt.Sprintf("%s -isysroot %s", cppflags, parts[i+1])
						break
					}
				}
			}

			// For C++: use -nostdinc++ to disable built-in paths, then add libc++ first,
			// then clang builtins (for stddef.h, stdarg.h), then SDK path
			if libcxxInclude != "" && cgoCflags != "" {
				cxxflags = fmt.Sprintf("%s -nostdinc++ -I%s %s", cxxflags, libcxxInclude, cgoCflags)
			} else if cgoCflags != "" {
				cxxflags = fmt.Sprintf("%s %s", cxxflags, cgoCflags)
			}
		}

		args = append(args,
			fmt.Sprintf("CFLAGS=%s", cflags),
			fmt.Sprintf("CPPFLAGS=%s", cppflags),
			fmt.Sprintf("CXXFLAGS=%s", cxxflags),
			fmt.Sprintf("LDFLAGS=%s", ldflags),
		)
	}

	// Run configure from source directory
	configurePath := "./configure"
	absConfigurePath := filepath.Join(srcPath, "configure")
	if !fileExists(absConfigurePath) {
		return fmt.Errorf("configure script not found at %s", absConfigurePath)
	}

	// Make configure executable
	if err := os.Chmod(absConfigurePath, 0755); err != nil {
		return fmt.Errorf("failed to make configure executable: %w", err)
	}

	return withBuildLog(buildDir, false, func(output io.Writer) error {
		return runCommand(srcPath, output, installDir, configurePath, args...)
	})
}

func (a *AutoconfBuild) Build(lib *Library, srcPath, buildDir string) error {
	// Touch automake files to prevent regeneration
	touchAutomakeFiles(srcPath)

	installDir := stagingDir(buildDir)

	return withBuildLog(buildDir, true, func(output io.Writer) error {
		// make
		if err := runCommand(srcPath, output, installDir, "make", "-j", fmt.Sprintf("%d", runtime.NumCPU())); err != nil {
			return err
		}
		// make install
		return runCommand(srcPath, output, installDir, "make", "install")
	})
}

// CMakeBuild implements the BuildSystem interface for CMake-based builds
type CMakeBuild struct {
	SourceSubdir string // Optional subdirectory containing source (e.g. "source" for x265)
}

func (c *CMakeBuild) Configure(lib *Library, srcPath, buildDir, installDir string) error {
	// Determine actual source path
	actualSrcPath := srcPath
	if c.SourceSubdir != "" {
		actualSrcPath = filepath.Join(srcPath, c.SourceSubdir)
	}

	args := []string{
		actualSrcPath,
		fmt.Sprintf("-DCMAKE_INSTALL_PREFIX=%s", installDir),
		"-DCMAKE_BUILD_TYPE=Release",
		"-DCMAKE_INSTALL_LIBDIR=lib",
	}

	if lib.ConfigureArgs != nil {
		args = append(args, lib.ConfigureArgs(runtime.GOOS)...)
	}

	return withBuildLog(buildDir, false, func(output io.Writer) error {
		return runCommand(buildDir, output, installDir, "cmake", args...)
	})
}

func (c *CMakeBuild) Build(lib *Library, srcPath, buildDir string) error {
	installDir := stagingDir(buildDir)

	return withBuildLog(buildDir, true, func(output io.Writer) error {
		// cmake --build . --target install
		if err := runCommand(buildDir, output, installDir, "cmake", "--build", ".", "--parallel", fmt.Sprintf("%d", runtime.NumCPU())); err != nil {
			return err
		}
		return runCommand(buildDir, output, installDir, "cmake", "--build", ".", "--target", "install")
	})
}

// MesonBuild implements the BuildSystem interface for Meson-based builds
type MesonBuild struct{}

func (m *MesonBuild) Configure(lib *Library, srcPath, buildDir, installDir string) error {
	args := []string{
		"setup",
		buildDir,
		srcPath,
		fmt.Sprintf("--prefix=%s", installDir),
		"--buildtype=release",
		"--default-library=static",
		"--libdir=lib",
	}

	if lib.ConfigureArgs != nil {
		args = append(args, lib.ConfigureArgs(runtime.GOOS)...)
	}

	return withBuildLog(buildDir, false, func(output io.Writer) error {
		return runCommand(".", output, installDir, "meson", args...)
	})
}

func (m *MesonBuild) Build(lib *Library, srcPath, buildDir string) error {
	installDir := stagingDir(buildDir)

	return withBuildLog(buildDir, true, func(output io.Writer) error {
		// meson compile
		if err := runCommand(buildDir, output, installDir, "meson", "compile"); err != nil {
			return err
		}
		// meson install
		return runCommand(buildDir, output, installDir, "meson", "install")
	})
}

// CargoBuild implements the BuildSystem interface for Cargo/Rust-based builds
type CargoBuild struct {
	InstallFunc func(srcPath, installDir string) error // Custom install function
}

func (c *CargoBuild) Configure(lib *Library, srcPath, buildDir, installDir string) error {
	// Cargo doesn't have a separate configure step
	return nil
}

func (c *CargoBuild) Build(lib *Library, srcPath, buildDir string) error {
	installDir := stagingDir(buildDir)

	// Custom install func handles the full cargo build process if provided
	if c.InstallFunc != nil {
		return c.InstallFunc(srcPath, installDir)
	}
	return nil
}

// MakefileBuild implements the BuildSystem interface for Makefile-based builds
type MakefileBuild struct {
	Targets     []string
	InstallFunc func(srcPath, installDir string) error
}

func (m *MakefileBuild) Configure(lib *Library, srcPath, buildDir, installDir string) error {
	// Makefile builds don't have a configure step
	return nil
}

func (m *MakefileBuild) Build(lib *Library, srcPath, buildDir string) error {
	installDir := stagingDir(buildDir)

	return withBuildLog(buildDir, false, func(output io.Writer) error {
		// Build the targets
		args := append([]string{"-j", fmt.Sprintf("%d", runtime.NumCPU())}, m.Targets...)
		if err := runCommand(srcPath, output, installDir, "make", args...); err != nil {
			return err
		}

		// If a custom install function is provided, use it
		if m.InstallFunc != nil {
			return m.InstallFunc(srcPath, installDir)
		}
		return nil
	})
}

// OpenSSLBuild implements the BuildSystem interface for OpenSSL's Configure/make
type OpenSSLBuild struct{}

func (o *OpenSSLBuild) Configure(lib *Library, srcPath, buildDir, installDir string) error {
	// OpenSSL uses 'Configure' (capital C) Perl script, not autoconf
	args := []string{
		fmt.Sprintf("--prefix=%s", installDir),
		fmt.Sprintf("--openssldir=%s", filepath.Join(installDir, "ssl")),
		"--libdir=lib",
		fmt.Sprintf("--with-zlib-include=%s/include", installDir),
		fmt.Sprintf("--with-zlib-lib=%s/lib", installDir),
		"zlib", // Enable zlib support
	}

	if lib.ConfigureArgs != nil {
		args = append(args, lib.ConfigureArgs(runtime.GOOS)...)
	}

	// Run Configure from source directory
	configurePath := "./Configure"
	absConfigurePath := filepath.Join(srcPath, "Configure")
	if !fileExists(absConfigurePath) {
		return fmt.Errorf("Configure script not found at %s", absConfigurePath)
	}

	// Make Configure executable
	if err := os.Chmod(absConfigurePath, 0755); err != nil {
		return fmt.Errorf("failed to make Configure executable: %w", err)
	}

	return withBuildLog(buildDir, false, func(output io.Writer) error {
		return runCommand(srcPath, output, installDir, configurePath, args...)
	})
}

func (o *OpenSSLBuild) Build(lib *Library, srcPath, buildDir string) error {
	installDir := stagingDir(buildDir)

	return withBuildLog(buildDir, true, func(output io.Writer) error {
		// make
		if err := runCommand(srcPath, output, installDir, "make", "-j", fmt.Sprintf("%d", runtime.NumCPU())); err != nil {
			return err
		}
		// make install_sw (install software only, skip docs)
		return runCommand(srcPath, output, installDir, "make", "install_sw")
	})
}
