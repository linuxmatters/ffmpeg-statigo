package ffmpeg

/*
#include <libavformat/avformat.h>
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
