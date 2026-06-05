package ffmpeg

/*
#include <stddef.h>
#include <stdlib.h>
#include <libavutil/log.h>
#include <libavutil/avstring.h>
#include <libavutil/bprint.h>
#include <libavutil/mem.h>

// CGO cannot call C variadic functions, so each FFmpeg variadic formatter is
// reached through a fixed-arity shim. Every shim hard-codes the literal "%s"
// format and passes the caller's text as the single string argument. The format
// string is therefore always a compile-time constant in C — never derived from
// caller content — so any '%' directives in the text (%n, %s, %x, ...) are
// emitted verbatim and never interpreted. This neutralises format-string
// injection from caller-supplied content.

static void ffg_av_log_s(void *avcl, int level, const char *msg) {
    av_log(avcl, level, "%s", msg);
}

static char *ffg_av_asprintf_s(const char *msg) {
    return av_asprintf("%s", msg);
}

static size_t ffg_av_strlcatf_s(char *dst, size_t size, const char *msg) {
    return av_strlcatf(dst, size, "%s", msg);
}

static void ffg_av_bprintf_s(AVBPrint *buf, const char *msg) {
    av_bprintf(buf, "%s", msg);
}

// Test-support allocators: _test.go files in this repo cannot host inline cgo,
// so heap-allocate/free an AVBPrint here for AVBprintf round-trip testing.
static AVBPrint *ffg_bprint_alloc(void) {
    return (AVBPrint *)calloc(1, sizeof(AVBPrint));
}

static void ffg_bprint_free(AVBPrint *buf) {
    char *unused = NULL;
    av_bprint_finalize(buf, &unused);
    av_free(unused);
    free(buf);
}
*/
import "C"

import "unsafe"

// Format-string injection note.
//
// FFmpeg's variadic logging/formatting functions (av_log, av_asprintf,
// av_strlcatf, av_bprintf) cannot be called from cgo, which has no way to pass
// a C variadic argument list. Each wrapper below calls a fixed-arity C shim
// (see the cgo preamble in this file) that hard-codes the format as "%s" and
// forwards the caller's text as the single string argument:
//
//	av_log(avcl, level, "%s", goText)
//
// The format string handed to C's printf parser is the literal "%s", baked into
// the shim, never caller-supplied text. Caller content is the single %s
// argument, so any '%' directives inside it (%n, %s, %x, ...) are copied out
// verbatim and never interpreted as conversions. This neutralises
// format-string injection: a caller passing "100% done" or a hostile "%n"
// reaches the log callback as plain bytes. Callers that want interpolation use
// fmt.Sprintf themselves and pass the result as a plain string.

// AVLog wraps av_log. It sends msg to the active FFmpeg log callback at the
// given level. ctx is the originating context, or nil for a general message.
// msg is treated as literal text: any '%' it contains is not interpreted as a
// printf directive (see the format-string injection note in this file).
func AVLog(ctx *LogCtx, level int, msg string) {
	var avcl unsafe.Pointer
	if ctx != nil {
		avcl = ctx.RawPtr()
	}

	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))

	C.ffg_av_log_s(avcl, C.int(level), cmsg)
}

// AVAsprintf wraps av_asprintf. It returns a newly allocated CStr holding msg.
// The caller owns the result and must call Free on it. msg is treated as
// literal text (see the format-string injection note in this file). A nil
// return indicates allocation failure.
func AVAsprintf(msg string) *CStr {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))

	cs := wrapCStr(C.ffg_av_asprintf_s(cmsg))
	if cs != nil {
		// av_asprintf allocates with FFmpeg's allocator, so Free must use av_free.
		cs.avFree = true
	}
	return cs
}

// AVStrlcatf wraps av_strlcatf. It appends msg to the C string in dst, never
// writing past size bytes (including the terminating NUL). It returns the
// length of the string dst would hold without truncation, mirroring the C
// return value; a value >= size means truncation occurred. msg is treated as
// literal text (see the format-string injection note in this file).
func AVStrlcatf(dst *CStr, size uint, msg string) uint {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))

	ret := C.ffg_av_strlcatf_s((*C.char)(dst.RawPtr()), C.size_t(size), cmsg)

	return uint(ret)
}

// AVBprintf wraps av_bprintf. It appends msg to the growable buffer buf. msg is
// treated as literal text (see the format-string injection note in this file).
func AVBprintf(buf *AVBPrint, msg string) {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))

	C.ffg_av_bprintf_s((*C.AVBPrint)(buf.RawPtr()), cmsg)
}

// newTestAVBPrint heap-allocates a zeroed AVBPrint for tests. _test.go files in
// this repo cannot host inline cgo, so the allocation lives here.
func newTestAVBPrint() *AVBPrint {
	return &AVBPrint{ptr: C.ffg_bprint_alloc()}
}

// freeTestAVBPrint finalises and frees an AVBPrint created by newTestAVBPrint.
func freeTestAVBPrint(buf *AVBPrint) {
	C.ffg_bprint_free((*C.AVBPrint)(buf.RawPtr()))
}

// av_sscanf is intentionally not wrapped.
//
// Unlike the four functions above, av_sscanf is a *parsing* function, not a
// formatting one. Its variadic arguments are output pointers that receive
// scanned values, so the "%s" indirection used elsewhere in this file does not
// apply: there is no caller text to neutralise, and replacing the format with
// "%s" would change the parse entirely. A sound binding would need per-call
// typed output pointers, which cgo cannot express through a single fixed-arity
// shim. Use Go's strconv or fmt.Sscanf on a string obtained from the relevant
// FFmpeg accessor instead. Left skipped with this documented reason.
