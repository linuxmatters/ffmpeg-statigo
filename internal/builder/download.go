package main

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/ulikunitz/xz"
)

// DownloadFile downloads a file using the grab library with resume support and retries
func DownloadFile(url, dest string, logger io.Writer) error {
	// Create downloads directory
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	// Check if file already exists and is complete
	if fileExists(dest) {
		fmt.Fprintf(logger, "Using cached %s\n", filepath.Base(dest))
		return nil
	}

	fmt.Fprintf(logger, "Downloading %s...\n", filepath.Base(dest))

	// Retry logic: 3 attempts with exponential backoff
	maxRetries := 3
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			backoff := time.Duration(attempt-1) * 5 * time.Second
			fmt.Fprintf(logger, "Retry %d/%d after %v...\n", attempt-1, maxRetries-1, backoff)
			time.Sleep(backoff)
		}

		// Create download client
		client := grab.NewClient()
		req, err := grab.NewRequest(dest, url)
		if err != nil {
			return err
		}

		// Start download
		resp := client.Do(req)

		// Monitor progress - update every 2 seconds to avoid log spam
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		lastProgress := 0.0
		for !resp.IsComplete() {
			select {
			case <-ticker.C:
				progress := resp.Progress() * 100
				if progress > lastProgress {
					fmt.Fprintf(logger, "  %.2f%% complete\n", progress)
					lastProgress = progress
				}
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}

		// Check for errors
		if err := resp.Err(); err != nil {
			lastErr = err
			fmt.Fprintf(logger, "Download attempt %d failed: %v\n", attempt, err)
			// Remove partial download
			os.Remove(dest)
			continue
		}

		fmt.Fprintf(logger, "  100.00%% complete\n")
		return nil
	}

	return fmt.Errorf("download failed after %d attempts: %w", maxRetries, lastErr)
}

// detectTarPrefix gets the prefix from the first entry in the tar archive
func detectTarPrefix(archivePath, archiveType string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var reader io.Reader = file

	// Handle compression
	switch {
	case strings.Contains(archiveType, ".gz"):
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return "", err
		}
		defer gzReader.Close()
		reader = gzReader

	case strings.Contains(archiveType, ".bz2"):
		reader = bzip2.NewReader(file)

	case strings.Contains(archiveType, ".xz"):
		xzReader, err := xz.NewReader(file)
		if err != nil {
			return "", err
		}
		reader = xzReader
	}

	tarReader := tar.NewReader(reader)

	// Find first real entry (skip PAX headers and other metadata)
	for {
		header, err := tarReader.Next()
		if err != nil {
			return "", err
		}

		// Skip PAX headers and other non-content entries
		if header.Typeflag == tar.TypeXGlobalHeader || header.Typeflag == tar.TypeXHeader {
			continue
		}

		name := header.Name

		// Skip if empty or contains .. (security)
		if name == "" || strings.Contains(name, "..") {
			continue
		}

		// Get the first directory component
		parts := strings.Split(strings.TrimSuffix(name, "/"), "/")
		if len(parts) == 0 {
			continue
		}

		return parts[0] + "/", nil
	}
}

// detectZipPrefix gets the prefix from the first entry in the zip archive
func detectZipPrefix(reader *zip.ReadCloser) string {
	if len(reader.File) == 0 {
		return ""
	}

	// Get first entry
	name := reader.File[0].Name

	// Skip if empty or contains .. (security)
	if name == "" || strings.Contains(name, "..") {
		return ""
	}

	// Get the first directory component
	parts := strings.Split(strings.TrimSuffix(name, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	return parts[0] + "/"
}

// ExtractArchive extracts an archive to the destination path
// Automatically detects and strips common parent directory
func ExtractArchive(archivePath, destPath, archiveType string, logger io.Writer) error {
	fmt.Fprintf(logger, "Extracting %s...\n", filepath.Base(archivePath))

	// Create destination directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}

	switch {
	case strings.HasSuffix(archiveType, ".zip"):
		return extractZip(archivePath, destPath, logger)
	case strings.Contains(archiveType, "tar"):
		return extractTar(archivePath, destPath, archiveType, logger)
	default:
		return fmt.Errorf("unsupported archive type: %s", archiveType)
	}
}

// extractTar extracts a tar archive (with optional compression)
func extractTar(archivePath, destPath, archiveType string, logger io.Writer) error {
	// First pass: detect common prefix
	stripPrefix, err := detectTarPrefix(archivePath, archiveType)
	if err != nil {
		return err
	}

	if stripPrefix != "" {
		fmt.Fprintf(logger, "  Auto-detected prefix: %s\n", stripPrefix)
	}

	// Second pass: extract with prefix stripped
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader io.Reader = file

	// Handle compression
	switch {
	case strings.Contains(archiveType, ".gz"):
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gzReader.Close()
		reader = gzReader

	case strings.Contains(archiveType, ".bz2"):
		reader = bzip2.NewReader(file)

	case strings.Contains(archiveType, ".xz"):
		xzReader, err := xz.NewReader(file)
		if err != nil {
			return err
		}
		reader = xzReader
	}

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Strip auto-detected prefix
		name := header.Name
		if stripPrefix != "" {
			if !strings.HasPrefix(name, stripPrefix) {
				continue
			}
			name = strings.TrimPrefix(name, stripPrefix)
		}

		if name == "" {
			continue
		}

		target := filepath.Join(destPath, name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()

		case tar.TypeSymlink:
			// Handle symlinks
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, target); err != nil && !os.IsExist(err) {
				return err
			}
		}
	}

	return nil
}

// extractZip extracts a zip archive
func extractZip(archivePath, destPath string, logger io.Writer) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Detect common prefix
	stripPrefix := detectZipPrefix(reader)
	if stripPrefix != "" {
		fmt.Fprintf(logger, "  Auto-detected prefix: %s\n", stripPrefix)
	}

	for _, file := range reader.File {
		// Strip auto-detected prefix
		name := file.Name
		if stripPrefix != "" {
			if !strings.HasPrefix(name, stripPrefix) {
				continue
			}
			name = strings.TrimPrefix(name, stripPrefix)
		}

		if name == "" {
			continue
		}

		target := filepath.Join(destPath, name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
