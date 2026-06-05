package ffmpeg_test

import (
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/assert"
)

// TestStreamGroupTileGridOffsetRoundTrip exercises the unsafe load/store casts
// in ToAVStreamGroupTileGridOffsetArray and the offset field accessors. It
// allocates a C buffer, writes distinct values through fresh wrappers, then reads
// them back via new wrappers over the same memory to prove the casts round-trip
// and that distinct indices do not alias.
func TestStreamGroupTileGridOffsetRoundTrip(t *testing.T) {
	const bufSize = 256

	buf := ffmpeg.AVMallocz(bufSize)
	if buf == nil {
		t.Fatal("AVMallocz returned nil")
	}
	defer ffmpeg.AVFree(buf)

	arr := ffmpeg.ToAVStreamGroupTileGridOffsetArray(buf)
	if arr == nil {
		t.Fatal("ToAVStreamGroupTileGridOffsetArray returned nil")
	}

	o0 := arr.Get(0)
	o0.SetIdx(7)
	o0.SetHorizontal(100)
	o0.SetVertical(200)

	o1 := arr.Get(1)
	o1.SetIdx(11)
	o1.SetHorizontal(-300)
	o1.SetVertical(-400)

	got0 := arr.Get(0)
	assert.Equal(t, uint(7), got0.Idx())
	assert.Equal(t, 100, got0.Horizontal())
	assert.Equal(t, 200, got0.Vertical())

	got1 := arr.Get(1)
	assert.Equal(t, uint(11), got1.Idx())
	assert.Equal(t, -300, got1.Horizontal())
	assert.Equal(t, -400, got1.Vertical())
}

// TestStreamGroupTileGridOffsetArrayNil confirms a nil pointer yields a nil Array.
func TestStreamGroupTileGridOffsetArrayNil(t *testing.T) {
	assert.Nil(t, ffmpeg.ToAVStreamGroupTileGridOffsetArray(nil))
}
