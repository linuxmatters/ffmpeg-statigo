package ffmpeg_test

import (
	"math"
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/assert"
)

// TestAVDisplayRotationRoundTrip sets a rotation into a matrix and reads it back,
// confirming the set/get pair agree for the cardinal angles.
func TestAVDisplayRotationRoundTrip(t *testing.T) {
	for _, angle := range []float64{0, 90, -90, 180} {
		var m ffmpeg.AVDisplayMatrix
		ffmpeg.AVDisplayRotationSet(&m, angle)

		got := ffmpeg.AVDisplayRotationGet(&m)
		// av_display_rotation_get returns counterclockwise, set is clockwise, so
		// the sign is inverted; 180 and -180 are equivalent.
		want := -angle
		if math.Abs(want) == 180 {
			assert.InDelta(t, 180, math.Abs(got), 0.5, "angle %v", angle)
			continue
		}
		assert.InDelta(t, want, got, 0.5, "angle %v", angle)
	}
}

// TestAVDisplayRotationSetOverwrites confirms the matrix is fully overwritten,
// not merely combined with prior contents.
func TestAVDisplayRotationSetOverwrites(t *testing.T) {
	m := ffmpeg.AVDisplayMatrix{1, 2, 3, 4, 5, 6, 7, 8, 9}
	ffmpeg.AVDisplayRotationSet(&m, 0)
	assert.InDelta(t, 0, ffmpeg.AVDisplayRotationGet(&m), 0.5)
}

// TestAVDisplayMatrixFlip verifies a horizontal flip alters the matrix and a
// second identical flip restores it.
func TestAVDisplayMatrixFlip(t *testing.T) {
	var m ffmpeg.AVDisplayMatrix
	ffmpeg.AVDisplayRotationSet(&m, 0)
	original := m

	ffmpeg.AVDisplayMatrixFlip(&m, true, false)
	assert.NotEqual(t, original, m, "horizontal flip should change the matrix")

	ffmpeg.AVDisplayMatrixFlip(&m, true, false)
	assert.Equal(t, original, m, "two identical flips should cancel out")
}

// TestAVExifOrientationRoundTrip converts each valid EXIF orientation to a
// matrix and back, confirming it survives the round-trip.
func TestAVExifOrientationRoundTrip(t *testing.T) {
	for orientation := 1; orientation <= 8; orientation++ {
		var m ffmpeg.AVDisplayMatrix
		err := ffmpeg.AVExifOrientationToMatrix(&m, orientation)
		assert.NoError(t, err, "orientation %d", orientation)

		got := ffmpeg.AVExifMatrixToOrientation(&m)
		assert.Equal(t, orientation, got, "orientation %d round-trip", orientation)
	}
}

// TestAVExifOrientationToMatrixInvalid confirms out-of-range orientations are
// rejected with an error.
func TestAVExifOrientationToMatrixInvalid(t *testing.T) {
	var m ffmpeg.AVDisplayMatrix
	assert.Error(t, ffmpeg.AVExifOrientationToMatrix(&m, 0))
	assert.Error(t, ffmpeg.AVExifOrientationToMatrix(&m, 9))
}

// TestAVExifMatrixToOrientationSingular confirms a singular (all-zero) matrix
// reports orientation 0.
func TestAVExifMatrixToOrientationSingular(t *testing.T) {
	var m ffmpeg.AVDisplayMatrix
	assert.Equal(t, 0, ffmpeg.AVExifMatrixToOrientation(&m))
}
