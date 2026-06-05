package ffmpeg

/*
#include <libavformat/avformat.h>

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
	"unsafe"
)

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
