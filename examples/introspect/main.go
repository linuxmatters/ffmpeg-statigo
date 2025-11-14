package main

import (
	"fmt"
	"sort"
	"strings"
	"unsafe"

	ffmpeg "github.com/csnewman/ffmpeg-go"
)

type codecInfo struct {
	name      string
	longName  string
	mediaType string
	isEncoder bool
	isDecoder bool
}

func main() {
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
	fmt.Printf("  Total decoders: %d (Video: %d, Audio: %d, Subtitle: %d, Other: %d)\n",
		totalDecoders, videoDecoders, audioDecoders, subtitleDecoders, otherDecoders)
	fmt.Printf("  Total encoders: %d (Video: %d, Audio: %d, Subtitle: %d, Other: %d)\n\n",
		totalEncoders, videoEncoders, audioEncoders, subtitleEncoders, otherEncoders)
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

	// Now list hardware-accelerated encoders and decoders
	fmt.Println("\n--------------------------------------------------\n")
	fmt.Printf(" %s  %-24s %-42s %s\n", "DE", "NAME", "DESCRIPTION", "TYPE")
	fmt.Println()

	type hwCodecInfo struct {
		name      string
		longName  string
		mediaType string
		isEncoder bool
		isDecoder bool
	}

	var hwCodecs []hwCodecInfo

	// Iterate through all codecs using av_codec_iterate
	var opaque unsafe.Pointer
	for {
		codec := ffmpeg.AVCodecIterate(&opaque)
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

		// Check if this is a hardware-accelerated codec
		// Hardware codecs typically have suffixes like _nvenc, _qsv, _vaapi, _videotoolbox, _vulkan, etc.
		isHWCodec := strings.Contains(name, "_nvenc") ||
			strings.Contains(name, "_nvdec") ||
			strings.Contains(name, "_qsv") ||
			strings.Contains(name, "_vaapi") ||
			strings.Contains(name, "_videotoolbox") ||
			strings.Contains(name, "_vulkan") ||
			strings.Contains(name, "_amf") ||
			strings.Contains(name, "_v4l2") ||
			strings.Contains(name, "_mediacodec") ||
			strings.Contains(name, "_mmal") ||
			strings.Contains(name, "_omx") ||
			strings.Contains(name, "_cuvid")

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

		hwCodecs = append(hwCodecs, hwCodecInfo{
			name:      name,
			longName:  longName,
			mediaType: mediaType,
			isEncoder: isEncoder,
			isDecoder: isDecoder,
		})
	}

	// Sort by name
	sort.Slice(hwCodecs, func(i, j int) bool {
		return hwCodecs[i].name < hwCodecs[j].name
	})

	// Count encoders and decoders
	hwEncoders := 0
	hwDecoders := 0

	// Display hardware codecs
	for _, info := range hwCodecs {
		flags := ""
		if info.isDecoder {
			flags += "D"
			hwDecoders++
		} else {
			flags += "."
		}
		if info.isEncoder {
			flags += "E"
			hwEncoders++
		} else {
			flags += "."
		}

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

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total hardware accelerators: %d\n", len(hwaccels))
	fmt.Printf("  Total hardware codecs: %d\n", len(hwCodecs))
	fmt.Printf("  Hardware decoders: %d\n", hwDecoders)
	fmt.Printf("  Hardware encoders: %d\n", hwEncoders)
}
