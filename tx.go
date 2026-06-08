package ffmpeg

/*
#cgo CFLAGS: -I${SRCDIR}/include

#include <libavutil/tx.h>

void ffg_tx_call(av_tx_fn fn, AVTXContext *s, void *out, void *in, ptrdiff_t stride);
*/
import "C"

import "unsafe"

// AVTxCall invokes the transform function returned by [AVTxInit] through the C
// shim in tx.c, so the C function pointer is never cast in Go. The caller owns
// buffer sizing, alignment, and stride: out and in must be large enough for the
// transform type and length, and aligned to the alignment av_tx_init() reports
// (unless AV_TX_UNALIGNED was set at init). A wrong size, alignment, or stride is
// a memory-safety bug, not a returned error; the shim cannot validate the buffers
// and will not fail gracefully. Windowing is the caller's responsibility.
func AVTxCall(fn AVTxFn, ctx *AVTXContext, out, in unsafe.Pointer, stride int) {
	C.ffg_tx_call(C.av_tx_fn(fn), ctx.ptr, out, in, C.ptrdiff_t(stride))
}
