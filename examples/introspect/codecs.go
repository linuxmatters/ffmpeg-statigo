package main

import (
	"fmt"
	"slices"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

func listCodecs() {
	fmt.Println("==================================================")
	fmt.Println("CODECS")
	fmt.Println("==================================================")
	fmt.Printf(" %s  %-24s %-42s %s\n", "DE", "NAME", "DESCRIPTION", "TYPE")
	fmt.Println()

	codecMap := make(map[string]*codecInfo)

	var desc *ffmpeg.AVCodecDescriptor
	for {
		desc = ffmpeg.AVCodecDescriptorNext(desc)
		if desc == nil {
			break
		}

		codecID := desc.Id()
		name := ""
		if desc.Name() != nil {
			name = desc.Name().String()
		}
		longName := ""
		if desc.LongName() != nil {
			longName = desc.LongName().String()
		}
		mediaType := getMediaTypeString(desc.Type())

		encoder := ffmpeg.AVCodecFindEncoder(codecID)
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

	names := make([]string, 0, len(codecMap))
	for name := range codecMap {
		names = append(names, name)
	}
	slices.Sort(names)

	videoEncoders, videoDecoders := 0, 0
	audioEncoders, audioDecoders := 0, 0
	subtitleEncoders, subtitleDecoders := 0, 0
	otherEncoders, otherDecoders := 0, 0

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
