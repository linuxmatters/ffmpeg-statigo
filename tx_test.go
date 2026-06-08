package ffmpeg

import (
	"math"
	"testing"
	"unsafe"
)

// complexFloat mirrors C's AVComplexFloat: two contiguous 32-bit floats with no
// padding. Indexing the C buffer through this layout avoids importing cgo into a
// test file (unsupported here) while matching the on-wire struct av_tx reads and
// writes. Stride for the FFT is therefore sizeof(complexFloat) = 8 bytes.
type complexFloat struct {
	re, im float32
}

// TestAVTxFloatFFTForwardDC drives a forward complex FFT over a constant real
// input and asserts the DC property: all energy lands in bin 0, every other bin
// is near zero. This exercises the av_tx_init → av_tx_fn (via AVTxCall) →
// av_tx_uninit path with caller-owned C buffers.
//
// Buffer alignment: in/out come from AVMalloc, which returns memory aligned to
// FFmpeg's max SIMD alignment (>= 32 bytes). That satisfies av_tx_fn's default
// requirement, so AV_TX_UNALIGNED is not set. The buffers are C-heap allocated,
// so the GC never moves them while their raw pointers are held by C.
func TestAVTxFloatFFTForwardDC(t *testing.T) {
	const (
		length  = 8
		realVal = float32(1.5)
		tol     = float32(1e-4)
	)

	elemSize := int(unsafe.Sizeof(complexFloat{}))
	bufBytes := uint64(elemSize * length) //nolint:gosec // small positive constant

	inBuf := AVMalloc(bufBytes)
	if inBuf == nil {
		t.Fatal("AVMalloc(in) returned nil")
	}
	defer AVFree(inBuf)

	outBuf := AVMalloc(bufBytes)
	if outBuf == nil {
		t.Fatal("AVMalloc(out) returned nil")
	}
	defer AVFree(outBuf)

	// scale is read at init only; FLOAT_FFT ignores it but the pointer must be valid.
	scale := float32(1.0)

	var ctx *AVTXContext
	var fn AVTxFn
	ret, err := AVTxInit(&ctx, &fn, AVTxFloatFFT, 0, length, unsafe.Pointer(&scale), 0)
	if err != nil {
		t.Fatalf("AVTxInit failed: %v (ret=%d)", err, ret)
	}
	if ret != 0 {
		t.Fatalf("AVTxInit returned non-zero: %d", ret)
	}
	if ctx == nil {
		t.Fatal("AVTxInit left context nil on success")
	}
	if fn == nil {
		t.Fatal("AVTxInit left transform function nil on success")
	}
	defer AVTxUninit(&ctx)

	// Fill input: constant real value, zero imaginary, across all samples.
	for i := range length {
		s := (*complexFloat)(unsafe.Add(inBuf, i*elemSize))
		s.re = realVal
		s.im = 0
	}

	// Stride for FFT is the size of one sample (complex float) in bytes.
	AVTxCall(fn, ctx, outBuf, inBuf, elemSize)

	// Constant real input of value V over N samples => bin 0 = N*V, others ~ 0.
	wantDC := realVal * float32(length)
	for i := range length {
		s := (*complexFloat)(unsafe.Add(outBuf, i*elemSize))
		re, im := s.re, s.im
		mag := float32(math.Hypot(float64(re), float64(im)))
		if i == 0 {
			if absf(re-wantDC) > tol || absf(im) > tol {
				t.Fatalf("bin 0 (DC) = (%g, %g), want (%g, 0)", re, im, wantDC)
			}
			continue
		}
		if mag > tol {
			t.Fatalf("bin %d magnitude = %g, want ~0 (re=%g im=%g)", i, mag, re, im)
		}
	}
}

func absf(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}
