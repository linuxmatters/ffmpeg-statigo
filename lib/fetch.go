package lib

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/linuxmatters/ffmpeg-statigo/internal/archiveextract"
)

// Version is the FFmpeg library version (major.minor.patch).
const Version = "8.1.2"

// LibRelease is the exact GitHub Release tag for the pre-built static libraries.
const LibRelease = "lib-8.1.2.1"

// DownloadLibs downloads the FFmpeg static libraries for the current platform
func DownloadLibs() error {
	return ensureLibrary()
}

func validatePlatformArch(platform, arch string) error {
	switch platform {
	case "linux", "darwin":
		if arch == "amd64" || arch == "arm64" {
			return nil
		}
	case "windows":
		if arch == "amd64" || arch == "386" {
			return nil
		}
	default:
		return fmt.Errorf("unsupported platform: %q", platform)
	}
	return fmt.Errorf("unsupported arch %q for platform %q", arch, platform)
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
	if err := validatePlatformArch(platform, arch); err != nil {
		return err
	}

	libDir := "lib"
	platArch := filepath.Join(libraryProfile, platform+"_"+arch)
	libPath := filepath.Join(libDir, platArch, "libffmpeg.a")

	// Library already exists
	if _, err := os.Stat(libPath); err == nil { //nolint:gosec // G703: platform and arch validated against allowlist above
		if err := stripCachedBitcodeSections(libPath, platform, arch); err != nil {
			return fmt.Errorf("stripping bitcode sections: %w", err)
		}
		return nil
	}

	assetPrefix := "ffmpeg-"
	if libraryProfile != "" {
		assetPrefix += libraryProfile + "-"
	}
	tarballName := fmt.Sprintf("%s%s-%s.tar.gz", assetPrefix, platform, arch)
	downloadURL := fmt.Sprintf(
		"https://github.com/linuxmatters/ffmpeg-statigo/releases/download/%s/%s",
		LibRelease, tarballName,
	)

	fmt.Printf("Downloading FFmpeg library %s for %s/%s...\n", LibRelease, platform, arch)

	// Fetch expected checksum before streaming download
	expectedChecksum, err := fetchExpectedChecksum(tarballName)
	if err != nil {
		// Treat a fetch failure as no checksum available, which is fatal downstream.
		expectedChecksum = ""
	}

	// Stage the payload so a wrong-profile archive cannot overwrite another profile.
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		return fmt.Errorf("creating library directory: %w", err)
	}
	stagingDir, err := os.MkdirTemp(libDir, ".download-*")
	if err != nil {
		return fmt.Errorf("creating download directory: %w", err)
	}
	defer os.RemoveAll(stagingDir)

	actualChecksum, err := streamDownloadAndExtract(downloadURL, stagingDir)
	if err != nil {
		return fmt.Errorf("download/extract: %w", err)
	}

	if err := checksumError(expectedChecksum, actualChecksum, tarballName); err != nil {
		return err
	}
	stagedLibrary := filepath.Join(stagingDir, platArch, "libffmpeg.a")
	info, err := os.Stat(stagedLibrary) //nolint:gosec // G703: profile is a build-time constant; platform and arch use allowlists.
	if err != nil {
		return fmt.Errorf("missing library in %s: %w", tarballName, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("invalid library in %s: not a regular file", tarballName)
	}
	if err := stripBitcodeSections(stagedLibrary, platform, arch); err != nil {
		return fmt.Errorf("stripping bitcode sections: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(libPath), 0o755); err != nil { //nolint:gosec // G703: libPath uses only validated platform, arch, and profile.
		return fmt.Errorf("creating library directory: %w", err)
	}
	if err := os.Rename(stagedLibrary, libPath); err != nil { //nolint:gosec // G703: both paths use only validated platform, arch, and profile.
		return fmt.Errorf("installing library: %w", err)
	}
	fmt.Printf("Checksum verified: %s\n", actualChecksum[:8])

	fmt.Printf("Successfully installed FFmpeg library to %s\n", libPath)
	return nil
}

// stripBitcodeSections removes Rust .llvmbc and .llvmcmd LTO sections from a
// native Linux ar archive so Go's cgo linker security check accepts libffmpeg.a.
func stripBitcodeSections(archivePath, platform, arch string) error {
	objcopyBin, arBin, ok := findBitcodeStripTools(platform, arch)
	if !ok {
		return nil
	}
	file, err := os.Open(archivePath) //nolint:gosec // G304: archivePath is a validated staged or existing library path.
	if err != nil {
		return err
	}
	var magic [8]byte
	_, readErr := io.ReadFull(file, magic[:])
	_ = file.Close()
	if readErr != nil || string(magic[:]) != "!<arch>\n" {
		return nil
	}
	return runBitcodeStripTools(archivePath, objcopyBin, arBin)
}

// stripCachedBitcodeSections strips a cached Linux archive via a staged copy in
// the same directory and replaces libPath only after objcopy and ar s succeed.
func stripCachedBitcodeSections(libPath, platform, arch string) error {
	objcopyBin, arBin, ok := findBitcodeStripTools(platform, arch)
	if !ok {
		return nil
	}
	tmpPath, mode, isArchive, err := copyArchiveForStrip(libPath)
	if err != nil || !isArchive {
		return err
	}
	keepTmp := false
	defer func() {
		if !keepTmp {
			_ = os.Remove(tmpPath) //nolint:gosec // G703: tmpPath is created by os.CreateTemp in a validated directory.
		}
	}()

	if err := runBitcodeStripTools(tmpPath, objcopyBin, arBin); err != nil {
		return err
	}
	if err := os.Chmod(tmpPath, mode); err != nil { //nolint:gosec // G703: tmpPath is created by os.CreateTemp in a validated directory.
		return err
	}
	if err := os.Rename(tmpPath, libPath); err != nil { //nolint:gosec // G703: both paths stay inside the validated library directory.
		return err
	}
	keepTmp = true
	return nil
}

func findBitcodeStripTools(platform, arch string) (objcopyBin, arBin string, ok bool) {
	if runtime.GOOS != "linux" || platform != "linux" || arch != runtime.GOARCH {
		return "", "", false
	}
	objcopyBin, err := exec.LookPath("objcopy")
	if err != nil {
		return "", "", false
	}
	arBin, err = exec.LookPath("ar")
	if err != nil {
		return "", "", false
	}
	return objcopyBin, arBin, true
}

func copyArchiveForStrip(libPath string) (string, os.FileMode, bool, error) {
	src, err := os.Open(libPath) //nolint:gosec // G304: libPath is a validated library path.
	if err != nil {
		return "", 0, false, err
	}
	defer src.Close()

	var magic [8]byte
	if _, err := io.ReadFull(src, magic[:]); err != nil || string(magic[:]) != "!<arch>\n" {
		return "", 0, false, nil
	}
	info, err := src.Stat()
	if err != nil {
		return "", 0, false, err
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(libPath), ".strip-*.a")
	if err != nil {
		return "", 0, false, err
	}
	tmpPath := tmpFile.Name()
	keepTmp := false
	defer func() {
		if !keepTmp {
			_ = os.Remove(tmpPath) //nolint:gosec // G703: tmpPath is created by os.CreateTemp in a validated directory.
		}
	}()

	if _, err := tmpFile.Write(magic[:]); err != nil {
		_ = tmpFile.Close()
		return "", 0, false, err
	}
	if _, err := io.Copy(tmpFile, src); err != nil {
		_ = tmpFile.Close()
		return "", 0, false, err
	}
	if err := tmpFile.Close(); err != nil {
		return "", 0, false, err
	}

	keepTmp = true
	return tmpPath, info.Mode().Perm(), true, nil
}

func runBitcodeStripTools(archivePath, objcopyBin, arBin string) error {
	output, err := exec.CommandContext(context.Background(), objcopyBin, "--remove-section=.llvmbc", "--remove-section=.llvmcmd", archivePath).CombinedOutput() //nolint:gosec // G204: objcopyBin is resolved from LookPath and archivePath is validated.
	if err != nil {
		return fmt.Errorf("objcopy: %w: %s", err, strings.TrimSpace(string(output)))
	}
	output, err = exec.CommandContext(context.Background(), arBin, "s", archivePath).CombinedOutput() //nolint:gosec // G204: arBin is resolved from LookPath and archivePath is validated.
	if err != nil {
		return fmt.Errorf("ar: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Digest             string `json:"digest"` // SHA256 digest in format "sha256:..."
}

type GitHubReleaseDetail struct {
	Assets []GitHubAsset `json:"assets"`
}

func newGitHubAPIRequest(apiURL string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "ffmpeg-statigo")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	if token := githubToken(); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req, nil
}

func githubToken() string {
	for _, name := range []string{"GITHUB_TOKEN", "GH_TOKEN"} {
		if token := strings.TrimSpace(os.Getenv(name)); token != "" {
			return token
		}
	}
	return ""
}

// fetchReleaseDetails retrieves asset metadata from GitHub API for LibRelease.
func fetchReleaseDetails() (*GitHubReleaseDetail, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/linuxmatters/ffmpeg-statigo/releases/tags/%s", LibRelease)

	req, err := newGitHubAPIRequest(apiURL)
	if err != nil {
		return nil, err
	}

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

	progressBody := &progressReader{
		reader: resp.Body,
		total:  resp.ContentLength,
	}

	hasher := sha256.New()

	// TeeReader: data flows to both hasher and gzip decompressor simultaneously
	teeReader := io.TeeReader(progressBody, hasher)

	gzr, err := gzip.NewReader(teeReader)
	if err != nil {
		return "", fmt.Errorf("gzip reader: %w", err)
	}
	defer gzr.Close()

	if err := archiveextract.ExtractTar(gzr, archiveextract.TarOptions{
		DestDir:          absDestDir,
		LinkPolicy:       archiveextract.SkipLinks,
		RemoveIncomplete: true,
		OnError: func() {
			fmt.Println() // Clear progress line
		},
	}); err != nil {
		return "", err
	}

	fmt.Println() // Clear progress line
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// checksumError reports whether an extracted library must be rejected for
// lack of, or a mismatched, checksum. A nil result means installation proceeds.
func checksumError(expected, actual, tarball string) error {
	switch {
	case expected == "":
		return fmt.Errorf("no checksum available to verify %s; refusing to install unverified", tarball)
	case actual != expected:
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actual)
	default:
		return nil
	}
}

// fetchExpectedChecksum retrieves the expected SHA256 checksum for a tarball from GitHub.
// It tries the asset digest first (newer releases), then falls back to SHA256SUMS file.
func fetchExpectedChecksum(tarballName string) (string, error) {
	releaseDetail, err := fetchReleaseDetails()
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
		parts := strings.Fields(line)
		if len(parts) == 2 && strings.TrimPrefix(parts[1], "*") == tarballName {
			return parts[0], nil
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

	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return fmt.Errorf("resolving destination directory: %w", err)
	}

	return archiveextract.ExtractTar(gzr, archiveextract.TarOptions{
		DestDir:          absDestDir,
		LinkPolicy:       archiveextract.SkipLinks,
		RemoveIncomplete: true,
	})
}
