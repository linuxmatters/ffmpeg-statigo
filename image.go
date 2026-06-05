package ffmpeg

/*
#include <stdint.h>
#include <libavutil/pixfmt.h>
#include <libavutil/imgutils.h>
*/
import "C"

import "unsafe"

// The plane/linesize family in FFmpeg passes image data as fixed C arrays:
// uint8_t *data[4] and int linesize[4]. cgocheck aborts if a Go pointer handed
// to C points at memory that itself holds Go pointers. The plane pointers used
// here always originate from C allocations (av_image_alloc, AVFrame.data) or are
// nil, so a Go-stack-local [4]*C.uint8_t holding only C pointers is safe to pass
// by address: its elements are never Go pointers. We therefore build the [4]
// arrays as Go-stack locals rather than heap C scratch (C.malloc), avoiding free
// bookkeeping while still satisfying the cgo pointer rules. Output-direction
// calls write into these locals; we copy the results back into the caller's Go
// slices after the call returns.

// imagePlanes is the fixed C uint8_t *data[4] array built on the Go stack.
type imagePlanes [4]*C.uint8_t

// imageLinesizes is the fixed C int linesize[4] array built on the Go stack.
type imageLinesizes [4]C.int

// toImagePlanes copies up to four plane pointers from a Go slice into a fixed C
// array. Each element must be a C-allocated pointer or nil, never a Go pointer.
func toImagePlanes(data []unsafe.Pointer) imagePlanes {
	var out imagePlanes
	for i := 0; i < 4 && i < len(data); i++ {
		out[i] = (*C.uint8_t)(data[i])
	}
	return out
}

// toImageLinesizes copies up to four linesizes from a Go slice into a fixed C array.
func toImageLinesizes(linesizes []int) imageLinesizes {
	var out imageLinesizes
	for i := 0; i < 4 && i < len(linesizes); i++ {
		out[i] = C.int(linesizes[i])
	}
	return out
}

// fromImagePlanes copies the four plane pointers from a filled C array back into
// a Go slice. The returned pointers reference C-owned memory.
func fromImagePlanes(planes imagePlanes) []unsafe.Pointer {
	out := make([]unsafe.Pointer, 4)
	for i := range planes {
		out[i] = unsafe.Pointer(planes[i])
	}
	return out
}

// fromImageLinesizes copies the four linesizes from a filled C array into a Go slice.
func fromImageLinesizes(linesizes imageLinesizes) []int {
	out := make([]int, 4)
	for i := range linesizes {
		out[i] = int(linesizes[i])
	}
	return out
}

// AVImageAlloc allocates an image with the given width, height and pixel format,
// and fills planes and linesizes accordingly. The returned planes slice holds
// the four C-allocated plane pointers and linesizes holds the four linesizes.
// The buffer must be freed with AVImageFreePlanes (av_freep on planes[0]).
//
// Returns the size in bytes required for the buffer, or a negative error code.
func AVImageAlloc(width, height int, pixFmt AVPixelFormat, align int) (planes []unsafe.Pointer, linesizes []int, size int, err error) {
	var cPlanes imagePlanes
	var cLinesizes imageLinesizes

	ret := C.av_image_alloc(
		&cPlanes[0], &cLinesizes[0],
		C.int(width), C.int(height),
		C.enum_AVPixelFormat(pixFmt), C.int(align),
	)

	return fromImagePlanes(cPlanes), fromImageLinesizes(cLinesizes), int(ret), WrapErr(int(ret))
}

// AVImageFreePlanes frees an image buffer allocated by AVImageAlloc. It frees the
// single backing allocation referenced by planes[0]; the remaining plane pointers
// point into the same allocation and must not be freed separately.
func AVImageFreePlanes(planes []unsafe.Pointer) {
	if len(planes) == 0 || planes[0] == nil {
		return
	}
	AVFree(planes[0])
}

// AVImageFillArrays fills the plane pointers and linesizes for an image whose data
// lives in the contiguous buffer src, laid out for the given pixel format and
// dimensions. Pass src == nil to compute the required buffer size only.
//
// Returns the size in bytes required for src, or a negative error code, alongside
// the filled planes and linesizes.
func AVImageFillArrays(src unsafe.Pointer, pixFmt AVPixelFormat, width, height, align int) (planes []unsafe.Pointer, linesizes []int, size int, err error) {
	var cPlanes imagePlanes
	var cLinesizes imageLinesizes

	ret := C.av_image_fill_arrays(
		&cPlanes[0], &cLinesizes[0],
		(*C.uint8_t)(src),
		C.enum_AVPixelFormat(pixFmt), C.int(width), C.int(height), C.int(align),
	)

	return fromImagePlanes(cPlanes), fromImageLinesizes(cLinesizes), int(ret), WrapErr(int(ret))
}

// AVImageFillPointers fills the plane pointers for an image whose data lives in the
// contiguous buffer ptr, using the supplied linesizes. Pass ptr == nil to compute
// the required buffer size only.
//
// Returns the size in bytes required for the buffer, or a negative error code,
// alongside the filled planes.
func AVImageFillPointers(ptr unsafe.Pointer, pixFmt AVPixelFormat, height int, linesizes []int) (planes []unsafe.Pointer, size int, err error) {
	var cPlanes imagePlanes
	cLinesizes := toImageLinesizes(linesizes)

	ret := C.av_image_fill_pointers(
		&cPlanes[0],
		C.enum_AVPixelFormat(pixFmt), C.int(height),
		(*C.uint8_t)(ptr),
		&cLinesizes[0],
	)

	return fromImagePlanes(cPlanes), int(ret), WrapErr(int(ret))
}

// AVImageCopy copies an image plane-by-plane from src to dst. The plane pointers
// must reference C-allocated memory (typically from AVImageAlloc or AVFrame.data).
func AVImageCopy(dstData []unsafe.Pointer, dstLinesizes []int, srcData []unsafe.Pointer, srcLinesizes []int, pixFmt AVPixelFormat, width, height int) {
	cDst := toImagePlanes(dstData)
	cDstLine := toImageLinesizes(dstLinesizes)
	cSrc := toImagePlanes(srcData)
	cSrcLine := toImageLinesizes(srcLinesizes)

	C.av_image_copy(
		&cDst[0], &cDstLine[0],
		&cSrc[0], &cSrcLine[0],
		C.enum_AVPixelFormat(pixFmt), C.int(width), C.int(height),
	)
}

// AVImageCopy2 is a wrapper around AVImageCopy provided for parity with the C API.
// It behaves identically to AVImageCopy.
func AVImageCopy2(dstData []unsafe.Pointer, dstLinesizes []int, srcData []unsafe.Pointer, srcLinesizes []int, pixFmt AVPixelFormat, width, height int) {
	cDst := toImagePlanes(dstData)
	cDstLine := toImageLinesizes(dstLinesizes)
	cSrc := toImagePlanes(srcData)
	cSrcLine := toImageLinesizes(srcLinesizes)

	C.av_image_copy2(
		&cDst[0], &cDstLine[0],
		&cSrc[0], &cSrcLine[0],
		C.enum_AVPixelFormat(pixFmt), C.int(width), C.int(height),
	)
}

// AVImageCopyToBuffer copies image data from the source planes into the contiguous
// buffer dst of size dstSize bytes. align is the assumed linesize alignment for dst.
//
// Returns the number of bytes written to dst, or a negative error code.
func AVImageCopyToBuffer(dst unsafe.Pointer, dstSize int, srcData []unsafe.Pointer, srcLinesizes []int, pixFmt AVPixelFormat, width, height, align int) (int, error) {
	cSrc := toImagePlanes(srcData)
	cSrcLine := toImageLinesizes(srcLinesizes)

	ret := C.av_image_copy_to_buffer(
		(*C.uint8_t)(dst), C.int(dstSize),
		&cSrc[0], &cSrcLine[0],
		C.enum_AVPixelFormat(pixFmt), C.int(width), C.int(height), C.int(align),
	)
	return int(ret), WrapErr(int(ret))
}
