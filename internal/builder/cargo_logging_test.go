package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func TestWindowsRav1eBuildLogsCargoFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("fake Cargo command requires a POSIX shell")
	}
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "cargo"), []byte("#!/bin/sh\necho cargo-stdout\necho cargo-stderr >&2\nexit 23\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))

	for _, arch := range []string{"amd64", "386"} {
		t.Run(arch, func(t *testing.T) {
			root := t.TempDir()
			lib := windowsLibrary(rav1e, arch)
			lib.URL = "https://example.test/rav1e.tar.gz"
			var archive bytes.Buffer
			gz := gzip.NewWriter(&archive)
			tw := tar.NewWriter(gz)
			if err := tw.WriteHeader(&tar.Header{Name: "rav1e/Cargo.toml", Mode: 0o644, Size: 7}); err != nil {
				t.Fatal(err)
			}
			if _, err := io.WriteString(tw, "fixture"); err != nil {
				t.Fatal(err)
			}
			if err := tw.Close(); err != nil {
				t.Fatal(err)
			}
			if err := gz.Close(); err != nil {
				t.Fatal(err)
			}
			withPinnedDigest(t, lib.URL, fmt.Sprintf("%x", sha256.Sum256(archive.Bytes())))
			downloads := filepath.Join(root, "downloads")
			if err := os.MkdirAll(downloads, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(downloads, "rav1e.tar.gz"), archive.Bytes(), 0o644); err != nil {
				t.Fatal(err)
			}
			buildDir := filepath.Join(root, "build", lib.Name)
			if err := os.MkdirAll(buildDir, 0o755); err != nil {
				t.Fatal(err)
			}
			logPath := filepath.Join(buildDir, "build.log")
			logFile, err := os.Create(logPath)
			if err != nil {
				t.Fatal(err)
			}
			defer logFile.Close()
			var output bytes.Buffer
			logger := io.MultiWriter(&output, logFile)
			fmt.Fprintln(logger, "before build")
			err = lib.Build(t.Context(), root, filepath.Join(root, "staging"), logger)
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) || exitErr.ExitCode() != 23 || !strings.Contains(err.Error(), "build failed: cargo failed:") {
				t.Fatalf("Build error = %v, want Cargo exit 23", err)
			}
			fmt.Fprintln(logger, "after build")
			for _, want := range []string{"before build\n", "Building rav1e...\n", "cargo-stdout\n", "cargo-stderr\n", "after build\n"} {
				if !strings.Contains(output.String(), want) {
					t.Errorf("logger output lacks %q: %s", want, &output)
				}
			}
			logged, err := os.ReadFile(logPath)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(logged, output.Bytes()) {
				t.Fatalf("build.log differs from logger output: %q", logged)
			}
		})
	}
}

func TestRav1eNonWindowsCargoCallbackRouting(t *testing.T) {
	for _, targetOS := range []string{"linux", "darwin"} {
		t.Run(targetOS, func(t *testing.T) {
			root := t.TempDir()
			libs := librariesForPlatform(false, filepath.Join(root, "staging"), targetOS, "amd64")
			index := slices.IndexFunc(libs, func(lib *Library) bool { return lib.Name == "rav1e" })
			if index < 0 || libs[index] != rav1e {
				t.Fatal("non-Windows profile changed the global rav1e definition")
			}
			build := *libs[index].BuildSystem.(*CargoBuild)
			if build.InstallFunc == nil || build.InstallLogFunc != nil {
				t.Fatal("non-Windows rav1e must keep the original install callback")
			}
			wantErr := errors.New("install failed")
			calls := 0
			build.InstallFunc = func(_ context.Context, srcPath, installDir string) error {
				calls++
				if srcPath != "source" || installDir != filepath.Join(root, "staging") {
					t.Fatalf("install paths = %q, %q", srcPath, installDir)
				}
				return wantErr
			}
			buildDir := filepath.Join(root, "build", "rav1e")
			if err := build.BuildWithLog(t.Context(), rav1e, "source", buildDir, io.Discard); !errors.Is(err, wantErr) {
				t.Fatalf("BuildWithLog error = %v, want install error", err)
			}
			if err := build.Build(t.Context(), rav1e, "source", buildDir); !errors.Is(err, wantErr) {
				t.Fatalf("Build error = %v, want install error", err)
			}
			if calls != 2 {
				t.Fatalf("install callback calls = %d, want 2", calls)
			}
		})
	}
}
