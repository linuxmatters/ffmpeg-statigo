// Package archiveextract provides shared archive extraction helpers.
package archiveextract

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/linuxmatters/ffmpeg-statigo/internal/pathsafe"
)

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
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return extractionError(opts, fmt.Errorf("reading tar header: %w", err))
		}

		name, ok := tarEntryName(header.Name, opts.StripPrefix)
		if !ok {
			continue
		}

		target, err := pathsafe.SanitizePath(opts.DestDir, name)
		if err != nil {
			return extractionError(opts, err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return extractionError(opts, fmt.Errorf("creating directory %s: %w", target, err))
			}
		case tar.TypeReg:
			if err := extractRegularFile(tarReader, header, target, opts); err != nil {
				return extractionError(opts, err)
			}
		case tar.TypeSymlink:
			if opts.LinkPolicy == PreserveSymlinks {
				if err := extractSymlink(opts.DestDir, header, target); err != nil {
					return extractionError(opts, err)
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

func extractRegularFile(tarReader *tar.Reader, header *tar.Header, target string, opts TarOptions) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("creating parent directory for %s: %w", target, err)
	}

	mode := os.FileMode(0o666)
	if opts.FileMode != nil {
		mode = opts.FileMode(header)
	}

	outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", target, err)
	}

	if err := pathsafe.CopyCapped(outFile, tarReader); err != nil {
		_ = outFile.Close()
		if opts.RemoveIncomplete {
			_ = os.Remove(target)
		}
		return fmt.Errorf("writing file %s: %w", target, err)
	}

	if err := outFile.Close(); err != nil {
		if opts.RemoveIncomplete {
			_ = os.Remove(target)
		}
		return fmt.Errorf("closing file %s: %w", target, err)
	}

	return nil
}

func extractSymlink(destDir string, header *tar.Header, target string) error {
	if err := symlinkTargetSafe(destDir, target, header.Linkname); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("creating parent directory for %s: %w", target, err)
	}
	if err := os.Symlink(header.Linkname, target); err != nil && !os.IsExist(err) {
		return fmt.Errorf("creating symlink %s: %w", target, err)
	}
	return nil
}

// symlinkTargetSafe rejects symlink targets that resolve outside destDir,
// preventing a malicious archive from planting a link to an absolute or
// parent path that later entries could write through.
func symlinkTargetSafe(destDir, linkPath, linkname string) error {
	if linkname == "" {
		return fmt.Errorf("symlink %s: empty link target", linkPath)
	}
	if filepath.IsAbs(linkname) {
		return fmt.Errorf("symlink %s: absolute target %q not allowed", linkPath, linkname)
	}

	cleanDest := filepath.Clean(destDir)
	resolved := filepath.Join(filepath.Dir(linkPath), linkname)
	if resolved != cleanDest && !strings.HasPrefix(resolved, cleanDest+string(filepath.Separator)) {
		return fmt.Errorf("symlink %s: target %q escapes destination directory", linkPath, linkname)
	}
	return nil
}

func extractionError(opts TarOptions, err error) error {
	if opts.OnError != nil {
		opts.OnError()
	}
	return err
}
