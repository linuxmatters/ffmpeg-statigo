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

func TestWindowsLibraries(t *testing.T) {
	for _, tc := range []struct {
		arch, ffmpegArch, vpxTarget, openh264Arch string
	}{
		{"amd64", "x86_64", "x86_64-win64-gcc", "x86_64"},
		{"386", "x86", "x86-win32-gcc", "i386"},
	} {
		for _, embedded := range []bool{false, true} {
			profile := "default"
			if embedded {
				profile = "embedded"
			}
			t.Run(profile+"/"+tc.arch, func(t *testing.T) {
				libs := librariesForPlatform(embedded, t.TempDir(), "windows", tc.arch)
				want := []string{"ffmpeg", "lame", "libvpx", "libwebp", "libxml2", "opus", "zimg", "zlib"}
				if embedded {
					want = append(want, "openh264")
				} else {
					want = append(want, "Vulkan-Headers", "dav1d", "glslang", "libsrt", "openssl", "rav1e", "x264", "x265")
				}
				slices.Sort(want)
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
					var args []string
					if lib.ConfigureArgs != nil {
						args = lib.ConfigureArgs("windows")
					}
					for _, arg := range args {
						if arg == "--enable-cross-compile" || strings.HasPrefix(arg, "--cross-prefix=") || (strings.HasPrefix(arg, "--host=") && lib.Name != "x264") {
							t.Errorf("%s gained cross-compilation argument %s", lib.Name, arg)
						}
					}
					var required []string
					switch lib.Name {
					case "ffmpeg":
						required = []string{"--target-os=mingw32", "--arch=" + tc.ffmpegArch, "--pkg-config-flags=--static", "--disable-autodetect"}
						if embedded {
							required = append(required, "--enable-libopenh264", "--disable-gpl", "--disable-version3", "--disable-nonfree", "--disable-shared", "--disable-hwaccels")
						} else {
							required = append(required, "--enable-gpl", "--enable-version3", "--enable-libx264", "--enable-libx265", "--enable-libdav1d", "--enable-librav1e", "--enable-libsrt", "--enable-openssl", "--enable-vulkan", "--enable-libglslang")
						}
						for _, flag := range []string{"gpl", "version3", "libx264", "libx265", "libdav1d", "librav1e", "libsrt", "openssl", "vulkan", "libglslang"} {
							if slices.Contains(args, "--enable-"+flag) == embedded {
								t.Errorf("%s enable flag differs from embedded=%v", flag, embedded)
							}
						}
						if slices.Contains(args, "--enable-libopenh264") != embedded {
							t.Error("OpenH264 must remain embedded-only")
						}
						for _, flag := range []string{"vaapi", "libvpl", "cuvid", "ffnvcodec", "nvdec", "nvenc", "avfoundation", "audiotoolbox", "videotoolbox", "libvvenc"} {
							if slices.Contains(args, "--enable-"+flag) {
								t.Errorf("Windows enabled unavailable library %s", flag)
							}
						}
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
}

func TestWindowsFFmpegPaths(t *testing.T) {
	for _, embedded := range []bool{false, true} {
		for _, targetOS := range []string{"windows", "linux", "darwin"} {
			for _, arch := range []string{"amd64", "386", "arm64"} {
				stage := `C:\work\staging`
				args := ffmpegConfigureArgsForPlatform(targetOS, arch, stage, embedded)
				incDir, libDir := filepath.Join(stage, "include"), filepath.Join(stage, "lib")
				windowsNative := targetOS == "windows" && arch != "arm64"
				if windowsNative {
					incDir, libDir = "C:/work/staging/include", "C:/work/staging/lib"
				}
				if !slices.Contains(args, "--extra-ldflags=-L"+libDir) {
					t.Errorf("%s/%s embedded=%v missing library path %s", targetOS, arch, embedded, libDir)
				}
				var cflags, targetFlags []string
				for _, arg := range args {
					if strings.HasPrefix(arg, "--extra-cflags=") {
						cflags = append(cflags, arg)
					}
					if strings.HasPrefix(arg, "--target-os=") || strings.HasPrefix(arg, "--arch=") {
						targetFlags = append(targetFlags, arg)
					}
				}
				wantCflags := []string{"--extra-cflags=-I" + incDir}
				var wantTargetFlags []string
				if windowsNative {
					ffmpegArch := "x86_64"
					if arch == "386" {
						ffmpegArch = "x86"
						wantCflags = append(wantCflags, "--extra-cflags=-mpreferred-stack-boundary=4", "--extra-cflags=-mstackrealign")
					}
					wantTargetFlags = []string{"--target-os=mingw32", "--arch=" + ffmpegArch}
				}
				if !slices.Equal(cflags, wantCflags) || !slices.Equal(targetFlags, wantTargetFlags) {
					t.Errorf("%s/%s embedded=%v flags = %v %v, want %v %v", targetOS, arch, embedded, cflags, targetFlags, wantCflags, wantTargetFlags)
				}
				if slices.Contains(args, "--disable-asm") || slices.Contains(args, "--disable-x86asm") {
					t.Errorf("%s/%s embedded=%v disabled assembly", targetOS, arch, embedded)
				}
			}
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
		for _, embedded := range []bool{false, true} {
			_ = librariesForPlatform(embedded, stage, "windows", "amd64")
			_ = librariesForPlatform(embedded, stage, "windows", "386")
		}
		for _, targetOS := range []string{"linux", "darwin", "windows"} {
			for _, arch := range []string{"amd64", "arm64", "386"} {
				for _, embedded := range []bool{false, true} {
					if targetOS == "windows" && (arch == "amd64" || arch == "386") {
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
							if slices.Contains([]string{"--without-iconv", "--target=x86_64-win64-gcc", "--target=x86-win32-gcc", "--target-os=mingw32", "--arch=x86_64", "--arch=x86", "--extra-cflags=-mpreferred-stack-boundary=4", "--extra-cflags=-mstackrealign", "OS=mingw_nt", "ARCH=x86_64", "ARCH=i386"}, arg) {
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
	lib := windowsLibrary(zimg, "amd64")
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
