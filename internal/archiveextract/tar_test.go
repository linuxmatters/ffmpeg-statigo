package archiveextract

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/linuxmatters/ffmpeg-statigo/internal/pathsafe"
)

func TestExtractTar(t *testing.T) {
	t.Run("strips_prefix_and_preserves_mode", func(t *testing.T) {
		destDir := t.TempDir()
		err := ExtractTar(tarStream(t, tarEntry{
			name: "src/bin/tool",
			mode: 0o755,
			body: "tool",
		}), TarOptions{
			DestDir:     destDir,
			StripPrefix: "src/",
			FileMode: func(header *tar.Header) os.FileMode {
				return os.FileMode(header.Mode & 0o777)
			},
		})
		if err != nil {
			t.Fatalf("ExtractTar() error = %v", err)
		}

		info, err := os.Stat(filepath.Join(destDir, "bin", "tool"))
		if err != nil {
			t.Fatalf("stat extracted file: %v", err)
		}
		if got := info.Mode().Perm(); got != 0o755 {
			t.Fatalf("mode = %o, want 755", got)
		}
	})

	t.Run("rejects_path_traversal", func(t *testing.T) {
		destDir := t.TempDir()
		err := ExtractTar(tarStream(t, tarEntry{
			name: "../escape",
			body: "bad",
		}), TarOptions{DestDir: destDir})
		if err == nil {
			t.Fatal("expected path traversal error, got nil")
		}
		if !strings.Contains(err.Error(), "path traversal") {
			t.Fatalf("error = %v, want path traversal", err)
		}
	})

	t.Run("skips_links_by_default", func(t *testing.T) {
		destDir := t.TempDir()
		err := ExtractTar(tarStream(t, tarEntry{
			name:     "link",
			linkname: "../target",
			typeflag: tar.TypeSymlink,
		}), TarOptions{DestDir: destDir})
		if err != nil {
			t.Fatalf("ExtractTar() error = %v", err)
		}
		if _, err := os.Lstat(filepath.Join(destDir, "link")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("link stat error = %v, want not exist", err)
		}
	})

	t.Run("preserves_symlinks_when_requested", func(t *testing.T) {
		destDir := t.TempDir()
		err := ExtractTar(tarStream(t, tarEntry{
			name:     "link",
			linkname: "target",
			typeflag: tar.TypeSymlink,
		}), TarOptions{
			DestDir:    destDir,
			LinkPolicy: PreserveSymlinks,
		})
		if err != nil {
			t.Fatalf("ExtractTar() error = %v", err)
		}

		got, err := os.Readlink(filepath.Join(destDir, "link"))
		if err != nil {
			t.Fatalf("readlink: %v", err)
		}
		if got != "target" {
			t.Fatalf("link target = %q, want %q", got, "target")
		}
	})

	t.Run("rejects_absolute_symlink_target", func(t *testing.T) {
		destDir := t.TempDir()
		err := ExtractTar(tarStream(t, tarEntry{
			name:     "link",
			linkname: "/etc/passwd",
			typeflag: tar.TypeSymlink,
		}), TarOptions{
			DestDir:    destDir,
			LinkPolicy: PreserveSymlinks,
		})
		if err == nil {
			t.Fatal("expected absolute symlink target error, got nil")
		}
		if !strings.Contains(err.Error(), "absolute target") {
			t.Fatalf("error = %v, want absolute target", err)
		}
		if _, statErr := os.Lstat(filepath.Join(destDir, "link")); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("link stat error = %v, want not exist", statErr)
		}
	})

	t.Run("rejects_escaping_symlink_target", func(t *testing.T) {
		destDir := t.TempDir()
		err := ExtractTar(tarStream(t, tarEntry{
			name:     "link",
			linkname: "../../etc/passwd",
			typeflag: tar.TypeSymlink,
		}), TarOptions{
			DestDir:    destDir,
			LinkPolicy: PreserveSymlinks,
		})
		if err == nil {
			t.Fatal("expected escaping symlink target error, got nil")
		}
		if !strings.Contains(err.Error(), "escapes destination directory") {
			t.Fatalf("error = %v, want escapes destination directory", err)
		}
	})

	t.Run("rejects_excess_entry_count", func(t *testing.T) {
		destDir := t.TempDir()
		entries := make([]tarEntry, pathsafe.MaxExtractEntries+1)
		for i := range entries {
			entries[i] = tarEntry{name: fmt.Sprintf("dir%d", i), typeflag: tar.TypeDir}
		}
		err := ExtractTar(tarStream(t, entries...), TarOptions{DestDir: destDir})
		if err == nil {
			t.Fatal("expected entry-count error, got nil")
		}
		if !strings.Contains(err.Error(), "entries") {
			t.Fatalf("error = %v, want entries cap", err)
		}
	})

	t.Run("removes_incomplete_file_on_copy_error", func(t *testing.T) {
		destDir := t.TempDir()
		err := ExtractTar(failingTarReader(t, "file"), TarOptions{
			DestDir:          destDir,
			RemoveIncomplete: true,
		})
		if err == nil {
			t.Fatal("expected copy error, got nil")
		}
		if _, statErr := os.Stat(filepath.Join(destDir, "file")); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("file stat error = %v, want not exist", statErr)
		}
	})
}

type tarEntry struct {
	name     string
	mode     int64
	body     string
	typeflag byte
	linkname string
}

func tarStream(t *testing.T, entries ...tarEntry) io.Reader {
	t.Helper()

	reader, writer := io.Pipe()
	go func() {
		tw := tar.NewWriter(writer)
		var err error
		for _, entry := range entries {
			mode := entry.mode
			if mode == 0 {
				mode = 0o644
			}
			typeflag := entry.typeflag
			if typeflag == 0 {
				typeflag = tar.TypeReg
			}
			header := &tar.Header{
				Name:     entry.name,
				Mode:     mode,
				Size:     int64(len(entry.body)),
				Typeflag: typeflag,
				Linkname: entry.linkname,
			}
			if typeflag != tar.TypeReg {
				header.Size = 0
			}
			if err = tw.WriteHeader(header); err != nil {
				break
			}
			if header.Size > 0 {
				_, err = tw.Write([]byte(entry.body))
				if err != nil {
					break
				}
			}
		}
		if closeErr := tw.Close(); err == nil {
			err = closeErr
		}
		_ = writer.CloseWithError(err)
	}()

	return reader
}

func failingTarReader(t *testing.T, name string) io.Reader {
	t.Helper()

	reader, writer := io.Pipe()
	go func() {
		tw := tar.NewWriter(writer)
		err := tw.WriteHeader(&tar.Header{
			Name: name,
			Mode: 0o644,
			Size: 1024,
		})
		if err == nil {
			_, err = tw.Write([]byte("partial"))
		}
		_ = writer.CloseWithError(err)
	}()

	return reader
}
