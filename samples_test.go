package ffmpeg

import (
	"testing"
)

// TestAVSamplesAllocSilenceRoundTrip exercises the alloc → silence → free
// path for a planar format, asserting a balanced free. Run under -race to exercise
// cgocheck on the plane-pointer array scratch.
func TestAVSamplesAllocSilenceRoundTrip(t *testing.T) {
	const (
		nbChannels = 2
		nbSamples  = 1024
		align      = 0
	)

	planes, linesize, size, err := AVSamplesAlloc(nbChannels, nbSamples, AVSampleFmtFltp, align)
	if err != nil {
		t.Fatalf("AVSamplesAlloc failed: %v", err)
	}
	if size <= 0 {
		t.Fatalf("AVSamplesAlloc returned non-positive size: %d", size)
	}
	if linesize <= 0 {
		t.Fatalf("AVSamplesAlloc returned non-positive linesize: %d", linesize)
	}
	// FLTP is planar: one plane per channel.
	if len(planes) != nbChannels {
		t.Fatalf("planar format should yield %d planes, got %d", nbChannels, len(planes))
	}
	if planes[0] == nil || planes[1] == nil {
		t.Fatal("planar planes must be non-nil")
	}

	if ret, err := AVSamplesSetSilence(planes, 0, nbSamples, nbChannels, AVSampleFmtFltp); err != nil {
		t.Fatalf("AVSamplesSetSilence failed: %v (ret=%d)", err, ret)
	}

	AVSamplesFreePlanes(planes)
}

// TestAVSamplesPacked verifies a packed (interleaved) format yields a single plane.
func TestAVSamplesPacked(t *testing.T) {
	const (
		nbChannels = 2
		nbSamples  = 512
		align      = 0
	)

	planes, _, size, err := AVSamplesAlloc(nbChannels, nbSamples, AVSampleFmtS16, align)
	if err != nil {
		t.Fatalf("AVSamplesAlloc failed: %v", err)
	}
	if size <= 0 {
		t.Fatalf("non-positive size: %d", size)
	}
	if len(planes) != 1 {
		t.Fatalf("packed format should yield 1 plane, got %d", len(planes))
	}

	if _, err := AVSamplesSetSilence(planes, 0, nbSamples, nbChannels, AVSampleFmtS16); err != nil {
		t.Fatalf("AVSamplesSetSilence failed: %v", err)
	}

	AVSamplesFreePlanes(planes)
}

// TestAVSamplesCopy allocates two planar buffers, sets silence on one, and copies
// it to the other, verifying no cgocheck abort and balanced frees.
func TestAVSamplesCopy(t *testing.T) {
	const (
		nbChannels = 2
		nbSamples  = 256
		align      = 0
	)

	src, _, _, err := AVSamplesAlloc(nbChannels, nbSamples, AVSampleFmtFltp, align)
	if err != nil {
		t.Fatalf("AVSamplesAlloc src failed: %v", err)
	}
	defer AVSamplesFreePlanes(src)

	dst, _, _, err := AVSamplesAlloc(nbChannels, nbSamples, AVSampleFmtFltp, align)
	if err != nil {
		t.Fatalf("AVSamplesAlloc dst failed: %v", err)
	}
	defer AVSamplesFreePlanes(dst)

	if _, err := AVSamplesSetSilence(src, 0, nbSamples, nbChannels, AVSampleFmtFltp); err != nil {
		t.Fatalf("AVSamplesSetSilence failed: %v", err)
	}

	ret, err := AVSamplesCopy(dst, src, 0, 0, nbSamples, nbChannels, AVSampleFmtFltp)
	if err != nil {
		t.Fatalf("AVSamplesCopy failed: %v (ret=%d)", err, ret)
	}
}

// TestAVSamplesAllocArrayAndSamples exercises the variant that allocates the
// pointer array itself, asserting a balanced free of both array and buffer.
func TestAVSamplesAllocArrayAndSamples(t *testing.T) {
	const (
		nbChannels = 2
		nbSamples  = 512
		align      = 0
	)

	planes, arrayPtr, _, size, err := AVSamplesAllocArrayAndSamples(nbChannels, nbSamples, AVSampleFmtFltp, align)
	if err != nil {
		t.Fatalf("AVSamplesAllocArrayAndSamples failed: %v", err)
	}
	if size <= 0 {
		t.Fatalf("non-positive size: %d", size)
	}
	if len(planes) != nbChannels {
		t.Fatalf("planar format should yield %d planes, got %d", nbChannels, len(planes))
	}

	// Free the samples buffer, then the FFmpeg-allocated pointer array.
	AVSamplesFreePlanes(planes)
	AVFree(arrayPtr)
}

// TestAVSamplesFillArrays fills plane pointers over a caller-owned buffer.
func TestAVSamplesFillArrays(t *testing.T) {
	const (
		nbChannels = 2
		nbSamples  = 128
		align      = 0
	)

	var linesize int
	bufSize, err := AVSamplesGetBufferSize(&linesize, nbChannels, nbSamples, AVSampleFmtFltp, align)
	if err != nil {
		t.Fatalf("AVSamplesGetBufferSize failed: %v", err)
	}
	if bufSize <= 0 {
		t.Fatalf("AVSamplesGetBufferSize returned non-positive: %d", bufSize)
	}

	buf := AVMalloc(uint64(bufSize)) //nolint:gosec // bufSize asserted positive above
	if buf == nil {
		t.Fatal("AVMalloc returned nil")
	}
	defer AVFree(buf)

	planes, ls, _, err := AVSamplesFillArrays(buf, nbChannels, nbSamples, AVSampleFmtFltp, align)
	if err != nil {
		t.Fatalf("AVSamplesFillArrays failed: %v", err)
	}
	if ls <= 0 {
		t.Fatalf("non-positive linesize: %d", ls)
	}
	if planes[0] != buf {
		t.Errorf("first plane should equal buffer: got %p, want %p", planes[0], buf)
	}
	if len(planes) != nbChannels {
		t.Fatalf("planar format should yield %d planes, got %d", nbChannels, len(planes))
	}
}
