package main

import (
	"cmp"
	"fmt"
	"slices"
	"unsafe"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

func listFilters() {
	fmt.Println("\n==================================================")
	fmt.Println("FILTERS")
	fmt.Println("==================================================")
	fmt.Printf(" %s  %-24s %-42s %s\n", "TSHM", "NAME", "DESCRIPTION", "TYPE")
	fmt.Println()

	type filterInfo struct {
		name                string
		description         string
		mediaType           string
		isHardware          bool
		isMetadataOnly      bool
		supportTimeline     bool
		supportSliceThreads bool
	}

	var filters []filterInfo

	// Iterate through all filters
	var opaque unsafe.Pointer
	for {
		filter := ffmpeg.AVFilterIterate(&opaque)
		if filter == nil {
			break
		}

		name := ""
		if filter.Name() != nil {
			name = filter.Name().String()
		}

		if name == "" {
			continue
		}

		description := ""
		if filter.Description() != nil {
			description = filter.Description().String()
		}

		// Check filter flags
		flags := filter.Flags()
		supportTimeline := (flags & ffmpeg.AVFilterFlagSupportTimelineGeneric) != 0
		supportSliceThreads := (flags & ffmpeg.AVFilterFlagSliceThreads) != 0
		isHardware := (flags & ffmpeg.AVFilterFlagHWDevice) != 0
		isMetadataOnly := (flags & ffmpeg.AVFilterFlagMetadataOnly) != 0

		// Determine filter media type by checking input pads
		mediaType := "UNKNOWN"
		inputs := filter.Inputs()
		if inputs != nil {
			// Get the type of the first input pad
			padType := ffmpeg.AVFilterPadGetType(inputs, 0)
			mediaType = getMediaTypeString(padType)
		} else {
			// No inputs means it's likely a source filter
			// Check outputs instead
			outputs := filter.Outputs()
			if outputs != nil {
				padType := ffmpeg.AVFilterPadGetType(outputs, 0)
				mediaType = getMediaTypeString(padType)
			}
		}

		filters = append(filters, filterInfo{
			name:                name,
			description:         description,
			mediaType:           mediaType,
			isHardware:          isHardware,
			isMetadataOnly:      isMetadataOnly,
			supportTimeline:     supportTimeline,
			supportSliceThreads: supportSliceThreads,
		})
	}

	// Sort filters by name
	slices.SortFunc(filters, func(a, b filterInfo) int {
		return cmp.Compare(a.name, b.name)
	})

	// Count filters by type
	timelineFilters := 0
	sliceThreadFilters := 0
	hwFilters := 0
	metadataFilters := 0
	videoFilters := 0
	audioFilters := 0
	subtitleFilters := 0
	dataFilters := 0

	// Display filters
	for _, f := range filters {
		timelineFlag := "."
		if f.supportTimeline {
			timelineFlag = "T"
			timelineFilters++
		}

		sliceFlag := "."
		if f.supportSliceThreads {
			sliceFlag = "S"
			sliceThreadFilters++
		}

		hwFlag := "."
		if f.isHardware {
			hwFlag = "H"
			hwFilters++
		}

		metadataFlag := "."
		if f.isMetadataOnly {
			metadataFlag = "M"
			metadataFilters++
		}

		// Count by media type
		switch f.mediaType {
		case "VIDEO":
			videoFilters++
		case "AUDIO":
			audioFilters++
		case "SUBTITLE":
			subtitleFilters++
		case "DATA":
			dataFilters++
		}

		filterName := f.name
		if len(filterName) > 24 {
			filterName = filterName[:24]
		}

		description := f.description
		if len(description) > 42 {
			description = description[:42]
		}

		fmt.Printf(" %s%s%s%s %-24s %-42s %s\n", timelineFlag, sliceFlag, hwFlag, metadataFlag, filterName, description, f.mediaType)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total filters: %d\n", len(filters))
	fmt.Printf("  Timeline support: %d\n", timelineFilters)
	fmt.Printf("  Slice threading: %d\n", sliceThreadFilters)
	fmt.Printf("  Hardware filters: %d\n", hwFilters)
	fmt.Printf("  Metadata-only filters: %d\n", metadataFilters)
	fmt.Printf("\nBy media type:\n")
	fmt.Printf("  Video filters: %d\n", videoFilters)
	fmt.Printf("  Audio filters: %d\n", audioFilters)
	fmt.Printf("  Subtitle filters: %d\n", subtitleFilters)
	fmt.Printf("  Data filters: %d\n", dataFilters)
	fmt.Println("\nFlags:")
	fmt.Println("  T - Timeline support")
	fmt.Println("  S - Slice threading")
	fmt.Println("  H - Hardware device required")
	fmt.Println("  M - Metadata only (does not modify frame data)")
}

// codecDependencies holds all the components needed for a codec
