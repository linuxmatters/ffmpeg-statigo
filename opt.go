package ffmpeg

import (
	"unsafe"
)

// optPrimitive constrains AVOptSetSlice to fixed-size, pointer-free element
// types. Passing a slice whose elements contain Go pointers to C would violate
// the cgo pointer rules, so the constraint excludes such types at compile time.
type optPrimitive interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 |
		~int64 | ~uint64 | ~int | ~uint | ~uintptr | ~float32 | ~float64 | ~bool
}

// AVOptSetSlice is a helper for storing a slice of primitive data to the named field. The element type is constrained
// to pointer-free primitives so the slice buffer is safe to hand to C.
//
// See AVOptSet for more information.
func AVOptSetSlice[T optPrimitive](obj unsafe.Pointer, name *CStr, val []T, searchFlags int) (int, error) {
	var ty T
	ptr := unsafe.SliceData(val)
	size := unsafe.Sizeof(ty)
	return AVOptSetBin(obj, name, unsafe.Pointer(ptr), int(size)*len(val), searchFlags)
}
