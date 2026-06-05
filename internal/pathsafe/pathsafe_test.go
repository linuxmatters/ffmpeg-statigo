package pathsafe

import (
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name      string
		destDir   string
		entryName string
		want      string
		wantErr   bool
	}{
		{
			name:      "valid entry with trailing separator on destDir",
			destDir:   "/tmp/out" + string(filepath.Separator),
			entryName: "file.txt",
			want:      filepath.Join("/tmp/out", "file.txt"),
		},
		{
			name:      "valid nested entry with clean destDir",
			destDir:   "/tmp/out",
			entryName: filepath.Join("sub", "file.txt"),
			want:      filepath.Join("/tmp/out", "sub", "file.txt"),
		},
		{
			name:      "parent traversal rejected",
			destDir:   "/tmp/out",
			entryName: filepath.Join("..", "escape.txt"),
			wantErr:   true,
		},
		{
			name:      "absolute entry rejected",
			destDir:   "/tmp/out",
			entryName: string(filepath.Separator) + "etc" + string(filepath.Separator) + "passwd",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SanitizePath(tt.destDir, tt.entryName)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (result %q)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCopyCapped(t *testing.T) {
	tests := []struct {
		name    string
		size    int64
		wantErr bool
	}{
		{
			name: "under cap succeeds",
			size: 1 << 20,
		},
		{
			name: "at cap succeeds",
			size: MaxExtractFileSize,
		},
		{
			name:    "over cap rejected",
			size:    MaxExtractFileSize + 1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := io.LimitReader(zeroReader{}, tt.size)
			n, err := CopyCapped(io.Discard, src)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), "exceeds") {
					t.Fatalf("error %q does not contain %q", err.Error(), "exceeds")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n != tt.size {
				t.Fatalf("bytes written = %d, want %d", n, tt.size)
			}
		})
	}
}

func TestBudgetEntries(t *testing.T) {
	var b Budget
	for i := range MaxExtractEntries {
		if err := b.AddEntry(); err != nil {
			t.Fatalf("AddEntry at %d: unexpected error %v", i, err)
		}
	}
	if err := b.AddEntry(); err == nil {
		t.Fatal("expected entry-count error, got nil")
	} else if !strings.Contains(err.Error(), "entries") {
		t.Fatalf("error %q does not contain %q", err.Error(), "entries")
	}
}

func TestBudgetBytes(t *testing.T) {
	var b Budget
	if err := b.AddBytes(MaxExtractTotalSize); err != nil {
		t.Fatalf("AddBytes at cap: unexpected error %v", err)
	}
	if err := b.AddBytes(1); err == nil {
		t.Fatal("expected aggregate-byte error, got nil")
	} else if !strings.Contains(err.Error(), "total bytes") {
		t.Fatalf("error %q does not contain %q", err.Error(), "total bytes")
	}
}

// zeroReader yields an endless stream of zero bytes without allocating a buffer.
type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}
