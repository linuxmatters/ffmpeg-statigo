package lib

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

// Version is the FFmpeg library version (major.minor.patch)
// The downloader will find the latest internal release (e.g., lib-8.0.1.0)
const Version = "8.0.1"

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

	// Use working directory for libraries (writable)
	// Libraries will be downloaded to ./ffmpeg-libs/<platform>_<arch>/
	// Libraries will be downloaded to lib/<platform>_<arch>/
	libDir := "lib"
	platArch := platform + "_" + arch
	libPath := filepath.Join(libDir, platArch, "libffmpeg.a")

	// Library already exists
	if _, err := os.Stat(libPath); err == nil {
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

	// Use unique temp file to avoid collision in high-concurrency CI/CD scenarios
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("ffmpeg-%s-%s-*.tar.gz", platform, arch))
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpTarball := tmpFile.Name()
	tmpFile.Close()

	fmt.Printf("Downloading FFmpeg library %s for %s/%s...\n", release, platform, arch)

	if err := downloadFile(downloadURL, tmpTarball); err != nil {
		os.Remove(tmpTarball)
		return fmt.Errorf("downloading: %w", err)
	}
	defer os.Remove(tmpTarball)

	// Verify checksum using GitHub's SHA256 checksums file
	if err := verifyChecksum(tmpTarball, release, tarballName); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}

	// Extract to lib directory
	if err := extractTarball(tmpTarball, libDir); err != nil {
		return fmt.Errorf("extracting: %w", err)
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

	req, err := http.NewRequest("GET", apiURL, nil)
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
	if resp.StatusCode == 403 {
		return "", fmt.Errorf("GitHub API rate limit exceeded")
	}

	if resp.StatusCode != 200 {
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
	sort.Strings(matchingReleases)

	// Return the last (highest) version
	return matchingReleases[len(matchingReleases)-1], nil
}

func downloadFile(url, dest string) error {
	client := grab.NewClient()
	req, err := grab.NewRequest(dest, url)
	if err != nil {
		return err
	}

	resp := client.Do(req)

	// Show progress for large downloads (~100MB libraries)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("\rDownloading... %.2f%%", resp.Progress()*100)
		case <-resp.Done:
			fmt.Println() // New line after progress
			return resp.Err()
		}
	}
}

type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Digest             string `json:"digest"` // SHA256 digest in format "sha256:..."
}

type GitHubReleaseDetail struct {
	Assets []GitHubAsset `json:"assets"`
}

func verifyChecksum(file, release, tarballName string) error {
	// Calculate file checksum
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	actualChecksum := hex.EncodeToString(h.Sum(nil))

	// Fetch the release details to get the digest
	apiURL := fmt.Sprintf("https://api.github.com/repos/linuxmatters/ffmpeg-statigo/releases/tags/%s", release)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "ffmpeg-statigo")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// If we can't verify, warn but don't fail (might be rate limited)
		fmt.Fprintf(os.Stderr, "WARNING: Could not fetch release details for checksum verification (status %d)\n", resp.StatusCode)
		return nil
	}

	var releaseDetail GitHubReleaseDetail
	if err := json.NewDecoder(resp.Body).Decode(&releaseDetail); err != nil {
		return err
	}

	// Find our tarball asset and get its digest
	var assetDigest string
	for _, asset := range releaseDetail.Assets {
		if asset.Name == tarballName {
			assetDigest = asset.Digest
			break
		}
	}

	if assetDigest == "" {
		// Fallback to SHA256SUMS file if digest not available (older releases)
		return verifyChecksumFromFile(releaseDetail.Assets, actualChecksum, tarballName)
	}

	// GitHub provides digests in "sha256:..." format
	if !strings.HasPrefix(assetDigest, "sha256:") {
		return fmt.Errorf("unexpected digest format: %s", assetDigest)
	}

	expectedChecksum := strings.TrimPrefix(assetDigest, "sha256:")
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	fmt.Printf("Checksum verified: %s\n", actualChecksum[:8])
	return nil
}

func verifyChecksumFromFile(assets []GitHubAsset, actualChecksum, tarballName string) error {
	// Find and download SHA256SUMS file (fallback for older releases)
	var sha256sumsURL string
	for _, asset := range assets {
		if asset.Name == "SHA256SUMS" {
			sha256sumsURL = asset.BrowserDownloadURL
			break
		}
	}

	if sha256sumsURL == "" {
		fmt.Fprintf(os.Stderr, "WARNING: No SHA256 verification available (no digest or SHA256SUMS file), skipping verification\n")
		return nil
	}

	// Download SHA256SUMS
	resp, err := http.Get(sha256sumsURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse checksums
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Find our file's checksum
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, tarballName) {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				expectedChecksum := parts[0]
				if actualChecksum != expectedChecksum {
					return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
				}
				fmt.Printf("Checksum verified: %s\n", actualChecksum[:8])
				return nil
			}
		}
	}

	fmt.Fprintf(os.Stderr, "WARNING: Could not find checksum for %s in SHA256SUMS\n", tarballName)
	return nil
}

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

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			outFile, err := os.Create(target)
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}
