package ffmpeg_test

import (
	"go/build"
	"path"
	"path/filepath"
	"slices"
	"testing"
)

func TestWindowsProfileLinkSelection(t *testing.T) {
	root, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, arch := range []string{"amd64", "386"} {
		for _, profile := range []string{"default", "embedded"} {
			t.Run(arch+"/"+profile, func(t *testing.T) {
				ctx := build.Default
				ctx.GOOS, ctx.GOARCH, ctx.CgoEnabled = "windows", arch, true
				ctx.BuildTags = nil
				otherProfile := "embedded"
				archiveDir := path.Join(filepath.ToSlash(root), "lib")
				if profile == "embedded" {
					ctx.BuildTags = []string{"embedded"}
					otherProfile = "default"
					archiveDir = path.Join(archiveDir, "embedded")
				}

				pkg, err := ctx.ImportDir(root, 0)
				if err != nil {
					t.Fatal(err)
				}
				if !slices.Contains(pkg.CgoFiles, "ffmpeg_link_"+profile+".go") ||
					slices.Contains(pkg.CgoFiles, "ffmpeg_link_"+otherProfile+".go") {
					t.Fatalf("wrong linker file selection: %v", pkg.CgoFiles)
				}
				if !slices.Contains(pkg.TestGoFiles, "profile_"+profile+"_test.go") ||
					slices.Contains(pkg.TestGoFiles, "profile_"+otherProfile+"_test.go") {
					t.Fatalf("wrong test profile selection: %v", pkg.TestGoFiles)
				}

				want := []string{path.Join(archiveDir, "windows_"+arch, "libffmpeg.a"), "-lstdc++"}
				if profile == "default" {
					want = append(want, "-lpthread")
				}
				want = append(want, "-lm", "-lws2_32", "-lbcrypt", "-luser32")
				if profile == "default" {
					want = append(want, "-lntdll", "-luserenv", "-lcrypt32")
				}
				if !slices.Equal(pkg.CgoLDFLAGS, want) {
					t.Errorf("link flags = %v, want %v", pkg.CgoLDFLAGS, want)
				}
			})
		}
	}
}
