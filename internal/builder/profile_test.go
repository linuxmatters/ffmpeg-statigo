package main

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func TestProfilePaths(t *testing.T) {
	for _, embedded := range []bool{false, true} {
		for _, targetOS := range []string{"linux", "darwin"} {
			for _, arch := range []string{"amd64", "arm64"} {
				root, output, err := profilePaths(embedded, targetOS, arch)
				if err != nil {
					t.Fatal(err)
				}
				rootSuffix, outputSuffix := ".build", "lib"
				if embedded {
					rootSuffix = filepath.Join(rootSuffix, "embedded")
					outputSuffix = filepath.Join(outputSuffix, "embedded")
				}
				wantRoot, _ := filepath.Abs(rootSuffix)
				wantOutput, _ := filepath.Abs(filepath.Join(outputSuffix, targetOS+"_"+arch, "libffmpeg.a"))
				if root != wantRoot || output != wantOutput {
					t.Fatalf("profilePaths(%v, %s, %s) = %s, %s", embedded, targetOS, arch, root, output)
				}
				if got := stagingDir(filepath.Join(root, "build", "openh264")); got != filepath.Join(root, "staging") {
					t.Fatalf("derived staging path = %s", got)
				}
			}
		}
	}
}

func TestEmbeddedLibrarySelection(t *testing.T) {
	excluded := []string{"dav1d", "glslang", "libdrm", "libsrt", "libva", "libvpl", "openssl", "rav1e", "x264", "x265", "nv-codec-headers", "Vulkan-Headers"}
	libs := librariesForProfile(true, filepath.Join(t.TempDir(), "staging"))
	seen := make(map[*Library]bool)
	for _, lib := range libs {
		if slices.Contains(excluded, lib.Name) {
			t.Errorf("excluded library %s remains", lib.Name)
		}
		for _, dep := range lib.Dependencies {
			if !seen[dep] {
				t.Errorf("%s depends on missing or later library %s", lib.Name, dep.Name)
			}
		}
		seen[lib] = true
		if lib.Enabled == nil || *lib.Enabled {
			if _, ok := expectedDigest(lib.URL); !ok {
				t.Errorf("missing digest for %s", lib.Name)
			}
		}
	}
	for _, lib := range AllLibraries {
		if lib != ffmpeg && !slices.Contains(excluded, lib.Name) && !seen[lib] {
			t.Errorf("remaining library %s was removed", lib.Name)
		}
	}
	if !seen[openh264] || libs[len(libs)-1].Name != "ffmpeg" {
		t.Fatal("OpenH264 must precede FFmpeg")
	}
	if vvenc.ShouldBuild() {
		t.Error("disabled vvenc was enabled")
	}
}

func featureEnabled(args []string, kind, name string) bool {
	for _, arg := range args {
		if value, ok := strings.CutPrefix(arg, "--enable-"+kind+"="); ok && slices.Contains(strings.Split(value, ","), name) {
			return true
		}
	}
	return false
}

func TestProfileConfigureArgs(t *testing.T) {
	stage, _ := filepath.Abs(".build/staging")
	original := ffmpeg.ConfigureArgs("linux")
	for _, targetOS := range []string{"linux", "darwin"} {
		for range 3 {
			libs := librariesForProfile(false, stage)
			want := ffmpeg.ConfigureArgs(targetOS)
			for _, lib := range AllLibraries {
				if lib == ffmpeg || (lib.Enabled != nil && !*lib.Enabled) || (len(lib.Platform) != 0 && !slices.Contains(lib.Platform, targetOS)) {
					continue
				}
				for _, enable := range lib.FFmpegEnables {
					want = append(want, "--enable-"+enable)
				}
			}
			if targetOS == "darwin" {
				want = append(want, "--enable-avfoundation", "--enable-audiotoolbox", "--enable-videotoolbox")
			}
			got := libs[len(libs)-1].ConfigureArgs(targetOS)
			if !slices.Equal(got, want) {
				t.Errorf("default %s configure arguments changed", targetOS)
			}

			embeddedStage, _ := filepath.Abs(".build/embedded/staging")
			embedded := librariesForProfile(true, embeddedStage)
			args := embedded[len(embedded)-1].ConfigureArgs(targetOS)
			for _, arg := range args {
				if !strings.HasPrefix(arg, "--enable-") {
					continue
				}
				for _, forbidden := range []string{"gpl", "version3", "nonfree", "av1", "dav1d", "avif", "obu", "vulkan", "nvenc", "nvdec", "cuvid", "ffnvcodec", "vaapi", "qsv", "vpl", "v4l2", "videotoolbox", "audiotoolbox", "avfoundation", "openssl", "libsrt", "libx264", "libx265", "glslang", "hwaccel"} {
					if strings.Contains(arg, forbidden) {
						t.Errorf("embedded %s enables forbidden feature: %s", targetOS, arg)
					}
				}
			}
			for _, required := range [][2]string{{"encoder", "libopenh264"}, {"encoder", "aac"}, {"muxer", "mp4"}, {"decoder", "hevc"}, {"encoder", "libvpx_vp9"}, {"encoder", "libopus"}, {"bsf", "extract_extradata"}} {
				if !featureEnabled(args, required[0], required[1]) {
					t.Errorf("missing %s %s", required[0], required[1])
				}
			}
			for _, required := range []string{"--enable-libopenh264", "--disable-hwaccels", "--disable-gpl", "--disable-version3", "--extra-cflags=-I" + filepath.Join(embeddedStage, "include"), "--extra-ldflags=-L" + filepath.Join(embeddedStage, "lib")} {
				if !slices.Contains(args, required) {
					t.Errorf("missing %s", required)
				}
			}
		}
	}
	if !slices.Equal(original, ffmpeg.ConfigureArgs("linux")) {
		t.Fatal("profile configuration mutated global FFmpeg arguments")
	}
}

func TestEmbeddedCleanIsolation(t *testing.T) {
	t.Chdir(t.TempDir())
	paths := []string{
		".build/src/ffmpeg/sentinel", ".build/build/ffmpeg/sentinel", ".build/staging/lib/libavcodec.a",
		".build/embedded/src/ffmpeg/sentinel", ".build/embedded/build/ffmpeg/sentinel", ".build/embedded/staging/lib/libavcodec.a",
	}
	for _, path := range paths {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("sentinel"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })
	os.Args = []string{"builder", "ffmpeg", "--clean", "--embedded"}
	if err := run(); err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		wantExists := !strings.Contains(path, "/embedded/")
		if fileExists(path) != wantExists {
			t.Errorf("%s existence differs from %v", path, wantExists)
		}
	}
}

func TestEmbeddedRejectsExcludedSelection(t *testing.T) {
	t.Chdir(t.TempDir())
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })
	os.Args = []string{"builder", "--embedded", "--clean", "x264"}
	if err := run(); err == nil || !strings.Contains(err.Error(), "excluded") {
		t.Fatalf("excluded selection error = %v", err)
	}
}

func TestCombineLibrariesRejectsMissingArchives(t *testing.T) {
	for _, makeDir := range []bool{false, true} {
		stage := t.TempDir()
		if makeDir {
			if err := os.MkdirAll(filepath.Join(stage, "lib"), 0o755); err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{"libpresent.a", "libunexpected.a"} {
				if err := os.WriteFile(filepath.Join(stage, "lib", name), []byte("archive"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
		}
		libs := []*Library{
			{Name: "test", LinkLibs: []string{"libpresent", "libmissing"}},
			{Name: "disabled", Enabled: Disabled(), LinkLibs: []string{"libdisabled"}},
			{Name: "other-platform", Platform: []string{"not-" + runtime.GOOS}, LinkLibs: []string{"libother"}},
		}
		output := filepath.Join(stage, "output.a")
		if err := os.WriteFile(output, []byte("original"), 0o644); err != nil {
			t.Fatal(err)
		}
		err := combineLibraries(context.Background(), libs, stage, output)
		if err == nil || !strings.Contains(err.Error(), "libmissing.a") || strings.Contains(err.Error(), "libdisabled") || strings.Contains(err.Error(), "libother") {
			t.Fatalf("combine error = %v", err)
		}
		data, err := os.ReadFile(output)
		if err != nil || string(data) != "original" {
			t.Fatal("failed combination changed existing output")
		}
	}
}

func TestOpenH264StaticBuild(t *testing.T) {
	build, ok := openh264.BuildSystem.(*MakefileBuild)
	if !ok || !slices.Equal(build.Targets, []string{"libopenh264.a"}) || build.InstallFunc == nil {
		t.Fatal("OpenH264 must build and install the static archive")
	}
	for _, targetOS := range []string{"linux", "darwin"} {
		args := openh264.ConfigureArgs(targetOS)
		for _, arg := range []string{"ENABLEPIC=Yes", "HAVE_GMP_API=No", "HAVE_GTEST=No"} {
			if !slices.Contains(args, arg) {
				t.Errorf("OpenH264 %s is missing %s", targetOS, arg)
			}
		}
		if targetOS == "darwin" && !slices.Contains(args, "STATIC_LDFLAGS=-lc++") {
			t.Fatal("macOS static pkg-config must link libc++")
		}
	}
}
