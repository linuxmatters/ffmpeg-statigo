package ffmpeg

import (
	"testing"
	"unsafe"
)

// fifoRoundTrip writes a known per-channel byte pattern through an AVAudioFifo
// and asserts the bytes read back match exactly, with balanced frees under -race.
func fifoRoundTrip(t *testing.T, sampleFmt AVSampleFormat, nbChannels, nbSamples int) {
	t.Helper()

	const align = 0

	fifo := AVAudioFifoAlloc(sampleFmt, nbChannels, nbSamples)
	if fifo == nil {
		t.Fatal("AVAudioFifoAlloc returned nil")
	}
	defer AVAudioFifoFree(fifo)

	src, srcLinesize, srcSize, err := AVSamplesAlloc(nbChannels, nbSamples, sampleFmt, align)
	if err != nil {
		t.Fatalf("AVSamplesAlloc src failed: %v", err)
	}
	defer AVSamplesFreePlanes(src)
	if srcSize <= 0 || srcLinesize <= 0 {
		t.Fatalf("non-positive src size/linesize: %d/%d", srcSize, srcLinesize)
	}

	nbPlanes := samplePlaneCount(nbChannels, sampleFmt)
	if len(src) != nbPlanes {
		t.Fatalf("expected %d planes, got %d", nbPlanes, len(src))
	}

	// Fill each source plane with a distinct, position-dependent pattern so a
	// dropped or swapped plane changes the read-back bytes.
	want := make([][]byte, nbPlanes)
	for p := 0; p < nbPlanes; p++ {
		plane := unsafe.Slice((*byte)(src[p]), srcLinesize)
		want[p] = make([]byte, srcLinesize)
		for i := range plane {
			b := byte((p*31 + i*7 + 1) & 0xff)
			plane[i] = b
			want[p][i] = b
		}
	}

	ret, err := AVAudioFifoWrite(fifo, src, nbSamples, nbChannels, sampleFmt)
	if err != nil {
		t.Fatalf("AVAudioFifoWrite failed: %v (ret=%d)", err, ret)
	}
	if ret != nbSamples {
		t.Fatalf("AVAudioFifoWrite wrote %d samples, want %d", ret, nbSamples)
	}

	if size, err := AVAudioFifoSize(fifo); err != nil {
		t.Fatalf("AVAudioFifoSize failed: %v", err)
	} else if size != nbSamples {
		t.Fatalf("AVAudioFifoSize returned %d, want %d", size, nbSamples)
	}

	dst, dstLinesize, dstSize, err := AVSamplesAlloc(nbChannels, nbSamples, sampleFmt, align)
	if err != nil {
		t.Fatalf("AVSamplesAlloc dst failed: %v", err)
	}
	defer AVSamplesFreePlanes(dst)
	if dstSize <= 0 || dstLinesize != srcLinesize {
		t.Fatalf("dst size/linesize mismatch: size=%d linesize=%d (src linesize=%d)", dstSize, dstLinesize, srcLinesize)
	}

	ret, err = AVAudioFifoRead(fifo, dst, nbSamples, nbChannels, sampleFmt)
	if err != nil {
		t.Fatalf("AVAudioFifoRead failed: %v (ret=%d)", err, ret)
	}
	if ret != nbSamples {
		t.Fatalf("AVAudioFifoRead read %d samples, want %d", ret, nbSamples)
	}

	for p := 0; p < nbPlanes; p++ {
		got := unsafe.Slice((*byte)(dst[p]), dstLinesize)
		for i := range got {
			if got[i] != want[p][i] {
				t.Fatalf("plane %d byte %d: read 0x%02x, wrote 0x%02x", p, i, got[i], want[p][i])
			}
		}
	}
}

// TestAVAudioFifoRoundTripPlanar exercises the FIFO data path for a planar
// format (FLTP, 2 channels → 2 planes).
func TestAVAudioFifoRoundTripPlanar(t *testing.T) {
	fifoRoundTrip(t, AVSampleFmtFltp, 2, 1024)
}

// TestAVAudioFifoRoundTripPacked exercises the FIFO data path for a packed
// format (S16 → 1 plane).
func TestAVAudioFifoRoundTripPacked(t *testing.T) {
	fifoRoundTrip(t, AVSampleFmtS16, 2, 1024)
}
