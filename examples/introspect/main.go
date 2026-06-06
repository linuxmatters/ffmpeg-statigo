package main

import (
	"fmt"
	"os"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

func main() {
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

	fmt.Println("ffmpeg-statigo")
	fmt.Println("==============")
	fmt.Println()

	fmt.Println("Configuration:")
	fmt.Printf("%s\n\n", ffmpeg.AVFormatConfiguration().String())

	fmt.Println("License:")
	fmt.Printf("%s\n\n", ffmpeg.AVFormatLicense().String())

	listCodecs()
	listHWAccels()
	listFormats()
	listFilters()
	listBSFs()
	listParsers()
	listProtocols()
}
