package main

import (
	"context"
	"os"
	"runtime"
)

var openh264 = &Library{
	Name:          "openh264",
	URL:           "https://github.com/cisco/openh264/archive/refs/tags/v2.6.0.tar.gz",
	FFmpegEnables: []string{"libopenh264"},
	BuildSystem:   openh264BuildSystem(openh264MakeArgs, runtime.GOOS),
	ConfigureArgs: openh264MakeArgs,
	LinkLibs:      []string{"libopenh264"},
}

func openh264BuildSystem(makeArgs func(string) []string, targetOS string) *MakefileBuild {
	return &MakefileBuild{
		Targets: []string{"libopenh264.a"},
		InstallFunc: func(ctx context.Context, srcPath, installDir string) error {
			args := append(makeArgs(targetOS), "PREFIX="+buildToolPath(installDir, targetOS), "install-static")
			return runCommand(ctx, srcPath, os.Stdout, installDir, "make", args...)
		},
	}
}

func windowsOpenH264() *Library {
	configured := *openh264
	configured.ConfigureArgs = func(targetOS string) []string {
		args := openh264MakeArgs(targetOS)
		if targetOS == "windows" {
			args = append(args, "OS=mingw_nt", "ARCH=x86_64")
		}
		return args
	}
	configured.BuildSystem = openh264BuildSystem(configured.ConfigureArgs, "windows")
	return &configured
}

func openh264MakeArgs(targetOS string) []string {
	args := []string{"ENABLEPIC=Yes", "HAVE_GMP_API=No", "HAVE_GTEST=No"}
	if targetOS == "darwin" {
		args = append(args, "STATIC_LDFLAGS=-lc++")
	}
	return args
}
