package ffmpeg_test

import (
	"testing"

	ffmpeg "github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/assert"
)

// TestAVAc3ParseHeaderGuards exercises the empty-buffer guard. A valid AC-3
// fixture is impractical to hand-craft, so only the guard path is covered.
func TestAVAc3ParseHeaderGuards(t *testing.T) {
	_, _, err := ffmpeg.AVAc3ParseHeader(nil)
	assert.Error(t, err)

	_, _, err = ffmpeg.AVAc3ParseHeader([]byte{})
	assert.Error(t, err)
}

// validADTSHeader returns a minimal 7-byte ADTS header: syncword 0xFFF,
// MPEG-4, layer 0, no CRC (protection absent), AAC-LC profile, 44100 Hz
// (sampling index 4), stereo (channel config 2), frame length 7 (header only).
func validADTSHeader() []byte {
	return []byte{
		0xFF, // syncword high
		0xF1, // syncword low, MPEG-4, layer 0, protection absent
		0x50, // profile (AAC-LC), sampling freq index 4, private 0, channel high bit 0
		0x80, // channel config 2, original/copy/home/copyright bits 0, frame len bits 0
		0x00, // frame length middle bits
		0xFF, // frame length low (7) + buffer fullness high
		0xFC, // buffer fullness low + 0 raw data blocks
	}
}

// TestAVAdtsHeaderParse covers the short-buffer guard and the success path with
// a hand-crafted valid 7-byte ADTS header.
func TestAVAdtsHeaderParse(t *testing.T) {
	_, _, err := ffmpeg.AVAdtsHeaderParse(nil)
	assert.Error(t, err)

	_, _, err = ffmpeg.AVAdtsHeaderParse(make([]byte, ffmpeg.AVAacAdtsHeaderSize-1))
	assert.Error(t, err)

	samples, frames, err := ffmpeg.AVAdtsHeaderParse(validADTSHeader())
	assert.NoError(t, err)
	assert.Equal(t, uint32(1024), samples)
	assert.Equal(t, uint8(1), frames)
}

// TestAVVorbisParseFrameFlagsGuards exercises the nil-context and empty-buffer
// guards. A valid Vorbis parser context plus frame fixture is impractical to
// hand-craft, so only the guard paths are covered.
func TestAVVorbisParseFrameFlagsGuards(t *testing.T) {
	_, _, err := ffmpeg.AVVorbisParseFrameFlags(nil, []byte{0x00})
	assert.Error(t, err)

	_, _, err = ffmpeg.AVVorbisParseFrameFlags(nil, nil)
	assert.Error(t, err)
}
