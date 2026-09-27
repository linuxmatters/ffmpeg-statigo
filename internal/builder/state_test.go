package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeArchive(t *testing.T, installDir, name string) {
	t.Helper()
	libDir := filepath.Join(installDir, "lib")
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		t.Fatalf("failed to create lib dir: %v", err)
	}
	path := filepath.Join(libDir, name+".a")
	if err := os.WriteFile(path, []byte("!<arch>\n"), 0o644); err != nil {
		t.Fatalf("failed to write archive %s: %v", path, err)
	}
}

// TestCanSkip covers all six branches of BuildState.CanSkip: empty URL,
// URL mismatch, config-hash mismatch, header-only library, missing output,
// and all outputs present.
func TestCanSkip(t *testing.T) {
	const url = "https://example.com/foo-1.0.tar.gz"

	tests := []struct {
		name string
		// setup returns the state under test and the install dir it checks.
		setup func(t *testing.T) (*BuildState, string)
		want  bool
	}{
		{
			name: "empty URL must build",
			setup: func(t *testing.T) (*BuildState, string) {
				lib := &Library{Name: "foo", URL: url, LinkLibs: []string{"libfoo"}}
				return NewBuildState(lib, t.TempDir()), t.TempDir()
			},
			want: false,
		},
		{
			name: "URL mismatch must rebuild",
			setup: func(t *testing.T) (*BuildState, string) {
				lib := &Library{Name: "foo", URL: url, LinkLibs: []string{"libfoo"}}
				s := NewBuildState(lib, t.TempDir())
				s.URL = "https://example.com/old-0.9.tar.gz"
				s.ConfigHash = lib.ConfigHash()
				return s, t.TempDir()
			},
			want: false,
		},
		{
			name: "config-hash mismatch must rebuild",
			setup: func(t *testing.T) (*BuildState, string) {
				lib := &Library{Name: "foo", URL: url, LinkLibs: []string{"libfoo"}}
				s := NewBuildState(lib, t.TempDir())
				s.URL = lib.URL
				s.ConfigHash = "stale-hash"
				return s, t.TempDir()
			},
			want: false,
		},
		{
			name: "header-only library skips",
			setup: func(t *testing.T) (*BuildState, string) {
				lib := &Library{Name: "foo", URL: url, LinkLibs: nil}
				s := NewBuildState(lib, t.TempDir())
				s.URL = lib.URL
				s.ConfigHash = lib.ConfigHash()
				return s, t.TempDir()
			},
			want: true,
		},
		{
			name: "missing output must rebuild",
			setup: func(t *testing.T) (*BuildState, string) {
				lib := &Library{Name: "foo", URL: url, LinkLibs: []string{"libfoo"}}
				s := NewBuildState(lib, t.TempDir())
				s.URL = lib.URL
				s.ConfigHash = lib.ConfigHash()
				return s, t.TempDir()
			},
			want: false,
		},
		{
			name: "all outputs present skips",
			setup: func(t *testing.T) (*BuildState, string) {
				lib := &Library{Name: "foo", URL: url, LinkLibs: []string{"libfoo", "libbar"}}
				s := NewBuildState(lib, t.TempDir())
				s.URL = lib.URL
				s.ConfigHash = lib.ConfigHash()
				installDir := t.TempDir()
				writeArchive(t, installDir, "libfoo")
				writeArchive(t, installDir, "libbar")
				return s, installDir
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, installDir := tt.setup(t)
			if got := s.CanSkip(installDir); got != tt.want {
				t.Errorf("CanSkip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanSkipInvalidArchives(t *testing.T) {
	tests := []struct {
		name   string
		create func(string) error
	}{
		{
			name: "empty archive",
			create: func(path string) error {
				return os.WriteFile(path, nil, 0o644)
			},
		},
		{
			name: "directory instead of archive",
			create: func(path string) error {
				return os.Mkdir(path, 0o755)
			},
		},
		{
			name: "broken symlink",
			create: func(path string) error {
				return os.Symlink(path+".missing", path)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lib := &Library{
				Name:     "foo",
				URL:      "https://example.com/foo-1.0.tar.gz",
				LinkLibs: []string{"libfoo", "libbar"},
			}
			buildRoot := t.TempDir()
			if err := NewBuildState(lib, buildRoot).Save(); err != nil {
				t.Fatalf("failed to save build state: %v", err)
			}
			installDir := t.TempDir()
			writeArchive(t, installDir, "libfoo")
			if err := tt.create(filepath.Join(installDir, "lib", "libbar.a")); err != nil {
				t.Fatalf("failed to create invalid archive: %v", err)
			}
			if NewBuildState(lib, buildRoot).CanSkip(installDir) {
				t.Error("CanSkip() = true, want false for invalid archive")
			}
		})
	}
}

// TestConfigHashBuildEnv verifies that changing BuildEnv busts the config hash,
// so a maintainer edit to build env values forces a rebuild.
func TestConfigHashBuildEnv(t *testing.T) {
	const url = "https://example.com/foo-1.0.tar.gz"

	base := &Library{Name: "foo", URL: url}
	withEnv := &Library{Name: "foo", URL: url, BuildEnv: func() []string { return []string{"AS=nasm"} }}
	changedEnv := &Library{Name: "foo", URL: url, BuildEnv: func() []string { return []string{"AS=clang"} }}
	sameEnv := &Library{Name: "foo", URL: url, BuildEnv: func() []string { return []string{"AS=nasm"} }}

	if base.ConfigHash() == withEnv.ConfigHash() {
		t.Error("adding BuildEnv did not change ConfigHash")
	}
	if withEnv.ConfigHash() == changedEnv.ConfigHash() {
		t.Error("changing BuildEnv value did not change ConfigHash")
	}
	if withEnv.ConfigHash() != sameEnv.ConfigHash() {
		t.Error("ConfigHash is not deterministic for identical BuildEnv")
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	present := filepath.Join(dir, "present")
	if err := os.WriteFile(present, []byte{}, 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if !fileExists(present) {
		t.Errorf("fileExists(%q) = false, want true", present)
	}
	if fileExists(filepath.Join(dir, "absent")) {
		t.Errorf("fileExists(absent) = true, want false")
	}
}
