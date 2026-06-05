package ffmpeg

/*
#include <libavformat/avformat.h>
#include <libavformat/avio.h>
#include <libavcodec/avcodec.h>
#include <libavcodec/bsf.h>
#include <libavfilter/avfilter.h>
#include <libavutil/channel_layout.h>

// Forward declarations for iteration functions
extern const AVCodec *av_codec_iterate(void **opaque);
extern const AVCodecParser *av_parser_iterate(void **opaque);
extern const char *avio_enum_protocols(void **opaque, int output);
extern const AVOutputFormat *av_muxer_iterate(void **opaque);
extern const AVInputFormat *av_demuxer_iterate(void **opaque);
extern const AVFilter *av_filter_iterate(void **opaque);
extern const AVBitStreamFilter *av_bsf_iterate(void **opaque);
extern const AVChannelLayout *av_channel_layout_standard(void **opaque);
*/
import "C"

import (
	"unsafe"
)

// iterateFFmpeg is a generic helper for FFmpeg iterator functions.
// It handles the common pattern of calling a C iterator, checking for nil,
// and wrapping the result in a Go type.
func iterateFFmpeg[T any](opaque *unsafe.Pointer, iterate func(*unsafe.Pointer) unsafe.Pointer, wrap func(unsafe.Pointer) *T) *T {
	ret := iterate(opaque)
	if ret == nil {
		return nil
	}
	return wrap(ret)
}

// AVMuxerIterate iterates over all registered muxers.
//
// @param opaque a pointer where libavformat will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered muxer or NULL when the iteration is finished
func AVMuxerIterate(opaque *unsafe.Pointer) *AVOutputFormat {
	return iterateFFmpeg(opaque,
		func(op *unsafe.Pointer) unsafe.Pointer {
			return unsafe.Pointer(C.av_muxer_iterate((*unsafe.Pointer)(unsafe.Pointer(op))))
		},
		func(p unsafe.Pointer) *AVOutputFormat {
			return &AVOutputFormat{ptr: (*C.AVOutputFormat)(p)}
		})
}

// AVDemuxerIterate iterates over all registered demuxers.
//
// @param opaque a pointer where libavformat will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered demuxer or NULL when the iteration is finished
func AVDemuxerIterate(opaque *unsafe.Pointer) *AVInputFormat {
	return iterateFFmpeg(opaque,
		func(op *unsafe.Pointer) unsafe.Pointer {
			return unsafe.Pointer(C.av_demuxer_iterate((*unsafe.Pointer)(unsafe.Pointer(op))))
		},
		func(p unsafe.Pointer) *AVInputFormat {
			return &AVInputFormat{ptr: (*C.AVInputFormat)(p)}
		})
}

// AVParserIterate iterates over all registered codec parsers.
//
// @param opaque a pointer where libavcodec will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered parser or NULL when the iteration is finished
func AVParserIterate(opaque *unsafe.Pointer) *AVCodecParser {
	return iterateFFmpeg(opaque,
		func(op *unsafe.Pointer) unsafe.Pointer {
			return unsafe.Pointer(C.av_parser_iterate((*unsafe.Pointer)(unsafe.Pointer(op))))
		},
		func(p unsafe.Pointer) *AVCodecParser {
			return &AVCodecParser{ptr: (*C.AVCodecParser)(p)}
		})
}

// AVCodecIterate iterates over all registered codecs.
//
// @param opaque a pointer where libavcodec will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered codec or NULL when the iteration is finished
func AVCodecIterate(opaque *unsafe.Pointer) *AVCodec {
	return iterateFFmpeg(opaque,
		func(op *unsafe.Pointer) unsafe.Pointer {
			return unsafe.Pointer(C.av_codec_iterate((*unsafe.Pointer)(unsafe.Pointer(op))))
		},
		func(p unsafe.Pointer) *AVCodec {
			return &AVCodec{ptr: (*C.AVCodec)(p)}
		})
}

// AVFilterIterate iterates over all registered filters.
//
// @param opaque a pointer where libavfilter will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered filter or NULL when the iteration is finished
func AVFilterIterate(opaque *unsafe.Pointer) *AVFilter {
	return iterateFFmpeg(opaque,
		func(op *unsafe.Pointer) unsafe.Pointer {
			return unsafe.Pointer(C.av_filter_iterate((*unsafe.Pointer)(unsafe.Pointer(op))))
		},
		func(p unsafe.Pointer) *AVFilter {
			return &AVFilter{ptr: (*C.AVFilter)(p)}
		})
}

// AVBSFIterate iterates over all registered bitstream filters.
//
// @param opaque a pointer where libavcodec will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered bitstream filter or NULL when the iteration is finished
func AVBSFIterate(opaque *unsafe.Pointer) *AVBitStreamFilter {
	return iterateFFmpeg(opaque,
		func(op *unsafe.Pointer) unsafe.Pointer {
			return unsafe.Pointer(C.av_bsf_iterate((*unsafe.Pointer)(unsafe.Pointer(op))))
		},
		func(p unsafe.Pointer) *AVBitStreamFilter {
			return &AVBitStreamFilter{ptr: (*C.AVBitStreamFilter)(p)}
		})
}

// AVChannelLayoutStandard iterates over all standard channel layouts.
//
// @param opaque a pointer where libavutil will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next standard channel layout or NULL when the iteration is finished
func AVChannelLayoutStandard(opaque *unsafe.Pointer) *AVChannelLayout {
	return iterateFFmpeg(opaque,
		func(op *unsafe.Pointer) unsafe.Pointer {
			return unsafe.Pointer(C.av_channel_layout_standard((*unsafe.Pointer)(unsafe.Pointer(op))))
		},
		func(p unsafe.Pointer) *AVChannelLayout {
			return &AVChannelLayout{ptr: (*C.AVChannelLayout)(p)}
		})
}

// AVIOEnumProtocols iterates through names of available protocols.
//
// @param opaque a pointer where libavformat will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @param output if set to 1, iterate over output protocols, otherwise over input protocols.
//
// @return a static string containing the name of current protocol or NULL
func AVIOEnumProtocols(opaque *unsafe.Pointer, output int) string {
	ret := C.avio_enum_protocols((*unsafe.Pointer)(unsafe.Pointer(opaque)), C.int(output))
	if ret == nil {
		return ""
	}
	return C.GoString(ret)
}
