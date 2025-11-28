package ffmpeg_test

import (
	"testing"

	"github.com/linuxmatters/ffmpeg-statigo"
	"github.com/stretchr/testify/assert"
)

func TestVersions(t *testing.T) {
	// FFmpeg 8.0.x: libavcodec version 62.11.100 (0x3E0B64 = 4066148)
	assert.Equal(t, 4066148, int(ffmpeg.AVCodecVersion()), "AVCodec version should match expected")
	assert.Equal(t, ffmpeg.LIBAVCodecVersionInt, int(ffmpeg.AVCodecVersion()), "AVCodec version func and const should match")
}
