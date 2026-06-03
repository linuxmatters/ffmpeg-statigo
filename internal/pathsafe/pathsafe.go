// Package pathsafe provides shared path-traversal guards for archive extraction.
package pathsafe

import (
	"fmt"
	"path/filepath"
	"strings"
)

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
	target := filepath.Join(destDir, entryName)

	// Defence in depth: confirm the resolved path is within destDir
	// This catches edge cases where filepath.Join might not prevent traversal
	if !strings.HasPrefix(target, destDir+string(filepath.Separator)) && target != destDir {
		return "", fmt.Errorf("path traversal detected: %q resolves outside destination directory", entryName)
	}

	return target, nil
}
