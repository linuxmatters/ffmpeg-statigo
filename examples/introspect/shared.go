package main

import (
	"fmt"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

type codecInfo struct {
	name      string
	longName  string
	mediaType string
	codecID   ffmpeg.AVCodecID
	isEncoder bool
	isDecoder bool
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
