package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"unsafe"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

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
	fmt.Println("==================================================")

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
	slices.Sort(hwaccels)

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
	fmt.Println("==================================================")
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
	slices.SortFunc(unifiedList, func(a, b unifiedHWEntry) int {
		return cmp.Compare(a.name, b.name)
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
