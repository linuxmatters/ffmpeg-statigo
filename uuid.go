package ffmpeg

/*
#include <libavutil/uuid.h>
*/
import "C"

import (
	"unsafe"
)

// AVUUID is a typedef to a 16-byte array in FFmpeg (uint8_t[16]).
// This represents a UUID as an opaque sequence of 16 unsigned bytes.
// Binary representation of a UUID per IETF RFC 4122.
type AVUUID = [16]uint8

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
//	@return        true if uu1 and uu2 are equal, false otherwise.
func AVUuidEqual(uu1 *AVUUID, uu2 *AVUUID) bool {
	return C.av_uuid_equal((*C.uint8_t)(unsafe.Pointer(&uu1[0])), (*C.uint8_t)(unsafe.Pointer(&uu2[0]))) != 0
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
