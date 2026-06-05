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

// MaxExtractTotalSize caps the combined size of every file in one archive.
// The largest input is an FFmpeg-dependency source tree, well under a gigabyte
// uncompressed; 8 GiB leaves generous headroom while still bounding a malicious
// archive that stays under the per-file cap but inflates in aggregate.
const MaxExtractTotalSize = 8 << 30 // 8 GiB

// MaxExtractEntries caps the number of entries in one archive, bounding a
// zip/tar bomb made of many tiny files. Real source trees hold a few tens of
// thousands of files; 100000 covers the largest with room to spare.
const MaxExtractEntries = 100000

// CopyCapped copies from src to dst, refusing more than MaxExtractFileSize bytes.
// It returns the number of bytes written so callers can track an archive-wide total.
func CopyCapped(dst io.Writer, src io.Reader) (int64, error) {
	n, err := io.Copy(dst, io.LimitReader(src, MaxExtractFileSize+1))
	if err != nil {
		return n, err
	}
	if n > MaxExtractFileSize {
		return n, fmt.Errorf("extracted file exceeds %d bytes", MaxExtractFileSize)
	}
	return n, nil
}

// Budget tracks running totals across an archive extraction and rejects archives
// that exceed the aggregate byte or entry-count caps. The zero value is ready to
// use and is not safe for concurrent use.
type Budget struct {
	entries int
	bytes   int64
}

// AddEntry records one archive entry, returning an error once the entry-count
// cap is exceeded.
func (b *Budget) AddEntry() error {
	b.entries++
	if b.entries > MaxExtractEntries {
		return fmt.Errorf("archive exceeds %d entries", MaxExtractEntries)
	}
	return nil
}

// AddBytes accumulates extracted bytes, returning an error once the aggregate
// byte cap is exceeded.
func (b *Budget) AddBytes(n int64) error {
	b.bytes += n
	if b.bytes > MaxExtractTotalSize {
		return fmt.Errorf("archive exceeds %d total bytes", MaxExtractTotalSize)
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
