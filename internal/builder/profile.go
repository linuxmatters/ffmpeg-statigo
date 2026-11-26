package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
)

func profilePaths(embedded bool, targetOS, arch string) (buildRoot, output string, err error) {
	buildRoot, output = ".build", "lib"
	if embedded {
		buildRoot = filepath.Join(buildRoot, "embedded")
		output = filepath.Join(output, "embedded")
	}
	if targetOS == "windows" && arch == "386" {
		buildRoot = filepath.Join(buildRoot, "windows_386_msvcrt")
	}
	buildRoot, err = filepath.Abs(buildRoot)
	if err != nil {
		return "", "", fmt.Errorf("get absolute path for build root: %w", err)
	}
	output, err = filepath.Abs(filepath.Join(output, targetOS+"_"+arch, "libffmpeg.a"))
	if err != nil {
		return "", "", fmt.Errorf("get absolute path for output: %w", err)
	}
	return buildRoot, output, nil
}

func programsOutputDir(embedded bool, targetOS, arch string) (string, error) {
	output := "bin"
	if embedded {
		output = filepath.Join(output, "embedded")
	}
	output, err := filepath.Abs(filepath.Join(output, targetOS+"_"+arch))
	if err != nil {
		return "", fmt.Errorf("get absolute path for programs output: %w", err)
	}
	return output, nil
}

func librariesForProfile(embedded bool, stagingDir string) []*Library {
	return librariesForProfileWithPrograms(embedded, false, stagingDir)
}

func librariesForProfileWithPrograms(embedded, programs bool, stagingDir string) []*Library {
	return librariesForPlatformWithPrograms(embedded, programs, stagingDir, runtime.GOOS, runtime.GOARCH)
}

func librariesForPlatform(embedded bool, stagingDir, targetOS, arch string) []*Library {
	return librariesForPlatformWithPrograms(embedded, false, stagingDir, targetOS, arch)
}

func librariesForPlatformWithPrograms(embedded, programs bool, stagingDir, targetOS, arch string) []*Library {
	windowsNative := targetOS == "windows" && (arch == "amd64" || arch == "386")
	configuredLibraries := make(map[*Library]*Library)
	var libs []*Library
	for _, lib := range AllLibraries {
		if lib == ffmpeg {
			continue
		}
		if embedded {
			switch lib.Name {
			case "dav1d", "glslang", "libdrm", "libsrt", "libva", "libvpl", "openssl", "rav1e", "x264", "x265", "nv-codec-headers", "Vulkan-Headers":
				continue
			}
		}
		if windowsNative {
			original := lib
			lib = windowsLibrary(lib, arch)
			if len(lib.Dependencies) != 0 {
				configured := *lib
				configured.Dependencies = slices.Clone(lib.Dependencies)
				for i, dep := range configured.Dependencies {
					if replacement, ok := configuredLibraries[dep]; ok {
						configured.Dependencies[i] = replacement
					}
				}
				lib = &configured
			}
			configuredLibraries[original] = lib
		}
		libs = append(libs, lib)
	}
	if embedded {
		lib := openh264
		if windowsNative {
			lib = windowsOpenH264(arch)
		}
		libs = append(libs, lib)
	}

	configuredFFmpeg := *ffmpeg
	configuredFFmpeg.Dependencies = slices.Clone(libs)
	configuredFFmpeg.ConfigureArgs = func(targetOS string) []string {
		args := ffmpegConfigureArgsWithPrograms(targetOS, arch, stagingDir, embedded, programs)
		for _, lib := range configuredFFmpeg.Dependencies {
			if lib.Enabled != nil && !*lib.Enabled {
				continue
			}
			if len(lib.Platform) != 0 && !slices.Contains(lib.Platform, targetOS) {
				continue
			}
			for _, flag := range lib.FFmpegEnables {
				args = append(args, "--enable-"+flag)
			}
		}
		if targetOS == "darwin" && !embedded {
			args = append(args, "--enable-avfoundation", "--enable-audiotoolbox", "--enable-videotoolbox")
		}
		return args
	}
	if programs {
		configuredFFmpeg.Binaries = []string{"ffmpeg", "ffprobe"}
		basePostExtract := ffmpeg.PostExtract
		configuredFFmpeg.PostExtract = func(ctx context.Context, srcPath string) error {
			if basePostExtract != nil {
				if err := basePostExtract(ctx, srcPath); err != nil {
					return err
				}
			}
			if err := stageFFmpegLicenses(srcPath, stagingDir); err != nil {
				return err
			}
			if targetOS == "linux" {
				if err := scrubLinuxPkgConfig(stagingDir); err != nil {
					return err
				}
			}
			return nil
		}
	}
	return append(libs, &configuredFFmpeg)
}

func stageFFmpegLicenses(srcPath, stagingDir string) error {
	destDir := filepath.Join(stagingDir, "share", "ffmpeg")
	for _, name := range []string{"COPYING.GPLv3", "COPYING.LGPLv2.1"} {
		src := filepath.Join(srcPath, name)
		data, err := os.ReadFile(src) //nolint:gosec // G304: fixed license filename under the local extracted FFmpeg source directory
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return fmt.Errorf("read %s: %w", src, err)
		}
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			return fmt.Errorf("create %s: %w", destDir, err)
		}
		dest := filepath.Join(destDir, name)
		if err := os.WriteFile(dest, data, 0o644); err != nil { //nolint:gosec // G703: fixed license filename under the local staging directory
			return fmt.Errorf("write %s: %w", dest, err)
		}
	}
	return nil
}

func ffmpegConfigureArgs(targetOS, stagingDir string, embedded bool) []string {
	return ffmpegConfigureArgsWithPrograms(targetOS, runtime.GOARCH, stagingDir, embedded, false)
}

func ffmpegConfigureArgsForPlatform(targetOS, arch, stagingDir string, embedded bool) []string {
	return ffmpegConfigureArgsWithPrograms(targetOS, arch, stagingDir, embedded, false)
}

func ffmpegConfigureArgsWithPrograms(targetOS, arch, stagingDir string, embedded, programs bool) []string {
	incDir, libDir := filepath.Join(stagingDir, "include"), filepath.Join(stagingDir, "lib")
	windowsNative := targetOS == "windows" && (arch == "amd64" || arch == "386")
	if windowsNative {
		incDir, libDir = buildToolPath(incDir, targetOS), buildToolPath(libDir, targetOS)
	}
	args := []string{
		"--pkg-config-flags=--static",
		"--extra-cflags=-I" + incDir,
		"--extra-ldflags=-L" + libDir,
	}
	if targetOS == "darwin" {
		args = append(args, "--cc=clang", "--cxx=clang++")
	}
	if programs && targetOS == "linux" {
		args = append(args, "--extra-ldexeflags=-static-libstdc++ -static-libgcc")
	}
	if windowsNative {
		ffmpegArch := "x86_64"
		if arch == "386" {
			ffmpegArch = "x86"
			// GCC 16 on Win32 can enter with a 4-byte stack; swscale needs 16-byte alignment.
			args = append(args, "--extra-cflags=-mpreferred-stack-boundary=4", "--extra-cflags=-mstackrealign")
		}
		if programs {
			args = append(args, "--extra-ldexeflags=-static -static-libgcc -static-libstdc++")
		}
		args = append(args, "--target-os=mingw32", "--arch="+ffmpegArch)
	}
	return append(args, ffmpegArgsCommonWithPrograms(targetOS, embedded, programs)...)
}
