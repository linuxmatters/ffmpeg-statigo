package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func TestWindowsEmbeddedLibraries(t *testing.T) {
	for _, tc := range []struct {
		arch, ffmpegArch, vpxTarget, openh264Arch string
	}{
		{"amd64", "x86_64", "x86_64-win64-gcc", "x86_64"},
		{"386", "x86", "x86-win32-gcc", "i386"},
	} {
		t.Run(tc.arch, func(t *testing.T) {
			libs := librariesForPlatform(true, t.TempDir(), "windows", tc.arch)
			want := []string{"ffmpeg", "lame", "libvpx", "libwebp", "libxml2", "openh264", "opus", "zimg", "zlib"}
			var names []string
			seen := make(map[*Library]bool)
			for _, lib := range libs {
				for _, dep := range lib.Dependencies {
					if !seen[dep] {
						t.Errorf("%s depends on missing or later library %s", lib.Name, dep.Name)
					}
				}
				seen[lib] = true
				if (lib.Enabled != nil && !*lib.Enabled) || (len(lib.Platform) != 0 && !slices.Contains(lib.Platform, "windows")) {
					continue
				}
				names = append(names, lib.Name)
				if _, ok := expectedDigest(lib.URL); !ok {
					t.Errorf("missing pinned digest for %s", lib.Name)
				}
				args := lib.ConfigureArgs("windows")
				for _, arg := range args {
					if arg == "--enable-cross-compile" || strings.HasPrefix(arg, "--cross-prefix=") || strings.HasPrefix(arg, "--host=") {
						t.Errorf("%s gained cross-compilation argument %s", lib.Name, arg)
					}
				}
				var required []string
				switch lib.Name {
				case "ffmpeg":
					required = []string{"--target-os=mingw32", "--arch=" + tc.ffmpegArch, "--pkg-config-flags=--static", "--enable-libopenh264", "--disable-autodetect", "--disable-gpl", "--disable-version3", "--disable-nonfree", "--disable-shared", "--disable-hwaccels"}
				case "libxml2":
					required = []string{"--without-iconv", "--with-zlib", "--disable-shared"}
				case "libvpx":
					required = []string{"--target=" + tc.vpxTarget, "--enable-static", "--disable-shared"}
				case "openh264":
					required = []string{"OS=mingw_nt", "ARCH=" + tc.openh264Arch, "HAVE_GMP_API=No", "HAVE_GTEST=No"}
				}
				for _, arg := range required {
					if !slices.Contains(args, arg) {
						t.Errorf("%s is missing %s", lib.Name, arg)
					}
				}
			}
			slices.Sort(names)
			if !slices.Equal(names, want) {
				t.Fatalf("active libraries = %v, want %v", names, want)
			}
			if libs[len(libs)-1].Name != "ffmpeg" {
				t.Fatal("FFmpeg must remain last")
			}
		})
	}
}

func TestWindowsFFmpegPaths(t *testing.T) {
	for _, arch := range []string{"amd64", "386"} {
		args := ffmpegConfigureArgsForPlatform("windows", arch, `C:\work\staging`, true)
		for _, want := range []string{"--extra-cflags=-IC:/work/staging/include", "--extra-ldflags=-LC:/work/staging/lib"} {
			if !slices.Contains(args, want) {
				t.Errorf("%s missing %s", arch, want)
			}
		}
		var cflags []string
		for _, arg := range args {
			if strings.HasPrefix(arg, "--extra-cflags=") {
				cflags = append(cflags, arg)
			}
		}
		want := []string{"--extra-cflags=-IC:/work/staging/include"}
		if arch == "386" {
			want = append(want, "--extra-cflags=-mpreferred-stack-boundary=4")
		}
		if !slices.Equal(cflags, want) {
			t.Errorf("%s extra C flags = %v, want %v", arch, cflags, want)
		}
	}
}

func TestWindowsProfileIsolation(t *testing.T) {
	stage := t.TempDir()
	originals := make(map[string]*Library)
	for _, lib := range append(slices.Clone(AllLibraries), openh264) {
		originals[lib.Name] = lib
	}
	for range 3 {
		_ = librariesForPlatform(true, stage, "windows", "amd64")
		_ = librariesForPlatform(true, stage, "windows", "386")
		for _, targetOS := range []string{"linux", "darwin", "windows"} {
			for _, arch := range []string{"amd64", "arm64", "386"} {
				for _, embedded := range []bool{false, true} {
					if embedded && targetOS == "windows" && (arch == "amd64" || arch == "386") {
						continue
					}
					for _, lib := range librariesForPlatform(embedded, stage, targetOS, arch) {
						if lib.Name != "ffmpeg" && lib != originals[lib.Name] {
							t.Errorf("%s %s embedded=%v changed %s", targetOS, arch, embedded, lib.Name)
						}
						if lib.ConfigureArgs == nil {
							continue
						}
						for _, arg := range lib.ConfigureArgs(targetOS) {
							if slices.Contains([]string{"--without-iconv", "--target=x86_64-win64-gcc", "--target=x86-win32-gcc", "--target-os=mingw32", "--arch=x86_64", "--arch=x86", "--extra-cflags=-mpreferred-stack-boundary=4", "OS=mingw_nt", "ARCH=x86_64", "ARCH=i386"}, arg) {
								t.Errorf("%s %s embedded=%v gained %s", targetOS, arch, embedded, arg)
							}
						}
					}
				}
			}
		}
	}
}

func TestWindowsOpenH264InstallTarget(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the command stub uses a POSIX executable script")
	}
	dir := t.TempDir()
	capture := filepath.Join(dir, "make-args")
	t.Setenv("MAKE_ARGS", capture)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	if err := os.WriteFile(filepath.Join(dir, "make"), []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$MAKE_ARGS\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, arch := range []string{"amd64", "386"} {
		lib := windowsOpenH264(arch)
		build := lib.BuildSystem.(*MakefileBuild)
		if !slices.Equal(build.Targets, []string{"libopenh264.a"}) {
			t.Fatalf("build targets = %v", build.Targets)
		}
		stage := `C:\work\staging`
		if err := build.InstallFunc(context.Background(), dir, stage); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(capture)
		if err != nil {
			t.Fatal(err)
		}
		want := append(lib.ConfigureArgs("windows"), "PREFIX=C:/work/staging", "install-static")
		if got := strings.Split(strings.TrimSpace(string(data)), "\n"); !slices.Equal(got, want) {
			t.Fatalf("install arguments = %v, want %v", got, want)
		}
	}
}

func TestWindowsZimgAutogen(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh is required for the script fixture")
	}
	dir := filepath.Join(t.TempDir(), "source with spaces")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string]string{
		"autogen.sh": "printf '%s' \"$0\" > invoked-script\n",
		"zimg.pc.in": "Libs.private: @STL_LIBS@\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	lib := windowsEmbeddedLibrary(zimg, "amd64")
	if err := lib.PostExtract(context.Background(), dir); err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]string{
		"invoked-script": "./autogen.sh",
		"zimg.pc.in":     "Libs.private: @STL_LIBS@ -lm\n",
	} {
		got, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil || string(got) != want {
			t.Errorf("%s = %q, %v, want %q", name, got, err, want)
		}
	}
}
