package main

import (
	"os"
	"runtime"
	"slices"
	"testing"
)

// TestCheckHostIncludes pins the dev-shell contract: empty include discovery
// fails inside the Nix shell and passes outside it.
func TestCheckHostIncludes(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		nixCC   string
		wantErr bool
	}{
		{name: "paths under nix", paths: []string{"/usr/include"}, nixCC: "/nix/store/cc", wantErr: false},
		{name: "paths without nix", paths: []string{"/usr/include"}, nixCC: "", wantErr: false},
		{name: "empty without nix", paths: nil, nixCC: "", wantErr: false},
		{name: "empty under nix", paths: nil, nixCC: "/nix/store/cc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkHostIncludes(tt.paths, tt.nixCC)
			if (err != nil) != tt.wantErr {
				t.Fatalf("checkHostIncludes(%v, %q) = %v, want error: %v", tt.paths, tt.nixCC, err, tt.wantErr)
			}
		})
	}
}

// TestNewCCConfigDiscoversHostIncludes asserts that host include discovery
// returns a non-empty list, so the parse sees the host's system headers. Under
// NIX_CC the dev-shell contract makes a configuration failure a test failure;
// off Nix a host without a C compiler skips instead.
func TestNewCCConfigDiscoversHostIncludes(t *testing.T) {
	nixCC := os.Getenv("NIX_CC")

	cfg, err := newCCConfig(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		if nixCC == "" {
			t.Skipf("no host C compiler configuration available: %v", err)
		}

		t.Fatalf("newCCConfig under NIX_CC=%q: %v", nixCC, err)
	}

	paths := hostIncludePaths(cfg)
	if len(paths) == 0 {
		t.Fatal("host include discovery returned no paths")
	}

	t.Logf("host include paths: %v", paths)

	if !cfg.Header {
		t.Error("cfg.Header is false; function bodies would be type checked")
	}

	if !slices.Contains(cfg.IncludePaths, AVLibPath) {
		t.Errorf("IncludePaths %v does not carry the header root %q", cfg.IncludePaths, AVLibPath)
	}

	if !slices.Contains(cfg.SysIncludePaths, AVLibPath) {
		t.Errorf("SysIncludePaths %v does not carry the header root %q", cfg.SysIncludePaths, AVLibPath)
	}
}
