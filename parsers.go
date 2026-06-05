package ffmpeg

/*
#cgo CFLAGS: -I${SRCDIR}/include -Wno-deprecated -Wno-deprecated-declarations

#include <stdint.h>
#include <stddef.h>
#include <libavcodec/ac3_parser.h>
#include <libavcodec/adts_parser.h>
#include <libavcodec/vorbis_parser.h>
*/
import "C"

import "unsafe"

// Hand-written wrappers for libavcodec parser entry points the generator skips.
// Each takes a const uint8_t *buf the C side only reads and never retains, plus
// primitive out-pointers (bitstream id, frame/sample counts, flags) the
// generator cannot classify as in or out. A flat Go []byte is a single pointer
// into pointer-free memory, so passing &buf[0] directly satisfies the cgo
// pointer rules: this is not the array-of-pointers case that samples.go needs
// C-heap scratch for. The out-params become Go return values, with the status
// code folded into error via WrapErr. Every entry guards bad input before
// touching &buf[0], which would panic on a zero-length slice.

// AVAc3ParseHeader parses an AC-3 frame header from the start of buf.
//
// buf must be non-empty. Returns the bitstream id, the frame size in bytes, and
// an error (nil on success, an AVError otherwise).
func AVAc3ParseHeader(buf []byte) (bitstreamID uint8, frameSize uint16, err error) {
	if len(buf) == 0 {
		return 0, 0, WrapErr(AVErrorUnknownConst)
	}

	var cBitstreamID C.uint8_t
	var cFrameSize C.uint16_t
	ret := C.av_ac3_parse_header(
		(*C.uint8_t)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)),
		&cBitstreamID, &cFrameSize,
	)
	return uint8(cBitstreamID), uint16(cFrameSize), WrapErr(int(ret))
}

// AVAdtsHeaderParse extracts the sample and frame counts from an ADTS header at
// the start of buf. buf must hold at least AVAacAdtsHeaderSize (7) bytes.
//
// Returns the number of samples, the number of frames, and an error (nil on
// success, an AVError otherwise).
func AVAdtsHeaderParse(buf []byte) (samples uint32, frames uint8, err error) {
	if len(buf) < AVAacAdtsHeaderSize {
		return 0, 0, WrapErr(AVErrorUnknownConst)
	}

	var cSamples C.uint32_t
	var cFrames C.uint8_t
	ret := C.av_adts_header_parse(
		(*C.uint8_t)(unsafe.Pointer(&buf[0])),
		&cSamples, &cFrames,
	)
	return uint32(cSamples), uint8(cFrames), WrapErr(int(ret))
}

// AVVorbisParseFrameFlags parses one Vorbis frame from buf, returning its
// duration in samples and the frame flags. The duration mirrors the generated
// AVVorbisParseFrame return shape: the non-negative C return value is surfaced
// as duration, with negative values folded into err.
//
// s must be non-nil and buf must be non-empty.
func AVVorbisParseFrameFlags(s *AVVorbisParseContext, buf []byte) (duration int, flags int, err error) {
	if s == nil || len(buf) == 0 {
		return 0, 0, WrapErr(AVErrorUnknownConst)
	}

	var cFlags C.int
	ret := C.av_vorbis_parse_frame_flags(
		s.ptr,
		(*C.uint8_t)(unsafe.Pointer(&buf[0])), C.int(len(buf)),
		&cFlags,
	)
	return int(ret), int(cFlags), WrapErr(int(ret))
}
