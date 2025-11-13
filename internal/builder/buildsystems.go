package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

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
		lib64Dir := filepath.Join(installDir, "lib64")

		args = append(args,
			fmt.Sprintf("CFLAGS=-I%s", incDir),
			fmt.Sprintf("CPPFLAGS=-I%s", incDir),
			fmt.Sprintf("LDFLAGS=-L%s -L%s", libDir, lib64Dir),
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

	logFile := filepath.Join(buildDir, "build.log")
	logger, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer logger.Close()

	multiWriter := io.MultiWriter(logger, os.Stdout)

	if err := runCommand(srcPath, multiWriter, installDir, configurePath, args...); err != nil {
		return err
	}

	return nil
}

func (a *AutoconfBuild) Build(lib *Library, srcPath, buildDir string) error {
	logFile := filepath.Join(buildDir, "build.log")
	logger, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer logger.Close()

	multiWriter := io.MultiWriter(logger, os.Stdout)

	// Touch automake files to prevent regeneration
	touchAutomakeFiles(srcPath)

	// Get installDir from buildDir path
	installDir := filepath.Join(filepath.Dir(filepath.Dir(buildDir)), "staging")

	// make
	if err := runCommand(srcPath, multiWriter, installDir, "make", "-j", fmt.Sprintf("%d", runtime.NumCPU())); err != nil {
		return err
	}

	// make install
	if err := runCommand(srcPath, multiWriter, installDir, "make", "install"); err != nil {
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
	}

	if lib.ConfigureArgs != nil {
		args = append(args, lib.ConfigureArgs(runtime.GOOS)...)
	}

	logFile := filepath.Join(buildDir, "build.log")
	logger, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer logger.Close()

	multiWriter := io.MultiWriter(logger, os.Stdout)

	if err := runCommand(buildDir, multiWriter, installDir, "cmake", args...); err != nil {
		return err
	}

	return nil
}

func (c *CMakeBuild) Build(lib *Library, srcPath, buildDir string) error {
	logFile := filepath.Join(buildDir, "build.log")
	logger, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer logger.Close()

	multiWriter := io.MultiWriter(logger, os.Stdout)

	// Get installDir from buildDir path
	installDir := filepath.Join(filepath.Dir(filepath.Dir(buildDir)), "staging")

	// cmake --build . --target install
	if err := runCommand(buildDir, multiWriter, installDir, "cmake", "--build", ".", "--parallel", fmt.Sprintf("%d", runtime.NumCPU())); err != nil {
		return err
	}

	if err := runCommand(buildDir, multiWriter, installDir, "cmake", "--build", ".", "--target", "install"); err != nil {
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
	}

	if lib.ConfigureArgs != nil {
		args = append(args, lib.ConfigureArgs(runtime.GOOS)...)
	}

	logFile := filepath.Join(buildDir, "build.log")
	logger, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer logger.Close()

	multiWriter := io.MultiWriter(logger, os.Stdout)

	if err := runCommand(".", multiWriter, installDir, "meson", args...); err != nil {
		return err
	}

	return nil
}

func (m *MesonBuild) Build(lib *Library, srcPath, buildDir string) error {
	logFile := filepath.Join(buildDir, "build.log")
	logger, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer logger.Close()

	multiWriter := io.MultiWriter(logger, os.Stdout)

	// Get installDir from buildDir path
	installDir := filepath.Join(filepath.Dir(filepath.Dir(buildDir)), "staging")

	// meson compile
	if err := runCommand(buildDir, multiWriter, installDir, "meson", "compile"); err != nil {
		return err
	}

	// meson install
	if err := runCommand(buildDir, multiWriter, installDir, "meson", "install"); err != nil {
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
	// Get installDir from buildDir path
	installDir := filepath.Join(filepath.Dir(filepath.Dir(buildDir)), "staging")

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
	logFile := filepath.Join(buildDir, "build.log")
	logger, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer logger.Close()

	multiWriter := io.MultiWriter(logger, os.Stdout)

	// Get installDir from buildDir path
	installDir := filepath.Join(filepath.Dir(filepath.Dir(buildDir)), "staging")

	// Build the targets
	args := append([]string{"-j", fmt.Sprintf("%d", runtime.NumCPU())}, m.Targets...)
	if err := runCommand(srcPath, multiWriter, installDir, "make", args...); err != nil {
		return err
	}

	// If a custom install function is provided, use it
	if m.InstallFunc != nil {
		return m.InstallFunc(srcPath, installDir)
	}

	return nil
}
