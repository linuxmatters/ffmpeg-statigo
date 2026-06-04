package lib

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/linuxmatters/ffmpeg-statigo/internal/pathsafe"
)

// Version is the FFmpeg library version (major.minor.patch)
// The downloader will find the latest internal release (e.g., lib-8.1.1.0)
const Version = "8.1.1"

// DownloadLibs downloads the FFmpeg static libraries for the current platform
func DownloadLibs() error {
	return ensureLibrary()
}

func ensureLibrary() error {
	// Support cross-compilation: use GOOS/GOARCH env vars if set
	platform := os.Getenv("GOOS")
	if platform == "" {
		platform = runtime.GOOS
	}

	arch := os.Getenv("GOARCH")
	if arch == "" {
		arch = runtime.GOARCH
	}

	// Restrict platform/arch to known values; these feed filesystem paths below.
	if !slices.Contains([]string{"linux", "darwin"}, platform) {
		return fmt.Errorf("unsupported platform: %s", platform)
	}
	if !slices.Contains([]string{"amd64", "arm64"}, arch) {
		return fmt.Errorf("unsupported arch: %s", arch)
	}

	// Use working directory for libraries (writable)
	// Libraries will be downloaded to lib/<platform>_<arch>/
	libDir := "lib"
	platArch := platform + "_" + arch
	libPath := filepath.Join(libDir, platArch, "libffmpeg.a")

	// Library already exists
	if _, err := os.Stat(libPath); err == nil { //nolint:gosec // G703: platform and arch validated against allowlist above
		return nil
	}

	// Determine latest compatible release
	release, err := findCompatibleRelease(Version)
	if err != nil {
		return fmt.Errorf("finding release: %w", err)
	}

	// Download tarball
	tarballName := fmt.Sprintf("ffmpeg-%s-%s.tar.gz", platform, arch)
	downloadURL := fmt.Sprintf(
		"https://github.com/linuxmatters/ffmpeg-statigo/releases/download/%s/%s",
		release, tarballName,
	)

	fmt.Printf("Downloading FFmpeg library %s for %s/%s...\n", release, platform, arch)

	// Fetch expected checksum before streaming download
	allowUnverified, _ := strconv.ParseBool(os.Getenv("FFMPEG_STATIGO_ALLOW_UNVERIFIED"))
	expectedChecksum, err := fetchExpectedChecksum(release, tarballName)
	if err != nil {
		// Treat a fetch failure as no checksum available; handled after extraction.
		if allowUnverified {
			fmt.Fprintf(os.Stderr, "WARNING: Could not fetch checksum: %v\n", err)
		}
		expectedChecksum = ""
	}

	// Stream download directly to extraction with concurrent checksum verification
	actualChecksum, err := streamDownloadAndExtract(downloadURL, libDir)
	if err != nil {
		// Clean up any partially extracted files so the next run re-downloads.
		_ = os.RemoveAll(filepath.Join(libDir, platArch)) //nolint:gosec // G703: platform and arch validated against allowlist above
		return fmt.Errorf("download/extract: %w", err)
	}

	// Verify checksum, refusing to install unverified libraries by default.
	if err := checksumError(expectedChecksum, actualChecksum, tarballName, allowUnverified); err != nil {
		// Clean up extracted files: the library is unverified or mismatched.
		_ = os.RemoveAll(filepath.Join(libDir, platArch)) //nolint:gosec // G703: platform and arch validated against allowlist above
		return err
	}
	if expectedChecksum == "" {
		fmt.Fprintf(os.Stderr, "WARNING: No checksum available for %s, installing unverified (FFMPEG_STATIGO_ALLOW_UNVERIFIED set)\n", tarballName)
	} else {
		fmt.Printf("Checksum verified: %s\n", actualChecksum[:8])
	}

	fmt.Printf("Successfully installed FFmpeg library to %s\n", libPath)
	return nil
}

func findCompatibleRelease(moduleVersion string) (string, error) {
	// Parse major.minor.patch from module version
	parts := strings.Split(moduleVersion, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid version format: %s", moduleVersion)
	}
	prefix := "lib-" + moduleVersion // "lib-8.0.0"

	// Try GitHub API first
	release, err := findViaAPI(prefix)
	if err == nil {
		return release, nil
	}

	// Fallback to predictable pattern if API fails (rate limit, network issue)
	// Assumes consistent lib-X.Y.Z.0 pattern for initial releases
	fmt.Fprintf(os.Stderr, "GitHub API unavailable, using fallback release pattern: %v\n", err)
	return fmt.Sprintf("lib-%s.0", moduleVersion), nil
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

func findViaAPI(prefix string) (string, error) {
	// Query GitHub API for releases
	apiURL := "https://api.github.com/repos/linuxmatters/ffmpeg-statigo/releases?per_page=100"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiURL, nil)
	if err != nil {
		return "", err
	}

	// GitHub recommends setting User-Agent
	req.Header.Set("User-Agent", "ffmpeg-statigo")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for rate limiting
	if resp.StatusCode == http.StatusForbidden {
		return "", fmt.Errorf("GitHub API rate limit exceeded")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	// Parse releases
	var releases []GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", fmt.Errorf("parsing releases: %w", err)
	}

	// Find all tags matching our prefix
	var matchingReleases []string
	for _, rel := range releases {
		if strings.HasPrefix(rel.TagName, prefix) {
			matchingReleases = append(matchingReleases, rel.TagName)
		}
	}

	if len(matchingReleases) == 0 {
		return "", fmt.Errorf("no releases found matching %s", prefix)
	}

	// Sort to find highest version (lib-8.0.0.0 < lib-8.0.0.1 < lib-8.0.0.3)
	slices.Sort(matchingReleases)

	// Return the last (highest) version
	return matchingReleases[len(matchingReleases)-1], nil
}

type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Digest             string `json:"digest"` // SHA256 digest in format "sha256:..."
}

type GitHubReleaseDetail struct {
	Assets []GitHubAsset `json:"assets"`
}

// fetchReleaseDetails retrieves asset metadata from GitHub API for a release.
func fetchReleaseDetails(release string) (*GitHubReleaseDetail, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/linuxmatters/ffmpeg-statigo/releases/tags/%s", release)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ffmpeg-statigo")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var detail GitHubReleaseDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// progressReader wraps an io.Reader to report download progress.
type progressReader struct {
	reader      io.Reader
	total       int64 // Total bytes expected (-1 if unknown)
	read        int64 // Bytes read so far
	lastPercent int   // Last reported percentage (to avoid spam)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.read += int64(n)

	if pr.total > 0 {
		percent := int(pr.read * 100 / pr.total)
		// Report every 10% to avoid flooding output
		if percent/10 > pr.lastPercent/10 {
			fmt.Printf("\rDownloading: %d%% (%d/%d MB)", percent, pr.read/(1024*1024), pr.total/(1024*1024))
			pr.lastPercent = percent
		}
	} else if pr.read%(10*1024*1024) == 0 {
		// Unknown size: report every 10MB
		fmt.Printf("\rDownloading: %d MB", pr.read/(1024*1024))
	}

	return n, err
}

// streamDownloadAndExtract downloads a tarball and extracts it in a single streaming pass.
// It returns the SHA256 checksum of the downloaded data for verification.
// This eliminates the need for temporary files and reduces total time by ~40%.
func streamDownloadAndExtract(url, destDir string) (string, error) {
	// URL is built from a fixed GitHub host and a validated release tag, not user input.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil) //nolint:gosec // G107: trusted, project-controlled release URL
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req) //nolint:gosec // G704: trusted, project-controlled release URL
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Resolve destDir to absolute path for security checks
	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return "", fmt.Errorf("resolving destination directory: %w", err)
	}

	// Wrap response body with progress reporting
	progressBody := &progressReader{
		reader: resp.Body,
		total:  resp.ContentLength,
	}

	// Create a hash writer to calculate checksum while streaming
	hasher := sha256.New()

	// TeeReader: data flows to both hasher and gzip decompressor simultaneously
	teeReader := io.TeeReader(progressBody, hasher)

	// Decompress gzip stream
	gzr, err := gzip.NewReader(teeReader)
	if err != nil {
		return "", fmt.Errorf("gzip reader: %w", err)
	}
	defer gzr.Close()

	// Extract tar entries
	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println() // Clear progress line
			return "", fmt.Errorf("reading tar header: %w", err)
		}

		// Security: Validate path to prevent path traversal attacks
		target, err := pathsafe.SanitizePath(absDestDir, header.Name)
		if err != nil {
			fmt.Println() // Clear progress line
			return "", err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				fmt.Println() // Clear progress line
				return "", fmt.Errorf("creating directory %s: %w", target, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				fmt.Println() // Clear progress line
				return "", fmt.Errorf("creating parent directory for %s: %w", target, err)
			}

			if err := extractFile(tr, target); err != nil {
				fmt.Println() // Clear progress line
				return "", err
			}
		case tar.TypeSymlink, tar.TypeLink:
			// Skip symlinks and hard links for security - they could point outside destDir
			continue
		}
	}

	fmt.Println() // Clear progress line
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// extractFile extracts a single file from a tar reader to the target path.
func extractFile(tr *tar.Reader, target string) error {
	outFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", target, err)
	}
	defer outFile.Close()

	if err := pathsafe.CopyCapped(outFile, tr); err != nil {
		_ = outFile.Close()
		_ = os.Remove(target)
		return fmt.Errorf("writing file %s: %w", target, err)
	}

	return nil
}

// checksumError reports whether an extracted library must be rejected for
// lack of, or a mismatched, checksum. A nil result means installation proceeds.
func checksumError(expected, actual, tarball string, allowUnverified bool) error {
	switch {
	case expected == "":
		if allowUnverified {
			return nil
		}
		return fmt.Errorf("no checksum available to verify %s; set FFMPEG_STATIGO_ALLOW_UNVERIFIED=1 to install without verification", tarball)
	case actual != expected:
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actual)
	default:
		return nil
	}
}

// fetchExpectedChecksum retrieves the expected SHA256 checksum for a tarball from GitHub.
// It tries the asset digest first (newer releases), then falls back to SHA256SUMS file.
func fetchExpectedChecksum(release, tarballName string) (string, error) {
	releaseDetail, err := fetchReleaseDetails(release)
	if err != nil {
		return "", err
	}

	// Try asset digest first (newer releases)
	for _, asset := range releaseDetail.Assets {
		if asset.Name == tarballName && asset.Digest != "" {
			if !strings.HasPrefix(asset.Digest, "sha256:") {
				return "", fmt.Errorf("unexpected digest format: %s", asset.Digest)
			}
			return strings.TrimPrefix(asset.Digest, "sha256:"), nil
		}
	}

	// Fallback to SHA256SUMS file for older releases
	return fetchChecksumFromFile(releaseDetail.Assets, tarballName)
}

// fetchChecksumFromFile downloads and parses the SHA256SUMS file to find a checksum.
func fetchChecksumFromFile(assets []GitHubAsset, tarballName string) (string, error) {
	var sha256sumsURL string
	for _, asset := range assets {
		if asset.Name == "SHA256SUMS" {
			sha256sumsURL = asset.BrowserDownloadURL
			break
		}
	}

	if sha256sumsURL == "" {
		return "", nil // No checksum available
	}

	// URL comes from the GitHub release asset list, not user input.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, sha256sumsURL, nil) //nolint:gosec // G107: trusted GitHub release asset URL
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	for line := range strings.SplitSeq(string(content), "\n") {
		if strings.Contains(line, tarballName) {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[0], nil
			}
		}
	}

	return "", nil // Checksum not found in file
}

// extractTarball extracts a gzipped tarball to a destination directory.
// This function is used for testing path traversal protection.
// Production code uses streamDownloadAndExtract which combines download and extraction.
func extractTarball(tarball, destDir string) error {
	f, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return fmt.Errorf("resolving destination directory: %w", err)
	}

	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		target, err := pathsafe.SanitizePath(absDestDir, header.Name)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			if err := extractFile(tr, target); err != nil {
				return err
			}
		case tar.TypeSymlink, tar.TypeLink:
			continue // Skip symlinks for security
		}
	}

	return nil
}
