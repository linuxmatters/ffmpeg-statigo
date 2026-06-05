package main

import (
	"maps"
	"slices"
	"unsafe"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
)

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

			cid := ffmpeg.AVCodecID(codecID) //nolint:gosec // G115: codec IDs are small enum values
			codecName := getCodecName(cid)
			parserMap[cid] = append(parserMap[cid], codecName)
		}
	}

	return parserMap
}

func buildFormatMap() map[ffmpeg.AVCodecID]*formatDeps {
	formatMap := make(map[ffmpeg.AVCodecID]*formatDeps)

	// AVInputFormat exposes no codec mapping, so demuxer support is inferred from
	// muxers that share a container name with a registered demuxer. A mux-only
	// format (no matching demuxer) is therefore never listed as a demuxer.
	demuxerNames := buildDemuxerNameSet()

	addFormat := func(codecID ffmpeg.AVCodecID, name string) {
		if codecID == ffmpeg.AVCodecIdNone {
			return
		}
		if _, exists := formatMap[codecID]; !exists {
			formatMap[codecID] = &formatDeps{}
		}
		formatMap[codecID].muxers = append(formatMap[codecID].muxers, name)
		if demuxerNames[name] {
			formatMap[codecID].demuxers = append(formatMap[codecID].demuxers, name)
		}
	}

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

		addFormat(muxer.VideoCodec(), name)
		addFormat(muxer.AudioCodec(), name)
		addFormat(muxer.SubtitleCodec(), name)
	}

	return formatMap
}

func buildDemuxerNameSet() map[string]bool {
	names := make(map[string]bool)

	var demuxerOpaque unsafe.Pointer
	for {
		demuxer := ffmpeg.AVDemuxerIterate(&demuxerOpaque)
		if demuxer == nil {
			break
		}
		if demuxer.Name() != nil {
			names[demuxer.Name().String()] = true
		}
	}

	return names
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

		name := bsf.Name().String()
		codecIDs := bsf.CodecIds()

		if codecIDs == nil {
			// Generic BSF - applies to all codecs
			genericBSFs = append(genericBSFs, name)
		} else {
			// Codec-specific BSF
			for i := uintptr(0); ; i++ {
				codecID := codecIDs.Get(i)
				if codecID == ffmpeg.AVCodecIdNone {
					break
				}
				bsfMap[codecID] = append(bsfMap[codecID], name)
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
	return slices.Sorted(maps.Keys(m))
}
