package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strings"
)

func combineMac(output string) {
	{
		var matches []string

		for _, lib := range libs {
			libPath := path.Join(stagingDir, "lib", fmt.Sprintf("%v.a", lib))
			// Check lib64 if not found in lib (aom and brotli may install there)
			if !exists(libPath) {
				libPath = path.Join(stagingDir, "lib64", fmt.Sprintf("%v.a", lib))
			}
			matches = append(matches, libPath)
		}

		log.Println("Running merge")
		cmd := exec.Command("libtool")
		cmd.Args = append(
			cmd.Args,
			"-static",
			"-o",
			output,
		)
		cmd.Args = append(
			cmd.Args,
			matches...,
		)
		run("[merge]", cmd)
	}

	{

		log.Println("Running strip")
		cmd := exec.Command("strip")
		cmd.Args = append(
			cmd.Args,
			"-S",
			output,
		)
		run("[strip]", cmd)
	}
}

func combineLinux(output string) {
	var matches []string

	for _, lib := range libs {
		libPath := path.Join(stagingDir, "lib", fmt.Sprintf("%v.a", lib))
		// Check lib64 if not found in lib (aom and brotli install there)
		if !exists(libPath) {
			libPath = path.Join(stagingDir, "lib64", fmt.Sprintf("%v.a", lib))
		}
		matches = append(matches, libPath)
	}

	var lines []string

	lines = append(lines, fmt.Sprintf("create %v", output))

	for _, match := range matches {
		log.Println(match)
		lines = append(lines, fmt.Sprintf("addlib %v", match))
	}

	lines = append(lines, "save", "end")
	mriAll := strings.Join(lines, "\n")

	log.Println(mriAll)

	{
		log.Println("Running merge")
		cmd := exec.Command("ar")
		cmd.Args = append(
			cmd.Args,
			"-M",
		)

		cmd.Stdin = bytes.NewBufferString(mriAll)

		run("[merge]", cmd)
	}

	{
		log.Println("Running strip")
		cmd := exec.Command("strip")
		cmd.Args = append(
			cmd.Args,
			"--strip-unneeded",
			output,
		)
		run("[strip]", cmd)
	}
}
