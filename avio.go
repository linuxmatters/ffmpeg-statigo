package ffmpeg

import (
	"runtime/cgo"
	"unsafe"
)

/*
#include <stdint.h>

#include <libavformat/avio.h>
#include <libavutil/buffer.h>
#include <libavutil/mem.h>

AVIOContext* ffg_avio_alloc(unsigned char* buffer, int buffer_size, int write_flag,
                            uintptr_t opaque, int want_read, int want_write, int want_seek);
void ffg_avio_free(AVIOContext* ctx);
AVBufferRef* ffg_avbuffer_create(uint8_t* data, size_t size, uintptr_t opaque, int flags);
*/
import "C"

// AVFmtFlagCustomIo mirrors AVFMT_FLAG_CUSTOM_IO. The macro is not generated, so
// it is defined here for callers wiring a custom AVIOContext into an
// AVFormatContext via SetFlags: it tells the (de)muxer the caller supplied the
// AVIOContext and must not avio_close it.
const AVFmtFlagCustomIo = 0x0080

func FFIOWFourCC(s *AVIOContext, a uint8, b uint8, c uint8, d uint8) {
	AVIOWl32(s, uint(a)|(uint(b)<<8)|(uint(c)<<16)|(uint(d)<<24))
}

func boolToCInt(b bool) C.int {
	if b {
		return 1
	}

	return 0
}

// AVIOReadFunc refills buf from the underlying source. It returns the number of
// bytes read, 0 or [AVErrorEofConst] at end of stream, or a negative
// [AVError]-style code on failure. buf aliases C memory valid only for the call.
type AVIOReadFunc func(buf []byte) int

// AVIOWriteFunc writes buf to the underlying sink. It returns the number of
// bytes consumed or a negative [AVError]-style code on failure. buf aliases C
// memory valid only for the call.
type AVIOWriteFunc func(buf []byte) int

// AVIOSeekFunc repositions the stream. whence is one of io.Seek* or carries the
// [AVSeekSize] flag, in which case it returns the total stream size. It returns
// the new offset or a negative [AVError]-style code on failure.
type AVIOSeekFunc func(offset int64, whence int) int64

// avioCallbacks holds one AVIOContext's Go closures. A cgo.Handle to this struct
// is passed as the AVIOContext opaque pointer, so each live context dispatches
// to its own closures with no cross-context aliasing.
type avioCallbacks struct {
	read  AVIOReadFunc
	write AVIOWriteFunc
	seek  AVIOSeekFunc
}

// AVIOCustomContext wraps a custom AVIOContext together with the cgo.Handle that
// keeps its Go closures alive. Close frees the context (and its possibly
// reallocated buffer) and deletes the handle. It is not safe to use the context
// after Close.
type AVIOCustomContext struct {
	ctx    *AVIOContext
	handle cgo.Handle
}

// Context returns the underlying AVIOContext. Do not free it directly; use Close.
func (c *AVIOCustomContext) Context() *AVIOContext {
	return c.ctx
}

// Close frees the AVIOContext and deletes the cgo.Handle backing its closures.
// After Close the handle is gone, so no late C callback can recover the
// closures. A second call is a no-op.
func (c *AVIOCustomContext) Close() {
	if c.ctx != nil {
		C.ffg_avio_free(c.ctx.ptr)
		c.ctx = nil
	}

	if c.handle != 0 {
		c.handle.Delete()
		c.handle = 0
	}
}

// AVIOAllocContext wraps avio_alloc_context with Go read/write/seek closures.
//
// It allocates an internal buffer of bufSize bytes with av_malloc and registers
// the closures behind a cgo.Handle passed as the C opaque pointer, so multiple
// custom contexts stay isolated. Any of read, write or seek may be nil to leave
// that operation unsupported; writeFlag selects write mode. The returned wrapper
// must be closed to free the context and delete the handle.
//
// The returned context can be attached to an AVFormatContext with SetPb; set
// AVFMT_FLAG_CUSTOM_IO on the format context so it does not close the context.
func AVIOAllocContext(bufSize int, writeFlag bool, read AVIOReadFunc, write AVIOWriteFunc, seek AVIOSeekFunc) *AVIOCustomContext {
	buffer := C.av_malloc(C.size_t(bufSize))
	if buffer == nil {
		return nil
	}

	cb := &avioCallbacks{read: read, write: write, seek: seek}
	handle := cgo.NewHandle(cb)

	wf := C.int(0)
	if writeFlag {
		wf = 1
	}

	ctx := C.ffg_avio_alloc(
		(*C.uchar)(buffer),
		C.int(bufSize),
		wf,
		C.uintptr_t(handle),
		boolToCInt(read != nil),
		boolToCInt(write != nil),
		boolToCInt(seek != nil),
	)
	if ctx == nil {
		C.av_free(buffer)
		handle.Delete()
		return nil
	}

	return &AVIOCustomContext{
		ctx:    &AVIOContext{ptr: ctx},
		handle: handle,
	}
}

// AVBufferFreeFunc frees a buffer's data. opaque is the value registered with
// [AVBufferCreate]; data is the buffer the callback created the AVBuffer over.
type AVBufferFreeFunc func(opaque unsafe.Pointer, data unsafe.Pointer)

// avBufferFreeReg pairs a free closure with its opaque so the dispatcher can
// recover both from one cgo.Handle.
type avBufferFreeReg struct {
	free   AVBufferFreeFunc
	opaque unsafe.Pointer
}

// AVBufferCustomRef wraps an AVBufferRef created with a Go free callback and the
// cgo.Handle backing it. The handle is deleted by the free dispatcher once
// FFmpeg drops the last reference, so callers must not delete it themselves.
type AVBufferCustomRef struct {
	ref *AVBufferRef
}

// Ref returns the underlying AVBufferRef. Unreference it with the usual
// av_buffer_unref path; the registered Go free callback runs on final release.
func (r *AVBufferCustomRef) Ref() *AVBufferRef {
	return r.ref
}

// AVBufferCreate wraps av_buffer_create with a Go free callback. data must point
// to size bytes that the AVBuffer takes ownership of. The free closure runs when
// FFmpeg releases the last reference; the cgo.Handle backing it is deleted at
// that point. flags is a combination of AV_BUFFER_FLAG_* values.
//
// On allocation failure it returns nil and the caller retains ownership of data.
func AVBufferCreate(data unsafe.Pointer, size int, free AVBufferFreeFunc, opaque unsafe.Pointer, flags int) *AVBufferCustomRef {
	reg := &avBufferFreeReg{free: free, opaque: opaque}
	handle := cgo.NewHandle(reg)

	ref := C.ffg_avbuffer_create(
		(*C.uint8_t)(data),
		C.size_t(size),
		C.uintptr_t(handle),
		C.int(flags),
	)
	if ref == nil {
		handle.Delete()
		return nil
	}

	return &AVBufferCustomRef{ref: &AVBufferRef{ptr: ref}}
}

//export ffgAVIOReadPacket
func ffgAVIOReadPacket(opaque C.uintptr_t, buf *C.uint8_t, bufSize C.int) C.int {
	cb := cgo.Handle(opaque).Value().(*avioCallbacks)
	if cb.read == nil {
		return C.int(AVErrorEofConst)
	}

	goBuf := unsafe.Slice((*byte)(unsafe.Pointer(buf)), int(bufSize))

	n := cb.read(goBuf)
	if n == 0 {
		// FFmpeg treats a zero-length read as EOF on a custom AVIOContext.
		return C.int(AVErrorEofConst)
	}

	return C.int(n)
}

//export ffgAVIOWritePacket
func ffgAVIOWritePacket(opaque C.uintptr_t, buf *C.uint8_t, bufSize C.int) C.int {
	cb := cgo.Handle(opaque).Value().(*avioCallbacks)
	if cb.write == nil {
		return C.int(AVErrorEofConst)
	}

	goBuf := unsafe.Slice((*byte)(unsafe.Pointer(buf)), int(bufSize))

	return C.int(cb.write(goBuf))
}

//export ffgAVIOSeek
func ffgAVIOSeek(opaque C.uintptr_t, offset C.int64_t, whence C.int) C.int64_t {
	cb := cgo.Handle(opaque).Value().(*avioCallbacks)
	if cb.seek == nil {
		return C.int64_t(AVErrorEofConst)
	}

	return C.int64_t(cb.seek(int64(offset), int(whence)))
}

//export ffgAVBufferFree
func ffgAVBufferFree(opaque C.uintptr_t, data *C.uint8_t) {
	handle := cgo.Handle(opaque)

	reg := handle.Value().(*avBufferFreeReg)
	if reg.free != nil {
		reg.free(reg.opaque, unsafe.Pointer(data))
	}

	handle.Delete()
}
