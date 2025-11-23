package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	// Initialize FFmpeg library enables after AllLibraries is fully defined
	CollectFFmpegEnables()

	// Derive output path based on architecture
	arch := runtime.GOARCH
	targetOutput := filepath.Join("lib", runtime.GOOS+"_"+arch, "libffmpeg.a")
	targetOutput, err := filepath.Abs(targetOutput)
	if err != nil {
		log.Fatalf("Failed to get absolute path for output: %v\n", err)
	}

	// Parse arguments
	selectedLibs := make(map[string]bool)
	cleanBuild := false
	listMode := false

	for _, arg := range os.Args[1:] {
		if arg == "--clean" {
			cleanBuild = true
		} else if arg == "--list" {
			listMode = true
		} else {
			selectedLibs[arg] = true
		}
	}

	// Handle --list mode: display library information and exit
	if listMode {
		printLibraryList(AllLibraries)
		return
	}

	// Setup directories
	buildRoot, err := filepath.Abs(".build")
	if err != nil {
		log.Fatalf("Failed to get absolute path for build root: %v\n", err)
	}
	stagingDir := filepath.Join(buildRoot, "staging")

	// Create directories (do NOT delete - incremental builds!)
	dirs := []string{
		filepath.Join(buildRoot, "downloads"),
		filepath.Join(buildRoot, "build"),
		filepath.Join(stagingDir, "lib"),
		filepath.Join(stagingDir, "include"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v\n", dir, err)
		}
	}

	// Filter libraries based on selection
	libs := AllLibraries
	if len(selectedLibs) > 0 {
		filtered := []*Library{}
		for _, lib := range AllLibraries {
			if selectedLibs[lib.Name] {
				filtered = append(filtered, lib)
			}
		}
		libs = filtered
	}

	// Handle --clean mode: clean and exit
	if cleanBuild {
		total := len(libs)
		for i, lib := range libs {
			fmt.Printf("[%d/%d] Cleaning %s...\n", i+1, total, lib.Name)

			libBuildDir := filepath.Join(buildRoot, "build", lib.Name)
			libSrcDir := filepath.Join(buildRoot, "src", lib.Name)

			os.RemoveAll(libBuildDir)
			os.RemoveAll(libSrcDir)

			// Also remove installed libraries from staging
			if lib.LinkLibs != nil {
				for _, libName := range lib.LinkLibs {
					for _, dir := range []string{"lib"} {
						libPath := filepath.Join(stagingDir, dir, libName+".a")
						if fileExists(libPath) {
							os.Remove(libPath)
						}
					}
				}
			}
		}
		fmt.Printf("\n✓ Cleaned %d libraries\n", total)
		return
	}

	// Build all libraries
	total := len(libs)
	built := 0

	for i, lib := range libs {
		fmt.Printf("\n[%d/%d] %s\n", i+1, total, lib.Name)
		fmt.Println(strings.Repeat("=", 60))

		// Create per-library logger
		logDir := filepath.Join(buildRoot, "build", lib.Name)
		os.MkdirAll(logDir, 0755)

		logFile := filepath.Join(logDir, "build.log")
		logFileWriter, err := os.Create(logFile)
		if err != nil {
			log.Fatalf("Failed to create log file for %s: %v\n", lib.Name, err)
		}

		// Use MultiWriter to send output to both stdout and log file
		logger := io.MultiWriter(os.Stdout, logFileWriter)

		if err := lib.Build(buildRoot, stagingDir, logger); err != nil {
			logFileWriter.Close()
			log.Fatalf("Build failed for %s: %v\nSee log: %s\n", lib.Name, err, logFile)
		}

		logFileWriter.Close()

		if lib.ShouldBuild() {
			built++
		}
	}

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Printf("Built %d/%d libraries\n", built, total)
	fmt.Println(strings.Repeat("=", 60))

	// Only combine libraries on a full build (no library filters)
	if len(selectedLibs) == 0 {
		if err := combineLibraries(libs, stagingDir, targetOutput); err != nil {
			log.Fatalf("Failed to combine libraries: %v\n", err)
		}
		fmt.Printf("\n✓ Success! Output: %s\n", targetOutput)
	} else {
		fmt.Printf("\n✓ Success! Built %d selected libraries\n", len(selectedLibs))
	}
}

// combineLibraries combines all built libraries into a single static library
func combineLibraries(libs []*Library, stagingDir, output string) error {
	// Collect library files from LinkLibs of all built libraries
	var libFiles []string
	linkLibsMap := make(map[string]bool) // Track which libs we need

	// Collect all LinkLibs from built libraries
	for _, lib := range libs {
		if !lib.ShouldBuild() {
			continue // Skip libraries that shouldn't be built on this platform
		}
		for _, linkLib := range lib.LinkLibs {
			linkLibsMap[linkLib] = true
		}
	}

	// Search for the specific .a files in lib
	libDirs := []string{
		filepath.Join(stagingDir, "lib"),
	}

	for _, libDir := range libDirs {
		entries, err := os.ReadDir(libDir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".a") {
				continue
			}
			// Check if this .a file matches one of our LinkLibs
			// Library files are named like "libfoo.a", so we check against the base name without .a
			baseName := strings.TrimSuffix(entry.Name(), ".a")
			if linkLibsMap[baseName] {
				libFiles = append(libFiles, filepath.Join(libDir, entry.Name()))
			}
		}
	}

	if len(libFiles) == 0 {
		return fmt.Errorf("no static libraries found in %s", stagingDir)
	}

	log.Printf("Combining %d libraries into %s\n", len(libFiles), output)

	// Platform-specific merge
	if runtime.GOOS == "darwin" {
		return combineMac(libFiles, output)
	}
	return combineLinux(libFiles, output)
}

// combineMac uses libtool to combine static libraries on macOS
// This is more efficient than ar as it doesn't require extracting all object files
func combineMac(libFiles []string, output string) error {
	log.Println("Using libtool -static approach (macOS)")

	// Ensure output directory exists
	outputDir := filepath.Dir(output)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Remove existing output file if present
	os.Remove(output)

	// Use libtool to combine libraries directly
	// libtool -static is Apple's tool specifically designed for this purpose
	args := append([]string{"-static", "-o", output}, libFiles...)
	libtoolCmd := exec.Command("libtool", args...)

	// Capture output for debugging
	var stdout, stderr bytes.Buffer
	libtoolCmd.Stdout = &stdout
	libtoolCmd.Stderr = &stderr

	if err := libtoolCmd.Run(); err != nil {
		return fmt.Errorf("libtool failed: %w\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	// Strip the library to reduce size
	if err := exec.Command("strip", "-S", output).Run(); err != nil {
		return fmt.Errorf("strip failed: %w", err)
	}

	return nil
}

// combineLinux uses ar to combine static libraries on Linux
func combineLinux(libFiles []string, output string) error {
	log.Println("Using thin archive approach to merge libraries (Linux)")

	// Stack Overflow solution: create thin archive (low memory), then convert to normal archive
	// Thin archives use pointers instead of copying data, avoiding memory exhaustion
	// Source: https://stackoverflow.com/a/23621751

	var stderr bytes.Buffer

	// Ensure output directory exists
	outputDir := filepath.Dir(output)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Remove existing output file if present
	os.Remove(output)

	// Step 1: Create thin archive with -T flag (pointers only, minimal memory)
	log.Println("  Creating thin archive...")
	args := append([]string{"cqT", output}, libFiles...)
	thinCmd := exec.Command("ar", args...)
	thinCmd.Stderr = &stderr
	stderr.Reset()

	if err := thinCmd.Run(); err != nil {
		return fmt.Errorf("thin archive creation failed: %w (stderr: %s)", err, stderr.String())
	}

	// Step 2: Convert thin archive to normal archive using MRI script
	// This extracts and rebuilds, but from a single thin archive (less memory than 32 separate libraries)
	log.Println("  Converting to normal archive...")
	mriScript := fmt.Sprintf("create %s.tmp\naddlib %s\nsave\nend", output, output)

	convertCmd := exec.Command("ar", "-M")
	convertCmd.Stdin = bytes.NewBufferString(mriScript)
	convertCmd.Stderr = &stderr
	stderr.Reset()

	if err := convertCmd.Run(); err != nil {
		return fmt.Errorf("thin archive conversion failed: %w (stderr: %s)", err, stderr.String())
	}

	// Replace original with converted archive
	if err := os.Rename(output+".tmp", output); err != nil {
		return fmt.Errorf("failed to rename converted archive: %w", err)
	}

	// Strip the library
	stripCmd := exec.Command("strip", "--strip-unneeded", output)
	stripCmd.Stderr = &stderr
	stderr.Reset()
	if err := stripCmd.Run(); err != nil {
		return fmt.Errorf("strip failed: %w (stderr: %s)", err, stderr.String())
	}

	return nil
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// printLibraryList displays a formatted table of all libraries
func printLibraryList(libs []*Library) {
	// Calculate column widths
	maxName := len("Library")
	maxPlatform := len("Platform")
	maxBuildSys := len("Build System")
	maxLinkLibs := len("Link Libraries")

	for _, lib := range libs {
		if len(lib.Name) > maxName {
			maxName = len(lib.Name)
		}
		platform := getPlatformString(lib)
		if len(platform) > maxPlatform {
			maxPlatform = len(platform)
		}
		buildSys := getBuildSystemString(lib)
		if len(buildSys) > maxBuildSys {
			maxBuildSys = len(buildSys)
		}
		linkLibs := getLinkLibsString(lib)
		if len(linkLibs) > maxLinkLibs {
			maxLinkLibs = len(linkLibs)
		}
	}

	// Print header
	fmt.Printf("\n%s%s╔═", colorBold, colorCyan)
	fmt.Printf("%s", strings.Repeat("═", maxName+6))
	fmt.Printf("╦═%s", strings.Repeat("═", maxPlatform+2))
	fmt.Printf("╦═%s", strings.Repeat("═", maxBuildSys+2))
	fmt.Printf("╦═%s", strings.Repeat("═", maxLinkLibs+2))
	fmt.Printf("╗%s\n", colorReset)

	fmt.Printf("%s║%s %s #%s   %s%-*s%s %s║%s %-*s %s ║%s %-*s %s ║%s %-*s %s ║%s\n",
		colorCyan, colorReset,
		colorBold+colorYellow, colorReset,
		colorBold+colorYellow, maxName, "Library", colorReset,
		colorCyan, colorReset,
		maxPlatform, "Platform",
		colorCyan, colorReset,
		maxBuildSys, "Build System",
		colorCyan, colorReset,
		maxLinkLibs, "Link Libraries",
		colorCyan, colorReset)

	fmt.Printf("%s╠═", colorCyan)
	fmt.Printf("%s", strings.Repeat("═", maxName+6))
	fmt.Printf("╬═%s", strings.Repeat("═", maxPlatform+2))
	fmt.Printf("╬═%s", strings.Repeat("═", maxBuildSys+2))
	fmt.Printf("╬═%s", strings.Repeat("═", maxLinkLibs+2))
	fmt.Printf("╣%s\n", colorReset)

	// Print rows
	for i, lib := range libs {
		num := fmt.Sprintf("%2d", i+1)
		platform := getPlatformString(lib)
		buildSys := getBuildSystemString(lib)

		// Get link libs display string (without embedded colors for now)
		var linkLibsDisplay string
		if lib.LinkLibs == nil {
			linkLibsDisplay = "(headers only)"
		} else if len(lib.LinkLibs) == 0 {
			linkLibsDisplay = "-"
		} else {
			linkLibsDisplay = strings.Join(lib.LinkLibs, ", ")
		}

		// Color code based on library type
		nameColor := colorGreen
		linkLibsColor := colorReset
		if lib.LinkLibs == nil {
			nameColor = colorGray // Header-only libraries in gray
			linkLibsColor = colorGray
		}

		fmt.Printf("%s║%s %s%s%s %s%-*s %s  ║%s %-*s %s ║%s %-*s %s ║%s %s%-*s%s %s ║%s\n",
			colorCyan, colorReset,
			colorBlue+colorBold, num, colorReset,
			nameColor, maxName, lib.Name,
			colorCyan, colorReset,
			maxPlatform, platform,
			colorCyan, colorReset,
			maxBuildSys, buildSys,
			colorCyan, colorReset,
			linkLibsColor, maxLinkLibs, linkLibsDisplay, colorReset,
			colorCyan, colorReset)
	}

	// Print footer
	fmt.Printf("%s╚═", colorCyan)
	fmt.Printf("%s", strings.Repeat("═", maxName+6))
	fmt.Printf("╩═%s", strings.Repeat("═", maxPlatform+2))
	fmt.Printf("╩═%s", strings.Repeat("═", maxBuildSys+2))
	fmt.Printf("╩═%s", strings.Repeat("═", maxLinkLibs+2))
	fmt.Printf("╝%s\n\n", colorReset)

	// Summary
	totalLibs := len(libs)
	headerOnly := 0
	for _, lib := range libs {
		if lib.LinkLibs == nil {
			headerOnly++
		}
	}
	fmt.Printf("%sTotal: %d libraries (%d build, %d header-only)%s\n",
		colorBold, totalLibs, totalLibs-headerOnly, headerOnly, colorReset)
}

// getPlatformString returns a formatted platform string
func getPlatformString(lib *Library) string {
	if len(lib.Platform) == 0 {
		return "-"
	}
	return strings.Join(lib.Platform, ", ")
}

// getBuildSystemString returns the build system type name
func getBuildSystemString(lib *Library) string {
	switch lib.BuildSystem.(type) {
	case *AutoconfBuild:
		return "autoconf"
	case *CMakeBuild:
		return "cmake"
	case *MesonBuild:
		return "meson"
	case *CargoBuild:
		return "cargo"
	case *MakefileBuild:
		return "makefile"
	default:
		return "unknown"
	}
}

// getLinkLibsString returns the display length of link libraries for column width calculation
func getLinkLibsString(lib *Library) string {
	if lib.LinkLibs == nil {
		return "(headers only)"
	}
	if len(lib.LinkLibs) == 0 {
		return "-"
	}
	return strings.Join(lib.LinkLibs, ", ")
}
