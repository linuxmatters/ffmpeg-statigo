package ffmpeg

/*
#include <stdint.h>
#include <libavutil/pixfmt.h>
#include <libswscale/swscale.h>
*/
import "C"

import "unsafe"

// sws_scale takes plane-pointer and stride arrays. Although the C signature is
// open-ended (uint8_t *const[]), swscale never addresses more than four planes,
// so we marshal the caller's slices into fixed Go-stack [4] arrays the same way
// image.go does. Plane pointers come from C allocations (AVImageAlloc /
// AVFrame.data), so the [4] arrays hold only C pointers and satisfy cgocheck.

// SwsGetContext allocates and returns a scaling context for converting an image
// of size srcW x srcH in srcFormat to dstW x dstH in dstFormat. flags selects the
// scaling algorithm (one of the SWS_* constants). srcFilter and dstFilter may be
// nil. param may be nil or supply algorithm-specific tuning values.
//
// Returns nil on failure. Free the context with SwsFreeContext.
func SwsGetContext(srcW, srcH int, srcFormat AVPixelFormat, dstW, dstH int, dstFormat AVPixelFormat, flags int, srcFilter, dstFilter *SwsFilter, param []float64) *SwsContext {
	var cParam *C.double
	if len(param) > 0 {
		cParam = (*C.double)(unsafe.Pointer(&param[0]))
	}

	ret := C.sws_getContext(
		C.int(srcW), C.int(srcH), C.enum_AVPixelFormat(srcFormat),
		C.int(dstW), C.int(dstH), C.enum_AVPixelFormat(dstFormat),
		C.int(flags),
		(*C.SwsFilter)(swsFilterPtr(srcFilter)),
		(*C.SwsFilter)(swsFilterPtr(dstFilter)),
		cParam,
	)
	if ret == nil {
		return nil
	}
	return &SwsContext{ptr: ret}
}

// SwsGetCachedContext checks whether context already matches the requested
// parameters and reuses it if so, otherwise frees it and allocates a new one. Pass
// nil context to allocate fresh. srcFilter and dstFilter are not checked and are
// assumed unchanged.
//
// Returns nil on failure. Free the returned context with SwsFreeContext.
func SwsGetCachedContext(context *SwsContext, srcW, srcH int, srcFormat AVPixelFormat, dstW, dstH int, dstFormat AVPixelFormat, flags int, srcFilter, dstFilter *SwsFilter, param []float64) *SwsContext {
	var cParam *C.double
	if len(param) > 0 {
		cParam = (*C.double)(unsafe.Pointer(&param[0]))
	}

	var cCtx *C.SwsContext
	if context != nil {
		cCtx = context.ptr
	}

	ret := C.sws_getCachedContext(
		cCtx,
		C.int(srcW), C.int(srcH), C.enum_AVPixelFormat(srcFormat),
		C.int(dstW), C.int(dstH), C.enum_AVPixelFormat(dstFormat),
		C.int(flags),
		(*C.SwsFilter)(swsFilterPtr(srcFilter)),
		(*C.SwsFilter)(swsFilterPtr(dstFilter)),
		cParam,
	)
	if ret == nil {
		return nil
	}
	return &SwsContext{ptr: ret}
}

// SwsScale scales the slice [srcSliceY, srcSliceY+srcSliceH) of the source image
// into dst. srcSlice and dst hold the plane pointers; srcStride and dstStride hold
// the per-plane strides. The plane pointers must reference C-allocated memory
// (typically AVImageAlloc or AVFrame.data).
//
// Returns the height of the output slice, or a negative error code.
func SwsScale(c *SwsContext, srcSlice []unsafe.Pointer, srcStride []int, srcSliceY, srcSliceH int, dst []unsafe.Pointer, dstStride []int) (int, error) {
	if c == nil {
		return 0, WrapErr(AVErrorUnknownConst)
	}
	cSrc := toImagePlanes(srcSlice)
	cSrcStride := toImageLinesizes(srcStride)
	cDst := toImagePlanes(dst)
	cDstStride := toImageLinesizes(dstStride)

	ret := C.sws_scale(
		c.ptr,
		&cSrc[0], &cSrcStride[0],
		C.int(srcSliceY), C.int(srcSliceH),
		&cDst[0], &cDstStride[0],
	)
	return int(ret), WrapErr(int(ret))
}

// SwsGetColorspaceDetails reads the YUV/RGB colorspace conversion parameters from
// the context. invTable and table return copies of the four-element inverse and
// forward colorspace tables.
//
// Returns a negative error code on failure (for example if the context is not
// initialised), non-negative otherwise.
func SwsGetColorspaceDetails(c *SwsContext) (invTable []int, srcRange int, table []int, dstRange, brightness, contrast, saturation int, err error) {
	if c == nil {
		return nil, 0, nil, 0, 0, 0, 0, WrapErr(AVErrorUnknownConst)
	}

	var cInvTable, cTable *C.int
	var cSrcRange, cDstRange, cBrightness, cContrast, cSaturation C.int

	//nolint:gocritic // dupSubExpr is a false positive on the cgo pointer-check expansion
	ret := C.sws_getColorspaceDetails(
		c.ptr,
		&cInvTable, &cSrcRange,
		&cTable, &cDstRange,
		&cBrightness, &cContrast, &cSaturation,
	)
	if int(ret) < 0 {
		return nil, 0, nil, 0, 0, 0, 0, WrapErr(int(ret))
	}

	invTable = copyIntTable(cInvTable, 4)
	table = copyIntTable(cTable, 4)
	return invTable, int(cSrcRange), table, int(cDstRange), int(cBrightness), int(cContrast), int(cSaturation), nil
}

// SwsSetColorspaceDetails sets the YUV/RGB colorspace conversion parameters on the
// context. invTable and table must each hold four elements (the SWS_CS_* derived
// coefficients).
//
// Returns a negative error code on failure, non-negative otherwise.
func SwsSetColorspaceDetails(c *SwsContext, invTable []int, srcRange int, table []int, dstRange, brightness, contrast, saturation int) (int, error) {
	if c == nil {
		return 0, WrapErr(AVErrorUnknownConst)
	}
	cInv := toImageLinesizes(invTable)
	cTable := toImageLinesizes(table)

	ret := C.sws_setColorspaceDetails(
		c.ptr,
		&cInv[0], C.int(srcRange),
		&cTable[0], C.int(dstRange),
		C.int(brightness), C.int(contrast), C.int(saturation),
	)
	return int(ret), WrapErr(int(ret))
}

// copyIntTable copies n ints out of a C int array into a Go slice. The C memory is
// owned by the context and must not be freed.
func copyIntTable(ptr *C.int, n int) []int {
	if ptr == nil {
		return nil
	}
	src := unsafe.Slice(ptr, n)
	out := make([]int, n)
	for i := range src {
		out[i] = int(src[i])
	}
	return out
}

// swsFilterPtr returns the underlying C pointer for a SwsFilter, or nil.
func swsFilterPtr(f *SwsFilter) unsafe.Pointer {
	if f == nil {
		return nil
	}
	return f.RawPtr()
}
