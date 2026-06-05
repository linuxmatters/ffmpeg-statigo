package ffmpeg_test

import (
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
