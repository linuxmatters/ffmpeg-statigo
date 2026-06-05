package ffmpeg

import (
	"sync"
	"unsafe"
)

/*
#cgo CFLAGS: -I${SRCDIR}/include -Wno-deprecated -Wno-deprecated-declarations

#include <libavcodec/avcodec.h>

void ffg_set_get_format(AVCodecContext *s);
void ffg_clear_get_format(AVCodecContext *s);
*/
import "C"

// GetFormatFunc chooses a decode pixel format from the formats offered by
// libavcodec. It must return one of the offered formats, e.g. a hardware format
// such as [AVPixFmtCuda]; return [AVPixFmtNone] to signal that none is usable.
type GetFormatFunc func(ctx *AVCodecContext, formats []AVPixelFormat) AVPixelFormat

// getFormatReg pairs the Go AVCodecContext wrapper with its callback so the
// dispatcher can hand the wrapper back to the caller's closure.
type getFormatReg struct {
	ctx *AVCodecContext
	fn  GetFormatFunc
}

// getFormatRegistry maps a C AVCodecContext pointer (as uintptr) to its
// registration. Keying on the pointer leaves AVCodecContext.opaque free for
// caller use, unlike the AVIO bridge which owns the opaque field. The context
// is C-allocated, so the Go garbage collector never moves it: the uintptr is a
// stable map key valid until ClearGetFormat removes the entry.
var (
	getFormatMu       sync.RWMutex
	getFormatRegistry = map[uintptr]getFormatReg{}
)

// SetGetFormat installs a Go callback as the context's get_format handler.
// libavcodec invokes it during decoder open with the offered pixel formats, even
// for pure software decode. The caller must call [AVCodecContext.ClearGetFormat]
// before the context is freed so the registry entry and trampoline are removed.
// SetGetFormat is a no-op if s, s.ptr or fn is nil.
func (s *AVCodecContext) SetGetFormat(fn GetFormatFunc) {
	if s == nil || s.ptr == nil || fn == nil {
		return
	}

	key := uintptr(unsafe.Pointer(s.ptr))

	getFormatMu.Lock()
	getFormatRegistry[key] = getFormatReg{ctx: s, fn: fn}
	getFormatMu.Unlock()

	C.ffg_set_get_format(s.ptr)
}

// ClearGetFormat removes the get_format callback and trampoline previously set
// by [AVCodecContext.SetGetFormat]. It is idempotent and a no-op if s or s.ptr
// is nil.
func (s *AVCodecContext) ClearGetFormat() {
	if s == nil || s.ptr == nil {
		return
	}

	key := uintptr(unsafe.Pointer(s.ptr))

	getFormatMu.Lock()
	delete(getFormatRegistry, key)
	getFormatMu.Unlock()

	C.ffg_clear_get_format(s.ptr)
}

//export ffgGetFormat
func ffgGetFormat(s *C.AVCodecContext, fmt *C.enum_AVPixelFormat) C.enum_AVPixelFormat {
	key := uintptr(unsafe.Pointer(s))

	getFormatMu.RLock()
	reg, ok := getFormatRegistry[key]
	getFormatMu.RUnlock()

	// No registration: return the first offered format. Never break decode by
	// returning NONE for an unregistered context.
	if !ok || reg.fn == nil {
		return *fmt
	}

	var formats []AVPixelFormat
	for p := fmt; *p != C.AV_PIX_FMT_NONE; p = (*C.enum_AVPixelFormat)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(*p))) {
		formats = append(formats, AVPixelFormat(*p))
	}

	return C.enum_AVPixelFormat(reg.fn(reg.ctx, formats))
}
