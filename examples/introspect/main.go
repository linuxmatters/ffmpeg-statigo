package main

import (
	"fmt"
	"os"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

func main() {
	// Check for --enable or --disable flags
	if len(os.Args) >= 3 {
		switch os.Args[1] {
		case "--enable":
			analyzeCodecDependencies(os.Args[2], false)
			return
		case "--disable":
			analyzeCodecDependencies(os.Args[2], true)
			return
		}
	}

	// Normal introspection output
	fmt.Println("ffmpeg-statigo")
	fmt.Println("==============")
	fmt.Println()

	// Get configuration
	fmt.Println("Configuration:")
	fmt.Printf("%s\n\n", ffmpeg.AVFormatConfiguration().String())

	// Get license
	fmt.Println("License:")
	fmt.Printf("%s\n\n", ffmpeg.AVFormatLicense().String())

	// List all components
	listCodecs()
	listHWAccels()
	listFormats()
	listFilters()
	listBSFs()
	listParsers()
	listProtocols()
}
