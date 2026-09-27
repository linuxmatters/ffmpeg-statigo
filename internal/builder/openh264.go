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
	BuildSystem: &MakefileBuild{
		Targets: []string{"libopenh264.a"},
		InstallFunc: func(ctx context.Context, srcPath, installDir string) error {
			args := append(openh264MakeArgs(runtime.GOOS), "PREFIX="+installDir, "install-static")
			return runCommand(ctx, srcPath, os.Stdout, installDir, "make", args...)
		},
	},
	ConfigureArgs: openh264MakeArgs,
	LinkLibs:      []string{"libopenh264"},
}

func openh264MakeArgs(targetOS string) []string {
	args := []string{"ENABLEPIC=Yes", "HAVE_GMP_API=No", "HAVE_GTEST=No"}
	if targetOS == "darwin" {
		args = append(args, "STATIC_LDFLAGS=-lc++")
	}
	return args
}
