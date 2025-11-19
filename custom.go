package ffmpeg

/*
#include <libavformat/avformat.h>
#include <libavformat/avio.h>
#include <libavcodec/avcodec.h>
#include <libavcodec/bsf.h>
#include <libavfilter/avfilter.h>

// Forward declarations for iteration functions
extern const AVCodec *av_codec_iterate(void **opaque);
extern const AVCodecParser *av_parser_iterate(void **opaque);
extern const char *avio_enum_protocols(void **opaque, int output);
extern const AVOutputFormat *av_muxer_iterate(void **opaque);
extern const AVInputFormat *av_demuxer_iterate(void **opaque);
extern const AVFilter *av_filter_iterate(void **opaque);
extern const AVBitStreamFilter *av_bsf_iterate(void **opaque);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// AVBitStreamFilter represents an FFmpeg bitstream filter.
// This is a minimal wrapper for the C struct.
type AVBitStreamFilter struct {
	ptr *C.AVBitStreamFilter
}

// Name returns the bitstream filter's name.
func (bsf *AVBitStreamFilter) Name() string {
	return C.GoString(bsf.ptr.name)
}

// CodecIds returns the array of codec IDs supported by this filter.
// Returns nil if the filter supports all codecs.
func (bsf *AVBitStreamFilter) CodecIds() *AVCodecID {
	if bsf.ptr.codec_ids == nil {
		return nil
	}
	return (*AVCodecID)(unsafe.Pointer(bsf.ptr.codec_ids))
}

// PrivClass returns the AVClass for private data options, or nil if none.
func (bsf *AVBitStreamFilter) PrivClass() *AVClass {
	if bsf.ptr.priv_class == nil {
		return nil
	}
	return &AVClass{ptr: bsf.ptr.priv_class}
}

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

// AVFilterIterate iterates over all registered filters.
//
// @param opaque a pointer where libavfilter will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered filter or NULL when the iteration is finished
func AVFilterIterate(opaque *unsafe.Pointer) *AVFilter {
	ret := C.av_filter_iterate((*unsafe.Pointer)(unsafe.Pointer(opaque)))
	if ret == nil {
		return nil
	}
	return &AVFilter{ptr: ret}
}

// AVBSFIterate iterates over all registered bitstream filters.
//
// @param opaque a pointer where libavcodec will store the iteration state. Must
//
//	point to NULL to start the iteration.
//
// @return the next registered bitstream filter or NULL when the iteration is finished
func AVBSFIterate(opaque *unsafe.Pointer) *AVBitStreamFilter {
	ret := C.av_bsf_iterate((*unsafe.Pointer)(unsafe.Pointer(opaque)))
	if ret == nil {
		return nil
	}
	return &AVBitStreamFilter{ptr: ret}
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

// AVWhitepointCoefficients is a typedef alias for AVCIExy in FFmpeg.
// This represents white point chromaticity coordinates.
type AVWhitepointCoefficients = AVCIExy

// AVAdler is a typedef alias for uint32_t in FFmpeg.
// This represents an Adler-32 checksum value.
type AVAdler = uint32

// AVCRC is a typedef alias for uint32_t in FFmpeg.
// This represents a CRC (Cyclic Redundancy Check) value.
type AVCRC = uint32
