package ffmpeg

/*
#include <stdint.h>
#include <libavutil/display.h>
#include <libavcodec/exif.h>
*/
import "C"

import "unsafe"

// Hand-written wrappers for the display-matrix and EXIF-orientation helpers the
// generator skips. Both families take a fixed int32_t[9] transformation matrix,
// which CGO cannot pass directly, so the matrix is modelled as a Go *[9]int32
// and handed to C as a pointer to its first element (the same pattern uuid.go
// uses for its 16-byte arrays). The display transformation matrix is a 3x3
// affine matrix stored as a 9-element array, compatible with the matrices stored
// in the ISO/IEC 14496-12 container format.

// AVDisplayMatrix is the fixed int32_t[9] display transformation matrix used by
// FFmpeg, stored row-major as (a, b, u, c, d, v, x, y, w).
type AVDisplayMatrix = [9]int32

// AVDisplayRotationGet extracts the rotation component of the transformation matrix.
//
//	@param matrix the transformation matrix
//	@return the angle (in degrees) by which the transformation rotates the frame
//	        counterclockwise. The angle will be in range [-180.0, 180.0],
//	        or NaN if the matrix is singular.
//
//	@note floating point numbers are inherently inexact, so callers are
//	      recommended to round the return value to nearest integer before use.
func AVDisplayRotationGet(matrix *AVDisplayMatrix) float64 {
	return float64(C.av_display_rotation_get((*C.int32_t)(unsafe.Pointer(&matrix[0]))))
}

// AVDisplayRotationSet initialises a transformation matrix describing a pure
// clockwise rotation by the specified angle (in degrees).
//
//	@param[out] matrix a transformation matrix (will be fully overwritten
//	                   by this function)
//	@param angle rotation angle in degrees.
func AVDisplayRotationSet(matrix *AVDisplayMatrix, angle float64) {
	C.av_display_rotation_set((*C.int32_t)(unsafe.Pointer(&matrix[0])), C.double(angle))
}

// AVDisplayMatrixFlip flips the input matrix horizontally and/or vertically.
//
//	@param[in,out] matrix a transformation matrix
//	@param hflip whether the matrix should be flipped horizontally
//	@param vflip whether the matrix should be flipped vertically
func AVDisplayMatrixFlip(matrix *AVDisplayMatrix, hflip, vflip bool) {
	C.av_display_matrix_flip(
		(*C.int32_t)(unsafe.Pointer(&matrix[0])),
		boolToCInt(hflip), boolToCInt(vflip),
	)
}

// AVExifMatrixToOrientation converts a display matrix used by
// AV_FRAME_DATA_DISPLAYMATRIX into an orientation constant used by EXIF's
// orientation tag.
//
// Returns an EXIF orientation between 1 and 8 (inclusive) depending on the
// rotation and flip factors. Returns 0 if the matrix is singular.
func AVExifMatrixToOrientation(matrix *AVDisplayMatrix) int {
	return int(C.av_exif_matrix_to_orientation((*C.int32_t)(unsafe.Pointer(&matrix[0]))))
}

// AVExifOrientationToMatrix converts an orientation constant used by EXIF's
// orientation tag into a display matrix used by AV_FRAME_DATA_DISPLAYMATRIX.
//
// Returns nil on success and an AVError if the orientation is invalid, i.e. not
// between 1 and 8 (inclusive).
func AVExifOrientationToMatrix(matrix *AVDisplayMatrix, orientation int) error {
	ret := C.av_exif_orientation_to_matrix(
		(*C.int32_t)(unsafe.Pointer(&matrix[0])), C.int(orientation),
	)
	return WrapErr(int(ret))
}
