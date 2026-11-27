package main

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestParseBuilderArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want builderOptions
	}{
		{
			name: "default empty args",
			args: nil,
			want: builderOptions{
				selectedLibs: map[string]bool{},
			},
		},
		{
			name: "programs flag only",
			args: []string{"--programs"},
			want: builderOptions{
				programs:     true,
				selectedLibs: map[string]bool{},
			},
		},
		{
			name: "embedded and programs",
			args: []string{"--embedded", "--programs"},
			want: builderOptions{
				embedded:     true,
				programs:     true,
				selectedLibs: map[string]bool{},
			},
		},
		{
			name: "programs with library names",
			args: []string{"--programs", "ffmpeg", "x264"},
			want: builderOptions{
				programs:     true,
				selectedLibs: map[string]bool{"ffmpeg": true, "x264": true},
			},
		},
		{
			name: "embedded programs clean with library",
			args: []string{"--embedded", "--programs", "--clean", "ffmpeg"},
			want: builderOptions{
				embedded:     true,
				programs:     true,
				clean:        true,
				selectedLibs: map[string]bool{"ffmpeg": true},
			},
		},
		{
			name: "programs with list",
			args: []string{"--programs", "--list"},
			want: builderOptions{
				programs:     true,
				list:         true,
				selectedLibs: map[string]bool{},
			},
		},
		{
			name: "embedded programs with update-digests",
			args: []string{"--embedded", "--programs", "--update-digests"},
			want: builderOptions{
				embedded:      true,
				programs:      true,
				updateDigests: true,
				selectedLibs:  map[string]bool{},
			},
		},
		{
			name: "other flags and libs without programs",
			args: []string{"--clean", "--list", "--update-digests", "zlib", "opus"},
			want: builderOptions{
				clean:         true,
				list:          true,
				updateDigests: true,
				selectedLibs:  map[string]bool{"zlib": true, "opus": true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseBuilderArgs(tt.args)
			if got.embedded != tt.want.embedded ||
				got.programs != tt.want.programs ||
				got.clean != tt.want.clean ||
				got.list != tt.want.list ||
				got.updateDigests != tt.want.updateDigests {
				t.Errorf("parseBuilderArgs(%v) flags = %+v, want %+v", tt.args, got, tt.want)
			}
			if !maps.Equal(got.selectedLibs, tt.want.selectedLibs) {
				t.Errorf("parseBuilderArgs(%v).selectedLibs = %v, want %v", tt.args, got.selectedLibs, tt.want.selectedLibs)
			}
		})
	}
}

func TestProgramsOutputDir(t *testing.T) {
	for _, embedded := range []bool{false, true} {
		for _, targetOS := range []string{"linux", "darwin", "windows"} {
			for _, arch := range []string{"amd64", "arm64", "386"} {
				got, err := programsOutputDir(embedded, targetOS, arch)
				if err != nil {
					t.Fatalf("programsOutputDir(%v, %s, %s) error = %v", embedded, targetOS, arch, err)
				}
				rel := filepath.Join("bin", targetOS+"_"+arch)
				if embedded {
					rel = filepath.Join("bin", "embedded", targetOS+"_"+arch)
				}
				want, err := filepath.Abs(rel)
				if err != nil {
					t.Fatalf("filepath.Abs(%q) error = %v", rel, err)
				}
				if got != want {
					t.Errorf("programsOutputDir(%v, %s, %s) = %q, want %q", embedded, targetOS, arch, got, want)
				}
			}
		}
	}
}

func TestFFmpegProgramsConfigureArgs(t *testing.T) {
	origLinuxArgs := slices.Clone(ffmpeg.ConfigureArgs("linux"))
	origBinaries := slices.Clone(ffmpeg.Binaries)

	platforms := []struct {
		targetOS        string
		arch            string
		wantLDExeFlags  string
		wantWin32CFlags bool
	}{
		{targetOS: "linux", arch: "amd64", wantLDExeFlags: "--extra-ldexeflags=-static-libstdc++ -static-libgcc"},
		{targetOS: "linux", arch: "arm64", wantLDExeFlags: "--extra-ldexeflags=-static-libstdc++ -static-libgcc"},
		{targetOS: "windows", arch: "amd64", wantLDExeFlags: "--extra-ldexeflags=-static -static-libgcc -static-libstdc++"},
		{targetOS: "windows", arch: "386", wantLDExeFlags: "--extra-ldexeflags=-static -static-libgcc -static-libstdc++", wantWin32CFlags: true},
		{targetOS: "darwin", arch: "arm64", wantLDExeFlags: ""},
		{targetOS: "darwin", arch: "amd64", wantLDExeFlags: ""},
	}

	hasLDExeFlags := func(args []string) bool {
		for _, arg := range args {
			if strings.HasPrefix(arg, "--extra-ldexeflags=") {
				return true
			}
		}
		return false
	}

	for _, embedded := range []bool{false, true} {
		profile := "default"
		if embedded {
			profile = "embedded"
		}
		for _, p := range platforms {
			t.Run(fmt.Sprintf("%s/%s_%s", profile, p.targetOS, p.arch), func(t *testing.T) {
				stage := filepath.Join(t.TempDir(), "staging")

				// programs=false keeps --disable-programs and emits no --extra-ldexeflags.
				noProgLibs := librariesForPlatformWithPrograms(embedded, false, stage, p.targetOS, p.arch)
				noProgFFmpeg := noProgLibs[len(noProgLibs)-1]
				noProgArgs := noProgFFmpeg.ConfigureArgs(p.targetOS)

				if !slices.Contains(noProgArgs, "--disable-programs") {
					t.Errorf("programs=false missing --disable-programs")
				}
				for _, forbidden := range []string{"--enable-ffmpeg", "--enable-ffprobe", "--disable-ffplay"} {
					if slices.Contains(noProgArgs, forbidden) {
						t.Errorf("programs=false unexpectedly contains %s", forbidden)
					}
				}
				if hasLDExeFlags(noProgArgs) {
					t.Errorf("programs=false unexpectedly contains --extra-ldexeflags")
				}
				if len(noProgFFmpeg.Binaries) != 0 {
					t.Errorf("programs=false Binaries = %v, want empty", noProgFFmpeg.Binaries)
				}

				// programs=true enables ffmpeg/ffprobe, disables ffplay, omits --disable-programs,
				// and sets per-OS --extra-ldexeflags.
				progLibs := librariesForPlatformWithPrograms(embedded, true, stage, p.targetOS, p.arch)
				progFFmpeg := progLibs[len(progLibs)-1]
				progArgs := progFFmpeg.ConfigureArgs(p.targetOS)

				for _, required := range []string{"--enable-ffmpeg", "--enable-ffprobe", "--disable-ffplay"} {
					if !slices.Contains(progArgs, required) {
						t.Errorf("programs=true missing %s", required)
					}
				}
				if slices.Contains(progArgs, "--disable-programs") {
					t.Errorf("programs=true must not contain --disable-programs")
				}
				if p.wantLDExeFlags != "" {
					if !slices.Contains(progArgs, p.wantLDExeFlags) {
						t.Errorf("programs=true missing %q", p.wantLDExeFlags)
					}
				} else if hasLDExeFlags(progArgs) {
					t.Errorf("programs=true on %s/%s must not emit --extra-ldexeflags", p.targetOS, p.arch)
				}
				if p.wantWin32CFlags {
					for _, required := range []string{"--extra-cflags=-mpreferred-stack-boundary=4", "--extra-cflags=-mstackrealign"} {
						if !slices.Contains(progArgs, required) {
							t.Errorf("windows/386 programs=true missing %s", required)
						}
					}
				}
				if !slices.Equal(progFFmpeg.Binaries, []string{"ffmpeg", "ffprobe"}) {
					t.Errorf("programs=true Binaries = %v, want [ffmpeg ffprobe]", progFFmpeg.Binaries)
				}

				// ConfigHash must differ between programs=false and programs=true.
				if noProgFFmpeg.ConfigHash() == progFFmpeg.ConfigHash() {
					t.Errorf("ConfigHash identical for programs=false and programs=true")
				}
			})
		}
	}

	if !slices.Equal(origLinuxArgs, ffmpeg.ConfigureArgs("linux")) {
		t.Error("global ffmpeg.ConfigureArgs was mutated")
	}
	if !slices.Equal(origBinaries, ffmpeg.Binaries) {
		t.Errorf("global ffmpeg.Binaries mutated to %v, want %v", ffmpeg.Binaries, origBinaries)
	}
}

func TestScrubPkgConfigDir(t *testing.T) {
	t.Run("scrubs -lgcc_s and preserves untouched files", func(t *testing.T) {
		pcDir := t.TempDir()

		vplContent := "prefix=${pcfiledir}/../..\n" +
			"exec_prefix=${prefix}\n" +
			"libdir=${exec_prefix}/lib\n" +
			"includedir=${prefix}/include\n\n" +
			"Name: vpl\n" +
			"Description: Intel Video Processing Library\n" +
			"Version: 2.16.0\n" +
			"Libs: -L${libdir} -lvpl\n" +
			"Libs.private: -lstdc++ -ldl\n" +
			"Cflags: -I${includedir}\n"

		fixtures := map[string]struct {
			before string
			after  string
		}{
			"haisrt.pc": {
				before: "prefix=/stage\nName: haisrt\nLibs: -L${prefix}/lib -lsrt -lgcc_s\nLibs.private: -lcrypto -lstdc++ -lgcc_s -lgcc -lpthread -lm\n",
				after:  "prefix=/stage\nName: haisrt\nLibs: -L${prefix}/lib -lsrt\nLibs.private: -lcrypto -lstdc++ -lgcc -lpthread -lm\n",
			},
			"srt.pc": {
				before: "prefix=/stage\nName: srt\nLibs: -L${prefix}/lib -lsrt\nLibs.private: -lssl -lcrypto -lstdc++ -lgcc_s -lgcc\n",
				after:  "prefix=/stage\nName: srt\nLibs: -L${prefix}/lib -lsrt\nLibs.private: -lssl -lcrypto -lstdc++ -lgcc\n",
			},
			"rav1e.pc": {
				before: "prefix=/stage\nName: rav1e\nLibs: -L${prefix}/lib -lrav1e\nLibs.private: -lgcc_s -lutil -lrt -lpthread -lm -ldl -lc\n",
				after:  "prefix=/stage\nName: rav1e\nLibs: -L${prefix}/lib -lrav1e\nLibs.private: -lutil -lrt -lpthread -lm -ldl -lc\n",
			},
			"x265.pc": {
				before: "prefix=/stage\nName: x265\nLibs: -L${prefix}/lib -lx265 -lgcc_s\nLibs.private: -lstdc++ -lm -lgcc_s -lgcc -lrt -ldl\n",
				after:  "prefix=/stage\nName: x265\nLibs: -L${prefix}/lib -lx265\nLibs.private: -lstdc++ -lm -lgcc -lrt -ldl\n",
			},
			"vpl.pc": {
				before: vplContent,
				after:  vplContent,
			},
			"zimg.pc": {
				before: "prefix=/stage\nName: zimg\nLibs: -L${prefix}/lib -lzimg\nLibs.private: -lstdc++ -lgcc_s -lgcc\n",
				after:  "prefix=/stage\nName: zimg\nLibs: -L${prefix}/lib -lzimg\nLibs.private: -lstdc++ -lgcc\n",
			},
		}

		for name, f := range fixtures {
			p := filepath.Join(pcDir, name)
			if err := os.WriteFile(p, []byte(f.before), 0o644); err != nil {
				t.Fatalf("write %s: %v", name, err)
			}
		}

		// Set an older mtime on vpl.pc so any rewrite is detectable.
		vplPath := filepath.Join(pcDir, "vpl.pc")
		oldTime := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
		if err := os.Chtimes(vplPath, oldTime, oldTime); err != nil {
			t.Fatalf("chtimes vpl.pc: %v", err)
		}

		if err := scrubPkgConfigDir(pcDir); err != nil {
			t.Fatalf("scrubPkgConfigDir() error = %v", err)
		}

		mtimesAfterFirst := make(map[string]time.Time, len(fixtures))
		for name, f := range fixtures {
			p := filepath.Join(pcDir, name)
			got, err := os.ReadFile(p)
			if err != nil {
				t.Fatalf("read %s: %v", name, err)
			}
			if string(got) != f.after {
				t.Errorf("%s after scrub =\n%q\nwant:\n%q", name, string(got), f.after)
			}
			if strings.Contains(string(got), "-lgcc_s") {
				t.Errorf("%s still contains -lgcc_s", name)
			}
			info, err := os.Stat(p)
			if err != nil {
				t.Fatalf("stat %s: %v", name, err)
			}
			mtimesAfterFirst[name] = info.ModTime()
		}

		if !mtimesAfterFirst["vpl.pc"].Equal(oldTime) {
			t.Errorf("vpl.pc mtime changed from %v to %v", oldTime, mtimesAfterFirst["vpl.pc"])
		}

		// Second run must be idempotent and leave all files and mtimes untouched.
		if err := scrubPkgConfigDir(pcDir); err != nil {
			t.Fatalf("second scrubPkgConfigDir() error = %v", err)
		}
		for name, f := range fixtures {
			p := filepath.Join(pcDir, name)
			got, err := os.ReadFile(p)
			if err != nil {
				t.Fatalf("read %s after second scrub: %v", name, err)
			}
			if string(got) != f.after {
				t.Errorf("%s changed on second scrub", name)
			}
			info, err := os.Stat(p)
			if err != nil {
				t.Fatalf("stat %s after second scrub: %v", name, err)
			}
			if !info.ModTime().Equal(mtimesAfterFirst[name]) {
				t.Errorf("%s mtime changed on second scrub: %v -> %v", name, mtimesAfterFirst[name], info.ModTime())
			}
		}
	})

	t.Run("missing pkgconfig directory succeeds", func(t *testing.T) {
		absent := filepath.Join(t.TempDir(), "does-not-exist", "pkgconfig")
		if err := scrubPkgConfigDir(absent); err != nil {
			t.Errorf("scrubPkgConfigDir(absent) error = %v, want nil", err)
		}
	})

	for _, target := range []struct {
		targetOS string
		arch     string
	}{
		{targetOS: "linux", arch: "amd64"},
		{targetOS: "windows", arch: "amd64"},
	} {
		t.Run("FFmpeg PostExtract scrubs pkgconfig and caches licenses on "+target.targetOS+"/"+target.arch, func(t *testing.T) {
			root := t.TempDir()
			stage := filepath.Join(root, "staging")
			pcDir := filepath.Join(stage, "lib", "pkgconfig")
			if err := os.MkdirAll(pcDir, 0o755); err != nil {
				t.Fatalf("mkdir pkgconfig: %v", err)
			}
			pcFile := filepath.Join(pcDir, "x265.pc")
			if err := os.WriteFile(pcFile, []byte("Name: x265\nLibs: -lx265\nLibs.private: -lstdc++ -lgcc_s -lgcc\n"), 0o644); err != nil {
				t.Fatalf("write x265.pc: %v", err)
			}

			srcPath := filepath.Join(root, "src", "ffmpeg")
			if err := os.MkdirAll(filepath.Join(srcPath, "libavformat"), 0o755); err != nil {
				t.Fatalf("mkdir libavformat: %v", err)
			}
			if err := os.WriteFile(filepath.Join(srcPath, "libavformat", "tls_openssl.c"), []byte("X509_gmtime_adj(X509_get_notBefore(*cert), 0);\n"), 0o644); err != nil {
				t.Fatalf("write tls_openssl.c: %v", err)
			}
			const (
				gplText  = "GNU GENERAL PUBLIC LICENSE Version 3 fixture\n"
				lgplText = "GNU LESSER GENERAL PUBLIC LICENSE Version 2.1 fixture\n"
			)
			if err := os.WriteFile(filepath.Join(srcPath, "COPYING.GPLv3"), []byte(gplText), 0o644); err != nil {
				t.Fatalf("write COPYING.GPLv3: %v", err)
			}
			if err := os.WriteFile(filepath.Join(srcPath, "COPYING.LGPLv2.1"), []byte(lgplText), 0o644); err != nil {
				t.Fatalf("write COPYING.LGPLv2.1: %v", err)
			}

			libs := librariesForPlatformWithPrograms(false, true, stage, target.targetOS, target.arch)
			configuredFFmpeg := libs[len(libs)-1]
			if configuredFFmpeg.PostExtract == nil {
				t.Fatal("FFmpeg PostExtract is nil")
			}
			if err := configuredFFmpeg.PostExtract(t.Context(), srcPath); err != nil {
				t.Fatalf("PostExtract() error = %v", err)
			}

			gotPC, err := os.ReadFile(pcFile)
			if err != nil {
				t.Fatalf("read scrubbed x265.pc: %v", err)
			}
			wantPC := "Name: x265\nLibs: -lx265\nLibs.private: -lstdc++ -lgcc\n"
			if string(gotPC) != wantPC {
				t.Errorf("x265.pc = %q, want %q", string(gotPC), wantPC)
			}

			for name, want := range map[string]string{
				"COPYING.GPLv3":    gplText,
				"COPYING.LGPLv2.1": lgplText,
			} {
				cached := filepath.Join(stage, "share", "ffmpeg", name)
				got, err := os.ReadFile(cached)
				if err != nil {
					t.Fatalf("read cached %s: %v", name, err)
				}
				if string(got) != want {
					t.Errorf("cached %s = %q, want %q", name, string(got), want)
				}
			}
		})
	}
}

func TestLibiconvDisableShared(t *testing.T) {
	args := libiconv.ConfigureArgs("darwin")
	if !slices.Contains(args, "--enable-static") {
		t.Errorf("libiconv.ConfigureArgs(darwin) = %v, missing --enable-static", args)
	}
	if !slices.Contains(args, "--disable-shared") {
		t.Errorf("libiconv.ConfigureArgs(darwin) = %v, missing --disable-shared", args)
	}
}

func TestPruneStagingDylibs(t *testing.T) {
	t.Run("missing lib directory succeeds", func(t *testing.T) {
		absent := filepath.Join(t.TempDir(), "does-not-exist")
		if err := pruneStagingDylibs(absent); err != nil {
			t.Errorf("pruneStagingDylibs(absent) error = %v, want nil", err)
		}
	})

	t.Run("darwin FFmpeg PostExtract removes stale dylibs and preserves static archives", func(t *testing.T) {
		root := t.TempDir()
		stage := filepath.Join(root, "staging")
		libDir := filepath.Join(stage, "lib")
		if err := os.MkdirAll(libDir, 0o755); err != nil {
			t.Fatalf("mkdir staging/lib: %v", err)
		}

		regularDylib := filepath.Join(libDir, "libiconv.2.dylib")
		if err := os.WriteFile(regularDylib, []byte("macho"), 0o644); err != nil {
			t.Fatalf("write libiconv.2.dylib: %v", err)
		}
		symlinkDylib := filepath.Join(libDir, "libiconv.dylib")
		if err := os.Symlink("libiconv.2.dylib", symlinkDylib); err != nil {
			t.Fatalf("symlink libiconv.dylib: %v", err)
		}
		danglingDylib := filepath.Join(libDir, "libstale.dylib")
		if err := os.Symlink("missing.dylib", danglingDylib); err != nil {
			t.Fatalf("symlink libstale.dylib: %v", err)
		}
		staticArchive := filepath.Join(libDir, "libiconv.a")
		if err := os.WriteFile(staticArchive, []byte("!<arch>\n"), 0o644); err != nil {
			t.Fatalf("write libiconv.a: %v", err)
		}
		subdir := filepath.Join(libDir, "nested.dylib")
		if err := os.MkdirAll(subdir, 0o755); err != nil {
			t.Fatalf("mkdir nested.dylib: %v", err)
		}

		srcPath := filepath.Join(root, "src", "ffmpeg")
		if err := os.MkdirAll(filepath.Join(srcPath, "libavformat"), 0o755); err != nil {
			t.Fatalf("mkdir libavformat: %v", err)
		}
		if err := os.WriteFile(filepath.Join(srcPath, "libavformat", "tls_openssl.c"), []byte("X509_gmtime_adj(X509_get_notBefore(*cert), 0);\n"), 0o644); err != nil {
			t.Fatalf("write tls_openssl.c: %v", err)
		}

		libs := librariesForPlatformWithPrograms(false, true, stage, "darwin", "arm64")
		configuredFFmpeg := libs[len(libs)-1]
		if configuredFFmpeg.PostExtract == nil {
			t.Fatal("FFmpeg PostExtract is nil")
		}
		if err := configuredFFmpeg.PostExtract(t.Context(), srcPath); err != nil {
			t.Fatalf("PostExtract() error = %v", err)
		}

		for _, removed := range []string{regularDylib, symlinkDylib, danglingDylib} {
			if _, err := os.Lstat(removed); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("%s still exists (err = %v), want removed", filepath.Base(removed), err)
			}
		}
		if got, err := os.ReadFile(staticArchive); err != nil || string(got) != "!<arch>\n" {
			t.Errorf("libiconv.a = %q (err = %v), want preserved", string(got), err)
		}
		if info, err := os.Stat(subdir); err != nil || !info.IsDir() {
			t.Errorf("nested.dylib directory err = %v, want preserved", err)
		}
	})
}

func TestBuildEnvMacOSDeploymentTarget(t *testing.T) {
	installDir := filepath.Join(t.TempDir(), "staging")

	tests := []struct {
		name      string
		targetOS  string
		baseEnv   []string
		wantSet   bool
		wantValue string
	}{
		{
			name:      "darwin unset defaults to 13.0",
			targetOS:  "darwin",
			baseEnv:   []string{"PATH=/usr/bin"},
			wantSet:   true,
			wantValue: "13.0",
		},
		{
			name:      "darwin empty defaults to 13.0",
			targetOS:  "darwin",
			baseEnv:   []string{"PATH=/usr/bin", "MACOSX_DEPLOYMENT_TARGET="},
			wantSet:   true,
			wantValue: "13.0",
		},
		{
			name:      "darwin preserves caller-set 14.0",
			targetOS:  "darwin",
			baseEnv:   []string{"PATH=/usr/bin", "MACOSX_DEPLOYMENT_TARGET=14.0"},
			wantSet:   true,
			wantValue: "14.0",
		},
		{
			name:     "linux does not add MACOSX_DEPLOYMENT_TARGET",
			targetOS: "linux",
			baseEnv:  []string{"PATH=/usr/bin"},
			wantSet:  false,
		},
		{
			name:     "windows does not add MACOSX_DEPLOYMENT_TARGET",
			targetOS: "windows",
			baseEnv:  []string{"PATH=C:\\Windows\\System32"},
			wantSet:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := buildEnvForOS(tt.baseEnv, installDir, tt.targetOS)
			var found bool
			var got string
			for _, entry := range env {
				if val, ok := strings.CutPrefix(entry, "MACOSX_DEPLOYMENT_TARGET="); ok {
					found = true
					got = val
				}
			}
			if found != tt.wantSet {
				t.Fatalf("MACOSX_DEPLOYMENT_TARGET present = %v, want %v (env=%v)", found, tt.wantSet, env)
			}
			if tt.wantSet && got != tt.wantValue {
				t.Errorf("MACOSX_DEPLOYMENT_TARGET = %q, want %q", got, tt.wantValue)
			}
		})
	}
}

func TestStripAndSignBinary(t *testing.T) {
	type recordedCall struct {
		name string
		args []string
	}

	t.Run("linux runs strip --strip-all", func(t *testing.T) {
		t.Setenv("STRIP", "")
		t.Setenv("CODESIGN", "")
		var calls []recordedCall
		runner := func(_ context.Context, name string, args ...string) ([]byte, error) {
			calls = append(calls, recordedCall{name: name, args: slices.Clone(args)})
			return nil, nil
		}

		if err := stripAndSignBinaryWithRunner(t.Context(), "/tmp/ffmpeg", "linux", runner); err != nil {
			t.Fatalf("stripAndSignBinaryWithRunner() error = %v", err)
		}
		if len(calls) != 1 {
			t.Fatalf("call count = %d, want 1 (%+v)", len(calls), calls)
		}
		if calls[0].name != "strip" || !slices.Equal(calls[0].args, []string{"--strip-all", "/tmp/ffmpeg"}) {
			t.Errorf("call[0] = %s %v, want strip [--strip-all /tmp/ffmpeg]", calls[0].name, calls[0].args)
		}
	})

	t.Run("windows runs strip --strip-all", func(t *testing.T) {
		t.Setenv("STRIP", "")
		t.Setenv("CODESIGN", "")
		var calls []recordedCall
		runner := func(_ context.Context, name string, args ...string) ([]byte, error) {
			calls = append(calls, recordedCall{name: name, args: slices.Clone(args)})
			return nil, nil
		}

		if err := stripAndSignBinaryWithRunner(t.Context(), `C:\bin\ffmpeg.exe`, "windows", runner); err != nil {
			t.Fatalf("stripAndSignBinaryWithRunner() error = %v", err)
		}
		if len(calls) != 1 {
			t.Fatalf("call count = %d, want 1 (%+v)", len(calls), calls)
		}
		if calls[0].name != "strip" || !slices.Equal(calls[0].args, []string{"--strip-all", `C:\bin\ffmpeg.exe`}) {
			t.Errorf("call[0] = %s %v, want strip [--strip-all C:\\bin\\ffmpeg.exe]", calls[0].name, calls[0].args)
		}
	})

	t.Run("darwin strips before codesign", func(t *testing.T) {
		t.Setenv("STRIP", "")
		t.Setenv("CODESIGN", "")
		var calls []recordedCall
		runner := func(_ context.Context, name string, args ...string) ([]byte, error) {
			calls = append(calls, recordedCall{name: name, args: slices.Clone(args)})
			return nil, nil
		}

		if err := stripAndSignBinaryWithRunner(t.Context(), "/tmp/ffmpeg", "darwin", runner); err != nil {
			t.Fatalf("stripAndSignBinaryWithRunner() error = %v", err)
		}
		if len(calls) != 2 {
			t.Fatalf("call count = %d, want 2 (%+v)", len(calls), calls)
		}
		if calls[0].name != "strip" || !slices.Equal(calls[0].args, []string{"/tmp/ffmpeg"}) {
			t.Errorf("call[0] = %s %v, want strip [/tmp/ffmpeg]", calls[0].name, calls[0].args)
		}
		if calls[1].name != "codesign" || !slices.Equal(calls[1].args, []string{"--force", "--sign", "-", "/tmp/ffmpeg"}) {
			t.Errorf("call[1] = %s %v, want codesign [--force --sign - /tmp/ffmpeg]", calls[1].name, calls[1].args)
		}
	})

	t.Run("strip failure short-circuits before codesign", func(t *testing.T) {
		t.Setenv("STRIP", "")
		t.Setenv("CODESIGN", "")
		stripErr := errors.New("exit status 1")
		var calls []recordedCall
		runner := func(_ context.Context, name string, args ...string) ([]byte, error) {
			calls = append(calls, recordedCall{name: name, args: slices.Clone(args)})
			return []byte("strip diagnostic"), stripErr
		}

		err := stripAndSignBinaryWithRunner(t.Context(), "/tmp/ffmpeg", "darwin", runner)
		if !errors.Is(err, stripErr) || !strings.Contains(err.Error(), "strip diagnostic") {
			t.Fatalf("error = %v, want wrapped stripErr with diagnostic", err)
		}
		if len(calls) != 1 {
			t.Errorf("call count = %d, want 1 (codesign must not run after strip failure)", len(calls))
		}
	})

	t.Run("codesign failure propagates error", func(t *testing.T) {
		t.Setenv("STRIP", "")
		t.Setenv("CODESIGN", "")
		signErr := errors.New("exit status 1")
		var calls []recordedCall
		runner := func(_ context.Context, name string, args ...string) ([]byte, error) {
			calls = append(calls, recordedCall{name: name, args: slices.Clone(args)})
			if name == "codesign" {
				return []byte("codesign diagnostic"), signErr
			}
			return nil, nil
		}

		err := stripAndSignBinaryWithRunner(t.Context(), "/tmp/ffmpeg", "darwin", runner)
		if !errors.Is(err, signErr) || !strings.Contains(err.Error(), "codesign diagnostic") {
			t.Fatalf("error = %v, want wrapped signErr with diagnostic", err)
		}
		if len(calls) != 2 {
			t.Errorf("call count = %d, want 2", len(calls))
		}
	})

	t.Run("STRIP and CODESIGN env overrides", func(t *testing.T) {
		t.Setenv("STRIP", "/custom/bin/my-strip")
		t.Setenv("CODESIGN", "/custom/bin/my-codesign")
		var calls []recordedCall
		runner := func(_ context.Context, name string, args ...string) ([]byte, error) {
			calls = append(calls, recordedCall{name: name, args: slices.Clone(args)})
			return nil, nil
		}

		if err := stripAndSignBinaryWithRunner(t.Context(), "/tmp/ffprobe", "darwin", runner); err != nil {
			t.Fatalf("stripAndSignBinaryWithRunner() error = %v", err)
		}
		if len(calls) != 2 {
			t.Fatalf("call count = %d, want 2 (%+v)", len(calls), calls)
		}
		if calls[0].name != "/custom/bin/my-strip" {
			t.Errorf("call[0].name = %q, want /custom/bin/my-strip", calls[0].name)
		}
		if calls[1].name != "/custom/bin/my-codesign" {
			t.Errorf("call[1].name = %q, want /custom/bin/my-codesign", calls[1].name)
		}
	})
}

func TestStageProgramsAndLicenseNotice(t *testing.T) {
	const (
		gplFixture  = "GNU GENERAL PUBLIC LICENSE\nVersion 3, 29 June 2007\n"
		lgplFixture = "GNU LESSER GENERAL PUBLIC LICENSE\nVersion 2.1, February 1999\n"
	)

	noopRunner := func(_ context.Context, _ string, _ ...string) ([]byte, error) {
		return nil, nil
	}

	targets := []struct {
		targetOS string
		arch     string
	}{
		{targetOS: "linux", arch: "amd64"},
		{targetOS: "darwin", arch: "arm64"},
		{targetOS: "windows", arch: "amd64"},
	}

	for _, embedded := range []bool{false, true} {
		profile := "default"
		wantLicense := gplFixture
		if embedded {
			profile = "embedded"
			wantLicense = lgplFixture
		}

		for _, target := range targets {
			for _, fromStagingFallback := range []bool{false, true} {
				sourceLabel := "src"
				if fromStagingFallback {
					sourceLabel = "staging-fallback"
				}
				testName := fmt.Sprintf("%s/%s_%s/%s", profile, target.targetOS, target.arch, sourceLabel)
				t.Run(testName, func(t *testing.T) {
					root := t.TempDir()
					buildRoot := filepath.Join(root, "buildroot")
					stagingDir := filepath.Join(buildRoot, "staging")
					binDir := filepath.Join(root, "out", target.targetOS+"_"+target.arch)

					stageBinDir := filepath.Join(stagingDir, "bin")
					if err := os.MkdirAll(stageBinDir, 0o755); err != nil {
						t.Fatalf("mkdir staging bin: %v", err)
					}
					for _, binName := range []string{"ffmpeg", "ffprobe"} {
						p := filepath.Join(stageBinDir, binaryFilename(binName, target.targetOS))
						if err := os.WriteFile(p, []byte("binary-"+binName), 0o755); err != nil { //nolint:gosec // G302: fixture executable in temp dir
							t.Fatalf("write fixture binary %s: %v", binName, err)
						}
					}

					licenseDir := filepath.Join(buildRoot, "src", "ffmpeg")
					if fromStagingFallback {
						licenseDir = filepath.Join(stagingDir, "share", "ffmpeg")
					}
					if err := os.MkdirAll(licenseDir, 0o755); err != nil {
						t.Fatalf("mkdir license dir: %v", err)
					}
					if err := os.WriteFile(filepath.Join(licenseDir, "COPYING.GPLv3"), []byte(gplFixture), 0o644); err != nil {
						t.Fatalf("write COPYING.GPLv3: %v", err)
					}
					if err := os.WriteFile(filepath.Join(licenseDir, "COPYING.LGPLv2.1"), []byte(lgplFixture), 0o644); err != nil {
						t.Fatalf("write COPYING.LGPLv2.1: %v", err)
					}

					libs := librariesForPlatformWithPrograms(embedded, true, stagingDir, target.targetOS, target.arch)
					err := stageProgramsWithRunner(t.Context(), buildRoot, stagingDir, binDir, embedded, target.targetOS, target.arch, "v1.2.3", libs, noopRunner)
					if err != nil {
						t.Fatalf("stageProgramsWithRunner() error = %v", err)
					}

					// Verify staged binaries exist with expected contents and Unix permissions.
					for _, binName := range []string{"ffmpeg", "ffprobe"} {
						dst := filepath.Join(binDir, binaryFilename(binName, target.targetOS))
						info, err := os.Stat(dst)
						if err != nil {
							t.Fatalf("stat staged %s: %v", dst, err)
						}
						if !info.Mode().IsRegular() {
							t.Errorf("staged %s is not a regular file", dst)
						}
						if runtime.GOOS != "windows" && info.Mode().Perm() != 0o755 {
							t.Errorf("staged %s mode = %o, want 755", dst, info.Mode().Perm())
						}
						gotBytes, err := os.ReadFile(dst)
						if err != nil {
							t.Fatalf("read staged %s: %v", dst, err)
						}
						if string(gotBytes) != "binary-"+binName {
							t.Errorf("staged %s content = %q, want %q", dst, string(gotBytes), "binary-"+binName)
						}
					}

					// Verify LICENSE matches the profile's FFmpeg copying file with zero MIT text.
					gotLicense, err := os.ReadFile(filepath.Join(binDir, "LICENSE"))
					if err != nil {
						t.Fatalf("read LICENSE: %v", err)
					}
					if string(gotLicense) != wantLicense {
						t.Errorf("LICENSE = %q, want %q", string(gotLicense), wantLicense)
					}
					if strings.Contains(string(gotLicense), "MIT") {
						t.Error("LICENSE must not contain MIT text")
					}

					// Verify NOTICE.txt contents.
					noticeBytes, err := os.ReadFile(filepath.Join(binDir, "NOTICE.txt"))
					if err != nil {
						t.Fatalf("read NOTICE.txt: %v", err)
					}
					notice := string(noticeBytes)

					ffmpegTag := ffmpegTagFromURL(ffmpeg.URL)
					for _, wantSubstr := range []string{
						"Profile: " + profile,
						fmt.Sprintf("Target: %s/%s", target.targetOS, target.arch),
						"FFmpeg: " + ffmpegTag,
						"Release: v1.2.3",
					} {
						if !strings.Contains(notice, wantSubstr) {
							t.Errorf("NOTICE.txt missing %q", wantSubstr)
						}
					}

					if embedded {
						archiveName := fmt.Sprintf("ffmpeg-embedded-%s-%s.tar.gz", target.targetOS, target.arch)
						for _, wantSubstr := range []string{
							"LGPLv2.1 Section 4 Notice",
							"LGPL-2.1-or-later",
							archiveName,
						} {
							if !strings.Contains(notice, wantSubstr) {
								t.Errorf("embedded NOTICE.txt missing %q", wantSubstr)
							}
						}
						if strings.Contains(notice, "GPLv3 Section 6 Notice") {
							t.Error("embedded NOTICE.txt must not include GPLv3 Section 6 Notice")
						}
					} else {
						for _, wantSubstr := range []string{
							"GPLv3 Section 6 Notice",
							"GPL-3.0-or-later",
						} {
							if !strings.Contains(notice, wantSubstr) {
								t.Errorf("default NOTICE.txt missing %q", wantSubstr)
							}
						}
						if strings.Contains(notice, "LGPLv2.1 Section 4 Notice") {
							t.Error("default NOTICE.txt must not include LGPLv2.1 Section 4 Notice")
						}
					}

					for _, lib := range libs {
						if !isLibraryActiveForTarget(lib, target.targetOS) {
							continue
						}
						digest, ok := expectedDigest(lib.URL)
						if !ok {
							t.Fatalf("missing pinned digest for active library %s (%s)", lib.Name, lib.URL)
						}
						if !strings.Contains(notice, lib.URL) {
							t.Errorf("NOTICE.txt missing URL for active library %s: %s", lib.Name, lib.URL)
						}
						if !strings.Contains(notice, "sha256:"+digest) {
							t.Errorf("NOTICE.txt missing SHA-256 digest for active library %s", lib.Name)
						}
					}
				})
			}
		}
	}

	t.Run("every active library in AllLibraries plus openh264 has a license entry", func(t *testing.T) {
		all := append(slices.Clone(AllLibraries), openh264)
		for _, lib := range all {
			if lib.Enabled != nil && !*lib.Enabled {
				continue
			}
			if _, ok := libraryLicenses[lib.Name]; !ok {
				t.Errorf("libraryLicenses missing entry for active library %q", lib.Name)
			}
		}
	})
}

func TestBuildStateCanSkipBinaries(t *testing.T) {
	const url = "https://example.com/ffmpeg-8.1.2.tar.gz"

	writeBinary := func(t *testing.T, installDir, name string, content []byte) {
		t.Helper()
		binDir := filepath.Join(installDir, "bin")
		if err := os.MkdirAll(binDir, 0o755); err != nil {
			t.Fatalf("mkdir bin dir: %v", err)
		}
		path := filepath.Join(binDir, binaryFilename(name, runtime.GOOS))
		if err := os.WriteFile(path, content, 0o755); err != nil { //nolint:gosec // G302: fixture executable in temp dir
			t.Fatalf("write binary %s: %v", path, err)
		}
	}

	tests := []struct {
		name  string
		setup func(t *testing.T, installDir string)
		want  bool
	}{
		{
			name: "missing declared binary returns false",
			setup: func(t *testing.T, installDir string) {
				writeArchive(t, installDir, "libavcodec")
				writeBinary(t, installDir, "ffmpeg", []byte("elf"))
				// ffprobe is missing
			},
			want: false,
		},
		{
			name: "empty declared binary returns false",
			setup: func(t *testing.T, installDir string) {
				writeArchive(t, installDir, "libavcodec")
				writeBinary(t, installDir, "ffmpeg", []byte("elf"))
				writeBinary(t, installDir, "ffprobe", nil)
			},
			want: false,
		},
		{
			name: "directory instead of binary returns false",
			setup: func(t *testing.T, installDir string) {
				writeArchive(t, installDir, "libavcodec")
				writeBinary(t, installDir, "ffmpeg", []byte("elf"))
				dirPath := filepath.Join(installDir, "bin", binaryFilename("ffprobe", runtime.GOOS))
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					t.Fatalf("mkdir fake binary: %v", err)
				}
			},
			want: false,
		},
		{
			name: "missing link lib with binaries present returns false",
			setup: func(t *testing.T, installDir string) {
				writeBinary(t, installDir, "ffmpeg", []byte("elf"))
				writeBinary(t, installDir, "ffprobe", []byte("elf"))
			},
			want: false,
		},
		{
			name: "both LinkLibs and Binaries present as non-empty regular files returns true",
			setup: func(t *testing.T, installDir string) {
				writeArchive(t, installDir, "libavcodec")
				writeBinary(t, installDir, "ffmpeg", []byte("elf"))
				writeBinary(t, installDir, "ffprobe", []byte("elf"))
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lib := &Library{
				Name:     "ffmpeg",
				URL:      url,
				LinkLibs: []string{"libavcodec"},
				Binaries: []string{"ffmpeg", "ffprobe"},
			}
			buildRoot := t.TempDir()
			if err := NewBuildState(lib, buildRoot).Save(); err != nil {
				t.Fatalf("save build state: %v", err)
			}
			installDir := t.TempDir()
			tt.setup(t, installDir)

			if got := NewBuildState(lib, buildRoot).CanSkip(installDir); got != tt.want {
				t.Errorf("CanSkip() = %v, want %v", got, tt.want)
			}
		})
	}
}
