package ffmpeg

/*
#cgo CFLAGS: -I${SRCDIR}/include -Wno-deprecated -Wno-deprecated-declarations

#include <libavcodec/avcodec.h>
#include <libavcodec/defs.h>
#include <libavutil/frame.h>
#include <libavutil/pixdesc.h>
#include <libavutil/mastering_display_metadata.h>
*/
import "C"

import "unsafe"

// Manual accessors for struct array fields the generator skips. Each field
// points into struct-owned memory, so every getter wraps the existing pointer
// and never frees it, mirroring the no-free AVStreamGroupTileGrid.Offsets
// convention in streamgroup.go. Getters return nil or a zero value when the backing
// pointer is nil.

// IntraMatrix gets the intra_matrix field from AVCodecContext.
//
// The field is a uint16_t* pointing at a 64-entry custom intra quant matrix, or
// nil when unset. The returned Array aliases struct-owned memory; do not free it.
func (s *AVCodecContext) IntraMatrix() *Array[uint16] {
	if s.ptr.intra_matrix == nil {
		return nil
	}
	return ToUint16Array(unsafe.Pointer(s.ptr.intra_matrix))
}

// InterMatrix gets the inter_matrix field from AVCodecContext.
//
// The field is a uint16_t* pointing at a 64-entry custom inter quant matrix, or
// nil when unset. The returned Array aliases struct-owned memory; do not free it.
func (s *AVCodecContext) InterMatrix() *Array[uint16] {
	if s.ptr.inter_matrix == nil {
		return nil
	}
	return ToUint16Array(unsafe.Pointer(s.ptr.inter_matrix))
}

// ChromaIntraMatrix gets the chroma_intra_matrix field from AVCodecContext.
//
// The field is a uint16_t* pointing at a 64-entry custom chroma intra quant
// matrix, or nil when unset. The returned Array aliases struct-owned memory; do
// not free it.
func (s *AVCodecContext) ChromaIntraMatrix() *Array[uint16] {
	if s.ptr.chroma_intra_matrix == nil {
		return nil
	}
	return ToUint16Array(unsafe.Pointer(s.ptr.chroma_intra_matrix))
}

// ExtendedData gets the extended_data field from AVFrame.
//
// The field is a uint8_t** array of plane pointers; for planar audio it holds
// one entry per channel, otherwise it aliases data. Returns nil when unset. The
// returned Array aliases struct-owned memory; do not free it.
func (s *AVFrame) ExtendedData() *Array[unsafe.Pointer] {
	if s.ptr.extended_data == nil {
		return nil
	}
	return ToUint8PtrArray(unsafe.Pointer(s.ptr.extended_data))
}

// Comp gets the component at index i from the AVPixFmtDescriptor comp field.
//
// The field is a fixed AVComponentDescriptor comp[4] inline array describing how
// each of up to four pixel components is laid out. Only the first NbComponents
// entries are meaningful. Returns nil when i is out of the [0,4) range. The
// returned descriptor aliases struct-owned memory; do not free it.
func (s *AVPixFmtDescriptor) Comp(i int) *AVComponentDescriptor {
	if i < 0 || i >= 4 {
		return nil
	}
	return &AVComponentDescriptor{ptr: &s.ptr.comp[i]}
}

// DisplayPrimaries gets row i of the display_primaries field from
// AVMasteringDisplayMetadata.
//
// The field is a fixed AVRational display_primaries[3][2] array holding the CIE
// 1931 xy chromaticity coordinates of the RGB display primaries: i selects the
// colour (0=R, 1=G, 2=B), and the returned two-element Array holds the x and y
// coordinates. Returns nil when i is out of the [0,3) range. The returned Array
// aliases struct-owned memory; do not free it.
func (s *AVMasteringDisplayMetadata) DisplayPrimaries(i int) *Array[*AVRational] {
	if i < 0 || i >= 3 {
		return nil
	}
	return ToAVRationalArray(unsafe.Pointer(&s.ptr.display_primaries[i][0]))
}

// Position gets row i of the position field from AVPanScan.
//
// The field is a fixed int16_t position[3][2] array of top-left corner
// coordinates in 1/16 pel for up to three frame areas: i selects the area, and
// the returned two-element Array holds the x and y coordinates. Returns nil when
// i is out of the [0,3) range. The returned Array aliases struct-owned memory;
// do not free it.
func (s *AVPanScan) Position(i int) *Array[int16] {
	if i < 0 || i >= 3 {
		return nil
	}
	return ToInt16Array(unsafe.Pointer(&s.ptr.position[i][0]))
}

// allocMasteringDisplayMetadata allocates a zeroed C AVMasteringDisplayMetadata
// and returns a wrapper over it, or nil on allocation failure. FFmpeg exposes no
// exported zero-value constructor for this struct; the caller must release it
// with freeMasteringDisplayMetadata.
func allocMasteringDisplayMetadata() *AVMasteringDisplayMetadata {
	raw := C.calloc(1, C.sizeof_AVMasteringDisplayMetadata)
	if raw == nil {
		return nil
	}
	return &AVMasteringDisplayMetadata{ptr: (*C.AVMasteringDisplayMetadata)(raw)}
}

// freeMasteringDisplayMetadata frees a struct allocated by
// allocMasteringDisplayMetadata.
func freeMasteringDisplayMetadata(m *AVMasteringDisplayMetadata) {
	if m == nil {
		return
	}
	C.free(unsafe.Pointer(m.ptr))
}

// allocPanScan allocates a zeroed C AVPanScan and returns a wrapper over it, or
// nil on allocation failure. FFmpeg exposes no exported zero-value constructor
// for this struct; the caller must release it with freePanScan.
func allocPanScan() *AVPanScan {
	raw := C.calloc(1, C.sizeof_AVPanScan)
	if raw == nil {
		return nil
	}
	return &AVPanScan{ptr: (*C.AVPanScan)(raw)}
}

// freePanScan frees a struct allocated by allocPanScan.
func freePanScan(p *AVPanScan) {
	if p == nil {
		return
	}
	C.free(unsafe.Pointer(p.ptr))
}
