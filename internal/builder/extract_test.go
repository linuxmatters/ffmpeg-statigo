package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/linuxmatters/ffmpeg-statigo/internal/pathsafe"
)

func TestExtractZipRejectsExcessEntryCount(t *testing.T) {
	archivePath := filepath.Join(t.TempDir(), "bomb.zip")
	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(file)
	for i := range pathsafe.MaxExtractEntries + 1 {
		if _, err := zw.Create(fmt.Sprintf("file%d.txt", i)); err != nil {
			t.Fatalf("add zip entry %d: %v", i, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close zip file: %v", err)
	}

	err = extractZip(archivePath, t.TempDir(), io.Discard)
	if err == nil {
		t.Fatal("expected entry-count error, got nil")
	}
	if !strings.Contains(err.Error(), "entries") {
		t.Fatalf("error = %v, want entries cap", err)
	}
}
