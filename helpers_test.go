package ffmpeg_test

import (
	"math"
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/assert"
)

// TestAVRationalString covers AVRational.String for both the defined and the
// undefined (den == 0) branches.
func TestAVRationalString(t *testing.T) {
	t.Run("defined", func(t *testing.T) {
		r := ffmpeg.AVMakeQ(16, 9)
		assert.Equal(t, "16/9 (1)", r.String())
	})

	t.Run("exact_division", func(t *testing.T) {
		r := ffmpeg.AVMakeQ(6, 2)
		assert.Equal(t, "6/2 (3)", r.String())
	})

	t.Run("undefined_zero_denominator", func(t *testing.T) {
		r := ffmpeg.AVMakeQ(1, 0)
		assert.Equal(t, "1/0 (undefined)", r.String())
		assert.Contains(t, r.String(), "undefined")
	})
}

// TestAVRescaleDelta checks the per-packet timestamp rescale against a known
// equivalence: rescaling from a time base to the same time base is a no-op, and
// the last state variable advances by the supplied duration.
func TestAVRescaleDelta(t *testing.T) {
	tb := ffmpeg.AVMakeQ(1, 44100)
	last := int64(ffmpeg.AVNoptsValue)

	out := ffmpeg.AVRescaleDelta(tb, 0, tb, 1024, &last, tb)
	assert.Equal(t, int64(0), out, "first timestamp scales to itself")

	out = ffmpeg.AVRescaleDelta(tb, 1024, tb, 1024, &last, tb)
	assert.Equal(t, int64(1024), out, "second timestamp scales to itself")
	assert.NotEqual(t, int64(ffmpeg.AVNoptsValue), last, "last state variable should be updated")
}

// TestAVSizeMult covers the success path and the overflow guard.
func TestAVSizeMult(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r, err := ffmpeg.AVSizeMult(6, 7)
		assert.NoError(t, err)
		assert.Equal(t, uint(42), r)
	})

	t.Run("zero", func(t *testing.T) {
		r, err := ffmpeg.AVSizeMult(0, math.MaxUint64)
		assert.NoError(t, err)
		assert.Equal(t, uint(0), r)
	})

	t.Run("overflow", func(t *testing.T) {
		_, err := ffmpeg.AVSizeMult(math.MaxUint64, 2)
		assert.Error(t, err)
	})
}
