// CGO bridge for FFmpeg custom-AVIO read/write/seek callbacks.
// This file is compiled inline with avio_custom.go via CGO. It mirrors the
// trampoline pattern in log.c: each C function below forwards to an exported
// Go symbol, passing the AVIOContext's opaque pointer through unchanged so the
// Go side can recover the per-context closures from its cgo.Handle.

#include <stdint.h>

#include <libavformat/avio.h>
#include <libavutil/buffer.h>
#include <libavutil/mem.h>

int ffgAVIOReadPacket(uintptr_t opaque, uint8_t* buf, int buf_size);
int ffgAVIOWritePacket(uintptr_t opaque, uint8_t* buf, int buf_size);
int64_t ffgAVIOSeek(uintptr_t opaque, int64_t offset, int whence);

// ffg_avio_read forwards to the Go read dispatcher. The opaque void pointer
// carries a cgo.Handle value, passed to Go as a uintptr to keep the race
// detector's checkptr instrumentation from rejecting a uintptr/pointer cast.
int ffg_avio_read(void* opaque, uint8_t* buf, int buf_size) {
    return ffgAVIOReadPacket((uintptr_t)opaque, buf, buf_size);
}

// ffg_avio_write forwards to the Go write dispatcher. FFmpeg declares the
// buffer const; the Go dispatcher only reads it, so casting away const is safe.
int ffg_avio_write(void* opaque, const uint8_t* buf, int buf_size) {
    return ffgAVIOWritePacket((uintptr_t)opaque, (uint8_t*)buf, buf_size);
}

// ffg_avio_seek forwards to the Go seek dispatcher.
int64_t ffg_avio_seek(void* opaque, int64_t offset, int whence) {
    return ffgAVIOSeek((uintptr_t)opaque, offset, whence);
}

// ffg_avio_alloc wraps avio_alloc_context, wiring only the trampolines the
// caller requested. Selecting the function pointers C-side avoids casting C
// function pointers to *[0]byte in Go, which the race detector's checkptr
// instrumentation rejects.
AVIOContext* ffg_avio_alloc(unsigned char* buffer, int buffer_size, int write_flag,
                            uintptr_t opaque, int want_read, int want_write, int want_seek) {
    int (*read_fn)(void*, uint8_t*, int) = want_read ? ffg_avio_read : NULL;
    int (*write_fn)(void*, const uint8_t*, int) = want_write ? ffg_avio_write : NULL;
    int64_t (*seek_fn)(void*, int64_t, int) = want_seek ? ffg_avio_seek : NULL;

    return avio_alloc_context(buffer, buffer_size, write_flag, (void*)opaque,
                              read_fn, write_fn, seek_fn);
}

// ffg_avbuffer_free forwards av_buffer_create's free callback to Go.
void ffg_avbuffer_free(void* opaque, uint8_t* data);

// ffg_avbuffer_create wraps av_buffer_create with the Go free trampoline,
// keeping the function-pointer reference C-side. opaque carries a cgo.Handle.
AVBufferRef* ffg_avbuffer_create(uint8_t* data, size_t size, uintptr_t opaque, int flags) {
    return av_buffer_create(data, size, ffg_avbuffer_free, (void*)opaque, flags);
}

// ffg_avio_free frees a custom AVIOContext. FFmpeg may reallocate the buffer
// passed to avio_alloc_context, so the current ctx->buffer is freed here (where
// the struct field is directly accessible) before the context itself. The Go
// teardown deletes the cgo.Handle separately.
void ffg_avio_free(AVIOContext* ctx) {
    if (ctx == NULL) {
        return;
    }

    av_freep(&ctx->buffer);
    avio_context_free(&ctx);
}

void ffgAVBufferFree(uintptr_t opaque, uint8_t* data);

// ffg_avbuffer_free forwards av_buffer_create's free callback to Go. opaque
// carries a cgo.Handle, passed to Go as a uintptr.
void ffg_avbuffer_free(void* opaque, uint8_t* data) {
    ffgAVBufferFree((uintptr_t)opaque, data);
}
