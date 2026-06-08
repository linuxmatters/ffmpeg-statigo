// Forward-call invoker for av_tx_fn, the inverse of get_format.c: nothing flows
// from C back into Go. av_tx_init() hands back a C function pointer in *tx; Go
// never casts it. Instead the pointer is passed back into C here and invoked on
// the C side, keeping the race detector's checkptr instrumentation from
// rejecting a C-function-pointer cast in Go.

#include <libavutil/tx.h>

// ffg_tx_call invokes the transform function pointer returned by av_tx_init().
void ffg_tx_call(av_tx_fn fn, AVTXContext *s, void *out, void *in, ptrdiff_t stride) {
    fn(s, out, in, stride);
}
