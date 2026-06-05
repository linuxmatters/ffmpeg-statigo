package main

import (
	"cmp"
	"fmt"
	"slices"
	"unsafe"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

func listProtocols() {
	fmt.Println("\n==================================================")
	fmt.Println("PROTOCOLS")
	fmt.Println("==================================================")
	fmt.Printf("%s  %-24s\n", "IO", "NAME")
	fmt.Println()

	type protocolInfo struct {
		name     string
		isInput  bool
		isOutput bool
	}

	protocolMap := make(map[string]*protocolInfo)

	// Iterate through input protocols
	var inputOpaque unsafe.Pointer
	for {
		name := ffmpeg.AVIOEnumProtocols(&inputOpaque, 0)
		if name == "" {
			break
		}

		if existing, exists := protocolMap[name]; exists {
			existing.isInput = true
		} else {
			protocolMap[name] = &protocolInfo{
				name:     name,
				isInput:  true,
				isOutput: false,
			}
		}
	}

	// Iterate through output protocols
	var outputOpaque unsafe.Pointer
	for {
		name := ffmpeg.AVIOEnumProtocols(&outputOpaque, 1)
		if name == "" {
			break
		}

		if existing, exists := protocolMap[name]; exists {
			existing.isOutput = true
		} else {
			protocolMap[name] = &protocolInfo{
				name:     name,
				isInput:  false,
				isOutput: true,
			}
		}
	}

	// Convert map to slice and sort
	var protocols []protocolInfo
	for _, p := range protocolMap {
		protocols = append(protocols, *p)
	}

	slices.SortFunc(protocols, func(a, b protocolInfo) int {
		return cmp.Compare(a.name, b.name)
	})

	// Count totals
	totalInput := 0
	totalOutput := 0

	// Display protocols
	for _, p := range protocols {
		inputFlag := "."
		if p.isInput {
			inputFlag = "I"
			totalInput++
		}

		outputFlag := "."
		if p.isOutput {
			outputFlag = "O"
			totalOutput++
		}

		// Truncate long protocol names to 24 chars
		protocolName := p.name
		if len(protocolName) > 24 {
			protocolName = protocolName[:24]
		}

		fmt.Printf("%s%s  %-24s\n", inputFlag, outputFlag, protocolName)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total protocols: %d\n", len(protocols))
	fmt.Printf("  Total input protocols: %d\n", totalInput)
	fmt.Printf("  Total output protocols: %d\n", totalOutput)
}

// testHardwareAvailable tests if a hardware device type is actually available
// by attempting to create a device context for it.
