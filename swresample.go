package ffmpeg

/*
#include <stdint.h>
#include <stddef.h>
#include <stdlib.h>
#include <libavutil/channel_layout.h>
#include <libswresample/swresample.h>
*/
import "C"

import "unsafe"

// swr_convert passes audio plane pointers as variable-length uint8_t* arrays, the
// same shape as the sample family in samples.go, so we reuse its C-heap pointer
// array scratch. The plane pointers come from C allocations (AVSamplesAlloc /
// AVFrame.data), satisfying the cgo pointer rules.

// SwrConvert converts inCount input samples per channel from in into up to
// outCount output samples per channel written to out. Pass in == nil to flush
// buffered samples. out and in hold the per-plane data pointers.
//
// Returns the number of samples output per channel, or a negative error code.
func SwrConvert(s *SwrContext, out []unsafe.Pointer, outCount int, in []unsafe.Pointer, inCount int) (int, error) {
	if s == nil {
		return 0, WrapErr(AVErrorUnknownConst)
	}

	outArr := newSamplePointerArray(out)
	defer outArr.free()
	inArr := newSamplePointerArray(in)
	defer inArr.free()

	ret := C.swr_convert(
		s.ptr,
		outArr.ptr, C.int(outCount),
		inArr.constPtr(), C.int(inCount),
	)
	return int(ret), WrapErr(int(ret))
}

// SwrBuildMatrix2 computes a rematrixing coefficient matrix mapping inLayout to
// outLayout and writes it into matrix, which must hold at least
// outLayout.nb_channels * stride doubles. stride is the row stride in elements.
// matrixEncoding selects an optional matrix-encoding downmix. logCtx may be nil.
//
// Returns 0 on success or a negative error code.
func SwrBuildMatrix2(inLayout, outLayout *AVChannelLayout, centerMixLevel, surroundMixLevel, lfeMixLevel, maxval, rematrixVolume float64, matrix []float64, stride int, matrixEncoding AVMatrixEncoding, logCtx unsafe.Pointer) (int, error) {
	if len(matrix) == 0 {
		return 0, WrapErr(AVErrorUnknownConst)
	}

	ret := C.swr_build_matrix2(
		(*C.AVChannelLayout)(channelLayoutPtr(inLayout)),
		(*C.AVChannelLayout)(channelLayoutPtr(outLayout)),
		C.double(centerMixLevel), C.double(surroundMixLevel),
		C.double(lfeMixLevel), C.double(maxval),
		C.double(rematrixVolume),
		(*C.double)(unsafe.Pointer(&matrix[0])),
		C.ptrdiff_t(stride),
		C.enum_AVMatrixEncoding(matrixEncoding),
		logCtx,
	)
	return int(ret), WrapErr(int(ret))
}

// SwrSetMatrix sets a custom rematrixing matrix on the context. matrix holds the
// coefficients laid out row-major with the given row stride in elements.
//
// Returns 0 on success or a negative error code.
func SwrSetMatrix(s *SwrContext, matrix []float64, stride int) (int, error) {
	if s == nil {
		return 0, WrapErr(AVErrorUnknownConst)
	}
	var cMatrix *C.double
	if len(matrix) > 0 {
		cMatrix = (*C.double)(unsafe.Pointer(&matrix[0]))
	}

	ret := C.swr_set_matrix(s.ptr, cMatrix, C.int(stride))
	return int(ret), WrapErr(int(ret))
}

// constPtr returns the plane-pointer array typed as const uint8_t *const * for
// swr_convert's input parameter.
func (s samplePointerArray) constPtr() **C.uint8_t {
	return s.ptr
}

// channelLayoutPtr returns the underlying C pointer for an AVChannelLayout, or nil.
func channelLayoutPtr(l *AVChannelLayout) unsafe.Pointer {
	if l == nil {
		return nil
	}
	return l.RawPtr()
}

// allocChannelLayout allocates a zeroed C AVChannelLayout and returns a wrapper
// over it. The caller must release it with AVChannelLayoutUninit followed by
// freeing the wrapper's RawPtr. Used internally to build layouts where no
// exported zero-value constructor exists.
func allocChannelLayout() *AVChannelLayout {
	raw := C.calloc(1, C.sizeof_AVChannelLayout)
	if raw == nil {
		return nil
	}
	return &AVChannelLayout{ptr: (*C.AVChannelLayout)(raw)}
}

// freeChannelLayout uninitialises and frees a layout allocated by
// allocChannelLayout.
func freeChannelLayout(l *AVChannelLayout) {
	if l == nil {
		return
	}
	AVChannelLayoutUninit(l)
	C.free(l.RawPtr())
}
