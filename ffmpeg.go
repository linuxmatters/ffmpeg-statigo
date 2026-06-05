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

var AVTimeBaseQ = &AVRational{value: C.AV_TIME_BASE_Q}

var (
	EAgain     = AVError{Code: -C.EAGAIN}
	AVErrorEOF = AVError{Code: AVErrorEofConst}
)

// AVWhitepointCoefficients is a typedef alias for AVCIExy in FFmpeg.
// This represents white point chromaticity coordinates.
type AVWhitepointCoefficients = AVCIExy

// AVAdler is a typedef alias for uint32_t in FFmpeg.
// This represents an Adler-32 checksum value.
type AVAdler = uint32

// AVCRC is a typedef alias for uint32_t in FFmpeg.
// This represents a CRC (Cyclic Redundancy Check) value.
type AVCRC = uint32

// cstrPtr returns the underlying C char pointer from a CStr, or nil if s is nil.
// This helper simplifies nil-safe access to CStr pointers in CGO calls.
func cstrPtr(s *CStr) *C.char {
	if s == nil {
		return nil
	}
	return s.ptr
}

// AVError represents a non-positive return code from FFmpeg.
type AVError struct {
	Code int
}

func (e AVError) Error() string {
	buf := AllocCStr(uint(AVErrorMaxStringSize))
	defer buf.Free()

	_, _ = AVStrerror(e.Code, buf, uint64(AVErrorMaxStringSize))

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

// wrapStaticCStr wraps a C string the caller does not own — e.g. a const char*
// returned by FFmpeg that points at static or struct-owned memory. The returned
// CStr is marked non-freeable so calling Free on it is a safe no-op.
func wrapStaticCStr(ptr *C.char) *CStr {
	cs := wrapCStr(ptr)
	if cs != nil {
		cs.dontFree = true
	}
	return cs
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
	s.ptr = nil
}

// RawPtr returns a raw reference to the underlying allocation.
func (s *CStr) RawPtr() unsafe.Pointer {
	return unsafe.Pointer(s.ptr)
}
