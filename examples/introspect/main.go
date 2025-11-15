package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unsafe"

	ffmpeg "github.com/csnewman/ffmpeg-go"
)

type codecInfo struct {
	name      string
	longName  string
	mediaType string
	codecID   ffmpeg.AVCodecID
	isEncoder bool
	isDecoder bool
}

func main() {
	// Check for --enable or --disable flags
	if len(os.Args) >= 3 {
		if os.Args[1] == "--enable" {
			analyzeCodecDependencies(os.Args[2], false)
			return
		} else if os.Args[1] == "--disable" {
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

func listCodecs() {
	fmt.Println("==================================================")
	fmt.Println("CODECS")
	fmt.Println("==================================================\n")
	fmt.Printf(" %s  %-24s %-42s %s\n", "DE", "NAME", "DESCRIPTION", "TYPE")
	fmt.Println()

	// Collect all codec information
	codecMap := make(map[string]*codecInfo)

	// Iterate through all codec descriptors
	var desc *ffmpeg.AVCodecDescriptor
	for {
		desc = ffmpeg.AVCodecDescriptorNext(desc)
		if desc == nil {
			break
		}

		codecID := desc.Id()
		name := desc.Name().String()
		longName := desc.LongName().String()
		mediaType := getMediaTypeString(desc.Type())

		// Check if encoder exists for this codec
		encoder := ffmpeg.AVCodecFindEncoder(codecID)
		// Check if decoder exists for this codec
		decoder := ffmpeg.AVCodecFindDecoder(codecID)

		if encoder != nil || decoder != nil {
			codecMap[name] = &codecInfo{
				name:      name,
				longName:  longName,
				mediaType: mediaType,
				isEncoder: encoder != nil,
				isDecoder: decoder != nil,
			}
		}
	}

	// Sort codecs by name for consistent output
	names := make([]string, 0, len(codecMap))
	for name := range codecMap {
		names = append(names, name)
	}
	sort.Strings(names)

	// Count encoders and decoders by type
	videoEncoders, videoDecoders := 0, 0
	audioEncoders, audioDecoders := 0, 0
	subtitleEncoders, subtitleDecoders := 0, 0
	otherEncoders, otherDecoders := 0, 0

	// Print all codecs
	for _, name := range names {
		info := codecMap[name]

		flags := ""
		if info.isDecoder {
			flags += "D"
			switch info.mediaType {
			case "VIDEO":
				videoDecoders++
			case "AUDIO":
				audioDecoders++
			case "SUBTITLE":
				subtitleDecoders++
			default:
				otherDecoders++
			}
		} else {
			flags += "."
		}
		if info.isEncoder {
			flags += "E"
			switch info.mediaType {
			case "VIDEO":
				videoEncoders++
			case "AUDIO":
				audioEncoders++
			case "SUBTITLE":
				subtitleEncoders++
			default:
				otherEncoders++
			}
		} else {
			flags += "."
		}

		// Truncate long codec names and descriptions to match format style
		codecName := info.name
		if len(codecName) > 24 {
			codecName = codecName[:24]
		}

		description := info.longName
		if len(description) > 42 {
			description = description[:42]
		}

		fmt.Printf(" %s  %-24s %-42s [%s]\n", flags, codecName, description, info.mediaType)
	}

	totalDecoders := videoDecoders + audioDecoders + subtitleDecoders + otherDecoders
	totalEncoders := videoEncoders + audioEncoders + subtitleEncoders + otherEncoders

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total codecs: %d\n", len(codecMap))
	fmt.Printf("  Decoders: %d (Video: %d, Audio: %d, Subtitle: %d, Other: %d)\n",
		totalDecoders, videoDecoders, audioDecoders, subtitleDecoders, otherDecoders)
	fmt.Printf("  Encoders: %d (Video: %d, Audio: %d, Subtitle: %d, Other: %d)\n",
		totalEncoders, videoEncoders, audioEncoders, subtitleEncoders, otherEncoders)
	fmt.Println("\nFlags:")
	fmt.Println("  D - Decoder available")
	fmt.Println("  E - Encoder available")
}

func getMediaTypeString(mediaType ffmpeg.AVMediaType) string {
	switch mediaType {
	case ffmpeg.AVMediaTypeVideo:
		return "VIDEO"
	case ffmpeg.AVMediaTypeAudio:
		return "AUDIO"
	case ffmpeg.AVMediaTypeSubtitle:
		return "SUBTITLE"
	case ffmpeg.AVMediaTypeData:
		return "DATA"
	case ffmpeg.AVMediaTypeAttachment:
		return "ATTACH"
	default:
		return "UNKNOWN"
	}
}

func getCodecName(codecID ffmpeg.AVCodecID) string {
	desc := ffmpeg.AVCodecDescriptorGet(codecID)
	if desc != nil && desc.Name() != nil {
		return desc.Name().String()
	}
	return fmt.Sprintf("codec_%d", codecID)
}

func listFormats() {
	fmt.Println("\n==================================================")
	fmt.Println("FORMATS")
	fmt.Println("==================================================\n")
	fmt.Printf("%s  %-24s %-42s %-35s %s\n", "DE", "NAME", "DESCRIPTION", "CODECS", "MIME TYPE")
	fmt.Println()

	type formatInfo struct {
		name          string
		longName      string
		exts          string
		mimeType      string
		videoCodec    string
		audioCodec    string
		subtitleCodec string
		hasMuxer      bool
		hasDemuxer    bool
	}

	formatMap := make(map[string]*formatInfo)

	// Iterate through all registered muxers
	var muxerOpaque unsafe.Pointer
	for {
		muxer := ffmpeg.AVMuxerIterate(&muxerOpaque)
		if muxer == nil {
			break
		}

		name := ""
		if muxer.Name() != nil {
			name = muxer.Name().String()
		}

		if name == "" {
			continue
		}

		longName := ""
		if muxer.LongName() != nil {
			longName = muxer.LongName().String()
		}

		extensions := ""
		if muxer.Extensions() != nil {
			extensions = muxer.Extensions().String()
		}

		mimeType := ""
		if muxer.MimeType() != nil {
			mimeType = muxer.MimeType().String()
		}

		videoCodec := ""
		if muxer.VideoCodec() != ffmpeg.AVCodecIdNone {
			videoCodec = getCodecName(muxer.VideoCodec())
		}

		audioCodec := ""
		if muxer.AudioCodec() != ffmpeg.AVCodecIdNone {
			audioCodec = getCodecName(muxer.AudioCodec())
		}

		subtitleCodec := ""
		if muxer.SubtitleCodec() != ffmpeg.AVCodecIdNone {
			subtitleCodec = getCodecName(muxer.SubtitleCodec())
		}

		if existing, exists := formatMap[name]; exists {
			existing.hasMuxer = true
			if existing.longName == "" {
				existing.longName = longName
			}
			if existing.mimeType == "" {
				existing.mimeType = mimeType
			}
			if existing.videoCodec == "" {
				existing.videoCodec = videoCodec
			}
			if existing.audioCodec == "" {
				existing.audioCodec = audioCodec
			}
			if existing.subtitleCodec == "" {
				existing.subtitleCodec = subtitleCodec
			}
		} else {
			formatMap[name] = &formatInfo{
				name:          name,
				longName:      longName,
				exts:          extensions,
				mimeType:      mimeType,
				videoCodec:    videoCodec,
				audioCodec:    audioCodec,
				subtitleCodec: subtitleCodec,
				hasMuxer:      true,
				hasDemuxer:    false,
			}
		}
	}

	// Iterate through all registered demuxers
	var demuxerOpaque unsafe.Pointer
	for {
		demuxer := ffmpeg.AVDemuxerIterate(&demuxerOpaque)
		if demuxer == nil {
			break
		}

		name := ""
		if demuxer.Name() != nil {
			name = demuxer.Name().String()
		}

		if name == "" {
			continue
		}

		longName := ""
		if demuxer.LongName() != nil {
			longName = demuxer.LongName().String()
		}

		extensions := ""
		if demuxer.Extensions() != nil {
			extensions = demuxer.Extensions().String()
		}

		if existing, exists := formatMap[name]; exists {
			existing.hasDemuxer = true
			if existing.longName == "" {
				existing.longName = longName
			}
		} else {
			formatMap[name] = &formatInfo{
				name:          name,
				longName:      longName,
				exts:          extensions,
				mimeType:      "",
				videoCodec:    "",
				audioCodec:    "",
				subtitleCodec: "",
				hasMuxer:      false,
				hasDemuxer:    true,
			}
		}
	} // Convert map to slice and sort
	var formats []formatInfo
	for _, f := range formatMap {
		formats = append(formats, *f)
	}

	sort.Slice(formats, func(i, j int) bool {
		return formats[i].name < formats[j].name
	})

	// Count totals
	totalMuxers := 0
	totalDemuxers := 0

	for _, f := range formats {
		demuxFlag := "."
		if f.hasDemuxer {
			demuxFlag = "D"
			totalDemuxers++
		}

		muxFlag := "."
		if f.hasMuxer {
			muxFlag = "E"
			totalMuxers++
		}

		// Build codec list
		codecs := []string{}
		if f.videoCodec != "" {
			codecs = append(codecs, fmt.Sprintf("video:%s", f.videoCodec))
		}
		if f.audioCodec != "" {
			codecs = append(codecs, fmt.Sprintf("audio:%s", f.audioCodec))
		}
		if f.subtitleCodec != "" {
			codecs = append(codecs, fmt.Sprintf("subtitle:%s", f.subtitleCodec))
		}

		codecList := strings.Join(codecs, ",")
		if len(codecList) > 35 {
			codecList = codecList[:35]
		}

		mimeType := f.mimeType
		if len(mimeType) > 20 {
			mimeType = mimeType[:20]
		}

		// Truncate long names and descriptions
		formatName := f.name
		if len(formatName) > 24 {
			formatName = formatName[:24]
		}

		description := f.longName
		if len(description) > 42 {
			description = description[:42]
		}

		// Display format on single line
		fmt.Printf("%s%s  %-24s %-42s %-35s %s\n", demuxFlag, muxFlag, formatName, description, codecList, mimeType)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total formats: %d\n", len(formats))
	fmt.Printf("  Total demuxers: %d\n", totalDemuxers)
	fmt.Printf("  Total muxers: %d\n", totalMuxers)
}

func listBSFs() {
	fmt.Println("\n==================================================")
	fmt.Println("BITSTREAM FILTERS")
	fmt.Println("==================================================\n")
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
		name := bsf.Name()

		// Truncate name if too long
		if len(name) > 24 {
			name = name[:24]
		}

		// Get supported codec IDs
		codecList := "all"
		codecIDs := bsf.CodecIds()
		if codecIDs != nil {
			codecSpecificCount++
			var codecs []string
			for i := uintptr(0); ; i++ {
				codecID := (*ffmpeg.AVCodecID)(unsafe.Pointer(uintptr(unsafe.Pointer(codecIDs)) + i*unsafe.Sizeof(*codecIDs)))
				if *codecID == ffmpeg.AVCodecIdNone {
					break
				}
				codecName := getCodecName(*codecID)
				codecs = append(codecs, codecName)
			}
			if len(codecs) > 0 {
				codecList = strings.Join(codecs, ", ")
				if len(codecList) > 64 {
					codecList = codecList[:61] + "..."
				}
			}
		}

		fmt.Printf("    %-24s %-42s\n", name, codecList)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total bitstream filters: %d\n", count)
	fmt.Printf("  Codec-specific filters: %d\n", codecSpecificCount)
	fmt.Printf("  Generic filters: %d\n", count-codecSpecificCount)
}

func listParsers() {
	fmt.Println("\n==================================================")
	fmt.Println("PARSERS")
	fmt.Println("==================================================\n")
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
			codecName := getCodecName(ffmpeg.AVCodecID(codecID))
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
	sort.Slice(parsers, func(i, j int) bool {
		return parsers[i].name < parsers[j].name
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

func listProtocols() {
	fmt.Println("\n==================================================")
	fmt.Println("PROTOCOLS")
	fmt.Println("==================================================\n")
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

	sort.Slice(protocols, func(i, j int) bool {
		return protocols[i].name < protocols[j].name
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
func testHardwareAvailable(deviceType ffmpeg.AVHWDeviceType) bool {
	// Save current log level and temporarily silence FFmpeg logs
	// to avoid error messages for unavailable hardware
	oldLevel, _ := ffmpeg.AVLogGetLevel()
	ffmpeg.AVLogSetLevel(ffmpeg.AVLogQuiet)
	defer ffmpeg.AVLogSetLevel(oldLevel)

	var hwDeviceCtx *ffmpeg.AVBufferRef
	ret, _ := ffmpeg.AVHWDeviceCtxCreate(&hwDeviceCtx, deviceType, nil, nil, 0)
	if ret == 0 && hwDeviceCtx != nil {
		// Successfully created - hardware is available
		ffmpeg.AVBufferUnref(&hwDeviceCtx)
		return true
	}
	return false
}

// getHWDeviceTypeFromCodec determines the hardware device type from codec's hardware config
func getHWDeviceTypeFromCodec(codec *ffmpeg.AVCodec) ffmpeg.AVHWDeviceType {
	// Iterate through hardware configurations to find device type
	for i := 0; ; i++ {
		hwConfig := ffmpeg.AVCodecGetHWConfig(codec, i)
		if hwConfig == nil {
			break
		}

		deviceType := hwConfig.DeviceType()
		if deviceType != ffmpeg.AVHWDeviceTypeNone {
			return deviceType
		}
	}

	// Fallback to name-based detection if no hardware config found
	// This handles some edge cases where the codec metadata might be incomplete
	codecName := ""
	if codec.Name() != nil {
		codecName = codec.Name().String()
	}

	switch {
	case strings.Contains(codecName, "_nvenc") || strings.Contains(codecName, "_nvdec") || strings.Contains(codecName, "_cuvid"):
		return ffmpeg.AVHWDeviceTypeCuda
	case strings.Contains(codecName, "_qsv"):
		return ffmpeg.AVHWDeviceTypeQsv
	case strings.Contains(codecName, "_vaapi"):
		return ffmpeg.AVHWDeviceTypeVaapi
	case strings.Contains(codecName, "_videotoolbox"):
		return ffmpeg.AVHWDeviceTypeVideotoolbox
	case strings.Contains(codecName, "_vulkan"):
		return ffmpeg.AVHWDeviceTypeVulkan
	case strings.Contains(codecName, "_amf"):
		return ffmpeg.AVHWDeviceTypeAmf
	case strings.Contains(codecName, "_mediacodec"):
		return ffmpeg.AVHWDeviceTypeMediacodec
	case strings.Contains(codecName, "_dxva2"):
		return ffmpeg.AVHWDeviceTypeDxva2
	case strings.Contains(codecName, "_d3d11va"):
		return ffmpeg.AVHWDeviceTypeD3D11Va
	case strings.Contains(codecName, "_vdpau"):
		return ffmpeg.AVHWDeviceTypeVdpau
	default:
		return ffmpeg.AVHWDeviceTypeNone
	}
}

func listHWAccels() {
	fmt.Println("\n==================================================")
	fmt.Println("HARDWARE ACCELERATORS")
	fmt.Println("==================================================\n")

	// First list hardware device types
	fmt.Printf("    %-24s\n", "NAME")
	fmt.Println()

	var hwaccels []string

	// Iterate through hardware device types
	deviceType := ffmpeg.AVHWDeviceTypeNone
	for {
		deviceType = ffmpeg.AVHWDeviceIterateTypes(deviceType)
		if deviceType == ffmpeg.AVHWDeviceTypeNone {
			break
		}

		name := ffmpeg.AVHWDeviceGetTypeName(deviceType)
		if name != nil {
			hwaccels = append(hwaccels, name.String())
		}
	}

	// Sort hwaccels by name
	sort.Strings(hwaccels)

	// Display hwaccels
	for _, name := range hwaccels {
		// Truncate long names to 24 chars
		if len(name) > 24 {
			name = name[:24]
		}
		fmt.Printf("    %-24s\n", name)
	}

	// Build a map of available hardware devices for quick lookup
	availableHW := make(map[ffmpeg.AVHWDeviceType]bool)
	testDeviceType := ffmpeg.AVHWDeviceTypeNone
	for {
		testDeviceType = ffmpeg.AVHWDeviceIterateTypes(testDeviceType)
		if testDeviceType == ffmpeg.AVHWDeviceTypeNone {
			break
		}
		availableHW[testDeviceType] = testHardwareAvailable(testDeviceType)
	}

	// Unified list of hardware acceleration (hwaccels + hardware codecs)
	fmt.Println("\n==================================================")
	fmt.Println("HARDWARE CODECS")
	fmt.Println("==================================================\n")
	fmt.Printf(" %s  %-24s %-42s %s %s\n", "DEH", "NAME", "DESCRIPTION", "TYPE", "PRESENT")
	fmt.Println()

	type unifiedHWEntry struct {
		name        string
		description string
		mediaType   string
		isDecoder   bool
		isEncoder   bool
		isHwaccel   bool
		deviceType  ffmpeg.AVHWDeviceType
	}

	var unifiedList []unifiedHWEntry

	// First, collect all hwaccels by iterating through decoders
	var hwaccelOpaque unsafe.Pointer
	for {
		codec := ffmpeg.AVCodecIterate(&hwaccelOpaque)
		if codec == nil {
			break
		}

		// Only check decoders
		isDecoderVal, _ := ffmpeg.AVCodecIsDecoder(codec)
		if isDecoderVal == 0 {
			continue
		}

		codecName := ""
		if codec.Name() != nil {
			codecName = codec.Name().String()
		}

		// Skip hardware-specific decoders (we want software decoders with hwaccel support)
		if strings.Contains(codecName, "_cuvid") || strings.Contains(codecName, "_qsv") ||
			strings.Contains(codecName, "_nvdec") || strings.Contains(codecName, "_vulkan") ||
			strings.Contains(codecName, "_vaapi") || strings.Contains(codecName, "_videotoolbox") {
			continue
		}

		mediaType := getMediaTypeString(codec.Type())

		// Get codec descriptor for better description
		codecDesc := ffmpeg.AVCodecDescriptorGet(codec.Id())
		baseDescription := ""
		if codecDesc != nil && codecDesc.LongName() != nil {
			baseDescription = codecDesc.LongName().String()
		}

		// Check for hardware acceleration configurations
		for i := 0; ; i++ {
			hwConfig := ffmpeg.AVCodecGetHWConfig(codec, i)
			if hwConfig == nil {
				break
			}

			deviceType := hwConfig.DeviceType()
			if deviceType == ffmpeg.AVHWDeviceTypeNone {
				continue
			}

			// Get device type name
			deviceTypeName := ""
			if name := ffmpeg.AVHWDeviceGetTypeName(deviceType); name != nil {
				deviceTypeName = name.String()
			}

			// Construct hwaccel name (e.g., h264_vulkan)
			hwaccelName := codecName + "_" + deviceTypeName

			// Construct description (e.g., "H.264 (Vulkan hardware acceleration)")
			description := baseDescription
			if description == "" {
				description = codecName
			}
			// Capitalize device type name for description
			deviceTypeCapitalized := strings.ToUpper(deviceTypeName[:1]) + deviceTypeName[1:]
			description = description + " (" + deviceTypeCapitalized + " hardware acceleration)"

			unifiedList = append(unifiedList, unifiedHWEntry{
				name:        hwaccelName,
				description: description,
				mediaType:   mediaType,
				isDecoder:   true,
				isEncoder:   false,
				isHwaccel:   true,
				deviceType:  deviceType,
			})
		}
	}

	// Now collect hardware codecs (encoders/decoders)
	var codecOpaque unsafe.Pointer
	for {
		codec := ffmpeg.AVCodecIterate(&codecOpaque)
		if codec == nil {
			break
		}

		name := ""
		if codec.Name() != nil {
			name = codec.Name().String()
		}

		if name == "" {
			continue
		}

		// Skip software decoders that only have hwaccel support (we already added those as hwaccels)
		// We only want dedicated hardware encoder/decoder implementations
		if !strings.Contains(name, "_cuvid") && !strings.Contains(name, "_qsv") &&
			!strings.Contains(name, "_nvenc") && !strings.Contains(name, "_nvdec") &&
			!strings.Contains(name, "_vulkan") && !strings.Contains(name, "_vaapi") &&
			!strings.Contains(name, "_videotoolbox") && !strings.Contains(name, "_amf") &&
			!strings.Contains(name, "_mediacodec") && !strings.Contains(name, "_dxva2") &&
			!strings.Contains(name, "_d3d11va") && !strings.Contains(name, "_vdpau") {
			continue
		}

		// Check if this is a hardware-accelerated codec by checking capabilities
		// and hardware configurations
		capabilities := codec.Capabilities()
		isHWCodec := (capabilities & ffmpeg.AVCodecCapHardware) != 0

		// Also check if codec has hardware configurations
		if !isHWCodec {
			hwConfig := ffmpeg.AVCodecGetHWConfig(codec, 0)
			isHWCodec = hwConfig != nil
		}

		if !isHWCodec {
			continue
		}

		longName := ""
		if codec.LongName() != nil {
			longName = codec.LongName().String()
		}

		mediaType := getMediaTypeString(codec.Type())

		isEncoderVal, _ := ffmpeg.AVCodecIsEncoder(codec)
		isEncoder := isEncoderVal != 0

		isDecoderVal, _ := ffmpeg.AVCodecIsDecoder(codec)
		isDecoder := isDecoderVal != 0

		// Determine hardware device type from codec's hardware config
		deviceType := getHWDeviceTypeFromCodec(codec)

		// Add to unified list
		unifiedList = append(unifiedList, unifiedHWEntry{
			name:        name,
			description: longName,
			mediaType:   mediaType,
			isDecoder:   isDecoder,
			isEncoder:   isEncoder,
			isHwaccel:   false,
			deviceType:  deviceType,
		})
	}

	// Sort by name
	sort.Slice(unifiedList, func(i, j int) bool {
		return unifiedList[i].name < unifiedList[j].name
	})

	// Count encoders, decoders, and hwaccels
	hwaccelCount := 0
	hwEncoders := 0
	hwDecoders := 0

	// Display unified list
	for _, info := range unifiedList {
		flags := ""
		if info.isDecoder {
			flags += "D"
			if !info.isHwaccel {
				hwDecoders++
			}
		} else {
			flags += "."
		}
		if info.isEncoder {
			flags += "E"
			if !info.isHwaccel {
				hwEncoders++
			}
		} else {
			flags += "."
		}
		if info.isHwaccel {
			flags += "H"
			hwaccelCount++
		} else {
			flags += "."
		}

		codecName := info.name
		if len(codecName) > 24 {
			codecName = codecName[:24]
		}

		description := info.description
		if len(description) > 42 {
			description = description[:42]
		}

		// Check if hardware is present
		present := "N"
		if info.deviceType != ffmpeg.AVHWDeviceTypeNone {
			if availableHW[info.deviceType] {
				present = "Y"
			}
		}

		fmt.Printf(" %s  %-24s %-42s [%s]  %s\n", flags, codecName, description, info.mediaType, present)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total hardware device types: %d\n", len(hwaccels))
	fmt.Printf("  Total hwaccels: %d\n", hwaccelCount)
	fmt.Printf("  Total hardware codecs: %d\n", len(unifiedList)-hwaccelCount)
	fmt.Printf("  Hardware decoders: %d\n", hwDecoders)
	fmt.Printf("  Hardware encoders: %d\n", hwEncoders)
}

func listFilters() {
	fmt.Println("\n==================================================")
	fmt.Println("FILTERS")
	fmt.Println("==================================================\n")
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
	sort.Slice(filters, func(i, j int) bool {
		return filters[i].name < filters[j].name
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
type codecDependencies struct {
	codecName    string
	longName     string
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
	// Build lookup maps for all components
	codecNameMap := buildCodecNameMap()
	parserMap := buildParserMap()
	formatMap := buildFormatMap()
	bsfMap := buildBSFMap()

	// Find all matching codecs using improved matching logic
	matches := findMatchingCodecs(codecName, codecNameMap, parserMap, formatMap, bsfMap)

	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No codec found matching '%s'\n", codecName)
		os.Exit(1)
	}

	// Sort codecs: exact match first, then software, then hardware
	sortedMatches := sortCodecsByPriority(matches, codecName)

	// Consolidate all codecs into single dependency set
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
	bsfMap map[ffmpeg.AVCodecID][]string) []*codecInfo {

	searchLower := strings.ToLower(search)
	matchedCodecs := make(map[string]*codecInfo)

	// 1. Direct codec name matching with improved logic
	for _, info := range codecNameMap {
		nameLower := strings.ToLower(info.name)

		// Exact match
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
				// Find all codecs with this codecID
				for _, info := range codecNameMap {
					if info.codecID == codecID {
						matchedCodecs[info.name] = info
					}
				}
			}
		}
	}

	// 3. Reverse lookup from formats (demuxers/muxers)
	for codecID := range formatMap {
		codecName := getCodecName(codecID)
		// Check if format name contains search term (e.g., "avif" for av1, "matroska" for various codecs)
		if strings.Contains(strings.ToLower(codecName), searchLower) {
			for _, info := range codecNameMap {
				if info.codecID == codecID {
					matchedCodecs[info.name] = info
				}
			}
		}
	}

	// 4. Reverse lookup from BSFs
	for codecID, bsfNames := range bsfMap {
		for _, bsfName := range bsfNames {
			if strings.Contains(strings.ToLower(bsfName), searchLower) {
				// Find all codecs with this codecID
				for _, info := range codecNameMap {
					if info.codecID == codecID {
						matchedCodecs[info.name] = info
					}
				}
			}
		}
	}

	// Convert map to slice
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

		// Exact match
		if nameLower == searchLower {
			exact = append(exact, codec)
		} else if isHardwareCodec(codec.name) {
			hardware = append(hardware, codec)
		} else {
			software = append(software, codec)
		}
	}

	// Combine: exact first, then software, then hardware
	result := make([]*codecInfo, 0, len(codecs))
	result = append(result, exact...)
	result = append(result, software...)
	result = append(result, hardware...)

	return result
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
	bsfMap map[ffmpeg.AVCodecID][]string) *codecDependencies {

	deps := &codecDependencies{}
	deps.descriptions = make([]string, 0)

	descriptionSet := make(map[string]bool)
	decoderSet := make(map[string]bool)
	encoderSet := make(map[string]bool)
	parserSet := make(map[string]bool)
	demuxerSet := make(map[string]bool)
	muxerSet := make(map[string]bool)
	bsfSet := make(map[string]bool)

	// Collect all components from all codecs
	for _, codec := range codecs {
		// Add description (deduplicate)
		if !descriptionSet[codec.longName] {
			deps.descriptions = append(deps.descriptions, codec.longName)
			descriptionSet[codec.longName] = true
		}

		codecID := codec.codecID

		// Get encoders and decoders
		if codec.isEncoder {
			encoderSet[codec.name] = true
		}
		if codec.isDecoder {
			decoderSet[codec.name] = true
		}

		// Get parsers
		if parsers, ok := parserMap[codecID]; ok {
			for _, p := range parsers {
				parserSet[p] = true
			}
		}

		// Get formats
		if formats, ok := formatMap[codecID]; ok {
			for _, d := range formats.demuxers {
				demuxerSet[d] = true
			}
			for _, m := range formats.muxers {
				muxerSet[m] = true
			}
		}

		// Get BSFs
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
type formatDeps struct {
	demuxers []string
	muxers   []string
}

func buildCodecNameMap() map[string]*codecInfo {
	codecMap := make(map[string]*codecInfo)

	var opaque unsafe.Pointer
	for {
		codec := ffmpeg.AVCodecIterate(&opaque)
		if codec == nil {
			break
		}

		codecID := codec.Id()
		name := ""
		if codec.Name() != nil {
			name = codec.Name().String()
		}

		if name == "" {
			continue
		}

		longName := ""
		if codec.LongName() != nil {
			longName = codec.LongName().String()
		}

		mediaType := getMediaTypeString(codec.Type())

		isEncoderVal, _ := ffmpeg.AVCodecIsEncoder(codec)
		isEncoder := isEncoderVal != 0

		isDecoderVal, _ := ffmpeg.AVCodecIsDecoder(codec)
		isDecoder := isDecoderVal != 0

		// Check if codec already exists in map (e.g., separate encoder/decoder entries)
		if existing, exists := codecMap[name]; exists {
			// Merge encoder/decoder flags
			existing.isEncoder = existing.isEncoder || isEncoder
			existing.isDecoder = existing.isDecoder || isDecoder
		} else {
			// Create new entry
			codecMap[name] = &codecInfo{
				name:      name,
				longName:  longName,
				mediaType: mediaType,
				codecID:   codecID,
				isEncoder: isEncoder,
				isDecoder: isDecoder,
			}
		}
	}

	return codecMap
}

func buildParserMap() map[ffmpeg.AVCodecID][]string {
	parserMap := make(map[ffmpeg.AVCodecID][]string)

	var opaque unsafe.Pointer
	for {
		parser := ffmpeg.AVParserIterate(&opaque)
		if parser == nil {
			break
		}

		codecIDArray := parser.CodecIds()
		for i := uintptr(0); ; i++ {
			codecID := codecIDArray.Get(i)
			if codecID == 0 {
				break
			}

			codecName := getCodecName(ffmpeg.AVCodecID(codecID))
			parserMap[ffmpeg.AVCodecID(codecID)] = append(parserMap[ffmpeg.AVCodecID(codecID)], codecName)
		}
	}

	return parserMap
}

func buildFormatMap() map[ffmpeg.AVCodecID]*formatDeps {
	formatMap := make(map[ffmpeg.AVCodecID]*formatDeps)

	// Note: Demuxers don't expose codec info directly via the API
	// We'll use muxer data and assume demuxers support the same codecs

	// Iterate through muxers
	var muxerOpaque unsafe.Pointer
	for {
		muxer := ffmpeg.AVMuxerIterate(&muxerOpaque)
		if muxer == nil {
			break
		}

		name := ""
		if muxer.Name() != nil {
			name = muxer.Name().String()
		}

		// Add format for video codec
		if muxer.VideoCodec() != ffmpeg.AVCodecIdNone {
			codecID := muxer.VideoCodec()
			if _, exists := formatMap[codecID]; !exists {
				formatMap[codecID] = &formatDeps{}
			}
			formatMap[codecID].muxers = append(formatMap[codecID].muxers, name)
			formatMap[codecID].demuxers = append(formatMap[codecID].demuxers, name)
		}

		// Add format for audio codec
		if muxer.AudioCodec() != ffmpeg.AVCodecIdNone {
			codecID := muxer.AudioCodec()
			if _, exists := formatMap[codecID]; !exists {
				formatMap[codecID] = &formatDeps{}
			}
			formatMap[codecID].muxers = append(formatMap[codecID].muxers, name)
			formatMap[codecID].demuxers = append(formatMap[codecID].demuxers, name)
		}

		// Add format for subtitle codec
		if muxer.SubtitleCodec() != ffmpeg.AVCodecIdNone {
			codecID := muxer.SubtitleCodec()
			if _, exists := formatMap[codecID]; !exists {
				formatMap[codecID] = &formatDeps{}
			}
			formatMap[codecID].muxers = append(formatMap[codecID].muxers, name)
			formatMap[codecID].demuxers = append(formatMap[codecID].demuxers, name)
		}
	}

	return formatMap
}

func buildBSFMap() map[ffmpeg.AVCodecID][]string {
	bsfMap := make(map[ffmpeg.AVCodecID][]string)
	var genericBSFs []string

	var opaque unsafe.Pointer
	for {
		bsf := ffmpeg.AVBSFIterate(&opaque)
		if bsf == nil {
			break
		}

		name := bsf.Name()
		codecIDs := bsf.CodecIds()

		if codecIDs == nil {
			// Generic BSF - applies to all codecs
			genericBSFs = append(genericBSFs, name)
		} else {
			// Codec-specific BSF
			for i := uintptr(0); ; i++ {
				codecID := (*ffmpeg.AVCodecID)(unsafe.Pointer(uintptr(unsafe.Pointer(codecIDs)) + i*unsafe.Sizeof(*codecIDs)))
				if *codecID == ffmpeg.AVCodecIdNone {
					break
				}
				bsfMap[*codecID] = append(bsfMap[*codecID], name)
			}
		}
	}

	// Add generic BSFs to all codec IDs
	for codecID := range bsfMap {
		bsfMap[codecID] = append(bsfMap[codecID], genericBSFs...)
	}

	return bsfMap
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
