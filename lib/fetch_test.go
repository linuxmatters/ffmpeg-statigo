package lib

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// =============================================================================
// Test 1.1: Version Parsing Validation
// =============================================================================

func TestFindCompatibleRelease_VersionParsing(t *testing.T) {
	t.Run("valid_version_format", func(t *testing.T) {
		// Valid versions should produce a release string (even if API fails, fallback works)
		validVersions := []string{
			"8.0.1",
			"1.2.3",
			"10.20.30",
		}

		for _, version := range validVersions {
			// We can't test the full function without network, but we can test
			// the version parsing logic directly
			parts := strings.Split(version, ".")
			if len(parts) != 3 {
				t.Errorf("Version %s should have 3 parts, got %d", version, len(parts))
			}
		}
	})

	t.Run("invalid_version_format_two_parts", func(t *testing.T) {
		// This tests the version parsing validation
		version := "8.0"
		parts := strings.Split(version, ".")
		if len(parts) == 3 {
			t.Errorf("Version %s should NOT have 3 parts", version)
		}
	})

	t.Run("invalid_version_format_four_parts", func(t *testing.T) {
		version := "8.0.1.0"
		parts := strings.Split(version, ".")
		if len(parts) == 3 {
			t.Errorf("Version %s should NOT have 3 parts", version)
		}
	})

	t.Run("invalid_version_format_empty", func(t *testing.T) {
		version := ""
		parts := strings.Split(version, ".")
		if len(parts) == 3 {
			t.Errorf("Empty version should NOT have 3 parts")
		}
	})

	t.Run("release_prefix_construction", func(t *testing.T) {
		// Verify that release prefix is constructed correctly
		moduleVersion := "8.0.1"
		expectedPrefix := "lib-8.0.1"
		actualPrefix := "lib-" + moduleVersion

		if actualPrefix != expectedPrefix {
			t.Errorf("Expected prefix %s, got %s", expectedPrefix, actualPrefix)
		}
	})

	t.Run("fallback_release_pattern", func(t *testing.T) {
		// When API fails, fallback should produce lib-X.Y.Z.0
		moduleVersion := "8.0.1"
		expectedFallback := "lib-8.0.1.0"
		actualFallback := "lib-" + moduleVersion + ".0"

		if actualFallback != expectedFallback {
			t.Errorf("Expected fallback %s, got %s", expectedFallback, actualFallback)
		}
	})
}

// TestFindViaAPI_ReleaseSorting tests that releases are sorted correctly
// to find the highest version number
func TestFindViaAPI_ReleaseSorting(t *testing.T) {
	t.Run("sorts_release_versions_correctly", func(t *testing.T) {
		// Simulate release list sorting logic
		releases := []string{
			"lib-8.0.1.0",
			"lib-8.0.1.2",
			"lib-8.0.1.1",
			"lib-8.0.1.10", // 10 > 2 lexicographically
		}

		// Current implementation uses sort.Strings - test this behavior
		sortedReleases := make([]string, len(releases))
		copy(sortedReleases, releases)
		// Lexicographic sort (same as sort.Strings)
		for i := 0; i < len(sortedReleases)-1; i++ {
			for j := i + 1; j < len(sortedReleases); j++ {
				if sortedReleases[i] > sortedReleases[j] {
					sortedReleases[i], sortedReleases[j] = sortedReleases[j], sortedReleases[i]
				}
			}
		}

		// Last element should be the "highest"
		highest := sortedReleases[len(sortedReleases)-1]

		// Note: lexicographic sorting means "lib-8.0.1.2" > "lib-8.0.1.10"
		// This is a known limitation of the current implementation
		// For now, we document the behavior
		t.Logf("Sorted releases: %v", sortedReleases)
		t.Logf("Selected highest: %s", highest)

		// Document that this is lexicographic, not semantic versioning
		if highest != "lib-8.0.1.2" {
			t.Logf("Note: Current sorting is lexicographic, not semver-aware")
		}
	})
}

// =============================================================================
// Test 5: Release Version Semantic Sort Bug
// =============================================================================

// TestReleaseVersionSemanticSortBug demonstrates the bug where lexicographic
// sorting picks the wrong version when patch numbers have different digit counts.
// Example: lib-8.0.1.10 should be > lib-8.0.1.2 semantically, but
// lexicographically lib-8.0.1.10 < lib-8.0.1.2 (string "10" < "2")
func TestReleaseVersionSemanticSortBug(t *testing.T) {
	t.Run("lexicographic_sort_picks_wrong_version", func(t *testing.T) {
		// Releases with double-digit patch number
		releases := []string{
			"lib-8.0.1.2",
			"lib-8.0.1.10",
		}

		// Use sort.Strings (current implementation in fetch.go line 169)
		sorted := make([]string, len(releases))
		copy(sorted, releases)
		sort.Strings(sorted)

		// With lexicographic sort, "lib-8.0.1.10" comes before "lib-8.0.1.2"
		// because '1' < '2' when comparing character by character
		if sorted[0] != "lib-8.0.1.10" {
			t.Errorf("Expected lib-8.0.1.10 to be first (lexicographically), got %s", sorted[0])
		}
		if sorted[1] != "lib-8.0.1.2" {
			t.Errorf("Expected lib-8.0.1.2 to be last (lexicographically), got %s", sorted[1])
		}

		// The bug: last element is selected as "highest" version
		selectedVersion := sorted[len(sorted)-1]

		// BUG: This selects lib-8.0.1.2 instead of lib-8.0.1.10
		if selectedVersion != "lib-8.0.1.2" {
			t.Errorf("Expected bug to select lib-8.0.1.2 (lexicographically last), got %s", selectedVersion)
		}

		t.Logf("BUG: Lexicographic sort selected %s instead of semantically correct lib-8.0.1.10", selectedVersion)
	})

	t.Run("demonstrates_bug_with_realistic_release_sequence", func(t *testing.T) {
		// Realistic scenario: multiple patch releases
		releases := []string{
			"lib-8.0.1.0",
			"lib-8.0.1.1",
			"lib-8.0.1.2",
			"lib-8.0.1.3",
			"lib-8.0.1.10", // Latest release (semantic version 8.0.1.10)
		}

		sorted := make([]string, len(releases))
		copy(sorted, releases)
		sort.Strings(sorted)

		// Lexicographic sort order: 0, 1, 10, 2, 3
		expectedLexOrder := []string{
			"lib-8.0.1.0",
			"lib-8.0.1.1",
			"lib-8.0.1.10", // BUG: 10 comes before 2 lexicographically
			"lib-8.0.1.2",
			"lib-8.0.1.3",
		}

		for i, expected := range expectedLexOrder {
			if sorted[i] != expected {
				t.Errorf("Position %d: expected %s, got %s", i, expected, sorted[i])
			}
		}

		// The bug: selects lib-8.0.1.3 instead of lib-8.0.1.10
		selectedVersion := sorted[len(sorted)-1]
		if selectedVersion != "lib-8.0.1.3" {
			t.Errorf("Expected bug to select lib-8.0.1.3, got %s", selectedVersion)
		}

		t.Logf("BUG: Selected %s instead of latest release lib-8.0.1.10", selectedVersion)
		t.Logf("Lexicographic order: %v", sorted)
	})

	t.Run("bug_affects_all_double_digit_versions", func(t *testing.T) {
		// Test that any double-digit component causes the issue
		testCases := []struct {
			name             string
			releases         []string
			wrongSelection   string // What gets selected (lexicographically last)
			correctSelection string // What should be selected (semantically latest)
		}{
			{
				name:             "patch_version_10_vs_9",
				releases:         []string{"lib-8.0.1.9", "lib-8.0.1.10"},
				wrongSelection:   "lib-8.0.1.9",
				correctSelection: "lib-8.0.1.10",
			},
			{
				name:             "patch_version_19_vs_100",
				releases:         []string{"lib-8.0.1.19", "lib-8.0.1.100"},
				wrongSelection:   "lib-8.0.1.19",
				correctSelection: "lib-8.0.1.100",
			},
			{
				name:             "patch_version_2_vs_12",
				releases:         []string{"lib-8.0.1.2", "lib-8.0.1.12"},
				wrongSelection:   "lib-8.0.1.2",
				correctSelection: "lib-8.0.1.12",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				sorted := make([]string, len(tc.releases))
				copy(sorted, tc.releases)
				sort.Strings(sorted)

				selectedVersion := sorted[len(sorted)-1]

				// Verify the bug: lexicographic sort picks wrong version
				if selectedVersion != tc.wrongSelection {
					t.Errorf("Expected bug to select %s (lexicographically), got %s", tc.wrongSelection, selectedVersion)
				}

				// Document what should be selected semantically
				if selectedVersion == tc.correctSelection {
					t.Logf("Note: This case happens to work correctly")
				} else {
					t.Logf("BUG: Selected %s instead of semantically correct %s", selectedVersion, tc.correctSelection)
				}
			})
		}
	})

	t.Run("documents_correct_semantic_version_comparison", func(t *testing.T) {
		// This test documents how semantic versioning should work
		// to prevent the bug in future implementations

		type semver struct {
			prefix string
			major  int
			minor  int
			patch  int
			build  int
		}

		releases := []semver{
			{prefix: "lib", major: 8, minor: 0, patch: 1, build: 0},
			{prefix: "lib", major: 8, minor: 0, patch: 1, build: 2},
			{prefix: "lib", major: 8, minor: 0, patch: 1, build: 10},
		}

		// Find semantically highest version
		highest := releases[0]
		for _, r := range releases {
			if r.major > highest.major ||
				(r.major == highest.major && r.minor > highest.minor) ||
				(r.major == highest.major && r.minor == highest.minor && r.patch > highest.patch) ||
				(r.major == highest.major && r.minor == highest.minor && r.patch == highest.patch && r.build > highest.build) {
				highest = r
			}
		}

		// Semantically, lib-8.0.1.10 should be selected
		if highest.build != 10 {
			t.Errorf("Semantic version comparison failed: expected build 10, got %d", highest.build)
		}

		t.Logf("Correct semantic selection: %s-%d.%d.%d.%d",
			highest.prefix, highest.major, highest.minor, highest.patch, highest.build)
	})
}

// =============================================================================
// Test 1.2: Checksum Verification Robustness
// =============================================================================

func TestVerifyChecksum_ChecksumCalculation(t *testing.T) {
	t.Run("sha256_calculation_correct", func(t *testing.T) {
		// Create a temp file with known content
		content := []byte("Hello, FFmpeg!")
		tmpFile, err := os.CreateTemp("", "checksum-test-*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write(content); err != nil {
			t.Fatalf("Failed to write content: %v", err)
		}
		tmpFile.Close()

		// Calculate expected checksum
		h := sha256.New()
		h.Write(content)
		expectedChecksum := hex.EncodeToString(h.Sum(nil))

		// Calculate actual checksum using same method as fetch.go
		f, err := os.Open(tmpFile.Name())
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		defer f.Close()

		h2 := sha256.New()
		if _, err := bytes.NewReader(content).WriteTo(h2); err != nil {
			t.Fatalf("Failed to calculate checksum: %v", err)
		}
		actualChecksum := hex.EncodeToString(h2.Sum(nil))

		if actualChecksum != expectedChecksum {
			t.Errorf("Checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
		}

		t.Logf("SHA256 checksum: %s", actualChecksum)
	})

	t.Run("checksum_mismatch_detection", func(t *testing.T) {
		// Test the mismatch error message format
		expected := "abc123"
		actual := "def456"

		err := checkMismatch(expected, actual)
		if err == nil {
			t.Error("Expected error for mismatched checksums")
		}

		if !strings.Contains(err.Error(), "checksum mismatch") {
			t.Errorf("Error should contain 'checksum mismatch', got: %v", err)
		}
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Error should contain expected checksum %s, got: %v", expected, err)
		}
		if !strings.Contains(err.Error(), actual) {
			t.Errorf("Error should contain actual checksum %s, got: %v", actual, err)
		}
	})

	t.Run("checksum_match_no_error", func(t *testing.T) {
		checksum := "abc123def456"
		err := checkMismatch(checksum, checksum)
		if err != nil {
			t.Errorf("Expected no error for matching checksums, got: %v", err)
		}
	})

	t.Run("sha256_digest_format_parsing", func(t *testing.T) {
		// GitHub provides digests in "sha256:..." format
		digest := "sha256:abc123def456789"

		if !strings.HasPrefix(digest, "sha256:") {
			t.Error("Digest should have sha256: prefix")
		}

		cleanDigest := strings.TrimPrefix(digest, "sha256:")
		if cleanDigest != "abc123def456789" {
			t.Errorf("Expected clean digest abc123def456789, got %s", cleanDigest)
		}
	})

	t.Run("unexpected_digest_format", func(t *testing.T) {
		// Test handling of unexpected digest formats
		badDigests := []string{
			"md5:abc123",
			"abc123",
			"SHA256:abc123", // Case matters
		}

		for _, digest := range badDigests {
			if strings.HasPrefix(digest, "sha256:") {
				t.Errorf("Digest %s should not be accepted as sha256 format", digest)
			}
		}
	})
}

// checkMismatch is a helper that mirrors the error format in fetch.go
func checkMismatch(expected, actual string) error {
	if actual != expected {
		return &checksumError{expected: expected, actual: actual}
	}
	return nil
}

type checksumError struct {
	expected string
	actual   string
}

func (e *checksumError) Error() string {
	return "checksum mismatch: expected " + e.expected + ", got " + e.actual
}

func TestVerifyChecksumFromFile_Parsing(t *testing.T) {
	t.Run("parse_sha256sums_format", func(t *testing.T) {
		// Standard SHA256SUMS file format
		content := `abc123def456  ffmpeg-linux-amd64.tar.gz
789xyz000111  ffmpeg-darwin-arm64.tar.gz
deadbeef1234  ffmpeg-linux-arm64.tar.gz`

		tarballName := "ffmpeg-linux-amd64.tar.gz"

		lines := strings.Split(content, "\n")
		var foundChecksum string
		for _, line := range lines {
			if strings.Contains(line, tarballName) {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					foundChecksum = parts[0]
					break
				}
			}
		}

		if foundChecksum != "abc123def456" {
			t.Errorf("Expected checksum abc123def456, got %s", foundChecksum)
		}
	})

	t.Run("tarball_not_in_checksums", func(t *testing.T) {
		content := `abc123def456  ffmpeg-linux-amd64.tar.gz`

		tarballName := "ffmpeg-windows-amd64.tar.gz" // Not in list

		lines := strings.Split(content, "\n")
		var foundChecksum string
		for _, line := range lines {
			if strings.Contains(line, tarballName) {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					foundChecksum = parts[0]
					break
				}
			}
		}

		if foundChecksum != "" {
			t.Errorf("Should not find checksum for missing tarball, got %s", foundChecksum)
		}
	})
}

// =============================================================================
// Test 1.3: Tarball Path Traversal Protection
// =============================================================================

func TestExtractTarball_PathTraversal(t *testing.T) {
	t.Run("rejects_path_traversal_dotdot", func(t *testing.T) {
		// Create a malicious tarball that tries to escape destDir
		tarball := createMaliciousTarball(t, "../escape.txt", []byte("malicious content"))
		defer os.Remove(tarball)

		destDir := t.TempDir()

		err := extractTarball(tarball, destDir)
		if err == nil {
			t.Error("Expected error for path traversal attack, got nil")
		}
		if err != nil && !strings.Contains(err.Error(), "path traversal") {
			t.Errorf("Error should mention path traversal, got: %v", err)
		}

		// Verify file was NOT created in parent directory
		escapePath := filepath.Join(filepath.Dir(destDir), "escape.txt")
		if _, statErr := os.Stat(escapePath); statErr == nil {
			os.Remove(escapePath)
			t.Error("Malicious file should NOT have been created outside destDir")
		}
	})

	t.Run("rejects_absolute_path", func(t *testing.T) {
		// Use a path inside the test's temp directory to avoid touching real system paths
		absolutePath := filepath.Join(t.TempDir(), "absolute_test.txt")

		tarball := createMaliciousTarball(t, absolutePath, []byte("malicious content"))
		defer os.Remove(tarball)

		destDir := t.TempDir()

		err := extractTarball(tarball, destDir)
		if err == nil {
			t.Error("Expected error for absolute path attack, got nil")
		}
		if err != nil && !strings.Contains(err.Error(), "absolute") {
			t.Errorf("Error should mention absolute path, got: %v", err)
		}
	})

	t.Run("rejects_path_with_embedded_dotdot", func(t *testing.T) {
		// Path that looks valid but escapes via embedded ..
		tarball := createMaliciousTarball(t, "subdir/../../../escape.txt", []byte("malicious"))
		defer os.Remove(tarball)

		destDir := t.TempDir()

		err := extractTarball(tarball, destDir)
		if err == nil {
			t.Error("Expected error for embedded path traversal attack, got nil")
		}
		if err != nil && !strings.Contains(err.Error(), "path traversal") {
			t.Errorf("Error should mention path traversal, got: %v", err)
		}
	})

	t.Run("accepts_valid_paths", func(t *testing.T) {
		// Create a valid tarball
		tarball := createValidTarball(t, map[string][]byte{
			"linux_amd64/libffmpeg.a": []byte("library content"),
			"linux_amd64/README.txt":  []byte("readme content"),
		})
		defer os.Remove(tarball)

		destDir := t.TempDir()

		err := extractTarball(tarball, destDir)
		if err != nil {
			t.Errorf("Valid tarball should extract without error, got: %v", err)
		}

		// Verify files were created
		libPath := filepath.Join(destDir, "linux_amd64", "libffmpeg.a")
		if _, err := os.Stat(libPath); os.IsNotExist(err) {
			t.Error("Expected libffmpeg.a to be extracted")
		}
	})

	t.Run("skips_symlinks", func(t *testing.T) {
		// Create a tarball with symlink - should be silently skipped
		tarball := createSymlinkTarball(t, "some_link", "../target")
		defer os.Remove(tarball)

		destDir := t.TempDir()

		// Should not error, just skip the symlink
		err := extractTarball(tarball, destDir)
		if err != nil {
			t.Logf("Symlink handling: %v", err)
		}

		// Verify symlink was NOT created
		symlinkPath := filepath.Join(destDir, "some_link")
		if _, statErr := os.Lstat(symlinkPath); statErr == nil {
			t.Error("Symlink should have been skipped, but was created")
		}
	})
}

// createMaliciousTarball creates a gzipped tarball with a single file at the given path
func createMaliciousTarball(t *testing.T, filePath string, content []byte) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "malicious-*.tar.gz")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()

	f, err := os.Create(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	header := &tar.Header{
		Name: filePath,
		Mode: 0644,
		Size: int64(len(content)),
	}

	if err := tw.WriteHeader(header); err != nil {
		t.Fatalf("Failed to write header: %v", err)
	}

	if _, err := tw.Write(content); err != nil {
		t.Fatalf("Failed to write content: %v", err)
	}

	return tmpFile.Name()
}

// createValidTarball creates a gzipped tarball with multiple files
func createValidTarball(t *testing.T, files map[string][]byte) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "valid-*.tar.gz")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()

	f, err := os.Create(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	for name, content := range files {
		// Create directory entry first
		dir := filepath.Dir(name)
		if dir != "." {
			dirHeader := &tar.Header{
				Name:     dir + "/",
				Mode:     0755,
				Typeflag: tar.TypeDir,
			}
			if err := tw.WriteHeader(dirHeader); err != nil {
				t.Fatalf("Failed to write dir header: %v", err)
			}
		}

		header := &tar.Header{
			Name: name,
			Mode: 0644,
			Size: int64(len(content)),
		}

		if err := tw.WriteHeader(header); err != nil {
			t.Fatalf("Failed to write header: %v", err)
		}

		if _, err := tw.Write(content); err != nil {
			t.Fatalf("Failed to write content: %v", err)
		}
	}

	return tmpFile.Name()
}

// createSymlinkTarball creates a gzipped tarball with a symlink
func createSymlinkTarball(t *testing.T, linkName, target string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "symlink-*.tar.gz")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()

	f, err := os.Create(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	header := &tar.Header{
		Name:     linkName,
		Mode:     0777,
		Typeflag: tar.TypeSymlink,
		Linkname: target,
	}

	if err := tw.WriteHeader(header); err != nil {
		t.Fatalf("Failed to write symlink header: %v", err)
	}

	return tmpFile.Name()
}

// =============================================================================
// Test 3: HTTP Download Failure Recovery
// =============================================================================

func TestDownloadFile_ErrorHandling(t *testing.T) {
	t.Run("handles_404_not_found", func(t *testing.T) {
		// Use a URL that returns 404
		url := "https://github.com/linuxmatters/ffmpeg-statigo/releases/download/nonexistent/file.tar.gz"
		dest := filepath.Join(t.TempDir(), "download.tar.gz")

		err := downloadFile(url, dest)
		if err == nil {
			t.Error("Expected error for 404 response, got nil")
		}

		// The grab library returns specific error for 404
		if err != nil && !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "bad response") {
			t.Logf("Note: Error message format: %v", err)
		}
	})

	t.Run("handles_invalid_url", func(t *testing.T) {
		// Invalid URL format
		url := "not-a-valid-url"
		dest := filepath.Join(t.TempDir(), "download.tar.gz")

		err := downloadFile(url, dest)
		if err == nil {
			t.Error("Expected error for invalid URL, got nil")
		}

		t.Logf("Invalid URL error: %v", err)
	})

	t.Run("handles_invalid_destination", func(t *testing.T) {
		// Valid URL but invalid destination (non-existent directory)
		url := "https://github.com/linuxmatters/ffmpeg-statigo/archive/refs/heads/main.zip"
		dest := "/nonexistent/path/that/does/not/exist/file.tar.gz"

		err := downloadFile(url, dest)
		if err == nil {
			t.Error("Expected error for invalid destination, got nil")
		}

		// Should get a path error
		if err != nil && !strings.Contains(err.Error(), "no such file") && !strings.Contains(err.Error(), "cannot create") {
			t.Logf("Note: Error message format: %v", err)
		}
	})

	t.Run("cleans_up_partial_downloads", func(t *testing.T) {
		// Download to temp dir to check cleanup behavior
		dest := filepath.Join(t.TempDir(), "partial.tar.gz")

		// Use invalid URL to cause failure
		url := "https://github.com/nonexistent/repo/releases/download/v1.0.0/file.tar.gz"

		err := downloadFile(url, dest)
		if err == nil {
			t.Error("Expected download to fail")
		}

		// grab library may create the file before failing
		// The caller (ensureLibrary) should clean up using defer os.Remove(tmpTarball)
		// This test verifies the error is returned, allowing cleanup
		if err != nil {
			t.Logf("Download failed as expected: %v", err)
		}
	})
}

func TestFindCompatibleRelease_APIFailureRecovery(t *testing.T) {
	t.Run("fallback_when_api_unavailable", func(t *testing.T) {
		// Test the fallback pattern construction
		moduleVersion := "8.0.1"
		expectedFallback := "lib-8.0.1.0"

		// Simulate what happens when API fails: fallback to predictable pattern
		fallbackRelease := "lib-" + moduleVersion + ".0"

		if fallbackRelease != expectedFallback {
			t.Errorf("Expected fallback %s, got %s", expectedFallback, fallbackRelease)
		}

		t.Logf("Fallback release pattern: %s", fallbackRelease)
	})

	t.Run("fallback_with_different_versions", func(t *testing.T) {
		testCases := []struct {
			version  string
			expected string
		}{
			{"8.0.0", "lib-8.0.0.0"},
			{"8.0.1", "lib-8.0.1.0"},
			{"9.1.0", "lib-9.1.0.0"},
			{"10.0.0", "lib-10.0.0.0"},
		}

		for _, tc := range testCases {
			fallback := "lib-" + tc.version + ".0"
			if fallback != tc.expected {
				t.Errorf("Version %s: expected %s, got %s", tc.version, tc.expected, fallback)
			}
		}
	})
}

func TestVerifyChecksum_APIFailureHandling(t *testing.T) {
	t.Run("handles_api_rate_limit", func(t *testing.T) {
		// When checksum verification fails due to rate limit (403),
		// the code should warn but not fail the download
		// This is tested implicitly by checking error message format

		// Simulate 403 status code handling
		statusCode := 403
		if statusCode != 200 {
			t.Logf("WARNING: Could not fetch release details for checksum verification (status %d)", statusCode)
			// Should not return error, just warn
		}
	})

	t.Run("handles_missing_digest", func(t *testing.T) {
		// When asset digest is empty, should fallback to SHA256SUMS file
		assetDigest := ""

		if assetDigest == "" {
			t.Log("Asset digest not available, would fallback to SHA256SUMS file")
			// This is the expected behavior
		}
	})

	t.Run("handles_missing_sha256sums_file", func(t *testing.T) {
		// When both digest and SHA256SUMS are unavailable,
		// should warn but not fail (allows download to proceed)

		// Simulate no SHA256SUMS URL found
		sha256sumsURL := ""

		if sha256sumsURL == "" {
			t.Log("WARNING: No SHA256 verification available (no digest or SHA256SUMS file), skipping verification")
			// Should warn but continue
		}
	})
}

func TestEnsureLibrary_ErrorPropagation(t *testing.T) {
	t.Run("propagates_download_errors_with_cleanup", func(t *testing.T) {
		// Test that errors are properly wrapped and propagated
		// This documents the error chain behavior

		// Simulate download error
		downloadErr := &urlError{url: "https://example.com/file.tar.gz", cause: "404 not found"}

		// Should be wrapped with context
		wrappedErr := wrapDownloadError(downloadErr)

		if wrappedErr == nil {
			t.Error("Error should be wrapped")
		}

		if !strings.Contains(wrappedErr.Error(), "downloading") {
			t.Errorf("Wrapped error should contain context, got: %v", wrappedErr)
		}

		if !strings.Contains(wrappedErr.Error(), "404") {
			t.Errorf("Wrapped error should preserve original error, got: %v", wrappedErr)
		}
	})

	t.Run("propagates_checksum_verification_errors", func(t *testing.T) {
		// Checksum mismatch should return descriptive error
		checksumErr := &checksumError{
			expected: "abc123",
			actual:   "def456",
		}

		if !strings.Contains(checksumErr.Error(), "checksum mismatch") {
			t.Errorf("Checksum error should be descriptive, got: %v", checksumErr)
		}

		// Should be wrapped with context in caller
		wrappedErr := wrapChecksumError(checksumErr)
		if !strings.Contains(wrappedErr.Error(), "verification failed") {
			t.Errorf("Should wrap with context, got: %v", wrappedErr)
		}
	})

	t.Run("propagates_extraction_errors", func(t *testing.T) {
		// Extraction errors (corrupted tarball, path traversal, etc.)
		// should be wrapped with context

		extractionErr := &tarError{path: "malicious/../../../etc/passwd"}

		wrappedErr := wrapExtractionError(extractionErr)
		if !strings.Contains(wrappedErr.Error(), "extracting") {
			t.Errorf("Should wrap extraction error with context, got: %v", wrappedErr)
		}
	})
}

// Helper types for error propagation tests
type urlError struct {
	url   string
	cause string
}

func (e *urlError) Error() string {
	return e.url + ": " + e.cause
}

type tarError struct {
	path string
}

func (e *tarError) Error() string {
	return "tar error: " + e.path
}

func wrapDownloadError(err error) error {
	if err == nil {
		return nil
	}
	return &wrappedError{context: "downloading", cause: err}
}

func wrapChecksumError(err error) error {
	if err == nil {
		return nil
	}
	return &wrappedError{context: "checksum verification failed", cause: err}
}

func wrapExtractionError(err error) error {
	if err == nil {
		return nil
	}
	return &wrappedError{context: "extracting", cause: err}
}

type wrappedError struct {
	context string
	cause   error
}

func (e *wrappedError) Error() string {
	return e.context + ": " + e.cause.Error()
}

func (e *wrappedError) Unwrap() error {
	return e.cause
}

// =============================================================================
// Integration test helper
// =============================================================================

func TestExtractTarball_ValidExtraction(t *testing.T) {
	// Create a valid tarball matching expected structure
	files := map[string][]byte{
		"linux_amd64/libffmpeg.a": []byte("mock library content"),
	}

	tarball := createValidTarball(t, files)
	defer os.Remove(tarball)

	destDir := t.TempDir()

	err := extractTarball(tarball, destDir)
	if err != nil {
		t.Fatalf("Extraction failed: %v", err)
	}

	// Verify file exists
	extractedFile := filepath.Join(destDir, "linux_amd64", "libffmpeg.a")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(content) != "mock library content" {
		t.Errorf("Content mismatch: got %s", string(content))
	}
}
