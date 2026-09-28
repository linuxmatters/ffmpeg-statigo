package main

import (
	"fmt"
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

func librariesForProfile(embedded bool, stagingDir string) []*Library {
	return librariesForPlatform(embedded, stagingDir, runtime.GOOS, runtime.GOARCH)
}

func librariesForPlatform(embedded bool, stagingDir, targetOS, arch string) []*Library {
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
		args := ffmpegConfigureArgsForPlatform(targetOS, arch, stagingDir, embedded)
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
	return append(libs, &configuredFFmpeg)
}

func ffmpegConfigureArgs(targetOS, stagingDir string, embedded bool) []string {
	return ffmpegConfigureArgsForPlatform(targetOS, runtime.GOARCH, stagingDir, embedded)
}

func ffmpegConfigureArgsForPlatform(targetOS, arch, stagingDir string, embedded bool) []string {
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
	if windowsNative {
		ffmpegArch := "x86_64"
		if arch == "386" {
			ffmpegArch = "x86"
			// GCC 16 on Win32 can enter with a 4-byte stack; swscale needs 16-byte alignment.
			args = append(args, "--extra-cflags=-mpreferred-stack-boundary=4", "--extra-cflags=-mstackrealign")
		}
		args = append(args, "--target-os=mingw32", "--arch="+ffmpegArch)
	}
	return append(args, ffmpegArgsCommon(targetOS, embedded)...)
}
