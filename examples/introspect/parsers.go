package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"unsafe"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

func listParsers() {
	fmt.Println("\n==================================================")
	fmt.Println("PARSERS")
	fmt.Println("==================================================")
	fmt.Printf("    %-24s %-42s\n", "NAME", "SUPPORTED CODECS")
	fmt.Println()

	type parserInfo struct {
		name     string
		codecIDs []string
	}

	var parsers []parserInfo

	// Iterate through all registered parsers
	var parserOpaque unsafe.Pointer
	for {
		parser := ffmpeg.AVParserIterate(&parserOpaque)
		if parser == nil {
			break
		}

		// Get codec IDs
		codecIDs := []string{}
		codecIDArray := parser.CodecIds()
		for i := uintptr(0); ; i++ {
			codecID := codecIDArray.Get(i)
			if codecID == 0 {
				break
			}
			// Get codec name from ID
			codecName := getCodecName(ffmpeg.AVCodecID(codecID)) //nolint:gosec // G115: codec IDs are small enum values
			codecIDs = append(codecIDs, codecName)
		}

		if len(codecIDs) > 0 {
			parsers = append(parsers, parserInfo{
				name:     codecIDs[0], // Use first codec as parser name
				codecIDs: codecIDs,
			})
		}
	}

	// Sort parsers by name
	slices.SortFunc(parsers, func(a, b parserInfo) int {
		return cmp.Compare(a.name, b.name)
	})

	// Display parsers
	for _, p := range parsers {
		// Truncate long parser names to 24 chars
		parserName := p.name
		if len(parserName) > 24 {
			parserName = parserName[:24]
		}

		codecList := strings.Join(p.codecIDs, ", ")
		if len(codecList) > 42 {
			codecList = codecList[:42]
		}

		fmt.Printf("    %-24s %-42s\n", parserName, codecList)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total parsers: %d\n", len(parsers))
}
