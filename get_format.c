// CGO bridge for AVCodecContext.get_format. libavcodec calls get_format with the
// list of offered pixel formats during decoder open; the caller returns the
// chosen one (e.g. a hardware format). This file mirrors avio.c: a C trampoline
// is installed into the function-pointer field so Go never casts a C function
// pointer, and the trampoline forwards to an exported Go dispatcher.

#include <libavcodec/avcodec.h>

// Defined in get_format.go via //export.
enum AVPixelFormat ffgGetFormat(AVCodecContext *s, enum AVPixelFormat *fmt);

// Trampoline installed into AVCodecContext.get_format. Casts away the const on
// fmt (the Go dispatcher only reads it), mirroring avio.c's const handling.
static enum AVPixelFormat ffg_get_format(struct AVCodecContext *s,
                                         const enum AVPixelFormat *fmt) {
    return ffgGetFormat(s, (enum AVPixelFormat *)fmt);
}

// ffg_set_get_format installs the trampoline. Selecting the function pointer
// C-side avoids casting a C function pointer in Go, which the race detector's
// checkptr instrumentation rejects.
void ffg_set_get_format(AVCodecContext *s) { s->get_format = ffg_get_format; }

// ffg_clear_get_format detaches the trampoline.
void ffg_clear_get_format(AVCodecContext *s) { s->get_format = NULL; }
