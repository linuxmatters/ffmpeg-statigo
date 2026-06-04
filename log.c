// CGO bridge for FFmpeg logging callbacks.
// This file is compiled inline with log.go via CGO.
// It bridges FFmpeg's C logging system to Go callbacks.

#include <libavutil/bprint.h>
#include <libavutil/avutil.h>
#include <libavutil/mem.h>

void ffgLogCallback(void* ctx, int level, void* msg);

void ffg_log_callback(void* avcl, int level, const char* fmt, va_list vl) {
    // Respect log level to save formatting cost
    if (level >= 0) {
        level &= 0xff;
    }

    if (level > av_log_get_level())
        return;

    AVBPrint msg;
    char *msg_buf;

    av_bprint_init(&msg, 0, AV_BPRINT_SIZE_UNLIMITED);
	av_vbprintf(&msg, fmt, vl);
    av_bprint_finalize(&msg, &msg_buf);

    ffgLogCallback(avcl, level, msg_buf);

    // msg_buf is allocated by av_bprint_finalize via FFmpeg's allocator, so it
    // must be freed with av_freep, not libc free. The Go callback copies the
    // string out before returning, so freeing here is safe.
    av_freep(&msg_buf);
}

void ffg_set_log() {
    av_log_set_callback(ffg_log_callback);
}

typedef const char* (*itemNameFunc) (void* ctx);

const char* invokeItemNameFunc(itemNameFunc f, void* ctx)
{
    return f(ctx);
}
