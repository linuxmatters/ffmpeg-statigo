// Package pathsafe provides shared path-traversal guards for archive extraction.
package pathsafe

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// MaxExtractFileSize caps a single extracted file to guard against decompression
// bombs. FFmpeg sources and prebuilt static libraries stay well below this.
const MaxExtractFileSize = 2 << 30 // 2 GiB

// CopyCapped copies from src to dst, refusing more than MaxExtractFileSize bytes.
func CopyCapped(dst io.Writer, src io.Reader) error {
	n, err := io.Copy(dst, io.LimitReader(src, MaxExtractFileSize+1))
	if err != nil {
		return err
	}
	if n > MaxExtractFileSize {
		return fmt.Errorf("extracted file exceeds %d bytes", MaxExtractFileSize)
	}
	return nil
}

// SanitizePath validates that an archive entry path is safe to extract.
// It prevents path traversal attacks by ensuring the resolved path
// stays within the destination directory.
func SanitizePath(destDir, entryName string) (string, error) {
	// Reject absolute paths (distinct branch so the error names the cause)
	if filepath.IsAbs(entryName) {
		return "", fmt.Errorf("path traversal detected: absolute path %q not allowed", entryName)
	}

	// Reject names that are not local: empty, absolute, or escaping via ..
	if !filepath.IsLocal(entryName) {
		return "", fmt.Errorf("path traversal detected: %q escapes destination directory", entryName)
	}

	// Construct the full target path
	cleanDest := filepath.Clean(destDir)
	target := filepath.Join(cleanDest, entryName)

	// Defence in depth: confirm the resolved path is within destDir
	// This catches edge cases where filepath.Join might not prevent traversal
	if target != cleanDest && !strings.HasPrefix(target, cleanDest+string(filepath.Separator)) {
		return "", fmt.Errorf("path traversal detected: %q resolves outside destination directory", entryName)
	}

	return target, nil
}
