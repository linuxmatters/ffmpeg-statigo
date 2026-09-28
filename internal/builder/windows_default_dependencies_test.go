package main

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestWindowsDefaultDependencyOptions(t *testing.T) {
	for _, arch := range []string{"amd64", "386"} {
		t.Run(arch, func(t *testing.T) {
			opensslTarget, alignment, host := "mingw64", "16", "x86_64-w64-mingw32"
			if arch == "386" {
				opensslTarget, alignment, host = "mingw", "4", "i686-w64-mingw32"
			}
			for _, tc := range []struct {
				lib  *Library
				want []string
			}{
				{dav1d, []string{"--default-library=static", "-Denable_tools=false", "-Denable_tests=false", "-Dstack_alignment=" + alignment}},
				{glslang, []string{"-DBUILD_SHARED_LIBS=OFF", "-DENABLE_PCH=OFF", "-DSPIRV_WERROR=OFF"}},
				{openssl, []string{"no-shared", "no-apps", "no-tests", opensslTarget}},
				{libsrt, []string{"-DENABLE_SHARED=OFF", "-DENABLE_STATIC=ON", "-DSRT_USE_OPENSSL_STATIC_LIBS=ON", "-DUSE_OPENSSL_PC=ON", "-DENABLE_STDCXX_SYNC=ON"}},
				{x264, []string{"--enable-static", "--disable-cli", "--host=" + host}},
				{x265, []string{"-DENABLE_SHARED=OFF", "-DENABLE_CLI=OFF"}},
			} {
				t.Run(tc.lib.Name, func(t *testing.T) {
					original := slices.Clone(tc.lib.ConfigureArgs("windows"))
					lib := windowsLibrary(tc.lib, arch)
					args := lib.ConfigureArgs("windows")
					for _, want := range tc.want {
						if !slices.Contains(args, want) {
							t.Errorf("missing %s in %v", want, args)
						}
					}
					if tc.lib == x265 && slices.Contains(args, "-DENABLE_ASSEMBLY=OFF") != (arch == "386") {
						t.Errorf("unexpected x265 assembly configuration: %v", args)
					}
					if lib.URL != tc.lib.URL || !slices.Equal(lib.LinkLibs, tc.lib.LinkLibs) || !slices.Equal(lib.FFmpegEnables, tc.lib.FFmpegEnables) {
						t.Error("adapter changed the source archive or codec selection")
					}
					if _, ok := expectedDigest(lib.URL); !ok {
						t.Error("missing source digest")
					}
					if !slices.Equal(original, tc.lib.ConfigureArgs("windows")) {
						t.Error("adapter changed the global definition")
					}
					for _, targetOS := range []string{"linux", "darwin"} {
						if !slices.Equal(lib.ConfigureArgs(targetOS), tc.lib.ConfigureArgs(targetOS)) {
							t.Errorf("Windows arguments reached %s", targetOS)
						}
					}
				})
			}
			lib := windowsLibrary(x264, arch)
			build := lib.BuildSystem.(*AutoconfBuild)
			name, args := configureCommand("windows", build.Shell, "./configure", lib.ConfigureArgs("windows"))
			if name != "bash" || args[0] != "./configure" {
				t.Fatalf("x264 configure command = %s %v", name, args)
			}
		})
	}
	for _, lib := range []*Library{libdrm, libva, libvpl, nvcodecheaders} {
		if windowsLibrary(lib, "amd64") != lib || !slices.Equal(lib.Platform, []string{"linux"}) {
			t.Errorf("changed the Linux-only gate for %s", lib.Name)
		}
	}
}

func TestDefaultGlslangShaderCompiler(t *testing.T) {
	baseArgs := []string{
		"-DBUILD_SHARED_LIBS=OFF",
		"-DENABLE_GLSLANG_BINARIES=OFF",
		"-DENABLE_HLSL=OFF",
		"-DGLSLANG_TESTS=OFF",
		"-DSPIRV_SKIP_EXECUTABLES=ON",
		"-DSPIRV_SKIP_TESTS=ON",
	}
	for _, targetOS := range []string{"windows", "linux", "darwin"} {
		for _, arch := range []string{"amd64", "386"} {
			t.Run(targetOS+"/"+arch, func(t *testing.T) {
				stage := filepath.Join(t.TempDir(), "staging")
				libs := librariesForPlatform(false, stage, targetOS, arch)
				index := slices.IndexFunc(libs, func(lib *Library) bool { return lib.Name == "glslang" })
				if index < 0 {
					t.Fatal("default profile lacks glslang")
				}
				lib := libs[index]
				want := slices.Clone(baseArgs)
				if targetOS == "windows" {
					want[1] = "-DENABLE_GLSLANG_BINARIES=ON"
					want = append(want, "-DGLSLANG_ENABLE_INSTALL=ON", "-DCMAKE_INSTALL_BINDIR=bin", "-DENABLE_PCH=OFF", "-DSPIRV_WERROR=OFF")
					if lib == glslang {
						t.Fatal("Windows shader compiler needs an isolated library adapter")
					}
					env := buildPathEnv([]string{"Path=system-bin"}, stage, targetOS)
					if !slices.Contains(env, "PATH="+buildToolPath(filepath.Join(stage, "bin"), targetOS)+";system-bin") {
						t.Errorf("shader compiler install directory is absent from PATH: %v", env)
					}
				}
				if args := lib.ConfigureArgs(targetOS); !slices.Equal(args, want) {
					t.Errorf("glslang arguments = %v, want %v", args, want)
				}
				if !slices.Equal(glslang.ConfigureArgs(targetOS), baseArgs) {
					t.Error("adapter changed the global glslang arguments")
				}
				if _, ok := lib.BuildSystem.(*CMakeBuild); !ok {
					t.Error("glslang must use the CMake build and install steps")
				}
				configuredFFmpeg := libs[len(libs)-1]
				args := configuredFFmpeg.ConfigureArgs(targetOS)
				if !slices.Contains(configuredFFmpeg.Dependencies, lib) {
					t.Error("FFmpeg does not depend on the configured glslang library")
				}
				for _, flag := range []string{"--enable-libglslang", "--enable-vulkan", "--enable-static"} {
					if !slices.Contains(args, flag) {
						t.Errorf("FFmpeg lacks %s", flag)
					}
				}
				if slices.Contains(args, "--enable-shared") {
					t.Error("default profile must not enable shared libraries")
				}
				if !featureEnabled(args, "encoder", "ffv1_vulkan") || featureConfigured(args, "disable", "encoder", "ffv1_vulkan") {
					t.Error("default profile must retain the ffv1_vulkan encoder")
				}
				for _, lib := range librariesForPlatform(true, stage, targetOS, arch) {
					if lib.Name == "glslang" {
						t.Error("embedded profile must exclude glslang")
					}
				}
			})
		}
	}
}

func TestRav1eWindowsInstallCommand(t *testing.T) {
	for _, tc := range []struct {
		arch, target, host, flags string
	}{
		{"amd64", "x86_64-pc-windows-gnu", "x86_64-w64-mingw32", "-C target-cpu=x86-64-v3 -C link-self-contained=no"},
		{"386", "i686-pc-windows-gnu", "i686-w64-mingw32", "-C link-self-contained=no"},
	} {
		t.Run(tc.arch, func(t *testing.T) {
			args, env := rav1eInstallCommand("windows", tc.arch, `C:\work\staging`, "-isysroot /unused")
			wantArgs := []string{"cinstall", "--prefix=C:/work/staging", "--libdir=lib", "--library-type=staticlib", "--crt-static", "--release", "--no-default-features", "--features=asm,threading", "--target=" + tc.target}
			wantEnv := []string{
				"CARGO_TARGET_" + strings.ToUpper(strings.ReplaceAll(tc.target, "-", "_")) + "_LINKER=" + tc.host + "-gcc",
				"WINAPI_NO_BUNDLED_LIBRARIES=1",
				"RUSTFLAGS=" + tc.flags,
				"CARGO_PROFILE_RELEASE_DEBUG=false",
			}
			if !slices.Equal(args, wantArgs) || !slices.Equal(env, wantEnv) {
				t.Fatalf("command = %v, env = %v; want %v, %v", args, env, wantArgs, wantEnv)
			}
			lib := windowsLibrary(rav1e, tc.arch)
			if lib == rav1e || lib.BuildSystem.(*CargoBuild).InstallFunc == nil || lib.URL != rav1e.URL {
				t.Fatal("missing isolated rav1e build adapter")
			}
		})
	}
}

func TestRav1eNonWindowsInstallCommand(t *testing.T) {
	for _, targetOS := range []string{"linux", "darwin"} {
		for _, arch := range []string{"amd64", "arm64", "386"} {
			args, env := rav1eInstallCommand(targetOS, arch, "/stage", "-isysroot /sdk")
			wantArgs := []string{"cinstall", "--prefix=/stage", "--libdir=lib", "--library-type=staticlib", "--crt-static", "--release", "--no-default-features", "--features=asm,threading"}
			flags := ""
			if arch == "amd64" {
				flags = "-C target-cpu=x86-64-v3"
			}
			if targetOS == "darwin" {
				flags += " -C link-arg=-L" + filepath.Join("/sdk", "usr", "lib")
			}
			wantEnv := []string{"RUSTFLAGS=" + flags, "CARGO_PROFILE_RELEASE_DEBUG=false"}
			if !slices.Equal(args, wantArgs) || !slices.Equal(env, wantEnv) {
				t.Errorf("%s/%s command = %v, env = %v; want %v, %v", targetOS, arch, args, env, wantArgs, wantEnv)
			}
		}
	}
}

func TestWindowsMesonCommandPaths(t *testing.T) {
	lib := windowsLibrary(dav1d, "386")
	args := mesonConfigureArgs(lib, `C:\work\src`, `C:\work\build`, `C:\work\staging`, "windows")
	want := []string{"setup", "C:/work/build", "C:/work/src", "--prefix=C:/work/staging", "--buildtype=release", "--default-library=static", "--libdir=lib"}
	want = append(want, lib.ConfigureArgs("windows")...)
	if !slices.Equal(args, want) {
		t.Fatalf("meson arguments = %v, want %v", args, want)
	}
	for _, targetOS := range []string{"linux", "darwin"} {
		args := mesonConfigureArgs(dav1d, "/src", "/build", "/stage", targetOS)
		if !slices.Equal(args[:4], []string{"setup", "/build", "/src", "--prefix=/stage"}) {
			t.Errorf("%s paths changed: %v", targetOS, args)
		}
	}
}
