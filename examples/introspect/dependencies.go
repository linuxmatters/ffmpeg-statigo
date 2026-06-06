package main

import (
	"fmt"
	"os"
	"strings"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

type codecDependencies struct {
	descriptions []string // Multiple codec descriptions when consolidating
	decoders     []string
	encoders     []string
	parsers      []string
	demuxers     []string
	muxers       []string
	bsfs         []string
}

// analyzeCodecDependencies finds all dependencies for a given codec and outputs configure flags
func analyzeCodecDependencies(codecName string, disable bool) {
	codecNameMap := buildCodecNameMap()
	parserMap := buildParserMap()
	formatMap := buildFormatMap()
	bsfMap := buildBSFMap()

	matches := findMatchingCodecs(codecName, codecNameMap, parserMap, formatMap, bsfMap)

	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No codec found matching '%s'\n", codecName)
		os.Exit(1)
	}

	sortedMatches := sortCodecsByPriority(matches, codecName)

	prefix := "--enable"
	if disable {
		prefix = "--disable"
	}

	deps := consolidateCodecDependencies(sortedMatches, parserMap, formatMap, bsfMap)
	outputConsolidatedDependencies(deps, prefix)
}

// findMatchingCodecs uses improved matching logic and reverse lookups
func findMatchingCodecs(search string, codecNameMap map[string]*codecInfo,
	parserMap map[ffmpeg.AVCodecID][]string, formatMap map[ffmpeg.AVCodecID]*formatDeps,
	bsfMap map[ffmpeg.AVCodecID][]string,
) []*codecInfo {
	searchLower := strings.ToLower(search)
	matchedCodecs := make(map[string]*codecInfo)

	// 1. Direct codec name matching with improved logic
	for _, info := range codecNameMap {
		nameLower := strings.ToLower(info.name)

		if nameLower == searchLower {
			matchedCodecs[info.name] = info
			continue
		}

		// Variant pattern: <codec>_<suffix> (e.g., av1_qsv, h264_nvenc)
		if strings.HasPrefix(nameLower, searchLower+"_") {
			matchedCodecs[info.name] = info
			continue
		}

		// Library pattern: lib*<codec>* (e.g., libdav1d, librav1e, libx264)
		if strings.HasPrefix(nameLower, "lib") && strings.Contains(nameLower, searchLower) {
			matchedCodecs[info.name] = info
			continue
		}

		// Starts with search term (e.g., searching "h26" finds "h264", "h265")
		if strings.HasPrefix(nameLower, searchLower) {
			matchedCodecs[info.name] = info
			continue
		}
	}

	// 2. Reverse lookup from parsers
	for codecID, parserNames := range parserMap {
		for _, parserName := range parserNames {
			if strings.Contains(strings.ToLower(parserName), searchLower) {
				for _, info := range codecNameMap {
					if info.codecID == codecID {
						matchedCodecs[info.name] = info
					}
				}
			}
		}
	}

	// 3. Reverse lookup from formats (demuxers/muxers)
	for codecID, formats := range formatMap {
		// Check if any muxer/demuxer name contains the search term
		// (e.g., "avif" for av1, "matroska" for various codecs)
		if !formatNameMatches(formats, searchLower) {
			continue
		}
		for _, info := range codecNameMap {
			if info.codecID == codecID {
				matchedCodecs[info.name] = info
			}
		}
	}

	// 4. Reverse lookup from BSFs
	for codecID, bsfNames := range bsfMap {
		for _, bsfName := range bsfNames {
			if strings.Contains(strings.ToLower(bsfName), searchLower) {
				for _, info := range codecNameMap {
					if info.codecID == codecID {
						matchedCodecs[info.name] = info
					}
				}
			}
		}
	}

	var result []*codecInfo
	for _, info := range matchedCodecs {
		result = append(result, info)
	}

	return result
}

// sortCodecsByPriority orders codecs: exact match, software, hardware
func sortCodecsByPriority(codecs []*codecInfo, searchTerm string) []*codecInfo {
	var exact []*codecInfo
	var software []*codecInfo
	var hardware []*codecInfo

	searchLower := strings.ToLower(searchTerm)

	for _, codec := range codecs {
		nameLower := strings.ToLower(codec.name)

		switch {
		case nameLower == searchLower:
			exact = append(exact, codec)
		case isHardwareCodec(codec.name):
			hardware = append(hardware, codec)
		default:
			software = append(software, codec)
		}
	}

	result := make([]*codecInfo, 0, len(codecs))
	result = append(result, exact...)
	result = append(result, software...)
	result = append(result, hardware...)

	return result
}

// formatNameMatches reports whether any muxer or demuxer name contains searchLower.
func formatNameMatches(formats *formatDeps, searchLower string) bool {
	for _, name := range formats.muxers {
		if strings.Contains(strings.ToLower(name), searchLower) {
			return true
		}
	}
	for _, name := range formats.demuxers {
		if strings.Contains(strings.ToLower(name), searchLower) {
			return true
		}
	}
	return false
}

// isHardwareCodec checks if codec uses hardware acceleration
func isHardwareCodec(name string) bool {
	nameLower := strings.ToLower(name)
	hwSuffixes := []string{"_qsv", "_nvenc", "_nvdec", "_vulkan", "_vaapi", "_videotoolbox", "_amf", "_v4l2m2m"}
	for _, suffix := range hwSuffixes {
		if strings.HasSuffix(nameLower, suffix) {
			return true
		}
	}
	return false
}

// consolidateCodecDependencies merges all codec dependencies into one set
func consolidateCodecDependencies(codecs []*codecInfo,
	parserMap map[ffmpeg.AVCodecID][]string, formatMap map[ffmpeg.AVCodecID]*formatDeps,
	bsfMap map[ffmpeg.AVCodecID][]string,
) *codecDependencies {
	deps := &codecDependencies{}
	deps.descriptions = make([]string, 0)

	descriptionSet := make(map[string]bool)
	decoderSet := make(map[string]bool)
	encoderSet := make(map[string]bool)
	parserSet := make(map[string]bool)
	demuxerSet := make(map[string]bool)
	muxerSet := make(map[string]bool)
	bsfSet := make(map[string]bool)

	for _, codec := range codecs {
		// Add description (deduplicate)
		if !descriptionSet[codec.longName] {
			deps.descriptions = append(deps.descriptions, codec.longName)
			descriptionSet[codec.longName] = true
		}

		codecID := codec.codecID

		if codec.isEncoder {
			encoderSet[codec.name] = true
		}
		if codec.isDecoder {
			decoderSet[codec.name] = true
		}

		if parsers, ok := parserMap[codecID]; ok {
			for _, p := range parsers {
				parserSet[p] = true
			}
		}

		if formats, ok := formatMap[codecID]; ok {
			for _, d := range formats.demuxers {
				demuxerSet[d] = true
			}
			for _, m := range formats.muxers {
				muxerSet[m] = true
			}
		}

		if bsfs, ok := bsfMap[codecID]; ok {
			for _, b := range bsfs {
				bsfSet[b] = true
			}
		}
	}

	// Convert sets to sorted slices (alphanumerically)
	deps.decoders = sortedKeys(decoderSet)
	deps.encoders = sortedKeys(encoderSet)
	deps.parsers = sortedKeys(parserSet)
	deps.demuxers = sortedKeys(demuxerSet)
	deps.muxers = sortedKeys(muxerSet)
	deps.bsfs = sortedKeys(bsfSet)

	return deps
}

// outputConsolidatedDependencies prints consolidated codec dependencies
func outputConsolidatedDependencies(deps *codecDependencies, prefix string) {
	// Print all codec descriptions as comments
	for _, desc := range deps.descriptions {
		fmt.Printf("// %s\n", desc)
	}

	// Print each component type with comma-delimited list (only if not empty)
	if len(deps.encoders) > 0 {
		fmt.Printf("%s-encoder=%s\n", prefix, strings.Join(deps.encoders, ","))
	}
	if len(deps.decoders) > 0 {
		fmt.Printf("%s-decoder=%s\n", prefix, strings.Join(deps.decoders, ","))
	}
	if len(deps.parsers) > 0 {
		fmt.Printf("%s-parser=%s\n", prefix, strings.Join(deps.parsers, ","))
	}
	if len(deps.demuxers) > 0 {
		fmt.Printf("%s-demuxer=%s\n", prefix, strings.Join(deps.demuxers, ","))
	}
	if len(deps.muxers) > 0 {
		fmt.Printf("%s-muxer=%s\n", prefix, strings.Join(deps.muxers, ","))
	}
	if len(deps.bsfs) > 0 {
		fmt.Printf("%s-bsf=%s\n", prefix, strings.Join(deps.bsfs, ","))
	}
}

// Helper functions for building lookup maps
