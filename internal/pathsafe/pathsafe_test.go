package pathsafe

import (
	"path/filepath"
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
