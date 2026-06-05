package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func run() error {
	// Root context cancelled on interrupt so in-flight build commands stop.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize FFmpeg library enables after AllLibraries is fully defined
	CollectFFmpegEnables()

	// Derive output path based on architecture
	arch := runtime.GOARCH
	targetOutput := filepath.Join("lib", runtime.GOOS+"_"+arch, "libffmpeg.a")
	targetOutput, err := filepath.Abs(targetOutput)
	if err != nil {
		return fmt.Errorf("get absolute path for output: %w", err)
	}

	// Parse arguments
	selectedLibs := make(map[string]bool)
	cleanBuild := false
	listMode := false
	updateDigests := false

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--clean":
			cleanBuild = true
		case "--list":
			listMode = true
		case "--update-digests":
			updateDigests = true
		default:
			selectedLibs[arg] = true
		}
	}

	// Handle --list mode: display library information and exit
	if listMode {
		printLibraryList(AllLibraries)
		return nil
	}

	// Handle --update-digests mode: trust-on-first-use bootstrap. Downloads every
	// enabled library's archive, computes its SHA-256, and prints the manifest
	// entries for review. Run once in a trusted environment, then verify each
	// digest against upstream-published checksums and commit digests.go.
	if updateDigests {
		buildRoot, err := filepath.Abs(".build")
		if err != nil {
			return fmt.Errorf("get absolute path for build root: %w", err)
		}
		return updateDigestsMode(ctx, buildRoot, AllLibraries)
	}

	// Setup directories
	buildRoot, err := filepath.Abs(".build")
	if err != nil {
		return fmt.Errorf("get absolute path for build root: %w", err)
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
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
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

			if err := os.RemoveAll(libBuildDir); err != nil {
				log.Printf("warning: failed to remove %s: %v", libBuildDir, err)
			}
			if err := os.RemoveAll(libSrcDir); err != nil {
				log.Printf("warning: failed to remove %s: %v", libSrcDir, err)
			}

			// Also remove installed libraries from staging
			if lib.LinkLibs != nil {
				for _, libName := range lib.LinkLibs {
					for _, dir := range []string{"lib"} {
						libPath := filepath.Join(stagingDir, dir, libName+".a")
						if fileExists(libPath) {
							if err := os.Remove(libPath); err != nil {
								log.Printf("warning: failed to remove %s: %v", libPath, err)
							}
						}
					}
				}
			}
		}
		fmt.Printf("\n✓ Cleaned %d libraries\n", total)
		return nil
	}

	// Build all libraries
	total := len(libs)
	built := 0

	for i, lib := range libs {
		fmt.Printf("\n[%d/%d] %s\n", i+1, total, lib.Name)
		fmt.Println(strings.Repeat("=", 60))

		// Create per-library logger
		logDir := filepath.Join(buildRoot, "build", lib.Name)
		if err := os.MkdirAll(logDir, 0o755); err != nil {
			return fmt.Errorf("create log directory for %s: %w", lib.Name, err)
		}

		logFile := filepath.Join(logDir, "build.log")
		logFileWriter, err := os.Create(logFile)
		if err != nil {
			return fmt.Errorf("create log file for %s: %w", lib.Name, err)
		}

		// Use MultiWriter to send output to both stdout and log file
		logger := io.MultiWriter(os.Stdout, logFileWriter)

		if err := lib.Build(ctx, buildRoot, stagingDir, logger); err != nil {
			logFileWriter.Close()
			return fmt.Errorf("build failed for %s: %w\nSee log: %s", lib.Name, err, logFile)
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
		if err := combineLibraries(ctx, libs, stagingDir, targetOutput); err != nil {
			return fmt.Errorf("combine libraries: %w", err)
		}
		fmt.Printf("\n✓ Success! Output: %s\n", targetOutput)
	} else {
		fmt.Printf("\n✓ Success! Built %d selected libraries\n", len(selectedLibs))
	}

	return nil
}

// updateDigestsMode downloads every enabled library's archive, computes its
// SHA-256, and prints the manifest entries for archiveDigests. This is the
// trust-on-first-use bootstrap; run it once in a trusted environment, verify the
// printed digests against upstream-published checksums, then paste the entries
// into internal/builder/digests.go and commit.
func updateDigestsMode(ctx context.Context, buildRoot string, libs []*Library) error {
	_ = ctx // grab download is not yet context-aware; retained for signature parity

	downloadsDir := filepath.Join(buildRoot, "downloads")
	if err := os.MkdirAll(downloadsDir, 0o755); err != nil {
		return fmt.Errorf("create downloads dir: %w", err)
	}

	digests := make(map[string]string)
	for _, lib := range libs {
		if !lib.ShouldBuild() {
			fmt.Fprintf(os.Stderr, "Skipping %s (not built on %s)\n", lib.Name, runtime.GOOS)
			continue
		}
		if _, seen := digests[lib.URL]; seen {
			continue // shared URL already hashed
		}

		archivePath := filepath.Join(downloadsDir, filepath.Base(lib.URL))
		if !fileExists(archivePath) {
			if err := downloadRaw(lib.URL, archivePath, os.Stderr); err != nil {
				return fmt.Errorf("download %s: %w", lib.Name, err)
			}
		}

		digest, err := hashFile(archivePath)
		if err != nil {
			return fmt.Errorf("hash %s: %w", lib.Name, err)
		}
		digests[lib.URL] = digest
		fmt.Fprintf(os.Stderr, "✓ %s %s\n", lib.Name, digest)
	}

	fmt.Print("\n// Paste into archiveDigests in internal/builder/digests.go after\n")
	fmt.Print("// verifying each digest against the upstream-published checksum.\n")
	fmt.Print("var archiveDigests = map[string]string{\n")
	fmt.Print(formatDigestEntries(digests))
	fmt.Print("}\n")
	return nil
}

// combineLibraries combines all built libraries into a single static library
func combineLibraries(ctx context.Context, libs []*Library, stagingDir, output string) error {
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
		return combineMac(ctx, libFiles, output)
	}
	return combineLinux(ctx, libFiles, output)
}

// combineMac uses Apple's libtool to combine static libraries on macOS
// This is more efficient than ar as it doesn't require extracting all object files
func combineMac(ctx context.Context, libFiles []string, output string) error {
	log.Println("Using Apple libtool -static approach (macOS)")

	// Ensure output directory exists
	outputDir := filepath.Dir(output)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Remove existing output file if present
	os.Remove(output)

	// Use Apple's libtool (not GNU libtool from Nix) to combine libraries directly
	// Apple's libtool -static is specifically designed for this purpose
	args := append([]string{"-static", "-o", output}, libFiles...)
	libtoolCmd := exec.CommandContext(ctx, "/usr/bin/libtool", args...)

	// Capture output for debugging
	var stdout, stderr bytes.Buffer
	libtoolCmd.Stdout = &stdout
	libtoolCmd.Stderr = &stderr

	if err := libtoolCmd.Run(); err != nil {
		return fmt.Errorf("libtool failed: %w\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	// Strip the library to reduce size
	if err := exec.CommandContext(ctx, "strip", "-S", output).Run(); err != nil {
		return fmt.Errorf("strip failed: %w", err)
	}

	return nil
}

// combineLinux uses ar to combine static libraries on Linux
func combineLinux(ctx context.Context, libFiles []string, output string) error {
	log.Println("Using thin archive approach to merge libraries (Linux)")

	// Stack Overflow solution: create thin archive (low memory), then convert to normal archive
	// Thin archives use pointers instead of copying data, avoiding memory exhaustion
	// Source: https://stackoverflow.com/a/23621751

	var stderr bytes.Buffer

	// Ensure output directory exists
	outputDir := filepath.Dir(output)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Remove existing output file if present
	os.Remove(output)

	// Step 1: Create thin archive with -T flag (pointers only, minimal memory)
	log.Println("  Creating thin archive...")
	args := append([]string{"cqT", output}, libFiles...)
	thinCmd := exec.CommandContext(ctx, "ar", args...)
	thinCmd.Stderr = &stderr
	stderr.Reset()

	if err := thinCmd.Run(); err != nil {
		return fmt.Errorf("thin archive creation failed: %w (stderr: %s)", err, stderr.String())
	}

	// Step 2: Convert thin archive to normal archive using MRI script
	// This extracts and rebuilds, but from a single thin archive (less memory than 32 separate libraries)
	log.Println("  Converting to normal archive...")
	mriScript := fmt.Sprintf("create %s.tmp\naddlib %s\nsave\nend", output, output)

	convertCmd := exec.CommandContext(ctx, "ar", "-M")
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
	stripCmd := exec.CommandContext(ctx, "strip", "--strip-unneeded", output)
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

// colorize wraps text with an ANSI colour code and reset.
func colorize(text, color string) string {
	return color + text + colorReset
}

// tableCell formats a left-aligned cell with the given width and colour.
func tableCell(text string, width int, color string) string {
	return fmt.Sprintf("%s%-*s%s", color, width, text, colorReset)
}

// tableBorder returns a horizontal border segment of the given width.
func tableBorder(width int) string {
	return strings.Repeat("═", width+2)
}

// libraryColumnWidths returns the four column widths (name, platform, build
// system, link libraries) sized to the widest entry or header label.
func libraryColumnWidths(libs []*Library) (maxName, maxPlatform, maxBuildSys, maxLinkLibs int) {
	maxName = len("Library")
	maxPlatform = len("Platform")
	maxBuildSys = len("Build System")
	maxLinkLibs = len("Link Libraries")

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
	return maxName, maxPlatform, maxBuildSys, maxLinkLibs
}

// renderBorder builds a horizontal table border. prefix and suffix bracket the
// line, and left/mid/right are the corner and junction glyphs. The first column
// border is widened by 3 to fit the "#" and name cells.
func renderBorder(prefix, left, mid, right, suffix string, maxName, maxPlatform, maxBuildSys, maxLinkLibs int) string {
	return prefix + left +
		tableBorder(maxName+3) + mid +
		tableBorder(maxPlatform) + mid +
		tableBorder(maxBuildSys) + mid +
		tableBorder(maxLinkLibs) + right + suffix
}

// renderRow builds a table data row from four pre-coloured cells, separating
// the "#" and name cells inside the first column with a space.
func renderRow(numCell, nameCell, platformCell, buildSysCell, linkLibsCell string) string {
	return fmt.Sprintf("%s║%s %s %s %s ║%s %s %s║%s %s %s║%s %s %s║%s\n",
		colorCyan, colorReset,
		numCell,
		nameCell,
		colorCyan, colorReset,
		platformCell,
		colorCyan, colorReset,
		buildSysCell,
		colorCyan, colorReset,
		linkLibsCell,
		colorCyan, colorReset)
}

// printLibraryList displays a formatted table of all libraries
func printLibraryList(libs []*Library) {
	maxName, maxPlatform, maxBuildSys, maxLinkLibs := libraryColumnWidths(libs)

	// Top border
	fmt.Print(renderBorder(
		"\n"+colorize("", colorBold+colorCyan), "╔═", "╦", "╗", colorReset+"\n",
		maxName, maxPlatform, maxBuildSys, maxLinkLibs,
	))

	// Header row
	fmt.Print(renderRow(
		tableCell("#", 2, colorBold+colorYellow),
		tableCell("Library", maxName, colorBold+colorYellow),
		tableCell("Platform", maxPlatform, colorReset),
		tableCell("Build System", maxBuildSys, colorReset),
		tableCell("Link Libraries", maxLinkLibs, colorReset),
	))

	// Separator
	fmt.Print(renderBorder(
		colorCyan, "╠═", "╬", "╣", colorReset+"\n",
		maxName, maxPlatform, maxBuildSys, maxLinkLibs,
	))

	// Print rows
	for i, lib := range libs {
		num := fmt.Sprintf("%2d", i+1)
		platform := getPlatformString(lib)
		buildSys := getBuildSystemString(lib)

		// Get link libs display string (without embedded colors for now)
		linkLibsDisplay := getLinkLibsString(lib)

		// Color code based on library type
		nameColor := colorGreen
		linkLibsColor := colorReset
		if lib.LinkLibs == nil {
			nameColor = colorGray // Header-only libraries in gray
			linkLibsColor = colorGray
		}

		fmt.Print(renderRow(
			tableCell(num, 2, colorBlue+colorBold),
			tableCell(lib.Name, maxName, nameColor),
			tableCell(platform, maxPlatform, colorReset),
			tableCell(buildSys, maxBuildSys, colorReset),
			tableCell(linkLibsDisplay, maxLinkLibs, linkLibsColor),
		))
	}

	// Footer
	fmt.Print(renderBorder(
		colorCyan, "╚═", "╩", "╝", colorReset+"\n\n",
		maxName, maxPlatform, maxBuildSys, maxLinkLibs,
	))

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
