package ffmpeg_test

import (
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/assert"
)

// TestAVParseRatio covers a plain num:den pair, the undefined "0:0" value, and a
// malformed string that should report an error.
func TestAVParseRatio(t *testing.T) {
	t.Run("num_den_pair", func(t *testing.T) {
		str := ffmpeg.ToCStr("16:9")
		defer str.Free()

		var q ffmpeg.AVRational
		err := ffmpeg.AVParseRatio(&q, str, 1000, 0, nil)
		assert.NoError(t, err)
		assert.Equal(t, 16, q.Num())
		assert.Equal(t, 9, q.Den())
	})

	t.Run("undefined", func(t *testing.T) {
		str := ffmpeg.ToCStr("0:0")
		defer str.Free()

		var q ffmpeg.AVRational
		err := ffmpeg.AVParseRatio(&q, str, 1000, 0, nil)
		assert.NoError(t, err)
		assert.Equal(t, 0, q.Num())
		assert.Equal(t, 0, q.Den())
	})

	t.Run("invalid", func(t *testing.T) {
		str := ffmpeg.ToCStr("not-a-ratio")
		defer str.Free()

		var q ffmpeg.AVRational
		err := ffmpeg.AVParseRatio(&q, str, 1000, 0, nil)
		assert.Error(t, err)
	})
}

// TestAVParseVideoRate covers a numeric rate, a named abbreviation, and a
// malformed string.
func TestAVParseVideoRate(t *testing.T) {
	t.Run("numeric", func(t *testing.T) {
		str := ffmpeg.ToCStr("25")
		defer str.Free()

		var rate ffmpeg.AVRational
		err := ffmpeg.AVParseVideoRate(&rate, str)
		assert.NoError(t, err)
		assert.Equal(t, 25, rate.Num())
		assert.Equal(t, 1, rate.Den())
	})

	t.Run("abbreviation", func(t *testing.T) {
		str := ffmpeg.ToCStr("pal")
		defer str.Free()

		var rate ffmpeg.AVRational
		err := ffmpeg.AVParseVideoRate(&rate, str)
		assert.NoError(t, err)
		assert.Equal(t, 25, rate.Num())
		assert.Equal(t, 1, rate.Den())
	})

	t.Run("invalid", func(t *testing.T) {
		str := ffmpeg.ToCStr("not-a-rate")
		defer str.Free()

		var rate ffmpeg.AVRational
		err := ffmpeg.AVParseVideoRate(&rate, str)
		assert.Error(t, err)
	})
}

// TestAVCodecGetTag2 looks up a known codec in the RIFF video tag table and
// confirms an unmatched lookup reports found == false.
func TestAVCodecGetTag2(t *testing.T) {
	tags := ffmpeg.AVFormatGetRiffVideoTags()
	if tags == nil {
		t.Skip("RIFF video tag table unavailable")
	}

	t.Run("known_codec", func(t *testing.T) {
		tag, found := ffmpeg.AVCodecGetTag2(&tags, ffmpeg.AVCodecIdMpeg4)
		assert.True(t, found, "MPEG-4 should have a RIFF tag")
		assert.NotZero(t, tag)
	})

	t.Run("absent_codec", func(t *testing.T) {
		// Opus is an audio codec, absent from the RIFF video table. A miss walks
		// the table list to its terminator, so this guards the wrapper's
		// internal NULL-termination against reading past the table pointer.
		_, found := ffmpeg.AVCodecGetTag2(&tags, ffmpeg.AVCodecIdOpus)
		assert.False(t, found, "audio codec should not match the video tag table")
	})
}
