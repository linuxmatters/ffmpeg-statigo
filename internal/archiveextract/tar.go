// Package archiveextract provides shared archive extraction helpers.
package archiveextract

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/linuxmatters/ffmpeg-statigo/internal/pathsafe"
)

const maxSymlinkHops = 40

// TarLinkPolicy controls how tar link entries are handled.
type TarLinkPolicy int

const (
	// SkipLinks ignores symlink and hard-link entries.
	SkipLinks TarLinkPolicy = iota
	// PreserveSymlinks creates symlink entries. Hard links are still skipped.
	PreserveSymlinks
)

// TarOptions configures tar extraction.
type TarOptions struct {
	DestDir          string
	StripPrefix      string
	LinkPolicy       TarLinkPolicy
	FileMode         func(*tar.Header) os.FileMode
	RemoveIncomplete bool
	OnError          func()
}

// ExtractTar extracts entries from a tar stream into DestDir.
func ExtractTar(reader io.Reader, opts TarOptions) error {
	if err := os.MkdirAll(opts.DestDir, 0o755); err != nil {
		return extractionError(opts, fmt.Errorf("creating parent directory %s: %w", opts.DestDir, err))
	}

	root, err := os.OpenRoot(opts.DestDir)
	if err != nil {
		return extractionError(opts, fmt.Errorf("opening destination directory %s: %w", opts.DestDir, err))
	}
	defer root.Close()

	cleanDest := filepath.Clean(opts.DestDir)
	tarReader := tar.NewReader(reader)
	var budget pathsafe.Budget
	var createdSymlinks []string

	fail := func(err error) error {
		_ = verifyAndCleanupSymlinks(root, createdSymlinks)
		return extractionError(opts, err)
	}

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			if err := verifyAndCleanupSymlinks(root, createdSymlinks); err != nil {
				return extractionError(opts, err)
			}
			return nil
		}
		if err != nil {
			return fail(fmt.Errorf("reading tar header: %w", err))
		}

		if err := budget.AddEntry(); err != nil {
			return fail(err)
		}

		name, ok := tarEntryName(header.Name, opts.StripPrefix)
		if !ok {
			continue
		}

		target, err := pathsafe.SanitizePath(cleanDest, name)
		if err != nil {
			return fail(err)
		}

		relPath, err := filepath.Rel(cleanDest, target)
		if err != nil {
			return fail(fmt.Errorf("resolving relative path for %s: %w", name, err))
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if relPath != "." {
				if err := root.MkdirAll(relPath, 0o755); err != nil {
					return fail(fmt.Errorf("creating directory %s: %w", target, err))
				}
			}
		case tar.TypeReg:
			if err := extractRegularFile(root, tarReader, header, relPath, target, opts, &budget); err != nil {
				return fail(err)
			}
		case tar.TypeSymlink:
			if opts.LinkPolicy == PreserveSymlinks {
				created, err := extractSymlink(root, cleanDest, header, relPath, target)
				if err != nil {
					return fail(err)
				}
				if created {
					createdSymlinks = append(createdSymlinks, relPath)
				}
			}
		case tar.TypeLink:
			continue
		}
	}
}

func tarEntryName(name, stripPrefix string) (string, bool) {
	if stripPrefix != "" {
		if !strings.HasPrefix(name, stripPrefix) {
			return "", false
		}
		name = strings.TrimPrefix(name, stripPrefix)
	}

	if name == "" {
		return "", false
	}

	return name, true
}

func extractRegularFile(root *os.Root, tarReader *tar.Reader, header *tar.Header, relPath, target string, opts TarOptions, budget *pathsafe.Budget) error {
	if parentDir := filepath.Dir(relPath); parentDir != "." {
		if err := root.MkdirAll(parentDir, 0o755); err != nil {
			return fmt.Errorf("creating parent directory for %s: %w", target, err)
		}
	}

	mode := os.FileMode(0o666)
	if opts.FileMode != nil {
		mode = opts.FileMode(header)
	}

	outFile, err := root.OpenFile(relPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode.Perm())
	if err != nil {
		return fmt.Errorf("creating file %s: %w", target, err)
	}

	written, err := pathsafe.CopyCapped(outFile, tarReader)
	if err == nil {
		err = budget.AddBytes(written)
	}
	if err != nil {
		_ = outFile.Close()
		if opts.RemoveIncomplete {
			_ = root.Remove(relPath)
		}
		return fmt.Errorf("writing file %s: %w", target, err)
	}

	if err := outFile.Close(); err != nil {
		if opts.RemoveIncomplete {
			_ = root.Remove(relPath)
		}
		return fmt.Errorf("closing file %s: %w", target, err)
	}

	return nil
}

func extractSymlink(root *os.Root, destDir string, header *tar.Header, relPath, target string) (bool, error) {
	if err := symlinkTargetSafe(destDir, target, header.Linkname); err != nil {
		return false, err
	}
	if parentDir := filepath.Dir(relPath); parentDir != "." {
		if err := root.MkdirAll(parentDir, 0o755); err != nil {
			return false, fmt.Errorf("creating parent directory for %s: %w", target, err)
		}
	}
	created := true
	if err := root.Symlink(header.Linkname, relPath); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			return false, fmt.Errorf("creating symlink %s: %w", target, err)
		}
		created = false
	}
	if err := verifySymlinkInRoot(root, relPath); err != nil {
		if created {
			_ = root.Remove(relPath)
		}
		return false, err
	}
	return created, nil
}

func verifyAndCleanupSymlinks(root *os.Root, relPaths []string) error {
	var firstErr error
	for _, relPath := range relPaths {
		if err := verifySymlinkInRoot(root, relPath); err != nil {
			_ = root.Remove(relPath)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

// verifySymlinkInRoot walks relPath from root, resolving every symlink along
// the path and its target chain to confirm the link stays within root.
func verifySymlinkInRoot(root *os.Root, relPath string) error {
	pending := splitPathSegments(relPath)
	current := make([]string, 0, len(pending))
	hops := 0

	for len(pending) > 0 {
		seg := pending[0]
		pending = pending[1:]

		switch seg {
		case "", ".":
			continue
		case "..":
			if len(current) == 0 {
				return fmt.Errorf("symlink %s: target escapes destination directory", relPath)
			}
			current = current[:len(current)-1]
			continue
		}

		current = append(current, seg)
		candidate := filepath.Join(current...)
		info, err := root.Lstat(candidate)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("symlink %s: inspecting %s: %w", relPath, candidate, err)
		}
		if info.Mode()&os.ModeSymlink == 0 {
			continue
		}

		hops++
		if hops > maxSymlinkHops {
			return fmt.Errorf("symlink %s: too many symlink hops (possible loop)", relPath)
		}

		linkTarget, err := root.Readlink(candidate)
		if err != nil {
			return fmt.Errorf("symlink %s: reading link %s: %w", relPath, candidate, err)
		}
		if err := validateLinkTargetSyntax(relPath, linkTarget); err != nil {
			return err
		}

		current = current[:len(current)-1]
		targetSegs := splitPathSegments(linkTarget)
		next := make([]string, 0, len(targetSegs)+len(pending))
		next = append(next, targetSegs...)
		next = append(next, pending...)
		pending = next
	}

	return nil
}

func splitPathSegments(p string) []string {
	return strings.Split(filepath.ToSlash(p), "/")
}

// symlinkTargetSafe rejects symlink targets that resolve outside destDir,
// preventing a malicious archive from planting a link to an absolute or
// parent path that later entries could write through.
func symlinkTargetSafe(destDir, linkPath, linkname string) error {
	if err := validateLinkTargetSyntax(linkPath, linkname); err != nil {
		return err
	}

	cleanDest := filepath.Clean(destDir)
	resolved := filepath.Join(filepath.Dir(linkPath), linkname)
	if resolved != cleanDest && !strings.HasPrefix(resolved, cleanDest+string(filepath.Separator)) {
		return fmt.Errorf("symlink %s: target %q escapes destination directory", linkPath, linkname)
	}
	return nil
}

func validateLinkTargetSyntax(linkPath, linkname string) error {
	if linkname == "" {
		return fmt.Errorf("symlink %s: empty link target", linkPath)
	}
	slashLinkname := strings.ReplaceAll(linkname, `\`, "/")
	driveDesignator := hasDriveDesignator(slashLinkname, runtime.GOOS)
	if filepath.IsAbs(linkname) || strings.HasPrefix(slashLinkname, "/") || driveDesignator {
		return fmt.Errorf("symlink %s: absolute target %q not allowed", linkPath, linkname)
	}
	return nil
}

func hasDriveDesignator(linkname, targetOS string) bool {
	if len(linkname) < 2 || linkname[1] != ':' {
		return false
	}
	return targetOS == "windows" ||
		(linkname[0] >= 'A' && linkname[0] <= 'Z') ||
		(linkname[0] >= 'a' && linkname[0] <= 'z')
}

func extractionError(opts TarOptions, err error) error {
	if opts.OnError != nil {
		opts.OnError()
	}
	return err
}
