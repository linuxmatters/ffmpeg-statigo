package ffmpeg

/*
#include <stdint.h>
#include <stdlib.h>
#include <libavutil/samplefmt.h>
*/
import "C"

import "unsafe"

// The sample plane family passes audio data as a variable-length uint8_t **
// (one plane per channel for planar formats, a single plane for packed). The
// same cgo pointer rule from image.go applies: the plane pointer array handed to
// C must not hold Go pointers. Sample plane pointers always come from C
// allocations (av_samples_alloc) or AVFrame.data, so we build the pointer array
// in heap C scratch via C.malloc, populate it with C pointers, and free it after
// the call. Heap C scratch (not a Go-stack local) is used here because the array
// length is dynamic (nb_channels), unlike the fixed [4] image arrays.

// samplePointerArray is a C-heap array of uint8_t* of dynamic length. It must be
// freed with free().
type samplePointerArray struct {
	ptr  **C.uint8_t
	n    int
	heap bool
}

// newSamplePointerArray allocates a C-heap uint8_t* array of n entries and copies
// the plane pointers from data into it. Each element must be a C-allocated pointer
// or nil, never a Go pointer. The caller must call free on the result.
func newSamplePointerArray(data []unsafe.Pointer) samplePointerArray {
	n := len(data)
	if n == 0 {
		return samplePointerArray{}
	}
	size := C.size_t(n) * C.size_t(unsafe.Sizeof(uintptr(0)))
	raw := C.malloc(size)
	arr := samplePointerArray{ptr: (**C.uint8_t)(raw), n: n, heap: true}
	slice := unsafe.Slice(arr.ptr, n)
	for i := range data {
		slice[i] = (*C.uint8_t)(data[i])
	}
	return arr
}

// newEmptySamplePointerArray allocates a zeroed C-heap uint8_t* array of n entries
// for the C side to fill. The caller must call free on the result.
func newEmptySamplePointerArray(n int) samplePointerArray {
	if n <= 0 {
		return samplePointerArray{}
	}
	raw := C.calloc(C.size_t(n), C.size_t(unsafe.Sizeof(uintptr(0))))
	return samplePointerArray{ptr: (**C.uint8_t)(raw), n: n, heap: true}
}

// planes copies the plane pointers out of the C array into a Go slice.
func (s samplePointerArray) planes() []unsafe.Pointer {
	if s.ptr == nil || s.n == 0 {
		return nil
	}
	out := make([]unsafe.Pointer, s.n)
	slice := unsafe.Slice(s.ptr, s.n)
	for i := range slice {
		out[i] = unsafe.Pointer(slice[i])
	}
	return out
}

// free releases the C-heap scratch.
func (s samplePointerArray) free() {
	if s.heap && s.ptr != nil {
		C.free(unsafe.Pointer(s.ptr))
	}
}

// AVSamplesAlloc allocates a samples buffer for nbSamples samples and fills the
// plane pointers and linesize accordingly. align of 0 selects the default
// alignment, 1 selects no alignment.
//
// The returned planes hold the C-allocated audio data pointers and must be freed
// with AVSamplesFreePlanes (av_freep on planes[0]).
//
// Returns the size in bytes of the allocated buffer, or a negative error code.
func AVSamplesAlloc(nbChannels, nbSamples int, sampleFmt AVSampleFormat, align int) (planes []unsafe.Pointer, linesize, size int, err error) {
	nbPlanes := samplePlaneCount(nbChannels, sampleFmt)
	arr := newEmptySamplePointerArray(nbPlanes)
	defer arr.free()

	var cLinesize C.int
	ret := C.av_samples_alloc(
		arr.ptr, &cLinesize,
		C.int(nbChannels), C.int(nbSamples),
		C.enum_AVSampleFormat(sampleFmt), C.int(align),
	)
	return arr.planes(), int(cLinesize), int(ret), WrapErr(int(ret))
}

// AVSamplesAllocArrayAndSamples allocates both the plane-pointer array and the
// samples buffer. This mirrors av_samples_alloc_array_and_samples: FFmpeg
// allocates the pointer array itself. The returned planes hold the audio data
// pointers; free the samples buffer with AVSamplesFreePlanes and the pointer
// array with AVFree on the returned arrayPtr.
//
// Returns the size in bytes of the allocated buffer, or a negative error code.
func AVSamplesAllocArrayAndSamples(nbChannels, nbSamples int, sampleFmt AVSampleFormat, align int) (planes []unsafe.Pointer, arrayPtr unsafe.Pointer, linesize, size int, err error) {
	var cAudioData **C.uint8_t
	var cLinesize C.int

	//nolint:gocritic // dupSubExpr is a false positive on the cgo pointer-check expansion
	ret := C.av_samples_alloc_array_and_samples(
		&cAudioData, &cLinesize,
		C.int(nbChannels), C.int(nbSamples),
		C.enum_AVSampleFormat(sampleFmt), C.int(align),
	)
	if ret < 0 || cAudioData == nil {
		return nil, nil, int(cLinesize), int(ret), WrapErr(int(ret))
	}

	nbPlanes := samplePlaneCount(nbChannels, sampleFmt)
	out := make([]unsafe.Pointer, nbPlanes)
	slice := unsafe.Slice(cAudioData, nbPlanes)
	for i := range slice {
		out[i] = unsafe.Pointer(slice[i])
	}
	return out, unsafe.Pointer(cAudioData), int(cLinesize), int(ret), WrapErr(int(ret))
}

// AVSamplesFreePlanes frees a samples buffer allocated by AVSamplesAlloc or
// AVSamplesAllocArrayAndSamples. It frees the single backing allocation
// referenced by planes[0].
func AVSamplesFreePlanes(planes []unsafe.Pointer) {
	if len(planes) == 0 || planes[0] == nil {
		return
	}
	AVFree(planes[0])
}

// AVSamplesFillArrays fills plane pointers and linesize for audio data stored in
// the contiguous buffer buf, laid out for the given channel count, sample count
// and sample format.
//
// Returns the size in bytes required for buf, or a negative error code, alongside
// the filled planes and linesize.
func AVSamplesFillArrays(buf unsafe.Pointer, nbChannels, nbSamples int, sampleFmt AVSampleFormat, align int) (planes []unsafe.Pointer, linesize, size int, err error) {
	nbPlanes := samplePlaneCount(nbChannels, sampleFmt)
	arr := newEmptySamplePointerArray(nbPlanes)
	defer arr.free()

	var cLinesize C.int
	ret := C.av_samples_fill_arrays(
		arr.ptr, &cLinesize,
		(*C.uint8_t)(buf),
		C.int(nbChannels), C.int(nbSamples),
		C.enum_AVSampleFormat(sampleFmt), C.int(align),
	)
	return arr.planes(), int(cLinesize), int(ret), WrapErr(int(ret))
}

// AVSamplesCopy copies nbSamples samples per channel from src to dst, applying the
// given per-channel sample offsets. The plane pointers must reference C-allocated
// memory.
//
// Returns 0 on success or a negative error code.
func AVSamplesCopy(dst, src []unsafe.Pointer, dstOffset, srcOffset, nbSamples, nbChannels int, sampleFmt AVSampleFormat) (int, error) {
	dstArr := newSamplePointerArray(dst)
	defer dstArr.free()
	srcArr := newSamplePointerArray(src)
	defer srcArr.free()

	ret := C.av_samples_copy(
		dstArr.ptr, srcArr.ptr,
		C.int(dstOffset), C.int(srcOffset),
		C.int(nbSamples), C.int(nbChannels),
		C.enum_AVSampleFormat(sampleFmt),
	)
	return int(ret), WrapErr(int(ret))
}

// AVSamplesSetSilence fills audioData with nbSamples of silence per channel,
// starting at the given per-channel sample offset.
//
// Returns 0 on success or a negative error code.
func AVSamplesSetSilence(audioData []unsafe.Pointer, offset, nbSamples, nbChannels int, sampleFmt AVSampleFormat) (int, error) {
	arr := newSamplePointerArray(audioData)
	defer arr.free()

	ret := C.av_samples_set_silence(
		arr.ptr,
		C.int(offset), C.int(nbSamples), C.int(nbChannels),
		C.enum_AVSampleFormat(sampleFmt),
	)
	return int(ret), WrapErr(int(ret))
}

// samplePlaneCount returns the number of data planes for the given channel count
// and sample format: nbChannels planes for planar formats, one plane for packed.
func samplePlaneCount(nbChannels int, sampleFmt AVSampleFormat) int {
	if planar, _ := AVSampleFmtIsPlanar(sampleFmt); planar != 0 {
		return nbChannels
	}
	return 1
}
