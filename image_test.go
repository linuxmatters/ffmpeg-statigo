package ffmpeg

import (
	"testing"
	"unsafe"
)

// TestImagePlaneRoundTrip proves the plane/linesize helper converts a
// four-plane set between Go slices and the fixed C [4] arrays without a cgocheck
// abort. The plane pointers come from a C allocation (AVMalloc), mirroring how
// real plane pointers originate from C-owned memory. Run under -race to exercise
// cgocheck on the slice-to-C-array conversion.
func TestImagePlaneRoundTrip(t *testing.T) {
	// Allocate one C-owned backing buffer and derive four plane pointers into it.
	const planeSize = 64
	buf := AVMalloc(uint64(planeSize * 4))
	if buf == nil {
		t.Fatal("AVMalloc returned nil")
	}
	defer AVFree(buf)

	planes := []unsafe.Pointer{
		unsafe.Add(buf, 0),
		unsafe.Add(buf, planeSize),
		unsafe.Add(buf, planeSize*2),
		unsafe.Add(buf, planeSize*3),
	}
	linesizes := []int{10, 20, 30, 40}

	cPlanes := toImagePlanes(planes)
	cLinesizes := toImageLinesizes(linesizes)

	gotPlanes := fromImagePlanes(cPlanes)
	gotLinesizes := fromImageLinesizes(cLinesizes)

	for i := range planes {
		if gotPlanes[i] != planes[i] {
			t.Errorf("plane %d round-trip mismatch: got %p, want %p", i, gotPlanes[i], planes[i])
		}
		if gotLinesizes[i] != linesizes[i] {
			t.Errorf("linesize %d round-trip mismatch: got %d, want %d", i, gotLinesizes[i], linesizes[i])
		}
	}
}

// TestAVImageAllocRoundTrip exercises the alloc → fill → free path,
// asserting dimensions and a balanced free.
func TestAVImageAllocRoundTrip(t *testing.T) {
	const (
		width  = 320
		height = 240
		align  = 32
	)

	planes, linesizes, size, err := AVImageAlloc(width, height, AVPixFmtYuv420P, align)
	if err != nil {
		t.Fatalf("AVImageAlloc failed: %v", err)
	}
	if size <= 0 {
		t.Fatalf("AVImageAlloc returned non-positive size: %d", size)
	}
	if planes[0] == nil {
		t.Fatal("AVImageAlloc returned nil first plane")
	}

	// YUV420P is planar with three planes; the Y plane stride must be at least the width.
	if linesizes[0] < width {
		t.Errorf("Y linesize %d should be >= width %d", linesizes[0], width)
	}
	// Chroma planes are half-width.
	if linesizes[1] < width/2 {
		t.Errorf("U linesize %d should be >= width/2 %d", linesizes[1], width/2)
	}
	if planes[1] == nil || planes[2] == nil {
		t.Error("YUV420P should have three non-nil planes")
	}

	AVImageFreePlanes(planes)
}

// TestAVImageFillArrays verifies plane/linesize population over a caller-owned
// contiguous buffer and a copy round-trip back to a packed buffer.
func TestAVImageFillArrays(t *testing.T) {
	const (
		width  = 16
		height = 16
		align  = 1
	)

	// Size-query mode: src == nil.
	_, _, bufSize, err := AVImageFillArrays(nil, AVPixFmtRgb24, width, height, align)
	if err != nil {
		t.Fatalf("AVImageFillArrays size query failed: %v", err)
	}
	if bufSize <= 0 {
		t.Fatalf("AVImageFillArrays size query returned non-positive: %d", bufSize)
	}

	src := AVMalloc(uint64(bufSize)) //nolint:gosec // bufSize asserted positive above
	if src == nil {
		t.Fatal("AVMalloc returned nil")
	}
	defer AVFree(src)

	planes, linesizes, _, err := AVImageFillArrays(src, AVPixFmtRgb24, width, height, align)
	if err != nil {
		t.Fatalf("AVImageFillArrays fill failed: %v", err)
	}
	if planes[0] != src {
		t.Errorf("first plane should equal src buffer: got %p, want %p", planes[0], src)
	}
	if linesizes[0] != width*3 {
		t.Errorf("RGB24 linesize should be width*3 = %d, got %d", width*3, linesizes[0])
	}

	// Copy the packed image into a separate buffer via AVImageCopyToBuffer.
	dst := AVMalloc(uint64(bufSize)) //nolint:gosec // bufSize asserted positive above
	if dst == nil {
		t.Fatal("AVMalloc returned nil")
	}
	defer AVFree(dst)

	written, err := AVImageCopyToBuffer(dst, bufSize, planes, linesizes, AVPixFmtRgb24, width, height, align)
	if err != nil {
		t.Fatalf("AVImageCopyToBuffer failed: %v", err)
	}
	if written != bufSize {
		t.Errorf("AVImageCopyToBuffer wrote %d bytes, want %d", written, bufSize)
	}
}

// TestAVImageCopy copies between two allocated images plane-by-plane and verifies
// no cgocheck abort and a balanced free.
func TestAVImageCopy(t *testing.T) {
	const (
		width  = 64
		height = 48
		align  = 16
	)

	srcPlanes, srcLines, _, err := AVImageAlloc(width, height, AVPixFmtYuv420P, align)
	if err != nil {
		t.Fatalf("AVImageAlloc src failed: %v", err)
	}
	defer AVImageFreePlanes(srcPlanes)

	dstPlanes, dstLines, _, err := AVImageAlloc(width, height, AVPixFmtYuv420P, align)
	if err != nil {
		t.Fatalf("AVImageAlloc dst failed: %v", err)
	}
	defer AVImageFreePlanes(dstPlanes)

	AVImageCopy(dstPlanes, dstLines, srcPlanes, srcLines, AVPixFmtYuv420P, width, height)
	AVImageCopy2(dstPlanes, dstLines, srcPlanes, srcLines, AVPixFmtYuv420P, width, height)
}

// TestAVImageFillPointers fills plane pointers over a contiguous buffer using
// caller-supplied linesizes.
func TestAVImageFillPointers(t *testing.T) {
	const (
		width  = 32
		height = 32
		align  = 1
	)

	_, _, bufSize, err := AVImageFillArrays(nil, AVPixFmtYuv420P, width, height, align)
	if err != nil {
		t.Fatalf("size query failed: %v", err)
	}
	if bufSize <= 0 {
		t.Fatalf("size query returned non-positive: %d", bufSize)
	}

	buf := AVMalloc(uint64(bufSize)) //nolint:gosec // bufSize asserted positive above
	if buf == nil {
		t.Fatal("AVMalloc returned nil")
	}
	defer AVFree(buf)

	// YUV420P linesizes: Y = width, U = V = width/2.
	linesizes := []int{width, width / 2, width / 2, 0}
	planes, size, err := AVImageFillPointers(buf, AVPixFmtYuv420P, height, linesizes)
	if err != nil {
		t.Fatalf("AVImageFillPointers failed: %v", err)
	}
	if size <= 0 {
		t.Fatalf("AVImageFillPointers returned non-positive size: %d", size)
	}
	if planes[0] != buf {
		t.Errorf("first plane should equal buffer: got %p, want %p", planes[0], buf)
	}
	if planes[1] == nil || planes[2] == nil {
		t.Error("YUV420P should fill three plane pointers")
	}
}
