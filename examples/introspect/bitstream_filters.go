package main

import (
	"fmt"
	"strings"
	"unsafe"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

func listBSFs() {
	fmt.Println("\n==================================================")
	fmt.Println("BITSTREAM FILTERS")
	fmt.Println("==================================================")
	fmt.Printf("    %-24s %-42s\n", "NAME", "SUPPORTED CODECS")
	fmt.Println()

	var opaque unsafe.Pointer
	count := 0
	codecSpecificCount := 0

	for {
		bsf := ffmpeg.AVBSFIterate(&opaque)
		if bsf == nil {
			break
		}
		count++

		// Get the bitstream filter name
		nameStr := bsf.Name().String()

		// Truncate name if too long
		if len(nameStr) > 24 {
			nameStr = nameStr[:24]
		}

		// Get supported codec IDs
		codecList := "all"
		codecIDs := bsf.CodecIds()
		if codecIDs != nil {
			codecSpecificCount++
			var codecs []string
			for i := uintptr(0); ; i++ {
				codecID := codecIDs.Get(i)
				if codecID == ffmpeg.AVCodecIdNone {
					break
				}
				codecName := getCodecName(codecID)
				codecs = append(codecs, codecName)
			}
			if len(codecs) > 0 {
				codecList = strings.Join(codecs, ", ")
				if len(codecList) > 64 {
					codecList = codecList[:61] + "..."
				}
			}
		}

		fmt.Printf("    %-24s %-42s\n", nameStr, codecList)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total bitstream filters: %d\n", count)
	fmt.Printf("  Codec-specific filters: %d\n", codecSpecificCount)
	fmt.Printf("  Generic filters: %d\n", count-codecSpecificCount)
}
