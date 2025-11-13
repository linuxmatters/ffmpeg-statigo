package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Library represents a third-party library to build
type Library struct {
	Name          string
	URL           string
	ArchiveType   string // "tar.gz", "tar.bz2", "tar.xz", "zip"
	StripPrefix   string
	Platform      []string // empty = all platforms, ["linux"], ["darwin"], etc.
	BuildSystem   BuildSystem
	ConfigureArgs func(os string) []string
	PostExtract   func(srcPath string) error // optional patches
	SkipAutoFlags bool                       // Skip automatic CFLAGS/LDFLAGS (for non-standard configure scripts like zlib)
	LinkLibs      []string                   // Libraries to link in final static lib (nil for header-only)
	Dependencies  []*Library
}

// BuildSystem defines the interface for different build systems
type BuildSystem interface {
	Configure(lib *Library, srcPath, buildDir, installDir string) error
	Build(lib *Library, srcPath, buildDir string) error
}

// ShouldBuild checks if this library should be built on the current platform
func (lib *Library) ShouldBuild() bool {
	if len(lib.Platform) == 0 {
		return true // no platform restriction = build everywhere
	}
	currentOS := runtime.GOOS
	for _, platform := range lib.Platform {
		if platform == currentOS {
			return true
		}
	}
	return false
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

	if err := ExtractArchive(archivePath, srcPath, lib.ArchiveType, lib.StripPrefix, logger); err != nil {
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
	h.Write([]byte(lib.ArchiveType))
	h.Write([]byte(lib.StripPrefix))
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
	// Include both lib and lib64 pkgconfig directories
	return filepath.Join(installDir, "lib", "pkgconfig") + ":" + filepath.Join(installDir, "lib64", "pkgconfig")
}

// buildEnv returns environment variables for building
func buildEnv(installDir string) []string {
	env := os.Environ()
	pkgConfigPath := pkgConfigPath(installDir)

	// Update or add PKG_CONFIG_PATH
	updated := false
	for i, e := range env {
		if strings.HasPrefix(e, "PKG_CONFIG_PATH=") {
			existing := strings.TrimPrefix(e, "PKG_CONFIG_PATH=")
			env[i] = "PKG_CONFIG_PATH=" + pkgConfigPath + ":" + existing
			updated = true
			break
		}
	}
	if !updated {
		env = append(env, "PKG_CONFIG_PATH="+pkgConfigPath)
	}

	return env
}
