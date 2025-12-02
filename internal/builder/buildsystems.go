package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
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
		ldflags := fmt.Sprintf("-L%s", libDir)

		args = append(args,
			fmt.Sprintf("CFLAGS=%s", cflags),
			fmt.Sprintf("CPPFLAGS=%s", cppflags),
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

	output, cleanup, err := openBuildLog(buildDir, false)
	if err != nil {
		return err
	}
	defer cleanup()

	if err := runCommand(srcPath, output, installDir, configurePath, args...); err != nil {
		return err
	}

	return nil
}

func (a *AutoconfBuild) Build(lib *Library, srcPath, buildDir string) error {
	output, cleanup, err := openBuildLog(buildDir, true)
	if err != nil {
		return err
	}
	defer cleanup()

	// Touch automake files to prevent regeneration
	touchAutomakeFiles(srcPath)

	installDir := stagingDir(buildDir)

	// make
	if err := runCommand(srcPath, output, installDir, "make", "-j", fmt.Sprintf("%d", runtime.NumCPU())); err != nil {
		return err
	}

	// make install
	if err := runCommand(srcPath, output, installDir, "make", "install"); err != nil {
		return err
	}

	return nil
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

	output, cleanup, err := openBuildLog(buildDir, false)
	if err != nil {
		return err
	}
	defer cleanup()

	if err := runCommand(buildDir, output, installDir, "cmake", args...); err != nil {
		return err
	}

	return nil
}

func (c *CMakeBuild) Build(lib *Library, srcPath, buildDir string) error {
	output, cleanup, err := openBuildLog(buildDir, true)
	if err != nil {
		return err
	}
	defer cleanup()

	installDir := stagingDir(buildDir)

	// cmake --build . --target install
	if err := runCommand(buildDir, output, installDir, "cmake", "--build", ".", "--parallel", fmt.Sprintf("%d", runtime.NumCPU())); err != nil {
		return err
	}

	if err := runCommand(buildDir, output, installDir, "cmake", "--build", ".", "--target", "install"); err != nil {
		return err
	}

	return nil
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

	output, cleanup, err := openBuildLog(buildDir, false)
	if err != nil {
		return err
	}
	defer cleanup()

	if err := runCommand(".", output, installDir, "meson", args...); err != nil {
		return err
	}

	return nil
}

func (m *MesonBuild) Build(lib *Library, srcPath, buildDir string) error {
	output, cleanup, err := openBuildLog(buildDir, true)
	if err != nil {
		return err
	}
	defer cleanup()

	installDir := stagingDir(buildDir)

	// meson compile
	if err := runCommand(buildDir, output, installDir, "meson", "compile"); err != nil {
		return err
	}

	// meson install
	if err := runCommand(buildDir, output, installDir, "meson", "install"); err != nil {
		return err
	}

	return nil
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
	output, cleanup, err := openBuildLog(buildDir, false)
	if err != nil {
		return err
	}
	defer cleanup()

	installDir := stagingDir(buildDir)

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

	output, cleanup, err := openBuildLog(buildDir, false)
	if err != nil {
		return err
	}
	defer cleanup()

	if err := runCommand(srcPath, output, installDir, configurePath, args...); err != nil {
		return err
	}

	return nil
}

func (o *OpenSSLBuild) Build(lib *Library, srcPath, buildDir string) error {
	output, cleanup, err := openBuildLog(buildDir, true)
	if err != nil {
		return err
	}
	defer cleanup()

	installDir := stagingDir(buildDir)

	// make
	if err := runCommand(srcPath, output, installDir, "make", "-j", fmt.Sprintf("%d", runtime.NumCPU())); err != nil {
		return err
	}

	// make install_sw (install software only, skip docs)
	if err := runCommand(srcPath, output, installDir, "make", "install_sw"); err != nil {
		return err
	}

	return nil
}
