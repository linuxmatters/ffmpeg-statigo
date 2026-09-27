package archiveextract

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
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
		wantMode := os.FileMode(0o755)
		if runtime.GOOS == "windows" {
			wantMode = 0o666 // Windows reports writable files without POSIX executable bits.
		}
		if got := info.Mode().Perm(); got != wantMode {
			t.Fatalf("mode = %o, want %o", got, wantMode)
		}
		data, err := os.ReadFile(filepath.Join(destDir, "bin", "tool"))
		if err != nil {
			t.Fatalf("read extracted file: %v", err)
		}
		if string(data) != "tool" {
			t.Fatalf("contents = %q, want %q", data, "tool")
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

func TestSymlinkTargetSafe(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "dest")
	linkPath := filepath.Join(destDir, "subdir", "link")
	numericColonErr := ""
	numericEscapeErr := "escapes destination directory"
	if runtime.GOOS == "windows" {
		numericColonErr = "absolute target"
		numericEscapeErr = "absolute target"
	}

	tests := []struct {
		name     string
		linkname string
		wantErr  string
	}{
		{name: "posix_root", linkname: "/etc/passwd", wantErr: "absolute target"},
		{name: "windows_root", linkname: `\Windows\system.ini`, wantErr: "absolute target"},
		{name: "drive_forward_slash", linkname: "C:/Windows/system.ini", wantErr: "absolute target"},
		{name: "drive_backslash", linkname: `C:\Windows\system.ini`, wantErr: "absolute target"},
		{name: "drive_relative", linkname: "C:outside.txt", wantErr: "absolute target"},
		{name: "drive_relative_parent", linkname: `C:..\outside.txt`, wantErr: "absolute target"},
		{name: "drive_only", linkname: "C:", wantErr: "absolute target"},
		{name: "numeric_colon", linkname: "1:target", wantErr: numericColonErr},
		{name: "numeric_colon_escape", linkname: "1:target/../../../outside", wantErr: numericEscapeErr},
		{name: "unc_forward_slash", linkname: "//server/share/file", wantErr: "absolute target"},
		{name: "unc_backslash", linkname: `\\server\share\file`, wantErr: "absolute target"},
		{name: "empty", wantErr: "empty link target"},
		{name: "escaping", linkname: "../../outside", wantErr: "escapes destination directory"},
		{name: "escaping_native_separator", linkname: filepath.Join("..", "..", "outside"), wantErr: "escapes destination directory"},
		{name: "sibling_prefix", linkname: "../../dest-other/file", wantErr: "escapes destination directory"},
		{name: "relative", linkname: "target"},
		{name: "relative_forward_slash", linkname: "nested/target"},
		{name: "relative_backslash", linkname: `nested\target`},
		{name: "relative_parent_inside", linkname: "../target"},
		{name: "destination_itself", linkname: ".."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := symlinkTargetSafe(destDir, linkPath, tt.linkname)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("symlinkTargetSafe() error = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("symlinkTargetSafe() error = %v, want %q", err, tt.wantErr)
			}
		})
	}
}

func TestHasDriveDesignator(t *testing.T) {
	tests := []struct {
		linkname string
		posix    bool
		windows  bool
	}{
		{linkname: "1:target", windows: true},
		{linkname: "1:", windows: true},
		{linkname: "!:target", windows: true},
		{linkname: "C:outside", posix: true, windows: true},
		{linkname: "c:outside", posix: true, windows: true},
		{linkname: "C:/outside", posix: true, windows: true},
		{linkname: `C:\outside`, posix: true, windows: true},
		{linkname: "C:", posix: true, windows: true},
		{linkname: "target"},
		{linkname: "nested/1:target"},
		{linkname: "long:target"},
		{linkname: ":"},
		{linkname: ""},
	}
	for _, targetOS := range []string{"linux", "darwin", "windows"} {
		for _, tt := range tests {
			t.Run(targetOS+"/"+tt.linkname, func(t *testing.T) {
				want := tt.posix
				if targetOS == "windows" {
					want = tt.windows
				}
				if got := hasDriveDesignator(tt.linkname, targetOS); got != want {
					t.Fatalf("hasDriveDesignator(%q, %q) = %v, want %v", tt.linkname, targetOS, got, want)
				}
			})
		}
	}
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
