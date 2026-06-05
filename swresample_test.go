package ffmpeg

import (
	"testing"
)

// newStereoLayout allocates a C AVChannelLayout set to the default stereo layout.
// The caller must call freeChannelLayout to release it.
func newStereoLayout(t *testing.T) *AVChannelLayout {
	t.Helper()
	layout := allocChannelLayout()
	if layout == nil {
		t.Fatal("allocChannelLayout returned nil")
	}
	AVChannelLayoutDefault(layout, 2)
	return layout
}

// TestSwrConvert resamples one buffer from 44100 Hz to 48000 Hz and frees the
// context, asserting a balanced free. Run under -race to exercise cgocheck on the
// plane-pointer arrays passed to swr_convert.
func TestSwrConvert(t *testing.T) {
	inLayout := newStereoLayout(t)
	defer freeChannelLayout(inLayout)
	outLayout := newStereoLayout(t)
	defer freeChannelLayout(outLayout)

	const (
		inRate     = 44100
		outRate    = 48000
		inSamples  = 1024
		nbChannels = 2
	)

	var swr *SwrContext
	ret, err := SwrAllocSetOpts2(
		&swr,
		outLayout, AVSampleFmtFltp, outRate,
		inLayout, AVSampleFmtFltp, inRate,
		0, nil,
	)
	if err != nil {
		t.Fatalf("SwrAllocSetOpts2 failed: %v (ret=%d)", err, ret)
	}
	if swr == nil {
		t.Fatal("SwrAllocSetOpts2 returned nil context")
	}
	defer SwrFree(&swr)

	if ret, err := SwrInit(swr); err != nil {
		t.Fatalf("SwrInit failed: %v (ret=%d)", err, ret)
	}

	inPlanes, _, _, err := AVSamplesAlloc(nbChannels, inSamples, AVSampleFmtFltp, 0)
	if err != nil {
		t.Fatalf("AVSamplesAlloc in failed: %v", err)
	}
	defer AVSamplesFreePlanes(inPlanes)
	if _, err := AVSamplesSetSilence(inPlanes, 0, inSamples, nbChannels, AVSampleFmtFltp); err != nil {
		t.Fatalf("AVSamplesSetSilence failed: %v", err)
	}

	// Output buffer sized for the upsample ratio plus headroom.
	outSamples := AVRescaleRnd(int64(inSamples), outRate, inRate, AVRoundUp) + 256
	outPlanes, _, _, err := AVSamplesAlloc(nbChannels, int(outSamples), AVSampleFmtFltp, 0)
	if err != nil {
		t.Fatalf("AVSamplesAlloc out failed: %v", err)
	}
	defer AVSamplesFreePlanes(outPlanes)

	got, err := SwrConvert(swr, outPlanes, int(outSamples), inPlanes, inSamples)
	if err != nil {
		t.Fatalf("SwrConvert failed: %v (got=%d)", err, got)
	}
	if got <= 0 {
		t.Fatalf("SwrConvert produced %d samples, want > 0", got)
	}
}

// TestSwrSetMatrix builds a stereo→stereo rematrix matrix and applies it,
// exercising SwrBuildMatrix2 and SwrSetMatrix.
func TestSwrSetMatrix(t *testing.T) {
	inLayout := newStereoLayout(t)
	defer freeChannelLayout(inLayout)
	outLayout := newStereoLayout(t)
	defer freeChannelLayout(outLayout)

	var swr *SwrContext
	if ret, err := SwrAllocSetOpts2(
		&swr,
		outLayout, AVSampleFmtFltp, 48000,
		inLayout, AVSampleFmtFltp, 48000,
		0, nil,
	); err != nil {
		t.Fatalf("SwrAllocSetOpts2 failed: %v (ret=%d)", err, ret)
	}
	defer SwrFree(&swr)

	// Stereo matrix is nb_out x nb_in = 2 x 2; stride is the input channel count.
	const stride = 2
	matrix := make([]float64, stride*2)

	ret, err := SwrBuildMatrix2(
		inLayout, outLayout,
		1.0, 1.0, 1.0, 1.0, 0.0,
		matrix, stride,
		AVMatrixEncodingNone, nil,
	)
	if err != nil {
		t.Fatalf("SwrBuildMatrix2 failed: %v (ret=%d)", err, ret)
	}

	if ret, err := SwrSetMatrix(swr, matrix, stride); err != nil {
		t.Fatalf("SwrSetMatrix failed: %v (ret=%d)", err, ret)
	}
}
