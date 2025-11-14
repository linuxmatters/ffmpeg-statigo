package ffmpeg

/*
#include <libavformat/avformat.h>
#include <libavformat/avio.h>
#include <libavcodec/avcodec.h>

// Forward declarations for iteration functions
extern const AVCodec *av_codec_iterate(void **opaque);
extern const AVCodecParser *av_parser_iterate(void **opaque);
extern const char *avio_enum_protocols(void **opaque, int output);
extern const AVOutputFormat *av_muxer_iterate(void **opaque);
extern const AVInputFormat *av_demuxer_iterate(void **opaque);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// AVMuxerIterate iterates over all registered muxers.
//
// @param opaque a pointer where libavformat will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered muxer or NULL when the iteration is finished
func AVMuxerIterate(opaque *unsafe.Pointer) *AVOutputFormat {
	ret := C.av_muxer_iterate((*unsafe.Pointer)(unsafe.Pointer(opaque)))
	if ret == nil {
		return nil
	}
	return &AVOutputFormat{ptr: ret}
}

// AVDemuxerIterate iterates over all registered demuxers.
//
// @param opaque a pointer where libavformat will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered demuxer or NULL when the iteration is finished
func AVDemuxerIterate(opaque *unsafe.Pointer) *AVInputFormat {
	ret := C.av_demuxer_iterate((*unsafe.Pointer)(unsafe.Pointer(opaque)))
	if ret == nil {
		return nil
	}
	return &AVInputFormat{ptr: ret}
}

// AVParserIterate iterates over all registered codec parsers.
//
// @param opaque a pointer where libavcodec will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered parser or NULL when the iteration is finished
func AVParserIterate(opaque *unsafe.Pointer) *AVCodecParser {
	ret := (*C.AVCodecParser)(C.av_parser_iterate((*unsafe.Pointer)(unsafe.Pointer(opaque))))
	if ret == nil {
		return nil
	}
	return &AVCodecParser{ptr: ret}
}

// AVCodecIterate iterates over all registered codecs.
//
// @param opaque a pointer where libavcodec will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered codec or NULL when the iteration is finished
func AVCodecIterate(opaque *unsafe.Pointer) *AVCodec {
	ret := C.av_codec_iterate((*unsafe.Pointer)(unsafe.Pointer(opaque)))
	if ret == nil {
		return nil
	}
	return &AVCodec{ptr: ret}
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

// AVOptSetSlice is a helper for storing a slice of primitive data to the named field. This function provides no
// guarantees for usage with Go wrapper types.
//
// See AVOptSet for more information.
func AVOptSetSlice[T any](obj unsafe.Pointer, name *CStr, val []T, searchFlags int) (int, error) {
	var ty T
	ptr := unsafe.SliceData(val)
	size := unsafe.Sizeof(ty)
	return AVOptSetBin(obj, name, unsafe.Pointer(ptr), int(size)*len(val), searchFlags)
}

func (s *AVRational) String() string {
	return fmt.Sprintf("%v/%v (%v)", s.Num(), s.Den(), s.Num()/s.Den())
}

func FFIOWFourCC(s *AVIOContext, a uint8, b uint8, c uint8, d uint8) {
	AVIOWl32(s, uint(a)|(uint(b)<<8)|(uint(c)<<16)|(uint(d)<<24))
}
