package main

import (
	"fmt"
	"path/filepath"
	"slices"
)

func profilePaths(embedded bool, targetOS, arch string) (buildRoot, output string, err error) {
	buildRoot, output = ".build", "lib"
	if embedded {
		buildRoot = filepath.Join(buildRoot, "embedded")
		output = filepath.Join(output, "embedded")
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
		libs = append(libs, lib)
	}
	if embedded {
		libs = append(libs, openh264)
	}

	configuredFFmpeg := *ffmpeg
	configuredFFmpeg.Dependencies = slices.Clone(libs)
	configuredFFmpeg.ConfigureArgs = func(targetOS string) []string {
		args := ffmpegConfigureArgs(targetOS, stagingDir, embedded)
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
	args := []string{
		"--pkg-config-flags=--static",
		"--extra-cflags=-I" + filepath.Join(stagingDir, "include"),
		"--extra-ldflags=-L" + filepath.Join(stagingDir, "lib"),
	}
	if targetOS == "darwin" {
		args = append(args, "--cc=clang", "--cxx=clang++")
	}
	return append(args, ffmpegArgsCommon(targetOS, embedded)...)
}
