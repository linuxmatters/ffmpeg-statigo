package main

import (
	"archive/tar"
	"cmp"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

const (
	repoURL             = "https://github.com/linuxmatters/ffmpeg-statigo"
	maxLicenseFileBytes = 1 << 20 // 1 MiB
)

type builderOptions struct {
	embedded      bool
	programs      bool
	clean         bool
	list          bool
	updateDigests bool
	selectedLibs  map[string]bool
}

func parseBuilderArgs(args []string) builderOptions {
	opts := builderOptions{
		selectedLibs: make(map[string]bool),
	}
	for _, arg := range args {
		switch arg {
		case "--embedded":
			opts.embedded = true
		case "--programs":
			opts.programs = true
		case "--clean":
			opts.clean = true
		case "--list":
			opts.list = true
		case "--update-digests":
			opts.updateDigests = true
		default:
			opts.selectedLibs[arg] = true
		}
	}
	return opts
}

var libraryLicenses = map[string]string{
	"Vulkan-Headers":   "Apache-2.0",
	"dav1d":            "BSD-2-Clause",
	"ffmpeg":           "GPL-3.0-or-later",
	"glslang":          "BSD-3-Clause",
	"lame":             "LGPL-2.0-or-later",
	"libdrm":           "MIT",
	"libiconv":         "LGPL-2.1-or-later",
	"libsrt":           "MPL-2.0",
	"libva":            "MIT",
	"libvpl":           "MIT",
	"libvpx":           "BSD-3-Clause",
	"libwebp":          "BSD-3-Clause",
	"libxml2":          "MIT",
	"nv-codec-headers": "MIT",
	"openh264":         "BSD-2-Clause",
	"openssl":          "Apache-2.0",
	"opus":             "BSD-3-Clause",
	"rav1e":            "BSD-2-Clause",
	"vvenc":            "BSD-3-Clause-Clear",
	"x264":             "GPL-2.0-or-later",
	"x265":             "GPL-2.0-or-later",
	"zimg":             "WTFPL",
	"zlib":             "Zlib",
}

func libraryLicense(name string, embedded bool) (string, error) {
	if name == "ffmpeg" {
		if embedded {
			return "LGPL-2.1-or-later", nil
		}
		return "GPL-3.0-or-later", nil
	}
	license, ok := libraryLicenses[name]
	if !ok {
		return "", fmt.Errorf("unknown license for library %q", name)
	}
	return license, nil
}

func ffmpegTagFromURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}
	p := rawURL
	if parsed, err := url.Parse(rawURL); err == nil && parsed.Path != "" {
		p = parsed.Path
	}
	if strings.HasSuffix(p, "/") {
		return ""
	}
	name := path.Base(p)
	if name == "." || name == "/" {
		return ""
	}
	lower := strings.ToLower(name)
	for _, ext := range []string{".tar.gz", ".tar.bz2", ".tar.xz", ".tgz", ".tbz2", ".txz", ".zip"} {
		if strings.HasSuffix(lower, ext) {
			return name[:len(name)-len(ext)]
		}
	}
	return name
}

func licenseFilenameForProfile(embedded bool) string {
	if embedded {
		return "COPYING.LGPLv2.1"
	}
	return "COPYING.GPLv3"
}

func releaseTagFromEnv() string {
	if tag := strings.TrimSpace(os.Getenv("FFMPEG_RELEASE_TAG")); tag != "" {
		return tag
	}
	return strings.TrimSpace(os.Getenv("GITHUB_REF_NAME"))
}

func stageLicense(buildRoot, stagingDir, binDir string, embedded bool) error {
	licenseFile := licenseFilenameForProfile(embedded)
	data, err := readFFmpegLicenseBytes(buildRoot, stagingDir, licenseFile)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return fmt.Errorf("create programs output directory %s: %w", binDir, err)
	}
	dest := filepath.Join(binDir, "LICENSE")
	if err := os.WriteFile(dest, data, 0o644); err != nil { //nolint:gosec // G703: dest is a fixed LICENSE filename inside binDir
		return fmt.Errorf("write %s: %w", dest, err)
	}
	return nil
}

func readFFmpegLicenseBytes(buildRoot, stagingDir, licenseFile string) ([]byte, error) {
	candidates := []string{
		filepath.Join(buildRoot, "src", "ffmpeg", licenseFile),
		filepath.Join(stagingDir, "share", "ffmpeg", licenseFile),
	}
	for _, candidate := range candidates {
		data, err := os.ReadFile(candidate) //nolint:gosec // G304: candidate is a fixed FFmpeg license path under buildRoot or stagingDir
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("read %s: %w", candidate, err)
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("license file %s is empty", candidate)
		}
		return data, nil
	}

	archivePath := filepath.Join(buildRoot, "downloads", filepath.Base(ffmpeg.URL))
	if _, err := os.Stat(archivePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("missing %s: checked %s, %s, and archive %s", licenseFile, candidates[0], candidates[1], archivePath)
		}
		return nil, fmt.Errorf("stat %s: %w", archivePath, err)
	}
	if err := verifyDigest(ffmpeg.URL, archivePath); err != nil {
		return nil, fmt.Errorf("verify FFmpeg archive before extracting %s: %w", licenseFile, err)
	}
	return extractLicenseFromTarGz(archivePath, licenseFile)
}

func extractLicenseFromTarGz(archivePath, licenseFile string) ([]byte, error) {
	file, err := os.Open(archivePath) //nolint:gosec // G304: archivePath is the SHA-256-verified FFmpeg download archive
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("open gzip stream %s: %w", archivePath, err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("%s not found in %s", licenseFile, archivePath)
		}
		if err != nil {
			return nil, fmt.Errorf("read tar %s: %w", archivePath, err)
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		clean := path.Clean(strings.TrimPrefix(header.Name, "./"))
		if !isTopLevelArchiveFile(clean, licenseFile) {
			continue
		}
		data, err := io.ReadAll(io.LimitReader(tarReader, maxLicenseFileBytes+1))
		if err != nil {
			return nil, fmt.Errorf("read %s from %s: %w", header.Name, archivePath, err)
		}
		if len(data) > maxLicenseFileBytes {
			return nil, fmt.Errorf("%s in %s exceeds %d bytes", header.Name, archivePath, maxLicenseFileBytes)
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("%s in %s is empty", header.Name, archivePath)
		}
		return data, nil
	}
}

func isTopLevelArchiveFile(cleanName, targetFile string) bool {
	if cleanName == "." || cleanName == ".." || strings.HasPrefix(cleanName, "/") || strings.HasPrefix(cleanName, "../") {
		return false
	}
	if cleanName == targetFile {
		return true
	}
	dir, base := path.Split(cleanName)
	dir = strings.TrimSuffix(dir, "/")
	return base == targetFile && dir != "" && !strings.Contains(dir, "/")
}

func isLibraryActiveForTarget(lib *Library, targetOS string) bool {
	if lib == nil {
		return false
	}
	if lib.Enabled != nil && !*lib.Enabled {
		return false
	}
	if len(lib.Platform) != 0 && !slices.Contains(lib.Platform, targetOS) {
		return false
	}
	return true
}

func renderNotice(embedded bool, targetOS, arch, releaseTag string, libs []*Library) (string, error) {
	profileName := "default"
	effectiveLicense := "GPL-3.0-or-later"
	if embedded {
		profileName = "embedded"
		effectiveLicense = "LGPL-2.1-or-later"
	}
	licenseFile := licenseFilenameForProfile(embedded)

	ffmpegURL := ffmpeg.URL
	var active []*Library
	for _, lib := range libs {
		if !isLibraryActiveForTarget(lib, targetOS) {
			continue
		}
		if lib.Name == "ffmpeg" && lib.URL != "" {
			ffmpegURL = lib.URL
		}
		active = append(active, lib)
	}
	slices.SortFunc(active, func(a, b *Library) int {
		return cmp.Compare(a.Name, b.Name)
	})

	var b strings.Builder
	b.WriteString("FFmpeg Static Binaries (ffmpeg, ffprobe)\n")
	fmt.Fprintf(&b, "Profile: %s\n", profileName)
	fmt.Fprintf(&b, "Target: %s/%s\n", targetOS, arch)
	fmt.Fprintf(&b, "FFmpeg: %s\n", ffmpegTagFromURL(ffmpegURL))
	fmt.Fprintf(&b, "License: %s (see LICENSE / %s)\n", effectiveLicense, licenseFile)
	fmt.Fprintf(&b, "Repository: %s\n", repoURL)
	if tag := strings.TrimSpace(releaseTag); tag != "" {
		fmt.Fprintf(&b, "Release: %s\n", tag)
	}
	b.WriteString("\n")

	if embedded {
		archiveName := fmt.Sprintf("ffmpeg-embedded-%s-%s.tar.gz", targetOS, arch)
		b.WriteString("LGPLv2.1 Section 4 Notice and Relinking Instructions:\n")
		b.WriteString("These binaries are distributed under the GNU Lesser General Public License,\n")
		b.WriteString("version 2.1 or later (LGPL-2.1-or-later). The complete license text is in\n")
		fmt.Fprintf(&b, "the accompanying LICENSE file (%s).\n\n", licenseFile)
		b.WriteString("To relink with a modified version of FFmpeg or another LGPL library:\n")
		fmt.Fprintf(&b, "1. Download the matching static library archive %s from the\n", archiveName)
		fmt.Fprintf(&b, "   release at %s, or build the embedded profile from\n", repoURL)
		b.WriteString("   source using the pinned archives listed below.\n")
		b.WriteString("2. Replace the modified LGPL library and rebuild or relink the binaries using\n")
		fmt.Fprintf(&b, "   the build scripts at %s\n", repoURL)
		b.WriteString("   (go run ./internal/builder --embedded --programs).\n")
	} else {
		b.WriteString("GPLv3 Section 6 Notice:\n")
		b.WriteString("These binaries are distributed under the GNU General Public License, version 3\n")
		b.WriteString("or later (GPL-3.0-or-later). The complete license text is in the accompanying\n")
		fmt.Fprintf(&b, "LICENSE file (%s). Corresponding source code and build scripts are\n", licenseFile)
		fmt.Fprintf(&b, "available at %s and from the pinned upstream source\n", repoURL)
		b.WriteString("archives listed below.\n")
	}

	b.WriteString("\nBundled Libraries:\n")
	for _, lib := range active {
		lic, err := libraryLicense(lib.Name, embedded)
		if err != nil {
			return "", err
		}
		digest, ok := expectedDigest(lib.URL)
		if !ok {
			return "", fmt.Errorf("missing pinned SHA-256 digest for %s (%s)", lib.Name, lib.URL)
		}
		fmt.Fprintf(&b, "- %s | %s | %s | sha256:%s\n", lib.Name, lic, lib.URL, digest)
	}

	return b.String(), nil
}

func stagePrograms(ctx context.Context, buildRoot, stagingDir, binDir string, embedded bool, targetOS, arch, releaseTag string, libs []*Library) error {
	return stageProgramsWithRunner(ctx, buildRoot, stagingDir, binDir, embedded, targetOS, arch, releaseTag, libs, nil)
}

func stageProgramsWithRunner(ctx context.Context, buildRoot, stagingDir, binDir string, embedded bool, targetOS, arch, releaseTag string, libs []*Library, run commandRunner) error {
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return fmt.Errorf("create programs output directory %s: %w", binDir, err)
	}

	for _, name := range []string{"ffmpeg", "ffprobe"} {
		filename := binaryFilename(name, targetOS)
		srcPath := filepath.Join(stagingDir, "bin", filename)
		dstPath := filepath.Join(binDir, filename)

		if err := copyExecutable(srcPath, dstPath); err != nil {
			return err
		}
		if run == nil {
			if err := stripAndSignBinary(ctx, dstPath, targetOS); err != nil {
				return err
			}
		} else {
			if err := stripAndSignBinaryWithRunner(ctx, dstPath, targetOS, run); err != nil {
				return err
			}
		}
	}

	if err := stageLicense(buildRoot, stagingDir, binDir, embedded); err != nil {
		return fmt.Errorf("stage license: %w", err)
	}

	notice, err := renderNotice(embedded, targetOS, arch, releaseTag, libs)
	if err != nil {
		return fmt.Errorf("render notice: %w", err)
	}
	noticePath := filepath.Join(binDir, "NOTICE.txt")
	if err := os.WriteFile(noticePath, []byte(notice), 0o644); err != nil { //nolint:gosec // G703: noticePath is a fixed NOTICE.txt filename inside binDir
		return fmt.Errorf("write %s: %w", noticePath, err)
	}

	return nil
}

func copyExecutable(srcPath, dstPath string) error {
	info, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("stat staged binary %s: %w", srcPath, err)
	}
	if !info.Mode().IsRegular() || info.Size() == 0 {
		return fmt.Errorf("expected non-empty regular binary %s", srcPath)
	}

	src, err := os.Open(srcPath) //nolint:gosec // G304: srcPath is a staged FFmpeg executable under stagingDir/bin
	if err != nil {
		return fmt.Errorf("open staged binary %s: %w", srcPath, err)
	}
	defer src.Close()

	// Remove any existing destination binary first so macOS allocates a fresh
	// inode and never retains a stale Mach-O vnode code-signature cache.
	if err := os.Remove(dstPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove existing binary %s: %w", dstPath, err)
	}

	dst, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o755) //nolint:gosec // G302,G304,G703: staged executable requires 0o755 permissions inside binDir
	if err != nil {
		return fmt.Errorf("create binary %s: %w", dstPath, err)
	}

	_, copyErr := io.Copy(dst, src)
	closeErr := dst.Close()
	if copyErr != nil {
		_ = os.Remove(dstPath)
		return fmt.Errorf("copy binary %s to %s: %w", srcPath, dstPath, copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(dstPath)
		return fmt.Errorf("close binary %s: %w", dstPath, closeErr)
	}

	return nil
}
