package ffmpeg

import (
	"fmt"
	"sync"
	"unsafe"
)

/*
#cgo CFLAGS: -I${SRCDIR}/include -Wno-deprecated -Wno-deprecated-declarations

#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/lib/linux_amd64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/lib/linux_arm64
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib/darwin_amd64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/lib/darwin_arm64

#cgo linux LDFLAGS: -lffmpeg -lm -ldl -lstdc++ -lpthread
#cgo darwin LDFLAGS: -lffmpeg -lstdc++ -lm -framework ApplicationServices -framework CoreVideo -framework CoreMedia -framework VideoToolbox -framework AudioToolbox

#include <errno.h>
#include <stdlib.h>

#include <libavutil/avutil.h>
*/
import "C"

const (
	ptrSize     = uintptr(C.sizeof_size_t)
	intSize     = uintptr(C.sizeof_int)
	int8Size    = uintptr(C.sizeof_int8_t)
	int16Size   = uintptr(C.sizeof_int16_t)
	int64Size   = uintptr(C.sizeof_int64_t)
	float64Size = uintptr(C.sizeof_double)
)

var AVTimeBaseQ = &AVRational{value: C.AV_TIME_BASE_Q}

var (
	EAgain     = AVError{Code: -C.EAGAIN}
	AVErrorEOF = AVError{Code: AVErrorEofConst}
)

// AVError represents a non-positive return code from FFmpeg.
type AVError struct {
	Code int
}

func (e AVError) Error() string {
	buf := AllocCStr(uint(AVErrorMaxStringSize))
	defer buf.Free()

	AVStrerror(e.Code, buf, uint64(AVErrorMaxStringSize))

	return fmt.Sprintf("averror %v: %v", e.Code, buf.String())
}

// WrapErr returns a AVError if the code is less than zero, otherwise nil.
func WrapErr(code int) error {
	if code >= 0 {
		return nil
	}

	return AVError{Code: code}
}

// CStr is a string allocated in the C memory space. You may need to call Free to clean up the string depending on the
// owner and use-case.
type CStr struct {
	ptr      *C.char
	dontFree bool
}

// AllocCStr allocates an empty string with the given length. The buffer will be initialised to 0.
func AllocCStr(len uint) *CStr {
	ptr := (*C.char)(C.calloc(C.ulong(len), C.sizeof_char))

	return &CStr{
		ptr: ptr,
	}
}

// ToCStr allocates a new CStr with the given content. The CStr will not be automatically garbage collected.
func ToCStr(val string) *CStr {
	return &CStr{
		ptr: C.CString(val),
	}
}

var (
	strMap  = map[string]*CStr{}
	strLock = sync.RWMutex{}
)

// GlobalCStr resolves the given string to a CStr. Multiple calls with the same input string will return the same CStr.
// You should not attempt to free the CStr returned. When passing to FFmpeg, you may need to call Dup to create a copy
// if the FFmpeg code expects to take ownership and will likely free the string.
func GlobalCStr(val string) *CStr {
	var (
		ptr *CStr
		ok  bool
	)

	strLock.RLock()
	ptr, ok = strMap[val]
	strLock.RUnlock()

	if ok {
		return ptr
	}

	strLock.Lock()
	defer strLock.Unlock()

	ptr, ok = strMap[val]
	if ok {
		return ptr
	}

	ptr = ToCStr(val)
	ptr.dontFree = true
	strMap[val] = ptr

	return ptr
}

func wrapCStr(ptr *C.char) *CStr {
	if ptr == nil {
		return nil
	}

	return &CStr{
		ptr: ptr,
	}
}

// Dup is a wrapper for AVStrdup.
func (s *CStr) Dup() *CStr {
	return AVStrdup(s)
}

// String converts the CStr to a Go string.
func (s *CStr) String() string {
	return C.GoString(s.ptr)
}

// Free frees the backing memory for this string. You should only call this function if you are the owner of the memory.
func (s *CStr) Free() {
	if s.dontFree {
		return
	}

	C.free(unsafe.Pointer(s.ptr))
}

// RawPtr returns a raw reference to the underlying allocation.
func (s *CStr) RawPtr() unsafe.Pointer {
	return unsafe.Pointer(s.ptr)
}

// Array is a helper utility for accessing arrays of FFmpeg types. You can not directly allocate this type, and you must
// use one of the inbuilt constructors, such as AllocAVCodecIDArray.
//
// Arrays have no inbuilt length, matching the behaviour of C code. Getting or setting an out of bound index will lead
// to undefined behaviour.
type Array[T any] struct {
	ptr      unsafe.Pointer
	elemSize uintptr
	loadPtr  func(pointer unsafe.Pointer) T
	storePtr func(pointer unsafe.Pointer, value T)
}

// Get returns the element at the ith offset.
func (a *Array[T]) Get(i uintptr) T {
	ptr := unsafe.Add(a.ptr, i*a.elemSize)
	return a.loadPtr(ptr)
}

// Set sets the element at the ith offset.
func (a *Array[T]) Set(i uintptr, value T) {
	ptr := unsafe.Add(a.ptr, i*a.elemSize)
	a.storePtr(ptr, value)
}

// Free deallocates the underlying memory. You should only call this if you own the array.
func (a *Array[T]) Free() {
	AVFree(a.ptr)
}

// RawPtr returns a raw handle the underlying allocation.
func (a *Array[T]) RawPtr() unsafe.Pointer {
	return a.ptr
}

// newArray creates an Array wrapper for a C array pointer.
// This is a factory function that encapsulates the common nil-check and Array construction pattern.
func newArray[T any](ptr unsafe.Pointer, elemSize uintptr, load func(unsafe.Pointer) T, store func(unsafe.Pointer, T)) *Array[T] {
	if ptr == nil {
		return nil
	}
	return &Array[T]{
		ptr:      ptr,
		elemSize: elemSize,
		loadPtr:  load,
		storePtr: store,
	}
}

func ToIntArray(ptr unsafe.Pointer) *Array[int] {
	return newArray(ptr, intSize,
		func(p unsafe.Pointer) int { return int(*(*C.int)(p)) },
		func(p unsafe.Pointer, v int) { *(*C.int)(p) = C.int(v) },
	)
}

func ToUintArray(ptr unsafe.Pointer) *Array[uint] {
	return newArray(ptr, intSize,
		func(p unsafe.Pointer) uint { return uint(*(*C.uint)(p)) },
		func(p unsafe.Pointer, v uint) { *(*C.uint)(p) = C.uint(v) },
	)
}

func ToUint8Array(ptr unsafe.Pointer) *Array[uint8] {
	return newArray(ptr, int8Size,
		func(p unsafe.Pointer) uint8 { return uint8(*(*C.uint8_t)(p)) },
		func(p unsafe.Pointer, v uint8) { *(*C.uint8_t)(p) = C.uint8_t(v) },
	)
}

func ToInt8Array(ptr unsafe.Pointer) *Array[int8] {
	return newArray(ptr, int8Size,
		func(p unsafe.Pointer) int8 { return int8(*(*C.int8_t)(p)) },
		func(p unsafe.Pointer, v int8) { *(*C.int8_t)(p) = C.int8_t(v) },
	)
}

func ToUint16Array(ptr unsafe.Pointer) *Array[uint16] {
	return newArray(ptr, int16Size,
		func(p unsafe.Pointer) uint16 { return uint16(*(*C.uint16_t)(p)) },
		func(p unsafe.Pointer, v uint16) { *(*C.uint16_t)(p) = C.uint16_t(v) },
	)
}

func ToInt16Array(ptr unsafe.Pointer) *Array[int16] {
	return newArray(ptr, int16Size,
		func(p unsafe.Pointer) int16 { return int16(*(*C.int16_t)(p)) },
		func(p unsafe.Pointer, v int16) { *(*C.int16_t)(p) = C.int16_t(v) },
	)
}

func ToUint32Array(ptr unsafe.Pointer) *Array[uint32] {
	return newArray(ptr, intSize,
		func(p unsafe.Pointer) uint32 { return uint32(*(*C.uint32_t)(p)) },
		func(p unsafe.Pointer, v uint32) { *(*C.uint32_t)(p) = C.uint32_t(v) },
	)
}

func ToInt32Array(ptr unsafe.Pointer) *Array[int32] {
	return newArray(ptr, intSize,
		func(p unsafe.Pointer) int32 { return int32(*(*C.int32_t)(p)) },
		func(p unsafe.Pointer, v int32) { *(*C.int32_t)(p) = C.int32_t(v) },
	)
}

func ToInt64Array(ptr unsafe.Pointer) *Array[int64] {
	return newArray(ptr, int64Size,
		func(p unsafe.Pointer) int64 { return int64(*(*C.int64_t)(p)) },
		func(p unsafe.Pointer, v int64) { *(*C.int64_t)(p) = C.int64_t(v) },
	)
}

func ToUint64Array(ptr unsafe.Pointer) *Array[uint64] {
	return newArray(ptr, int64Size,
		func(p unsafe.Pointer) uint64 { return uint64(*(*C.uint64_t)(p)) },
		func(p unsafe.Pointer, v uint64) { *(*C.uint64_t)(p) = C.uint64_t(v) },
	)
}

func ToFloat64Array(ptr unsafe.Pointer) *Array[float64] {
	return newArray(ptr, float64Size,
		func(p unsafe.Pointer) float64 { return float64(*(*C.double)(p)) },
		func(p unsafe.Pointer, v float64) { *(*C.double)(p) = C.double(v) },
	)
}

func ToUint8PtrArray(ptr unsafe.Pointer) *Array[unsafe.Pointer] {
	return newArray(ptr, ptrSize,
		func(p unsafe.Pointer) unsafe.Pointer { return *(*unsafe.Pointer)(p) },
		func(p unsafe.Pointer, v unsafe.Pointer) { *(*unsafe.Pointer)(p) = v },
	)
}
