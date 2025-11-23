// download-lib is a utility to download the FFmpeg static libraries
// Run this program to trigger the automatic download and SHA256 verification
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/linuxmatters/ffmpeg-statigo/lib"
)

func main() {
	fmt.Printf("Downloading FFmpeg libraries for %s/%s...\n", runtime.GOOS, runtime.GOARCH)

	if err := lib.DownloadLibs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Download complete!")
}
