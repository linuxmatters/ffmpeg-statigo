package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"unsafe"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

func listFormats() {
	fmt.Println("\n==================================================")
	fmt.Println("FORMATS")
	fmt.Println("==================================================")
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
	}
	var formats []formatInfo
	for _, f := range formatMap {
		formats = append(formats, *f)
	}

	slices.SortFunc(formats, func(a, b formatInfo) int {
		return cmp.Compare(a.name, b.name)
	})

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

		formatName := f.name
		if len(formatName) > 24 {
			formatName = formatName[:24]
		}

		description := f.longName
		if len(description) > 42 {
			description = description[:42]
		}

		fmt.Printf("%s%s  %-24s %-42s %-35s %s\n", demuxFlag, muxFlag, formatName, description, codecList, mimeType)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total formats: %d\n", len(formats))
	fmt.Printf("  Total demuxers: %d\n", totalDemuxers)
	fmt.Printf("  Total muxers: %d\n", totalMuxers)
}
