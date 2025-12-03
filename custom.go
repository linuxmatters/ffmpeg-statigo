package ffmpeg

/*
#include <libavformat/avformat.h>
#include <libavformat/avio.h>
#include <libavcodec/avcodec.h>
#include <libavcodec/bsf.h>
#include <libavfilter/avfilter.h>
#include <libavutil/uuid.h>

// Forward declarations for iteration functions
extern const AVCodec *av_codec_iterate(void **opaque);
extern const AVCodecParser *av_parser_iterate(void **opaque);
extern const char *avio_enum_protocols(void **opaque, int output);
extern const AVOutputFormat *av_muxer_iterate(void **opaque);
extern const AVInputFormat *av_demuxer_iterate(void **opaque);
extern const AVFilter *av_filter_iterate(void **opaque);
extern const AVBitStreamFilter *av_bsf_iterate(void **opaque);

// Typedef for anonymous struct inside AVStreamGroupTileGrid.
// The 'offsets' field is an array of these structs containing tile position info.
// We create a typedef so CGO can reference it by a stable name.
typedef struct {
    unsigned int idx;
    int horizontal;
    int vertical;
} AVStreamGroupTileGridOffset;
*/
import "C"
import (
	"fmt"
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

// AVUUID is a typedef to a 16-byte array in FFmpeg (uint8_t[16]).
// This represents a UUID as an opaque sequence of 16 unsigned bytes.
// Binary representation of a UUID per IETF RFC 4122.
type AVUUID = [16]uint8

// cstrPtr returns the underlying C char pointer from a CStr, or nil if s is nil.
// This helper simplifies nil-safe access to CStr pointers in CGO calls.
func cstrPtr(s *CStr) *C.char {
	if s == nil {
		return nil
	}
	return s.ptr
}

// --- Manual UUID function wrappers (arrays need pointer conversion in CGO) ---

// AVUuidParse parses a string representation of a UUID formatted according to IETF RFC 4122
// into an AVUUID. The parsing is case-insensitive. The string must be 37
// characters long, including the terminating NUL character.
//
// Example string representation: "2fceebd0-7017-433d-bafb-d073a7116696"
//
//	@param[in]  in  String representation of a UUID
//	@param[out] uu  AVUUID
//	@return         A non-zero value in case of an error.
func AVUuidParse(in *CStr, uu *AVUUID) (int, error) {
	ret := C.av_uuid_parse(cstrPtr(in), (*C.uint8_t)(unsafe.Pointer(&uu[0])))
	return int(ret), WrapErr(int(ret))
}

// AVUuidUrnParse parses a URN representation of a UUID, as specified at IETF RFC 4122,
// into an AVUUID. The parsing is case-insensitive. The string must be 46
// characters long, including the terminating NUL character.
//
// Example string representation: "urn:uuid:2fceebd0-7017-433d-bafb-d073a7116696"
//
//	@param[in]  in  URN UUID
//	@param[out] uu  AVUUID
//	@return         A non-zero value in case of an error.
func AVUuidUrnParse(in *CStr, uu *AVUUID) (int, error) {
	ret := C.av_uuid_urn_parse(cstrPtr(in), (*C.uint8_t)(unsafe.Pointer(&uu[0])))
	return int(ret), WrapErr(int(ret))
}

// AVUuidParseRange parses a string representation of a UUID formatted according to IETF RFC 4122
// into an AVUUID. The parsing is case-insensitive.
//
//	@param[in]  inStart Pointer to the first character of the string representation
//	@param[in]  inEnd   Pointer to the character after the last character of the
//	                    string representation. That memory location is never
//	                    accessed. It is an error if `inEnd - inStart != 36`.
//	@param[out] uu      AVUUID
//	@return             A non-zero value in case of an error.
func AVUuidParseRange(inStart *CStr, inEnd *CStr, uu *AVUUID) (int, error) {
	ret := C.av_uuid_parse_range(cstrPtr(inStart), cstrPtr(inEnd), (*C.uint8_t)(unsafe.Pointer(&uu[0])))
	return int(ret), WrapErr(int(ret))
}

// AVUuidUnparse serializes a AVUUID into a string representation according to IETF RFC 4122.
// The string is lowercase and always 37 characters long, including the terminating NUL character.
//
//	@param[in]  uu  AVUUID
//	@param[out] out Pointer to an array of no less than 37 characters.
func AVUuidUnparse(uu *AVUUID, out *CStr) {
	C.av_uuid_unparse((*C.uint8_t)(unsafe.Pointer(&uu[0])), cstrPtr(out))
}

// AVUuidEqual compares two UUIDs for equality.
//
//	@param[in] uu1 AVUUID
//	@param[in] uu2 AVUUID
//	@return        Nonzero if uu1 and uu2 are equal, 0 otherwise.
func AVUuidEqual(uu1 *AVUUID, uu2 *AVUUID) (int, error) {
	ret := C.av_uuid_equal((*C.uint8_t)(unsafe.Pointer(&uu1[0])), (*C.uint8_t)(unsafe.Pointer(&uu2[0])))
	return int(ret), WrapErr(int(ret))
}

// AVUuidCopy copies the bytes of src into dest.
//
//	@param[out] dest AVUUID
//	@param[in]  src  AVUUID
func AVUuidCopy(dest *AVUUID, src *AVUUID) {
	C.av_uuid_copy((*C.uint8_t)(unsafe.Pointer(&dest[0])), (*C.uint8_t)(unsafe.Pointer(&src[0])))
}

// AVUuidNil sets a UUID to the nil UUID, i.e. a UUID with have all
// its 128 bits set to zero.
//
//	@param[out] uu AVUUID
func AVUuidNil(uu *AVUUID) {
	C.av_uuid_nil((*C.uint8_t)(unsafe.Pointer(&uu[0])))
}

// --- Manual binding for anonymous struct in AVStreamGroupTileGrid ---

// AVStreamGroupTileGridOffset wraps the anonymous struct used in AVStreamGroupTileGrid.offsets.
// This is a tile grid offset structure containing the stream index and pixel offsets.
type AVStreamGroupTileGridOffset struct {
	ptr *C.AVStreamGroupTileGridOffset
}

// RawPtr returns the underlying C pointer.
func (s *AVStreamGroupTileGridOffset) RawPtr() unsafe.Pointer {
	return unsafe.Pointer(s.ptr)
}

// ToAVStreamGroupTileGridOffsetArray creates an Array accessor for a slice of offsets.
func ToAVStreamGroupTileGridOffsetArray(ptr unsafe.Pointer) *Array[*AVStreamGroupTileGridOffset] {
	if ptr == nil {
		return nil
	}

	return &Array[*AVStreamGroupTileGridOffset]{
		elemSize: unsafe.Sizeof(C.AVStreamGroupTileGridOffset{}),
		loadPtr: func(pointer unsafe.Pointer) *AVStreamGroupTileGridOffset {
			cValue := (*C.AVStreamGroupTileGridOffset)(pointer)
			return &AVStreamGroupTileGridOffset{ptr: cValue}
		},
		ptr: ptr,
		storePtr: func(pointer unsafe.Pointer, value *AVStreamGroupTileGridOffset) {
			dest := (*C.AVStreamGroupTileGridOffset)(pointer)
			if value != nil {
				*dest = *value.ptr
			}
		},
	}
}

// Idx gets the idx field.
// Index of the stream in the group this tile references.
// Must be < AVStreamGroup.nb_streams.
func (s *AVStreamGroupTileGridOffset) Idx() uint {
	return uint(s.ptr.idx)
}

// SetIdx sets the idx field.
func (s *AVStreamGroupTileGridOffset) SetIdx(value uint) {
	s.ptr.idx = (C.uint)(value)
}

// Horizontal gets the horizontal field.
// Offset in pixels from the left edge of the canvas where the tile should be placed.
func (s *AVStreamGroupTileGridOffset) Horizontal() int {
	return int(s.ptr.horizontal)
}

// SetHorizontal sets the horizontal field.
func (s *AVStreamGroupTileGridOffset) SetHorizontal(value int) {
	s.ptr.horizontal = (C.int)(value)
}

// Vertical gets the vertical field.
// Offset in pixels from the top edge of the canvas where the tile should be placed.
func (s *AVStreamGroupTileGridOffset) Vertical() int {
	return int(s.ptr.vertical)
}

// SetVertical sets the vertical field.
func (s *AVStreamGroupTileGridOffset) SetVertical(value int) {
	s.ptr.vertical = (C.int)(value)
}

// Offsets gets the offsets field from AVStreamGroupTileGrid.
// Returns an array of tile grid offsets containing stream index and pixel positions.
// The offsets field is an anonymous struct in C, so we cast to our compatible typedef.
func (s *AVStreamGroupTileGrid) Offsets() *AVStreamGroupTileGridOffset {
	// The C field is an anonymous struct pointer; cast to our compatible typedef
	value := unsafe.Pointer(s.ptr.offsets)
	if value == nil {
		return nil
	}
	return &AVStreamGroupTileGridOffset{ptr: (*C.AVStreamGroupTileGridOffset)(value)}
}

// SetOffsets sets the offsets field on AVStreamGroupTileGrid.
func (s *AVStreamGroupTileGrid) SetOffsets(value *AVStreamGroupTileGridOffset) {
	if value != nil {
		// Use unsafe pointer to bypass CGO type checking between compatible struct types
		*(*unsafe.Pointer)(unsafe.Pointer(&s.ptr.offsets)) = unsafe.Pointer(value.ptr)
	} else {
		*(*unsafe.Pointer)(unsafe.Pointer(&s.ptr.offsets)) = nil
	}
}

// ToAVHWFramesContext converts an unsafe.Pointer (typically from AVBufferRef.Data())
// to an AVHWFramesContext wrapper. This is needed for configuring hardware frames
// contexts returned by AVHWFrameCtxAlloc().
func ToAVHWFramesContext(ptr unsafe.Pointer) *AVHWFramesContext {
	if ptr == nil {
		return nil
	}
	return &AVHWFramesContext{ptr: (*C.AVHWFramesContext)(ptr)}
}
